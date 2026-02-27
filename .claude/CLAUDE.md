# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

llcppg is an LLGo autogen tool for C/C++ libraries that automatically generates LLGo bindings for C/C++ libraries. The tool enhances integration between LLGo and the C ecosystem.

**Key Technologies:**
- Go 1.23.0 
- LLGo (Go+ compiler for low-level programming)
- LLVM 19 with clang/libclang for C/C++ parsing
- Dependencies: gogen, llgo, mod, qiniu/x

## Build and Development Commands

### Installation and Setup
```bash
# Install system dependencies (required before building)
# For macOS:
brew install cjson llvm@19 bdw-gc openssl libffi libuv lld@19 zlib

# Install all tools (REQUIRED for full functionality)
bash ./install.sh
```

### Build Commands
```bash
# Build all Go components
go build -v ./...

# Install individual tools
go install ./cmd/llcppcfg    # Configuration generator
go install ./cmd/gogensig    # Go signature generator  
go install ./cmd/llcppg      # Main binding generator
llgo install ./_xtool/llcppsymg      # Symbol generator (LLGo)
llgo install ./_xtool/llcppsigfetch  # Signature fetcher (LLGo)
```

### Test Commands
```bash
# Basic Go tests (fast, 2 seconds)
go test -v ./config ./internal/name ./internal/arg ./internal/unmarshal

# Full test suite (8-12 minutes, requires LLGo tools)
go test -timeout=10m -v ./...

# LLGo-specific tests (3-5 minutes)
llgo test ./_xtool/internal/...
llgo test ./_xtool/llcppsigfetch/internal/...
llgo test ./_xtool/llcppsymg/internal/...

# Demo validation (5-10 minutes)
bash .github/workflows/test_demo.sh
```

### Usage Commands
```bash
# Generate configuration file
llcppcfg [libname]

# Generate LLGo bindings
llcppg [config-file]  # Uses llcppg.cfg if not specified

# Test a complete workflow
llcppcfg cjson
# Edit llcppg.cfg to add include files
llcppg llcppg.cfg
```

## High-Level Architecture

### Three-Stage Pipeline
1. **llcppsymg** (`_xtool/llcppsymg/`) - Symbol table generator
   - Analyzes C/C++ libraries and header files
   - Extracts symbols using nm tool
   - Matches library symbols with header declarations
   - Outputs: `llcppg.symb.json`

2. **llcppsigfetch** (`_xtool/llcppsigfetch/`) - Signature fetcher  
   - Uses libclang to parse C/C++ header files
   - Extracts type information and function signatures
   - Outputs: JSON-formatted package information structure

3. **gogensig** (`cmd/gogensig/`) - Go code generator
   - Converts C/C++ declarations to Go code
   - Generates functions, methods, types, and constants
   - Uses symbol table to determine what to generate

### Core Components

**Parser Layer** (`parser/`, `_xtool/internal/parser/`)
- C/C++ header file parsing using libclang
- AST generation and traversal
- Type extraction and analysis

**Conversion Layer** (`cl/`, `cl/internal/convert/`)
- Type mapping from C/C++ to Go
- Function signature conversion
- Method generation logic
- Name transformation rules

**Configuration** (`config/`)
- Configuration file parsing (`llcppg.cfg`)
- Dependency management
- Type mapping configuration

**Tool Layer** (`tool/`)
- High-level APIs for binding generation
- File composition and output management

### Key File Types

**Configuration Files:**
- `llcppg.cfg` - Main configuration (cflags, libs, includes, deps)
- `llcppg.pub` - Type mapping for dependencies  
- `llcppg.symb.json` - Symbol table mapping

**Generated Files:**
- `{header}.go` - Go bindings for each header file
- `{name}_autogen.go` - Implementation files
- `{name}_autogen_link.go` - Linking information

### Critical Dependencies

**LLGo Tools** (compiled with LLGo, not Go):
- `llcppsymg` - Symbol analysis
- `llcppsigfetch` - Header parsing
- These MUST be installed via `install.sh` for testing

**System Dependencies:**
- LLVM 19 development libraries
- libclang-19-dev
- Platform-specific C libraries (libcjson-dev, etc.)

## Development Guidelines

### Before Making Changes
1. Run `bash ./install.sh` to ensure all tools are available
2. Test with a known working example from `_llcppgtest/` or `_demo/`
3. Verify the three-stage pipeline works end-to-end

### After Making Changes  
1. Build: `go build -v ./...`
2. Install: `bash ./install.sh` (essential for testing)
3. Test: `go test -timeout=10m -v ./...`
4. Validate with demo: `bash .github/workflows/test_demo.sh`

### Testing Strategy
- Quick validation: `go test -v ./config ./internal/name ./internal/arg`
- Full validation: Install tools + run full test suite
- Real-world validation: Test with cjson, sqlite, or zlib examples

### Common Patterns
- Look at `_llcppgtest/` for real configuration examples
- Follow existing type mapping patterns in `cl/internal/convert/`
- Use naming conventions from `internal/name/`
- Reference design documentation in `doc/en/dev/llcppg.md`

## Important Notes

- Never attempt to modify LLGo tools (`_xtool/llcppsymg`, `_xtool/llcppsigfetch`) without LLGo installed
- Always run `install.sh` after changes - many tests depend on installed tools
- Build times: 15 seconds for Go, 2-3 minutes for complete install
- Test times: 2 seconds basic, 8-15 minutes full suite
- The project requires specific LLVM 19 and LLGo versions - do not change these dependencies