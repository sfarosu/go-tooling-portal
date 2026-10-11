package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
)

// WrapHTMXValidationErrors converts Huma schema-validation responses into
// fragments that HTMX can display, without changing regular API responses.
func WrapHTMXValidationErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" || !strings.HasPrefix(r.URL.Path, "/api/") || !strings.HasSuffix(r.URL.Path, "/htmx") {
			next.ServeHTTP(w, r)
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)
		if recorder.Code != http.StatusUnprocessableEntity {
			copyResponse(w, recorder)
			return
		}

		var validationError struct {
			Detail string `json:"detail"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &validationError); err != nil || validationError.Detail == "" {
			copyResponse(w, recorder)
			return
		}

		messages := []string{validationError.Detail}
		for _, detail := range validationError.Errors {
			if detail.Message != "" {
				messages = append(messages, detail.Message)
			}
		}

		copyHeaders(w.Header(), recorder.Header())
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(generateHTMXError(errors.New(strings.Join(messages, ": ")))))
	})
}

func copyResponse(w http.ResponseWriter, recorder *httptest.ResponseRecorder) {
	copyHeaders(w.Header(), recorder.Header())
	w.WriteHeader(recorder.Code)
	_, _ = w.Write(recorder.Body.Bytes())
}

func copyHeaders(dst, src http.Header) {
	for name, values := range src {
		dst[name] = append([]string(nil), values...)
	}
}

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
