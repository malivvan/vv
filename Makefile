.PHONY: generate lint test fmt build docs
default: build

COMMIT = $(shell git rev-parse HEAD)
ifeq ($(shell git status --porcelain),)
	VERSION = $(shell git describe --tags --abbrev=0)
endif


TEST_FORMAT ?= pkgname

define build
	@mkdir -p build
   	$(eval COMPILED := $(shell date -u +'%Y-%m-%dT%H:%M:%S.%9N'))
	$(eval OUTPUT := $(if $(filter windows,$(1)),vv-$(1)-$(2).exe,vv-$(1)-$(2)))
	$(eval URL := $(shell if [ -z "$(VERSION)" ]; then echo -n "" ; else echo -n https://github.com/malivvan/vv/releases/download/$(VERSION)/$(OUTPUT); fi))
	$(eval SERIAL := $(shell if [ -z "$(VERSION)" ]; then uuidgen --random ; else uuidgen --sha1 --namespace @url --name $(URL); fi))
	@echo "$(OUTPUT)"
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) GOFLAGS=-tags="$(4)" cyclonedx-gomod \
      app -json -packages -licenses \
      -serial=$(SERIAL) \
      -output build/$(OUTPUT).json -main ./cmd > /dev/null 2>&1
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go \
	build -trimpath -tags="$(4)" \
	  -ldflags="$(3) \
	  -buildid=$(SERIAL) \
	  -X main.serial=$(SERIAL) \
	  -X main.commit=$(COMMIT) \
	  -X main.version=$(VERSION) \
	  -X main.compiled=$(COMPILED)" \
	  -o build/$(OUTPUT) ./cmd
	@if [ ! -f build/release.md ]; then \
	  echo "| filename | serial |" > build/release.md; \
	  echo "|----------|--------|" >> build/release.md; \
	fi
	@if [ -z "$(VERSION)" ]; then \
	  echo "| $(OUTPUT) | $(SERIAL) |" >> build/release.md; \
	else \
	  echo "| [$(OUTPUT)]($(URL)) | [$(SERIAL)]($(URL).json) |" >> build/release.md; \
	fi
endef

install/build:
	@go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest

install/test:
	@go install golang.org/x/lint/golint@latest
	@go install gotest.tools/gotestsum@latest

install/release:
	@go install go.mills.io/zs@latest

install: install/build install/test install/release

lint:
	@golint -set_exit_status ./vvm/...

docs:
	@rm -rf build/docs
	@mkdir -p build/docs
	@find ./docs/ -type f -name "*.md" -exec sh -c 'zs -p build $$1 > build/$${1%.md}.html' _ {} \;
	@cd build/docs && find . -type f -name "*.html" -exec sh -c 'echo "- [$${1%.html}](./$$1)" >> index.md' _ {} \;
	@ZS_TITLE="Documentation" zs -p build build/docs/index.md > build/docs/index.html
	@ZS_TITLE="Release" zs -p build build/release.md > build/docs/release.html
	@cp -rf .zs/js build/docs/js
	@cp -rf .zs/css build/docs/css
	@cp -rf .zs/img build/docs/img

test: generate lint
	@GODEBUG=randseednop=0 gotestsum --format $(TEST_FORMAT) --format-hide-empty-pkg --hide-summary skipped --raw-command -- go test -json -race -cover ./vvm/...
	@gotestsum --format $(TEST_FORMAT) --format-hide-empty-pkg --hide-summary skipped --raw-command -- go test -json -cover ./pkg/...
	@go run ./cmd ./vvm/testdata/cli/test.vv > /dev/null 2>&1 || (echo "END TO END TEST FAILED" && exit 1)

fmt:
	@go fmt ./...

generate:
	@go generate ./vvm/...

build: clean
	$(call build,$(shell go env GOOS),$(shell go env GOARCH),-s -w,)

release: clean
	$(call build,linux,386,-s -w,)
	$(call build,linux,amd64,-s -w,)
	$(call build,linux,arm,-s -w,)
	$(call build,linux,arm64,-s -w,)
	$(call build,darwin,amd64,-s -w,)
	$(call build,darwin,arm64,-s -w,)
	$(call build,windows,amd64,-s -w,)
	$(call build,windows,386,-s -w,)
	$(call build,windows,arm,-s -w,)
	$(call build,windows,arm64,-s -w,)

clean:
	@rm -rf ./build
