# Repository Guide

## Project
Go Tooling Portal is a self-hosted web app with developer utilities and an HTTP API. It uses Go and Huma for the API, HTMX for interactive pages, and Bootstrap for the UI. Web pages and assets are embedded in the Go binary.

## Layout
- `internal/service/`: tool implementations
- `internal/apis/`: HTTP and HTMX handlers
- `internal/httpserver/`: HTTP routing, middleware, and server setup
- `internal/mcpserver/`: MCP tools that proxy to the JSON API
- `internal/util/generalutil/`: shared general-purpose utilities
- `web/`: HTML, CSS, JavaScript, and embed declaration
- `docs/`: developer and release guides

## Development
- Use the Go version declared in `go.mod`.
- Run locally with `go run .`; the default address is `:8100`.
- Run relevant package tests while iterating, then run all tests with `go test ./...`.
- Add or update focused tests when changing behavior.

See `docs/developer-guide.md` for setup, testing, and API documentation.

## Safety
- Never commit secrets, credentials, tokens, or real user data. Do not print secrets in logs or command output.
- Inspect the worktree before editing and preserve changes you did not make.
- Do not create, delete, move, or push Git tags, or push commits, unless explicitly asked.
- Do not run `make all`, `make gitpushtags`, or `goreleaser release` unless explicitly asked to perform a release. These can change Git refs or publish artifacts.
- Do not run destructive Docker, Kubernetes, or filesystem commands unless their exact targets and effects are clear and authorized.
- Keep changes focused; change dependencies, release configuration, or deployment manifests only when required for the task.

See `docs/release-guide.md` before performing release-related work.
