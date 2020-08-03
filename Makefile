GIT_VER := $(shell git describe --tags)
export GO111MODULE := on

.PHONY: test binary install clean

cmd/switchbot/switchbot: *.go cmd/switchbot/*.go go.*
	cd cmd/switchbot && go build -ldflags "-s -w -X main.Version=${GIT_VER}" -gcflags="-trimpath=${PWD}"

test:
	go test -v .

# https://goreleaser.com/install/#running-with-docker
goreleaser/build:
	docker run --rm --privileged \
		-v ${PWD}:/go/src/github.com/user/repo \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/user/repo \
		goreleaser/goreleaser \
		build \
		--rm-dist --skip-validate

clean:
	rm -f cmd/switchbot/switchbot
