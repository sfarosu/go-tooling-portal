# Developer Guide

## Technology

- [Go](https://go.dev/) and [Huma](https://huma.rocks/) provide the API.
- [HTMX](https://htmx.org/) powers dynamic page interactions.
- [Bootstrap](https://getbootstrap.com/) provides responsive layout and UI components.

## Run locally

Install the Go version required by [`go.mod`](../go.mod) (currently Go 1.27.1), then clone and run the app:

```sh
git clone https://github.com/sfarosu/go-tooling-portal.git
cd go-tooling-portal
go run .
```

Alternatively, build and run the binary:

```sh
go build .
./go-tooling-portal
```

By default, the app listens on port `8100`. Open [http://localhost:8100](http://localhost:8100). To choose another address or port, pass `-addr`, for example `go run . -addr :8080`.

The web pages, stylesheets, scripts, and favicon are embedded in the Go binary, so a built executable can run without a separate `web/` directory.

## Tests and coverage

Run all Go tests with:

```sh
go test ./...
```

To run tests without cached results and write a coverage profile:

```sh
go test -coverprofile=coverage.out ./... -count=1
```

Summarize coverage by function with `go tool cover -func=coverage.out`, or create a browsable report with `go tool cover -html=coverage.out`.

## API documentation

With the app running, open the interactive API documentation at [http://localhost:8100/docs](http://localhost:8100/docs). The OpenAPI specification is available at [http://localhost:8100/openapi.json](http://localhost:8100/openapi.json).

The stateless Streamable HTTP MCP endpoint is available at `http://localhost:8100/mcp`. It exposes tools for the JSON API operations and proxies calls through the app's HTTP API; HTMX endpoints are not exposed as MCP tools.

Prometheus metrics are exposed at `http://localhost:8100/metrics`.

### Test the MCP server with MCP Inspector

With the app running locally, start the [MCP Inspector](https://github.com/modelcontextprotocol/inspector) in another terminal:

```sh
HOST=127.0.0.1 npx @modelcontextprotocol/inspector
```

Open the Inspector at [http://localhost:6274](http://localhost:6274), select **Streamable HTTP**, and connect to `http://localhost:8100/mcp`. Binding Inspector to `127.0.0.1` keeps its local servers on loopback and avoids relying on a hostname such as `desktop` that may not resolve on the developer's machine.
