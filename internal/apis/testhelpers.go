package apis

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// setupTestAPI is a helper function for unit testing to set up a test API with the given registration function.
func setupTestAPI(register func(huma.API)) *http.ServeMux {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Test API", "1.0.0"))
	register(api)
	return router
}
