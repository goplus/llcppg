# Update LLGo Version

## Description

Update the `github.com/goplus/llgo` dependency to a new version across all configuration and CI files in the repository.

## When to Use

When a new version of LLGo is released and the project needs to be updated to use it.

## Files to Update

The following files must all be updated consistently to the same new version:

1. **`go.mod`** — Update the `github.com/goplus/llgo` require directive to the new version.
2. **`go.sum`** — Run `go mod tidy` to update checksums after changing `go.mod`.
3. **`.github/workflows/go.yml`** — Update the `llgo: [vX.Y.Z]` entry in the strategy matrix.
4. **`.github/workflows/end2end.yml`** — Update the `llgo: [vX.Y.Z]` entry in the strategy matrix.
5. **`.github/workflows/gentest.yml`** — Update the `llgo: [vX.Y.Z]` entry in the strategy matrix.
6. **`.github/actions/setup-llcppg/action.yml`** — Update the default value for the `llgo` input and its description.

## Steps

1. Determine the current llgo version by reading `go.mod`.
2. Update `go.mod` to reference the new version:
   - Change `github.com/goplus/llgo vOLD` to `github.com/goplus/llgo vNEW`.
3. Run `go mod tidy` to update `go.sum`.
4. Update all three CI workflow files (`.github/workflows/go.yml`, `.github/workflows/end2end.yml`, `.github/workflows/gentest.yml`):
   - Change `llgo: [vOLD]` to `llgo: [vNEW]` in the strategy matrix.
5. Update `.github/actions/setup-llcppg/action.yml`:
   - Change the `default:` value for the `llgo` input from `"vOLD"` to `"vNEW"`.
   - Update the `description:` to reference the new version (e.g., `"LLGo version to download (e.g. vNEW)"`).
6. Run validation:
   ```bash
   go build -v ./...
   go vet ./...
   go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal
   ```
7. Commit all changes together with a message like:
   ```
   build(deps): bump github.com/goplus/llgo from vOLD to vNEW
   ```
