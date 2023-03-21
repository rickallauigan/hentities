################################################################################
# Tool versions
################################################################################
GO_VERSION            := 1.18.1
GOTESTSUM_VERSION     := v1.8.0
GOLANGCI_LINT_VERSION := v1.44.2

################################################################################
# Tool binaries
################################################################################
GO            := @$(shell command -v go${GO_VERSION} || echo "go")
GOTESTSUM     := @$(shell echo "gotestsum")
GOLANGCI_LINT := @$(shell echo "golangci-lint")

################################################################################
# Common functions
################################################################################
define create-banner
	@echo '==================================================='
	@echo '>>  '$(1)
	@echo '==================================================='
endef

define create-subtext
	@echo ">>>> $(1)"
endef

################################################################################
# Initialization
################################################################################

.PHONY: init
init: install
	$(call create-banner, 'Initializing project')
	@cp .misc/git-hooks/* .git/hooks/
	$(call install-go)

################################################################################
# Tool installation
################################################################################

.PHONY: install
install: install-go install-gotestsum install-golangci-lint

.PHONY: install-go
install-go:
	$(call create-banner, "Ensure go is at version ${GO_VERSION}")
	@go install golang.org/dl/go${GO_VERSION}@latest
	@go${GO_VERSION} download

.PHONY: install-gotestsum
install-gotestsum:
	$(call create-banner, "Ensure gotestsum is at version ${GOTESTSUM_VERSION}")
	@go install gotest.tools/gotestsum@${GOTESTSUM_VERSION}

.PHONY: install-golangci-lint
install-golangci-lint:
	$(call create-banner, "Ensure golanglint-ci is at version ${GOLANGCI_LINT_VERSION}")
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
