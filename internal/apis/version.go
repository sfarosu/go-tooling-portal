package apis

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

type VersionOutput struct {
	Version string `json:"version"`
}

func RegisterVersion(api huma.API) {
	huma.Get(api, "/api/version", func(ctx context.Context, input *struct{}) (*VersionOutput, error) {
		return &VersionOutput{Version: version.Version}, nil
	})
}
