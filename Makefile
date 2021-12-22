GIT_VER := $(shell git describe --tags)

GO111MODULE := on

GO_FILES := $(shell find . -type f -name '*.go' -print)

.PHONY: build
build: $(GO_FILES)
	go build -trimpath -ldflags "-s -w -X main.Version=${GIT_VER}" -o tmp/switchbot cmd/switchbot/main.go

.PHONY: clean
clean:
	rm -f switchbot

.PHONY: lint
lint:
	go run honnef.co/go/tools/cmd/staticcheck@2021.1.2 -f stylish ./...

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
