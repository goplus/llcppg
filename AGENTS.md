# AGENTS.md

Guidance for AI coding agents and new contributors working in the **llcppg** repository.

## Project overview

llcppg (LLGo autogen tool) automatically generates [LLGo](https://github.com/goplus/llgo)
bindings for C/C++ libraries, enhancing the experience of integrating LLGo with
C/C++. It parses C/C++ headers via libclang and emits Go source that maps the
library's public API onto LLGo.

- Module: `github.com/goplus/llcppg`
- Language: Go
- Build/test toolchain: [LLGo](https://github.com/goplus/llgo) with LLVM/Clang
  (llcppg depends on libclang and is compiled/tested with `llgo`, not plain `go`).

The core functionality lives in two packages: `cl` and `tool`. `cl` is the
kernel that loads clang translation units and generates the Go package; `tool`
orchestrates and integrates the components (loading configuration, parsing
sources, and driving `cl` to produce the final package). These two packages are
the most critical; changes to either must be tested with `llgo test`.

> llcppg is the first production-grade application built with LLGo, and it can
> only be compiled and executed with LLGo. The official Go compiler can build
> the project and everything proceeds normally until the final linking stage,
> where it fails with missing symbols. Therefore every package in this
> repository must be tested and verified with `llgo test`, never plain
> `go test`.

## Repository layout

| Path | Description |
| --- | --- |
| `cl/` | Core compiler kernel: loads clang translation units and generates the Go package. |
| `cl/_testc/`, `cl/_testcpp/`, `cl/_testpp/` | Fixture inputs (C, C++, and package tests) consumed by `cl` tests. Directories are `_`-prefixed so the Go toolchain ignores them as packages. |
| `tool/` | Orchestration layer that integrates the components: loads configuration, parses sources, and drives `cl` to generate a package (`Config.NewPackage`). Used by `cmd/llcppg`. |
| `tool/_testc/`, `tool/_testcpp/` | Fixture inputs consumed by `tool` tests. Unlike the flat `in.h` + `out.go` `cl` fixtures, these are config/package-driven: nested `include/` header trees plus `llcppg.cfg`/`llcppg.pub`, with a golden `out.go` per generated sub-package. |
| `tool/pputil/` | Preprocessor utilities (header listing, include resolution) used by `tool`. |
| `cmd/llcppg/` | Main `llcppg` command-line entry point. |
| `clang/` | Higher-level clang helpers (based on package `github.com/llarhub/clang-c`). |
| `.github/` | CI workflow (`workflows/llgo.yml`) validates the setup-llgo action and its standalone installer. |

## Build and test

### Install the llgo toolchain

**Installing llgo is a required, supported step. "I cannot install llgo" is not
an acceptable reason to skip `llgo test`.** Use setup-llgo's standalone installer
for local/agent work; it shares the action's version resolver, source builder and
dependency setup. CI tests both entrypoints (`.github/workflows/llgo.yml`).

During this draft's validation, use the candidate from
[setup-llgo #51](https://github.com/xgo-dev/setup-llgo/pull/51), pinned to the same
commit as CI below. Follow its
[local installation guide](https://github.com/cpunion/setup-llgo/blob/66efbc5bac45a7810acc79717e8bee5ac62f8e77/README.md#local-development-and-agents);
the entrypoint is `scripts/install.sh`.
For released LLGo and platform package-manager instructions, see the
[upstream LLGo README](https://raw.githubusercontent.com/xgo-dev/llgo/refs/heads/main/README.md).

Prerequisites: Linux (Ubuntu/Debian with `apt`) or macOS (Homebrew), Bash, Git,
Node.js 20+, an existing Go 1.21+ launcher, network access, and `sudo` (or root)
on Linux. If Go or Node is absent, install it with the platform package manager
first. The installer then selects the exact Go 1.27.0 toolchain.

1. Check out the candidate into a new sibling directory, then install:

   ```bash
   git clone https://github.com/cpunion/setup-llgo.git ../setup-llgo
   git -C ../setup-llgo checkout --detach 66efbc5bac45a7810acc79717e8bee5ac62f8e77
   LLGO_VERSION=main GO_VERSION=1.27.0 LLVM_VERSION=22 bash ../setup-llgo/scripts/install.sh
   ```

   Reuse that installer checkout on subsequent invocations; do not reset another
   user's checkout. `LLGO_VERSION` accepts a branch, tag, or commit. Each install
   owns a new directory under `~/.cache/setup-llgo` by default, without changing
   shell profiles. No `npm install` is needed.

2. Run the exact `source .../env.sh` command printed at the end, and repeat it
   in each new shell. It sets `LLGO_ROOT`, `GOTOOLCHAIN=local` and `PATH`, with
   the selected Go/LLVM/LLGo executables before older installations. Do not use
   the old `~/.llgo/bin` or `~/.llgo-src` paths for this installer.

3. Confirm success:

   ```bash
   GOTOOLCHAIN=local go version     # must report go1.27.0
   llvm-config --version           # must report LLVM 22.x
   llgo version
   llgo test ./tool/pputil/...
   ```

Common failure modes:

- `go: unknown GOEXPERIMENT dwarf5`: an older `go` (e.g. 1.24) is first on
  `PATH`. Source the generated `env.sh` to activate the actual Go 1.27 binary
  with `GOTOOLCHAIN=local`.
- `llvm-config`/`clang` not found or wrong version: add
  `/usr/lib/llvm-22/bin` (Linux) or `$(brew --prefix llvm@22)/bin` (macOS) to
  `PATH`.
- `llgo: command not found`: source the generated `env.sh` from step 2.
- apt/brew/network errors: re-run the script; it is idempotent.
- Link errors with missing symbols: you ran plain `go test`; use `llgo test`.

### Run the tests

Once `llgo` is installed and on `PATH`, run the same command as CI. Always use
`llgo test`, never plain `go test` — the latter fails at link time with missing
symbols:

```bash
llgo test -v ./...
```

Both the `cl` kernel and the `tool` orchestration layer have fixture-driven
tests. `cl` drives every fixture under `cl/_testc` (C, via `TestC`),
`cl/_testcpp` (C++, via `TestCpp`), and `cl/_testpp` (`TestPreprocessor`).
`tool` drives fixtures under `tool/_testc` (via `TestC`) and `tool/_testcpp`
(the `TestCpp_*` / `TestLLVM_*` tests). To iterate on one fixture, filter by
name and scope it to a package:

```bash
llgo test -v -run 'TestC/union_struct' ./cl/
llgo test -v -run 'TestLLVM_String' ./tool/
```

#### Fixtures and golden files

The two packages use different fixture layouts:

- **`cl`** (`cl/compile_test.go`): each fixture directory has a flat input
  header (`in.h`) and a golden `out.go`. The harness parses `in.h` with
  libclang, generates Go, and diffs it against `out.go`.
- **`tool`** (`tool/gen_test.go`): fixtures are config/package-driven — a
  nested `include/` header tree plus `llcppg.cfg`/`llcppg.pub`. The harness
  loads the config, drives `cl` via `Config.NewPackage`, and diffs each
  generated sub-package against its own `out.go` (e.g.
  `tool/_testc/clang-c-22.1.8/CXString/out.go`); there is no `in.h`.

In both cases the output is compared against a golden `out.go`. When adding or
changing a fixture:

1. Implement the generator change, then run the fixture test.
2. When the generated Go differs from the golden (including a missing or stale
   `out.go`), the harness writes the actual output to `out.go.txt` (gitignored
   via the `*.txt` rule) **and fails the test** — it does not silently accept a
   mismatch. Do not hand-write or hand-edit `out.go`; always let the harness
   produce it so the golden matches the generator byte-for-byte.
3. Inspect `out.go.txt`, and once correct promote it: `mv out.go.txt out.go`.
4. Re-run until the fixture passes with no `out.go.txt` produced.

Make the tests pass before submitting a change.

## Issue title convention

- Format: `feat(pkg): xxx`, `fix(pkg): xxx`, or `Proposal: xxx`.
- `pkg` is the package(s) affected by the issue. Multiple packages are separated
  by `,`, e.g. `fix(cl,parser): xxx`.
- `Proposal: xxx` issues usually affect the `cl` package, so they are equivalent
  to `feat(cl): xxx`.
