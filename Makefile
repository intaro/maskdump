GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOVULNCHECK ?= govulncheck

TOOLS_GOLANGCI_LINT := github.com/golangci/golangci-lint/cmd/golangci-lint
TOOLS_GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck

.PHONY: fmt fmt-check vet lint vuln tools check

fmt:
	gofmt -w .

fmt-check:
	gofmt -l .

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run

vuln:
	$(GOVULNCHECK) ./...

tools:
	cd tools && $(GO) install $(TOOLS_GOLANGCI_LINT)@$$($(GO) list -m -f '{{.Version}}' github.com/golangci/golangci-lint)
	cd tools && $(GO) install $(TOOLS_GOVULNCHECK)@$$($(GO) list -m -f '{{.Version}}' golang.org/x/vuln)

check: fmt-check vet lint vuln
