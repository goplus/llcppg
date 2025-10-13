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
- `parser/` - C/C++ header file parsing
- `config/` - Configuration file handling
- `_llcppgtest/` - Real-world binding examples (cjson, sqlite, lua, etc.)
- `_demo/` - Simple demonstration projects
- `_cmptest/` - Comparison and end-to-end tests

## Development Setup

### Prerequisites

llcppg has strict dependencies that MUST be installed in the correct order:

1. **LLVM 19** - Exactly version 19 (not compatible with other versions)
   ```bash
   # Ubuntu/Linux
   sudo apt-get install -y llvm-19-dev clang-19 libclang-19-dev lld-19 libunwind-19-dev libc++-19-dev
   export PATH="/usr/lib/llvm-19/bin:$PATH"
   
   # macOS
   brew install llvm@19 lld@19
   echo "$(brew --prefix llvm@19)/bin" >> $GITHUB_PATH
   ```

2. **System Dependencies**
   ```bash
   # Ubuntu/Linux
   sudo apt-get install -y pkg-config libgc-dev libssl-dev zlib1g-dev libffi-dev libuv1-dev libcjson-dev
   
   # macOS
   brew install bdw-gc openssl libffi libuv zlib cjson
   ```

3. **LLGo** - CRITICAL dependency, specific commit required
   ```bash
   git clone https://github.com/goplus/llgo.git .llgo
   cd .llgo
   git checkout f0728c4fe028fbc72455c1242cef638ebdf60454
   go install -v ./cmd/llgo/...
   export LLGO_ROOT=$(pwd)
   cd ..
   ```

### Installation

**CRITICAL**: Always run the installation script to build all tools:

```bash
bash ./install.sh
```

**NEVER CANCEL**: This takes 2-3 minutes including LLGo compilation. Set timeouts to 10+ minutes.

This script installs all five core tools. Without it, tests will fail.

## Building and Testing

### Build Commands

```bash
go build -v ./...
```

Timing: ~15 seconds. Build always succeeds without LLGo dependencies.

### Testing Strategy

#### Quick Tests (Standard Go)
```bash
go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal
```
Timing: ~2 seconds. Always run these first for quick validation.

#### Full Test Suite
```bash
go test -v ./...
```
Timing: 8-12 minutes. Some tests require LLGo tools installed via `install.sh`.

#### LLGo-Dependent Tests
```bash
llgo test ./_xtool/internal/...
llgo test ./_xtool/llcppsigfetch/internal/...
llgo test ./_xtool/llcppsymg/internal/...
```
**NEVER CANCEL**: Takes 3-5 minutes. Set timeout to 10+ minutes.

#### Demo Validation
```bash
bash .github/workflows/test_demo.sh
```
**NEVER CANCEL**: Takes 5-10 minutes. Set timeout to 15+ minutes.

### Pre-Commit Validation

**ALWAYS** run before committing:

```bash
go fmt ./...
go vet ./...
go test -timeout=10m ./...
```

Timing: 8-12 minutes total. **NEVER CANCEL** - set timeout to 20+ minutes.

## Usage Workflow

### 1. Generate Configuration

```bash
llcppcfg [options] libname
```

Examples:
- `llcppcfg cjson` - Basic configuration
- `llcppcfg -cpp libname` - For C++ libraries
- `llcppcfg -deps "c/os,github.com/author/llpkg" libname` - With dependencies

Timing: Instant (< 1 second).

### 2. Edit Configuration

Edit the generated `llcppg.cfg` to specify:
- `include`: Header files to process
- `cflags`: Compiler flags
- `libs`: Library flags
- `trimPrefixes`: Prefixes to remove from names
- `deps`: Dependencies on other packages
- `typeMap`: Custom type name mappings
- `symMap`: Custom function name mappings

### 3. Generate Bindings

```bash
llcppg [config-file]
```

**NEVER CANCEL**: Takes 30 seconds to 5 minutes. Set timeout to 10+ minutes.

### 4. Validate Output

Check the generated Go package:
- Type definitions match C structures
- Functions are correctly mapped
- Dependencies are properly imported

## Architecture and Design

### Type Mapping System

llcppg converts C/C++ types to Go types following strict rules:

#### Basic Types
- `int` → `c.Int`
- `unsigned int` → `c.Uint`
- `char` → `c.Char`
- `void*` → `c.Pointer`
- `float` → `c.Float`
- `double` → `c.Double`

