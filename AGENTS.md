# AGENTS Instructions

These instructions apply to any automated or assisted work in this repository.

## Source of Truth

- Follow the rules from `CONTRIBUTING.md`.
- If this file and `CONTRIBUTING.md` overlap, treat `CONTRIBUTING.md` as the baseline and this file as an operational extension.

## Language and Documentation

- Keep `README.md` bilingual if it is updated.
- Write all other documentation and repository text in English.
- Write commit messages in English.

## Release and Commit Conventions

- Follow Conventional Commits for commit messages and PR titles.
- Prefer these types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`, `chore`, `revert`.
- Prefer an optional scope when it improves clarity, for example: `fix(cli): handle empty input`.
- Use `type!:` or a `BREAKING CHANGE:` footer for backward-incompatible changes.
- Assume squash merge is the default strategy. In that model, the PR title is the authoritative message for release automation.
- Keep commit subjects concise, imperative, and release-readable. Avoid vague subjects such as `updates`, `misc fixes`, or `changes`.
- Do not mix unrelated change types in one commit or one PR when that would obscure release notes.
- If a change is internal-only, avoid presenting it as a user-facing `feat` or `fix` unless it should affect versioning and release notes.
- For release process details, see `docs/releasing.md`.

## Required Validation After Changes

After every code change, run the full local verification flow from `CONTRIBUTING.md`:

1. `make fmt`
2. `make check`
3. Ensure tests pass

Do not consider the task complete until these checks have been run successfully, unless the environment is missing required tools or dependencies. In that case, explicitly report what could not be verified and why.

`make check` includes:

- `gofmt -l .`
- `go vet ./...`
- `golangci-lint run`
- `govulncheck ./...`

## Testing Rules

- Every new feature, branch of logic, or bug fix must include unit tests.
- If functionality changes, review the existing tests first and update them to match the new behavior.
- If functionality changes without adequate test coverage, add the missing tests in the same change.
- Do not leave behavior changes covered only by manual reasoning.
- Prefer focused unit tests over broad end-to-end coverage when validating internal logic.

## Go Development Practices

- Keep functions small, readable, and easy to test.
- Handle errors explicitly; do not silently ignore failures unless there is a documented reason.
- Avoid global mutable state. Prefer dependency injection and explicit inputs.
- Keep packages cohesive and APIs narrow.
- Prefer clear, idiomatic Go over clever abstractions.
- Use table-driven tests when multiple scenarios share the same logic.
- Keep tests deterministic: avoid time, randomness, network, filesystem, and environment coupling unless the test explicitly targets that behavior.
- Minimize side effects and make them visible at package boundaries.
- Pass `context.Context` where cancellation, deadlines, or request-scoped values matter.
- Return wrapped errors with useful context when propagating failures.
- Keep exported names and package layout consistent with Go conventions.
- Avoid unnecessary allocations, reflection, and interface indirection on hot paths.
- Prefer standard library solutions unless a dependency provides clear, justified value.
- When changing CLI behavior, configuration loading, parsing, or masking logic, add regression tests for the affected paths.

## Change Scope

- Keep changes focused on the requested task.
- Do not mix unrelated refactors into feature or bug-fix work.
- If you notice adjacent problems that should not be fixed in the same change, document them separately instead of expanding scope.

## Tooling Baseline

- Use Go `1.26.0` or newer.
- Use the repository-provided tooling versions installed via `make tools` when local tools are missing.

## GitHub Workflow Guidance

- Prefer GitHub CLI `gh` for repository inspection when it is available and authenticated.
- If `gh` is not available, use the GitHub web UI or local git state as fallback.
- Do not assume repository admin access. First determine whether the current account can view:
  - repository settings
  - Actions runs and logs
  - releases and tags
  - pull requests
- If admin-only settings are needed and the current account cannot access them, explicitly report the limitation instead of guessing.

### PR and Merge Rules

- Treat the PR title as the future squash commit message on `main`.
- Prefer `Squash and merge`.
- Before confirming a squash merge, ensure the final commit title matches the PR title exactly.
- Avoid merge methods that change the release-significant commit message unexpectedly.

### Release Please Rules

- Release Please workflow file: `.github/workflows/release-please.yml`
- Release Please config file: `release-please-config.json`
- Release Please manifest file: `.release-please-manifest.json`
- When working on Release Please configuration, preserve these workflow permissions unless there is a strong documented reason to change them:
  - `contents: write`
  - `issues: write`
  - `pull-requests: write`
- Do not enable `skip-labeling: true` for this repository unless the release flow is deliberately redesigned and revalidated.
- If a merged release PR does not produce the expected tag or GitHub Release:
  1. inspect GitHub Releases and tags
  2. inspect the Release Please workflow log
  3. confirm the merge commit landed on `main`
  4. recover missing release objects only on the real merge commit
  5. re-run Release Please after recovery
- Do not create placeholder tags on arbitrary commits.

### Standard GitHub CLI Commands

- Repository:
  - `gh repo view intaro/maskdump`
- Pull requests:
  - `gh pr list --repo intaro/maskdump --limit 20`
  - `gh pr view <number> --repo intaro/maskdump`
  - `gh pr diff <number> --repo intaro/maskdump`
- Actions:
  - `gh run list --repo intaro/maskdump --limit 20`
  - `gh run view <run-id> --repo intaro/maskdump`
  - `gh run view <run-id> --repo intaro/maskdump --log`
  - `gh workflow run release-please.yml --repo intaro/maskdump`
- Releases:
  - `gh release list --repo intaro/maskdump --limit 20`
  - `gh release view <tag> --repo intaro/maskdump`
