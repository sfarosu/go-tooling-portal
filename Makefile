.PHONY: prechecks tests bumpversion gitpushtags goreleaser help

BUMP_LEVEL := patch
TOOLS := make git docker bumpversion go goreleaser
DOCKER_REPO := docker.io/sfarosu/go-tooling-portal
DEFAULT_GIT_BRANCH := master

prechecks:
	@echo "#################################"
	@echo "####### Running prechecks #######"

	@for tool in $(TOOLS); do \
		if ! command -v $$tool >/dev/null 2>&1; then \
			echo "Verify if [$$tool] is installed ... ERROR, not installed or not in PATH"; \
			exit 1; \
		else \
			echo "Verify if [$$tool] is installed ...  found it in path [$$(command -v $$tool)]"; \
		fi; \
	done

	@if [ -z "$$GITHUB_TOKEN" ]; then \
		echo "Verify if [GITHUB_TOKEN] env var is set ... ERROR, not found"; \
		exit 1; \
	else \
		echo "Verify if [GITHUB_TOKEN] env var is set ... found it"; \
	fi

	@response_code=$$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: token $$GITHUB_TOKEN" https://api.github.com/user); \
	if [ $$response_code -eq 200 ]; then \
		echo "Verify if [GITHUB_TOKEN is valid] ... success"; \
	else \
		echo "Verify if [GITHUB_TOKEN is valid] ... ERROR, token is invalid"; \
		exit 1; \
	fi

	@if ! goreleaser check >/dev/null 2>&1; then \
		echo "Verify if [.goreleaser.yaml] file is valid ... ERROR, not valid"; \
		exit 1; \
	else \
		echo "Verify if [.goreleaser.yaml file] is valid ... success"; \
	fi

	@if [ -d "dist" ]; then \
		echo "Verify if [dist] folder exists ... ERROR, it must be deleted before running goreleaser"; \
		exit 1; \
	else \
		echo "Verify if [dist] folder exists ... success, does not exist"; \
	fi

	@if ! docker login $(DOCKER_REPO) >/dev/null 2>&1; then \
		echo "Verify docker login ... ERROR, docker login failed"; \
		exit 1; \
	else \
		echo "Verify docker login ... docker login successful"; \
	fi
	
	@if [ $$(git rev-parse --abbrev-ref HEAD) != $(DEFAULT_GIT_BRANCH) ]; then \
		echo "Verify current branch ... ERROR, not on default branch [$(DEFAULT_GIT_BRANCH)]"; \
		exit 1; \
	else \
		echo "Verify current branch ... success, on default branch [$(DEFAULT_GIT_BRANCH)]"; \
	fi

	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Verify git repo state ... ERROR, repo is in dirty state"; \
		exit 1; \
	else \
		echo "Verify git repo state ... repository is clean"; \
	fi

tests:
	@echo "####################################"
	@echo "######### Running go tests #########"
	@if ! go test ./... -count=1; then \
		echo "Go tests failed"; \
		exit 1; \
	else \
		echo "All go tests passed successfully"; \
	fi

bumpversion:
	@echo "################################"
	@echo "##### Running bumpversion #####"
	@if ! bumpversion $(BUMP_LEVEL); then \
		echo "Bump version failed"; \
		exit 1; \
	else \
		echo "Bump version successful"; \
	fi

gitpushtags:
	@echo "#################################"
	@echo "###### Running gitpushtags ######"
	@if ! git push origin $(DEFAULT_GIT_BRANCH) --follow-tags --porcelain; then \
		echo "Git push tags failed"; \
		exit 1; \
	else \
		echo "Git push tags successful"; \
	fi

goreleaser:
	@echo "################################"
	@echo "###### Running goreleaser ######"
	@if ! goreleaser release; then \
		echo "goreleaser release failed"; \
		exit 1; \
	else \
		echo "goreleaser release successful"; \
	fi

all: prechecks tests bumpversion gitpushtags goreleaser

help:
	@echo "Available targets:"
	@echo "  prechecks:       run a set of verifications before doing anything"
	@echo "  tests:           run a set of unit tests"
	@echo "  bumpversion:     increment the program version"
	@echo "  gitpushtags:     publish the commit and tags created by bumpversion"
	@echo "  goreleaser:      create a GitHub release"
	@echo "  all:             run all targets"


# If no target is provided, display the help message
.DEFAULT_GOAL := help
