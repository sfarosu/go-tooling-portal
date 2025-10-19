package apis

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// SSHGenHTMXInput is the request structure for the /api/sshgen/htmx endpoint.
type SSHGenHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`

	// Body contains the main input fields for the ssh key generation.
	Body struct {
		Algorithm string `json:"algorithm" example:"ed25519" doc:"SSH key algorithm (ed25519, ecdsa, rsa)" enum:"ed25519,ecdsa,rsa"`
		EcdsaBits string `json:"ecdsaBits,omitempty" example:"256" doc:"ECDSA key size (256, 384, 521)"`
		RsaBits   string `json:"rsaBits,omitempty" example:"2048" doc:"RSA key size (1024, 3072, 4096)"`
		Email     string `json:"email" example:"alice@example.com" doc:"Email/comment for the SSH public key"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *SSHGenHTMXInput) Resolve(ctx huma.Context) []error {
	val := ctx.Header("HX-Request")
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		m.HtmxHeader = false
		return nil
	}
	m.HtmxHeader = parsed
	return nil
}

// SSHGenHTMXOutput is the response structure for the /api/sshgen/htmx endpoint.
type SSHGenHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var sshgenResultTmpl = template.Must(template.New("sshgen-result").Parse(`
<div class="row mb-3">
  <div class="col-md-6 d-flex flex-column align-items-stretch">
    <div class="d-flex align-items-center mb-2 justify-content-center">
      <span class="form-label fw-bold text-center">Private Key</span>
    </div>
    <textarea class="form-control custom-output flex-grow-1" id="sshgen-privatekey" rows="8" readonly style="resize: none; min-height: 180px;">{{.PrivateKey}}</textarea>
    <div class="d-flex justify-content-center mt-2">
      <button class="btn btn-graphite" type="button" onclick="copyToClipboard('sshgen-privatekey')" title="Copy Private Key">
        <i class="bi bi-clipboard"></i>
      </button>
    </div>
	<div class="usage-tip mt-3">
      <i class="bi"></i>
      Place the private key under "~/.ssh/" folder and make sure to run "chmod 400" on it.
    </div>
  </div>
  <div class="col-md-6 d-flex flex-column align-items-stretch">
    <div class="d-flex align-items-center mb-2 justify-content-center">
      <span class="form-label fw-bold text-center">Public Key</span>
    </div>
	<textarea class="form-control custom-output flex-grow-1" id="sshgen-publickey" rows="8" readonly style="resize: none; min-height: 180px;">{{.PublicKey}}</textarea>
    <div class="d-flex justify-content-center mt-2">
      <button class="btn btn-graphite" type="button" onclick="copyToClipboard('sshgen-publickey')" title="Copy Public Key">
        <i class="bi bi-clipboard"></i>
      </button>
    </div>
	<div class="usage-tip mt-3">
      <i class="bi"></i>
      Paste the public key on the target machine in "~/.ssh/authorized_keys".
    </div>
  </div>
</div>
`))

// RegisterSSHGenHtmx registers the /api/sshgen/htmx endpoint with the given Huma API
func RegisterSSHGenHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-ssh-key-htmx",
		Summary:       "generate ssh key - htmx",
		Description:   "Returns an HTML fragment suitable for HTMX injection containing an SSH key pair for the selected algorithm.",
		Method:        http.MethodPost,
		Path:          "/api/sshgen/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"SSHGen"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully generated an SSH key pair, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request, invalid input (including missing HX-Request header) or unsupported algorithm.",
			},
		},
	}, func(ctx context.Context, input *SSHGenHTMXInput) (*SSHGenHTMXOutput, error) {
		// Validate input; notice we don't validate the ecdsaBits and rsaBits in SSHGenInput using enums as the user will select a single algorithm
		algo := input.Body.Algorithm
		ecdsaBits := input.Body.EcdsaBits
		rsaBits := input.Body.RsaBits

		switch algo {
		case "ecdsa":
			if ecdsaBits == "" {
				return nil, huma.Error400BadRequest("ecdsaBits is required for ECDSA algorithm")
			}
		case "rsa":
			if rsaBits == "" {
				return nil, huma.Error400BadRequest("rsaBits is required for RSA algorithm")
			}
			// ed25519 needs neither
		}

		// Generate the SSH key pair
		privKey, pubKey, err := service.GenerateSSHKeyPair(input.Body.Algorithm, input.Body.EcdsaBits, input.Body.RsaBits, input.Body.Email)
		if err != nil {
			logger.Logger.Error(
				"failed to generate ssh key pair",
				"algorithm", input.Body.Algorithm,
				"ecdsaBits", input.Body.EcdsaBits,
				"rsaBits", input.Body.RsaBits,
				"email", input.Body.Email,
				"error", err,
			)
			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				return &SSHGenHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("failed to generate ssh key pair: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			sshgenResultTmpl.Execute(&buf, map[string]string{
				"PrivateKey": privKey,
				"PublicKey":  pubKey,
			})
			return &SSHGenHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/sshgen for JSON responses")
	})
}
