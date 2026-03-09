# Skill: Update LLGo Version

## Description
Update the `github.com/goplus/llgo` dependency to a new version across the project. This includes the Go module, CI workflows, and the setup action.

## When to Use
When asked to bump, update, or upgrade the llgo version (e.g., "update llgo to v0.14.0").

## Instructions

Given a target version (e.g., `v0.14.0`), update all of the following files:

### 1. `go.mod`
Replace the `github.com/goplus/llgo` version in the `require` block with the new version.

### 2. `go.sum`
Run `go mod tidy` to update the checksums automatically.

### 3. `.github/workflows/go.yml`
Update the `llgo` entry in the `strategy.matrix` section:
```yaml
llgo: [v0.14.0]
```

### 4. `.github/workflows/end2end.yml`
Update the `llgo` entry in the `strategy.matrix` section:
```yaml
llgo: [v0.14.0]
```

### 5. `.github/workflows/gentest.yml`
Update the `llgo` entry in the `strategy.matrix` section:
```yaml
llgo: [v0.14.0]
```

### 6. `.github/actions/setup-llcppg/action.yml`
Update the default value for the `llgo` input:
```yaml
llgo:
  description: "LLGo version to download (e.g. v0.14.0)"
  default: "v0.14.0"
```

## Verification

After making all changes, run:

```bash
go mod tidy
go build -v ./...
go vet ./...
go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal
```

## Commit Convention

Use this commit message format:
```
build(deps): bump github.com/goplus/llgo from <old_version> to <new_version>
```