All basic types are imported from `github.com/goplus/lib/c`.

#### Function Pointers
C function pointers become Go function types with `llgo:type C` tag:

```c
typedef int (*CallBack)(void *L);
```

Becomes:

```go
// llgo:type C
type CallBack func(c.Pointer) c.Int
```

#### Method Generation

When a C function's first parameter matches a converted type, it becomes a Go method:

```c
int sqlite3_close(sqlite3*);
```

Becomes:

```go
// llgo:link (*Sqlite3).Close C.sqlite3_close
func (recv_ *Sqlite3) Close() c.Int {
    return 0
}
```

### Name Conversion Rules

1. **Type Names**: Convert to PascalCase after removing configured prefixes
   - `cJSON_Hooks` → `CJSONHooks` (or `Hooks` with `trimPrefixes: ["cJSON_"]`)
   - `sqlite3_destructor_type` → `Sqlite3DestructorType`

2. **Field Names**: Convert to PascalCase for export
   - `value_string` → `Valuestring`

3. **Parameters**: Preserve original case, add `_` suffix for Go keywords
   - `func` → `func_`
   - Variadic params always named `__llgo_va_list`

### Dependency System

llcppg handles cross-package dependencies through:

1. **llcppg.pub** - Type mapping table (C type → Go type name)
2. **deps field** - List of dependency packages in `llcppg.cfg`
3. **Special aliases** - `c/` prefix maps to `github.com/goplus/lib/c/`

Example: `c/os` → `github.com/goplus/lib/c/os`

### File Generation Rules

- **Interface headers** (in `include`): Each generates a `.go` file
- **Implementation headers** (same directory): Content goes in `{name}_autogen.go`
- **Link file**: `{name}_autogen_link.go` contains linking info
- **Type mapping**: `llcppg.pub` for dependency resolution

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

**Parser tests failing**
- Install llclang: `llgo install ./_xtool/llclang`

**Demo tests failing**
- Verify library dependencies (libcjson-dev, etc.) are installed

### Performance Expectations

- Build time: 15 seconds (Go only), 2-3 minutes (with install.sh)
- Test time: 2 seconds (basic), 8-15 minutes (full suite)
- Demo time: 5-10 minutes
- Memory usage: 2-4GB during LLVM compilation

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

### Validation Workflow

After making changes, ALWAYS:

1. Build: `go build -v ./...`
2. Install tools: `bash ./install.sh` (**ESSENTIAL**)
3. Generate test config: `llcppcfg sqlite`
4. Edit config to add proper headers
5. Run binding generation: `llcppg llcppg.cfg`
6. Verify Go files are generated
7. Test with example from `_demo/` or `_llcppgtest/`

## Important Constraints

### Version Requirements
- Go 1.23+
- LLGo commit: `f0728c4fe028fbc72455c1242cef638ebdf60454`
- LLVM 19 (exact version)

**NEVER** use different versions without updating the entire toolchain.

### Timeout Settings

Many operations are CPU-intensive. Set appropriate timeouts:

- `install.sh`: 10+ minutes
- Full tests: 20+ minutes
- Demo tests: 15+ minutes
- Binding generation: 10+ minutes

**NEVER CANCEL** long-running operations - they will complete successfully with proper timeout settings.

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

- `.github/workflows/go.yml` - Main test suite (8-15 minutes)
- `.github/workflows/end2end.yml` - End-to-end validation (15-30 minutes)
- `.github/workflows/test_demo.sh` - Demo validation script

These run automatically on PR and provide validation feedback.

## Getting Help

- Check `README.md` for comprehensive usage documentation
- Review design documentation in `doc/en/dev/`
- Study working examples in `_llcppgtest/`
- Reference `.github/copilot-instructions.md` for detailed build instructions

## Key Principles

1. **Always install tools first** - Run `bash ./install.sh` before testing
2. **Set generous timeouts** - LLGo compilation takes time
3. **Follow dependency order** - LLGo requires specific LLVM and commit versions
4. **Validate thoroughly** - Run full test suite and demos
5. **Study examples** - Real-world bindings in `_llcppgtest/` are the best reference
6. **Never cancel long operations** - They complete successfully with proper timeouts

This guide provides the foundation for working effectively with llcppg. For detailed technical specifications, always reference the design documentation in `doc/en/dev/`.
