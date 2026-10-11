package apis

import (
	"bytes"
	"html/template"
)

func renderHTMXTemplate(tmpl *template.Template, data any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func generateHTMXError(err error) string {
	content := `
<div class="d-flex justify-content-center">
  <div style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output is-invalid"
	  readonly aria-label="error" rows="2">ERROR: ` + template.HTMLEscapeString(err.Error()) + `</textarea>
  </div>
</div>`
	return content
}
