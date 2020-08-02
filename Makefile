GIT_VER := $(shell git describe --tags)
export GO111MODULE := on

.PHONY: test binary install clean

cmd/switchbot/switchbot: *.go cmd/switchbot/*.go go.*
	cd cmd/switchbot && go build -ldflags "-s -w -X main.Version=${GIT_VER}" -gcflags="-trimpath=${PWD}"

test:
	go test -v .

clean:
	rm -f cmd/switchbot/switchbot
