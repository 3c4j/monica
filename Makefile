$(eval GIT_HASH := $(shell git rev-parse --short HEAD))
$(eval GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD))
$(eval BUILD_TIME := $(shell date +'%Y-%m-%d %H:%M:%S'))
$(eval LAST_COMMIT_HASH := $(shell git rev-parse --short HEAD))
$(eval LAST_COMMIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD))
$(eval LAST_COMMIT_TIME := $(shell git show -s --format=%ci HEAD))
$(eval LAST_COMMIT_AUTHOR := $(shell git show -s --format='%an <%ae>' HEAD))
$(eval LAST_COMMIT_MESSAGE := $(shell git show -s --format='%s' HEAD))

LDFLAGS := -X 'main.GitHash=$(GIT_HASH)' \
           -X 'main.GitBranch=$(GIT_BRANCH)' \
           -X 'main.BuildTime=$(BUILD_TIME)' \
           -X 'main.LastCommitHash=$(LAST_COMMIT_HASH)' \
           -X 'main.LastCommitBranch=$(LAST_COMMIT_BRANCH)' \
           -X 'main.LastCommitTime=$(LAST_COMMIT_TIME)' \
           -X 'main.LastCommitAuthor=$(LAST_COMMIT_AUTHOR)' \
           -X 'main.LastCommitMessage=$(LAST_COMMIT_MESSAGE)'

.PHONY: wire
wire:
	@go run -mod=mod github.com/google/wire/cmd/wire ./...

.PHONY: build
build:
	@mkdir -p bin
	@make wire
	go build -ldflags "$(LDFLAGS)" -o bin/ ./...
