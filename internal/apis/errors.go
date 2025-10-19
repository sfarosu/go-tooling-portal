package apis

func generateHTMXError(err error) string {
	content := `
<div class="d-flex justify-content-center">
  <div style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output is-invalid"
	  readonly aria-label="error" rows="2">` + "ERROR: " + err.Error() + `</textarea>
  </div>
</div>`
	return content
}
