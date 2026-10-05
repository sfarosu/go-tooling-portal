package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// SSHGenInput is the request structure for the /api/sshgen endpoint.
type SSHGenInput struct {
	// Body contains the main input fields for the ssh key generation.
	Body struct {
		Algorithm string `json:"algorithm" example:"ed25519" doc:"SSH key algorithm (ed25519, ecdsa, rsa)" enum:"ed25519,ecdsa,rsa"`
		EcdsaBits string `json:"ecdsaBits" example:"256" doc:"ECDSA key size (256, 384, 521)"`
		RsaBits   string `json:"rsaBits" example:"2048" doc:"RSA key size (1024, 3072, 4096)"`
		Email     string `json:"email" example:"alice@example.com" doc:"Email/comment for the SSH public key"`
	}
}

// SSHGenOutput is the response structure for the /api/sshgen endpoint.
type SSHGenOutput struct {
	Body struct {
		PrivateKey string `json:"privateKey"`
		PublicKey  string `json:"publicKey"`
	}
}

// RegisterSSHGen registers the /api/sshgen endpoint with the given Huma API
func RegisterSSHGen(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-ssh-key",
		Summary:       "generate ssh key - json",
		Description:   "Returns a JSON object containg a ssh key-pair depending on the selected algorithm.",
		Method:        http.MethodPost,
		Path:          "/api/sshgen",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"SSHGen"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully generated SSH key pair, responding with JSON content.",
			},
			"400": {
				Description: "Bad request - invalid input.",
			},
		},
	}, func(ctx context.Context, input *SSHGenInput) (*SSHGenOutput, error) {
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

		// Generate SSH key pair
		privKey, pubKey, err := service.GenerateSSHKeyPair(input.Body.Algorithm, input.Body.EcdsaBits, input.Body.RsaBits, input.Body.Email)
		if err != nil {
			logger.Logger.Error(
				"failed to generate SSH key pair",
				"algorithm", input.Body.Algorithm,
				"ecdsaBits", input.Body.EcdsaBits,
				"rsaBits", input.Body.RsaBits,
				"email", input.Body.Email,
				"error", err,
			)
			return nil, huma.Error400BadRequest("SSH key generation failed: " + err.Error())
		}

		resp := &SSHGenOutput{}
		resp.Body.PrivateKey = privKey
		resp.Body.PublicKey = pubKey
		return resp, nil
	})
}
