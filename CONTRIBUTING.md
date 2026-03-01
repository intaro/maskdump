# Contributing

Thank you for your interest in MaskDump. This document describes the contribution workflow and the quality tools we use.

## Language Policy

- The main `README.md` is bilingual (English and Russian).
- All other project documents and files are in English.
- Commit messages must be in English.

## Requirements

- Go **1.26.0** or newer
- `git`, `make`
- Installed quality tools (see below)

## Code Quality Tools

We pin tool versions in the `tools/` module (`tools/go.mod` and `tools/tools.go`).

Install tools:
```bash
make tools
```

Run all checks:
```bash
make check
```

`make check` includes:
- `gofmt -l .` — formatting check
- `go vet ./...` — basic analysis
- `golangci-lint run` — static analysis
- `govulncheck ./...` — vulnerability scan

## Formatting

Auto-format code:
```bash
make fmt
```

## Useful Commands

Check tool versions:
```bash
golangci-lint version
govulncheck -version
```

## Code Guidelines

- Keep functions small and readable.
- Avoid global state when possible.
- New functionality should include tests.
- Handle errors explicitly.

## Before Opening a PR

1. `make fmt`
2. `make check`
3. Ensure tests pass

Thanks for your contributions.
