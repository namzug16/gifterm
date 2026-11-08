BINARY_NAME=gifterm
RELEASE_TARGETS=darwin/arm64 darwin/amd64 linux/arm64 linux/amd64

_: build install

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Builds the archives attached to a GitHub release into dist/, one per OS/arch.
release-artifacts:
	rm -rf dist
	mkdir -p dist
	for target in $(RELEASE_TARGETS); do \
		os=$${target%/*}; arch=$${target#*/}; \
		mkdir -p dist/$$os-$$arch; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags="-s -w" -o dist/$$os-$$arch/$(BINARY_NAME) . || exit 1; \
		tar -czf dist/$(BINARY_NAME)-$$os-$$arch.tar.gz -C dist/$$os-$$arch $(BINARY_NAME); \
		rm -rf dist/$$os-$$arch; \
	done

.PHONY: _ build install release-artifacts
.DEFAULT_GOAL := _
