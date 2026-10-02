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
| `tool/_testc/`, `tool/_testcpp/` | Fixture inputs consumed by `tool` tests, mirroring the `cl` fixture layout. |
| `tool/pputil/` | Preprocessor utilities (header listing, include resolution) used by `tool`. |
| `cmd/llcppg/` | Main `llcppg` command-line entry point. |
| `clang/` | Higher-level clang helpers (based on package `github.com/llarhub/clang-c`). |
| `.github/` | CI workflow (`workflows/llgo.yml`) that installs llgo via the `xgo-dev/setup-llgo` action. |

## Build and test

### Install the llgo toolchain

The single most important setup step is installing `llgo` correctly; AI agents
frequently get this wrong. Do **not** try to build llgo from source by hand.
Instead, mirror exactly what CI does. The `.github/workflows/llgo.yml` workflow
installs llgo through the [`xgo-dev/setup-llgo`](https://github.com/xgo-dev/setup-llgo)
action:

```yaml
    - name: Setup llgo
      uses: xgo-dev/setup-llgo@v0.2.0
      with:
        go-version: ${{ matrix.go }}
        llvm-version: ${{ matrix.llvm }}
        llgo-version: ${{ matrix.llgo }}
```

The versions CI pins (see `workflows/llgo.yml`) are Go `1.27`, LLVM `22`, and
llgo `main`. Use the same action/versions to reproduce the toolchain; after it
runs, `llgo` is on `PATH`.

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

Each fixture directory has an input header (`in.h`) and a golden `out.go`. The
test harness (`cl/compile_test.go` for `cl`, `tool/gen_test.go` for `tool`)
parses `in.h` with libclang, generates Go, and diffs it against `out.go`. When
adding or changing a fixture:

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
