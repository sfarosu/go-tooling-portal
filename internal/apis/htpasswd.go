package apis

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	// import your service package if you separate business logic
)

type HtpasswdInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type HtpasswdOutput struct {
	Htpasswd string `json:"htpasswd"`
}

func RegisterHtpasswd(api huma.API) {
	huma.Post(api, "/api/htpasswd", func(ctx context.Context, input *HtpasswdInput) (*HtpasswdOutput, error) {
		// Call your business logic here
		htpasswd := input.Username + ":" + input.Password // Replace with real logic!
		return &HtpasswdOutput{Htpasswd: htpasswd}, nil
	})
}
