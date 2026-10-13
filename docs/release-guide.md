# Release Guide

The `make all` target runs checks and tests, bumps the version, pushes the default branch and version tag, and publishes a release with GoReleaser. It changes Git history and publishes artifacts, so run it only when you intend to release.

## Prerequisites

Install and configure:

- `make`, `git`, `docker`, `go`, `curl`, `bumpversion`, and `goreleaser`.
- A GitHub personal access token with permission to publish releases, exported as `GITHUB_TOKEN`.
- Docker Hub credentials with Read & Write access to `sfarosu/go-tooling-portal`. Authenticate with `docker login -u <dockerhub-username> docker.io`, using a Docker Hub access token with Read & Write permissions as the password.
- On Debian/Ubuntu systems building the ARM64 Docker image, `qemu-user-static` may also be required.

The release prechecks require the `main` branch to be checked out, a clean working tree, and no existing `dist/` directory. Docker must be running, and the GitHub token and Docker login must be valid.

## Publish a release

First review the available targets:

```sh
make help
```

Then, from a clean checkout of `main`, pull the latest changes and run:

```sh
git checkout main
git pull origin main
make all
```

The version bump defaults to `patch`. To select another bump level, pass it as a Make variable, for example `make all BUMP_LEVEL=minor` or `make all BUMP_LEVEL=major`.
