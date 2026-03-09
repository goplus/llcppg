---
name: update-llgo-version
description: Update the github.com/goplus/llgo dependency version across the llcppg project — Go module, CI workflows, and setup action. Use this skill whenever the user asks to bump, update, or upgrade the llgo version, or mentions changing the llgo dependency, even if they don't use the exact phrase "update llgo version".
---

# Update LLGo Version

Update `github.com/goplus/llgo` to a new version across the entire project. Six files reference the llgo version and they all need to stay in sync.

## Instructions

Given a target version (e.g., `v0.14.0`):

### Step 1: Detect the current version

Read `go.mod` and note the current `github.com/goplus/llgo` version — you'll need it for the commit message.

### Step 2: Update Go module dependency

Use `go get` to update the module dependency, then tidy:

```bash
go get github.com/goplus/llgo@<target_version>
go mod tidy
```

This updates both `go.mod` and `go.sum` automatically. Using `go get` is safer than manually editing `go.mod` because it resolves transitive dependencies correctly.

### Step 3: Update CI workflow matrices

Three workflow files have a `strategy.matrix.llgo` entry that pins the llgo version. Update each one:

- `.github/workflows/go.yml` — line with `llgo: [...]`
- `.github/workflows/end2end.yml` — line with `llgo: [...]`
- `.github/workflows/gentest.yml` — line with `llgo: [...]`

Replace the version in the `llgo` array, e.g., `llgo: [v0.14.0]`.

### Step 4: Update the setup action default

In `.github/actions/setup-llcppg/action.yml`, update the `default` value for the `llgo` input to match the new version.

## Verification

After making all changes, run these commands to confirm nothing is broken:

```bash
go mod tidy
go build -v ./...
go vet ./...
go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal
```

These are the fast unit tests. If they pass, the version bump is safe to commit.

## Commit Convention

```
build(deps): bump github.com/goplus/llgo from <old_version> to <new_version>
```

Replace `<old_version>` with the version you noted in Step 1 and `<new_version>` with the target version.
