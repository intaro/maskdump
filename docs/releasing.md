# Releases on GitHub

## Goal

This repository uses a GitHub-native semi-automated release flow based on:

- GitHub Actions
- Release Please
- Conventional Commits

The goal is to make the next version tag predictable from the merged change description and to keep release notes consistent without manual tag bookkeeping.

## Release Model

- Release branch: `main`
- Tag format: `vX.Y.Z`
- Current baseline: `v0.1.5`
- Changelog source: `CHANGELOG.md`
- Release artifact page: GitHub Releases

Release Please monitors pushes to `main`, opens or updates a Release PR, and after that PR is merged creates:

- the next Git tag
- a GitHub Release
- changelog updates in `CHANGELOG.md`

## How the Next Version Is Calculated

The next version is derived from Conventional Commits in merged changes:

- `feat:` -> minor bump
- `fix:` -> patch bump
- `perf:` -> patch bump
- `refactor:` -> patch bump when it changes shipped behavior or fixes a defect
- `type!:` or `BREAKING CHANGE:` footer -> major bump

Examples from the current baseline `v0.1.5`:

- next merged releasable change is `fix(cli): preserve empty values` -> `v0.1.6`
- next merged releasable change is `feat(config): support yaml config` -> `v0.2.0`
- next merged releasable change is `feat!: remove legacy config fields` -> `v1.0.0`

Documentation-only or maintenance-only changes such as `docs:`, `test:`, `ci:`, and `chore:` should still use Conventional Commits, but they should not be relied on as the primary trigger for a release.

## What Developers Must Write

The repository should prefer squash merges. In that model, the PR title becomes the canonical commit message on `main`, so the PR title controls release automation.

Use these PR title and commit subject formats:

- `feat(scope): add new capability`
- `fix(scope): correct broken behavior`
- `perf(scope): reduce allocations in parser`
- `refactor(scope): simplify masking pipeline`
- `docs(scope): update release process`
- `test(scope): cover invalid config handling`
- `ci(scope): align workflow with local checks`
- `chore(scope): update dependencies`

Scope is optional but recommended when it clarifies impact.

Breaking changes:

- `feat!: change config schema`
- or add a footer/body line: `BREAKING CHANGE: config field X was removed`

## Recommended Merge Strategy

- Enable squash merge.
- Make the PR title the source of truth for release classification.
- Avoid merge commits for routine feature work.
- Keep PRs focused so the release notes stay readable.

## Release Lifecycle

1. A developer opens a PR with a Conventional Commit title.
2. CI validates code and the PR title workflow validates the title format.
3. The PR is squash-merged into `main`.
4. Release Please opens or updates the Release PR.
5. A maintainer reviews the generated version and changelog.
6. The Release PR is merged.
7. GitHub Actions creates the tag `vX.Y.Z` and publishes the GitHub Release.

No manual tag creation is needed in the normal flow.

## Files Used by the Release System

- `.github/workflows/release-please.yml`
- `.github/workflows/pr-title-conventional.yml`
- `release-please-config.json`
- `.release-please-manifest.json`
- `CHANGELOG.md`

## Repository Settings Required on GitHub

Configure these once in the GitHub repository:

1. `Settings -> Actions -> General`
   - allow GitHub Actions for the repository
2. `Settings -> Actions -> General -> Workflow permissions`
   - set `Read and write permissions`
   - enable `Allow GitHub Actions to create and approve pull requests`
3. `Settings -> Pull Requests`
   - enable `Squash merge`
   - optionally disable merge methods you do not want to support
4. `Settings -> Branches`
   - protect `main`
   - require PRs before merge
   - require status checks from `CI` and `Conventional PR Title`

No personal access token is required for the default setup. `GITHUB_TOKEN` is enough if the repository permissions above are enabled.

## Maintainer Notes

- Review the generated Release PR before merging it.
- If the generated version is intentionally different from the default bump, use a `Release-As: X.Y.Z` footer in the merged change that should control the version.
- If an internal-only change should be excluded from release notes, avoid turning it into a user-facing `feat` or `fix`.

## Migration Notes

- Existing tags already use the `vX.Y.Z` format, so the automated flow continues from `v0.1.5`.
- The manifest baseline is pinned to `0.1.5` to make the release starting point explicit in the repository.
