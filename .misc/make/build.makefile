.DEFAULT_GOAL := all
.PHONY: all
all: init clean deps gen build test lint

.PHONY: clean
clean:
	$(call create-banner, 'Cleaning project')
	$(GO) clean
	@rm -rf build | true

.PHONY: deps
deps:
	$(call create-banner, 'Getting Dependencies')
	$(call create-subtext, 'Cleaning up dependency list...')
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: gen
gen:
	$(call create-banner, 'Generating dependencies')
	$(GO) generate

.PHONY: build
build:
	$(call create-banner, 'Building library')
	@$(GO) build .

.PHONY: test
test:
	$(call create-banner, 'Running tests')
	@mkdir -p build/tests
	@gotestsum --format standard-verbose -- -coverprofile=build/tests/cover.out ./...
	$(GO) tool cover -html build/tests/cover.out -o build/tests/cover.html
	$(call create-subtext, '========================================')
	$(GO) tool cover -func build/tests/cover.out

.PHONY: lint
lint:
	$(call create-banner, 'Running lint')
	$(GOLANGCI_LINT) run
	$(call create-subtext, 'Done.')

