# AGENTS.md

Guidance for AI coding agents and new contributors working in the **llcppg** repository.

## Project overview

llcppg (LLGo autogen tool) automatically generates [LLGo](https://github.com/goplus/llgo)
bindings for C/C++ libraries, enhancing the experience of integrating LLGo with
C/C++. It parses C/C++ headers via libclang and emits Go source that maps the
library's public API onto LLGo.

- Module: `github.com/goplus/llcppg`
- Language: Go (see `go.mod` for the required Go version)
- Build/test toolchain: [LLGo](https://github.com/goplus/llgo) with LLVM/Clang
  (llcppg depends on libclang and is compiled/tested with `llgo`, not plain `go`).

## Repository layout

| Path | Description |
| --- | --- |
| `cl/` | Core compiler: loads clang translation units and generates the Go package (`Package`, `Config`, `NewPackage`, type/enum/func/class/namespace handling, inline-function wrapping). |
| `cl/cltest/` | Test helpers for `cl` (config loading, comparison utilities). |
| `cl/_testc/`, `cl/_testcpp/`, `cl/_testpp/` | Fixture inputs (C, C++, and package tests) consumed by `cl` tests. Directories are `_`-prefixed so the Go toolchain ignores them as packages. |
| `clang/` | Higher-level clang cursor helpers and name mangling. |
| `lib/clang/` | Low-level libclang bindings used by llcppg. |
| `xtool/clang/preprocessor/` | Preprocessor tooling built on clang. |
| `cmd/llcppg/` | Main `llcppg` command-line entry point. |
| `chore/` | Auxiliary tools: `llcppdump`, `llgogen`, `llimport`. |
| `.github/` | CI workflow (`workflows/llgo.yml`) and the `setup-llgo` composite action. |

To see the full package list: `go list ./...`.

## Build and test

Tests run under LLGo, mirroring CI (`.github/workflows/llgo.yml`).

### Set up the toolchain

CI uses the `.github/actions/setup-llgo` composite action, which:

1. Installs the pinned Go version.
2. Installs LLVM/Clang and related system libraries
   (`clang-<llvm>`, `libgc-dev`, `libssl-dev`, `zlib1g-dev`, `libffi-dev`,
   `libuv1-dev` on Linux; the Homebrew equivalents on macOS).
3. Downloads a pinned LLGo release via
   `.github/actions/setup-llgo/download-llgo.sh <llgo-version> .llgo` and puts
   `llgo` on `PATH` (setting `LLGO_ROOT`).

Match the versions pinned in the workflow matrix (Go, LLVM, and LLGo) when
reproducing CI locally.

### Run the tests

Once `llgo` is installed and on `PATH`, run the same command as CI:

```bash
llgo test -v ./...
```

Make the tests pass before submitting a change.

## Code style

- Format all Go code with `gofmt` (e.g. `gofmt -w .`) before committing.
- Keep changes focused and minimal; avoid unrelated refactors in the same change.
- Follow the conventions of the surrounding code (naming, structure, and the
  existing Apache-2.0 license header at the top of new Go files).

## Issue title convention

- Format: `feat(pkg): xxx`, `fix(pkg): xxx`, or `Proposal: xxx`.
- `pkg` is the package(s) affected by the issue. Multiple packages are separated
  by `,`, e.g. `fix(cl,parser): xxx`.
- `Proposal: xxx` issues usually affect the `cl` package, so they are equivalent
  to `feat(cl): xxx`.

## Pull requests

- Use the same title convention as issues (`feat(pkg): xxx`, `fix(pkg): xxx`,
  or `Proposal: xxx`).
- Link the related issue in the PR description (e.g. `Fixes #123`).
- Ensure `gofmt` is clean and `llgo test -v ./...` passes before requesting review.
- Keep the change scoped to a single concern to make review straightforward.
