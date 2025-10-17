# llcppg Project AI Assistant Guide

This guide helps AI assistants like Claude Code understand and work effectively with the llcppg project - a tool for automatically generating LLGo bindings for C/C++ libraries.

## Project Overview

llcppg is a binding generator that bridges C/C++ libraries to LLGo (a Go-based compiler). It processes C/C++ header files and generates idiomatic Go bindings, enabling Go code to seamlessly interact with C/C++ libraries.

### Core Components

1. **llcppcfg** - Configuration file generator (Go)
2. **llcppg** - Main binding generator (Go) 
3. **gogensig** - Go signature generator (Go)
4. **llcppsymg** - Symbol generator (LLGo, requires LLGo compilation)
5. **llcppsigfetch** - Signature fetcher (LLGo, requires LLGo compilation)

### Key Directories

- `cmd/` - Main executables (llcppg, llcppcfg, gogensig)
- `_xtool/` - LLGo-compiled tools (llcppsymg, llcppsigfetch)
- `cl/` - Core conversion logic and AST processing
- `parser/` - C/C++ header file parsing interface. The core C AST to internal AST conversion logic is in `_xtool/internal/parser`, compiled with LLGo
- `config/` - Configuration file handling
- `_llcppgtest/` - Real-world binding examples (cjson, sqlite, lua, etc.)
- `_demo/` - Simple demonstration projects
- `_cmptest/` - Comparison and end-to-end tests

## Development Setup

For detailed setup instructions including prerequisites, dependencies, and installation steps, see [README.md](README.md).

## Building and Testing

### Build Commands

```bash
go build -v ./...
```

### Testing Strategy

#### Unit Tests (Priority)
```bash
go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal
```
Always run these unit tests first for quick validation before running the full test suite.

#### LLGo-Dependent Tests
```bash
llgo test ./_xtool/internal/...
llgo test ./_xtool/llcppsigfetch/internal/...
llgo test ./_xtool/llcppsymg/internal/...
```

#### Full Test Suite
```bash
go test -v ./...
```
Some tests require LLGo tools installed via `install.sh`.

#### Demo Validation
```bash
bash .github/workflows/test_demo.sh
```

### Pre-Commit Validation

**ALWAYS** run before committing:

```bash
go fmt ./...
go vet ./...
go test -timeout=10m ./...
```

## Architecture and Design

For detailed technical specifications, see [llcppg Design Documentation](doc/en/dev/llcppg.md).

### Type Mapping System

For complete details, see [Type Mapping](doc/en/dev/llcppg.md#type-mapping).

### Name Conversion Rules

For complete details, see [Name Mapping Rules](doc/en/dev/llcppg.md#name-mapping-rules).

### Dependency System

For complete details, see [Dependency](doc/en/dev/llcppg.md#dependency).

### File Generation Rules

For complete details, see [File Generation Rules](doc/en/dev/llcppg.md#file-generation-rules).

## Common Issues and Solutions

### Build Failures

**"llgo: command not found"**
- LLGo not installed or not in PATH
- Solution: Install LLGo with correct commit hash

**"llcppsymg: executable file not found"**
- **CRITICAL**: MUST run `bash ./install.sh`
- This is absolutely essential for testing

**"BitReader.h: No such file or directory"**
- Install LLVM 19 development packages
- Ensure LLVM 19 is in PATH

### Test Failures

**Tests requiring llcppsigfetch or llcppsymg**
- MUST install via `bash ./install.sh`
- Do not skip this step

**Demo tests failing**
- Verify library dependencies (libcjson-dev, etc.) are installed

## Working with Examples

### Study Real Examples

Look at `_llcppgtest/` subdirectories for working configurations:

- `cjson/` - JSON library binding
- `sqlite/` - Database binding
- `lua/` - Scripting language binding
- `zlib/` - Compression library binding
- `libxml2/` - XML parsing with dependencies

Each contains:
- `llcppg.cfg` - Configuration
- Generated `.go` files
- `demo/` - Usage examples

### Unit Test Verification for New Features

When adding a new feature to llcppg, follow this workflow to verify your changes:

1. **Create a test case** in `cl/internal/convert/_testdata/` with:
   - A `conf/` directory containing `llcppg.cfg` and `llcppg.symb.json`
   - An `hfile/` directory with your test header files
   - Configuration that exercises your new feature

2. **Generate the expected output** using the `testFrom` function:
   - Temporarily set `gen:true` in the test call to generate `gogensig.expect`
   - Run the test: `go test -v ./cl/internal/convert -run TestFromTestdata`
   - This creates the expected output file that future test runs will compare against

3. **Verify the test passes** with `gen:false`:
   - Change back to `gen:false` (or remove the gen parameter)
   - Run the test again to ensure it passes
   - The test compares generated output against `gogensig.expect`

4. **Do NOT commit the test case** to the repository unless it's a permanent regression test
   - Use test cases for verification during development
   - Only add to `_testdata/` if it should be part of the test suite

**Example test structure:**
```
cl/internal/convert/_testdata/yourfeature/
├── conf/
│   ├── llcppg.cfg
│   └── llcppg.symb.json
├── hfile/
│   └── test.h
└── gogensig.expect  (generated with gen:true)
```

## Important Constraints

### Version Requirements

For detailed version requirements and installation instructions, see [README.md](README.md).

### Header File Order

Header files in `include` must be in dependency order. If `filter.h` uses types from `vli.h`, then `vli.h` must appear first in the `include` array.

## Code Style and Conventions

### Configuration Files
- Use JSON format for `llcppg.cfg`
- Follow examples in `_llcppgtest/` for structure
- Comment complex configurations

### Generated Code
- Do not manually edit generated `.go` files
- Regenerate bindings after config changes
- Use `typeMap` and `symMap` for customization

### Testing
- Add test cases for new features
- Run full test suite before PR
- Validate with real library examples

## CI/CD

The project uses GitHub Actions workflows:

- `.github/workflows/go.yml` - Main test suite
- `.github/workflows/end2end.yml` - End-to-end validation
- `.github/workflows/test_demo.sh` - Demo validation script

These run automatically on PR and provide validation feedback.

## Getting Help

- Check [README.md](README.md) for comprehensive usage documentation
- Review design documentation in `doc/en/dev/`
- Study working examples in `_llcppgtest/`

## Key Principles

1. **Always install tools first** - Run `bash ./install.sh` before testing
2. **Follow dependency order** - LLGo requires specific LLVM and commit versions
3. **Validate thoroughly** - Run full test suite and demos
4. **Study examples** - Real-world bindings in `_llcppgtest/` are the best reference

This guide provides the foundation for working effectively with llcppg. For detailed technical specifications, always reference the design documentation in `doc/en/dev/`.
