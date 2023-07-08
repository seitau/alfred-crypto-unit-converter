GO_ENV ?= GOBIN=$(BIN) GOFLAGS=-mod=mod
WORKFLOW_FILE := build/alfred-crypto-unit-converter.alfredworkflow

.PHONY: build-workflow
build-workflow: clean build
	zip $(WORKFLOW_FILE) \
	coins.json \
	info.plist \
	icon.png \
	icons \
	./bin/coverter

.PHONY: clean
clean:
	rm ./bin/converter $(WORKFLOW_FILE) || :

.PHONY: build
build: BUILD_ENV ?= GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
build: BUILD_OPTS ?= -ldflags='-s -w' -trimpath
build:
	$(BUILD_ENV) go build $(BUILD_OPTS) -o ./bin/converter ./cmd/converter
