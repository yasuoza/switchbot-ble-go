GIT_VER := $(shell git describe --tags)

GO111MODULE := on

GO_FILES := $(shell find . -type f -name '*.go' -print)

.PHONY: tools
tools:
	cd tools && cat main.go | awk '/_/ {print $$2}' | xargs -tI {} go install {}

.PHONY: build
build: $(GO_FILES)
	go build -trimpath -ldflags "-s -w -X main.Version=${GIT_VER}" -o tmp/switchbot cmd/switchbot/main.go

.PHONY: clean
clean:
	rm -f switchbot

.PHONY: lint
lint:
	staticcheck -f stylish ./...

.PHONY: test
test:
	go test -v -shuffle=1 ./...

# https://goreleaser.com/install/#running-with-docker
goreleaser/build:
	docker run --rm --privileged \
		-v ${PWD}:/go/src/github.com/user/repo \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/user/repo \
		goreleaser/goreleaser \
		build \
		--rm-dist --skip-validate
