GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOVULNCHECK ?= govulncheck

TOOLS_GOLANGCI_LINT := github.com/golangci/golangci-lint/cmd/golangci-lint
TOOLS_GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: fmt-check
fmt-check:
	gofmt -l .

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

.PHONY: vuln
vuln:
	$(GOVULNCHECK) ./...

.PHONY: test
test:
	$(GO) test ./...

.PHONY: tools
tools:
	cd tools && $(GO) install $(TOOLS_GOLANGCI_LINT)@$$($(GO) list -m -f '{{.Version}}' github.com/golangci/golangci-lint)
	cd tools && $(GO) install $(TOOLS_GOVULNCHECK)@$$($(GO) list -m -f '{{.Version}}' golang.org/x/vuln)

.PHONY: check
check: fmt-check vet lint vuln
