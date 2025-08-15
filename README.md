llcppg - LLGo autogen tool for C/C++ libraries
====

[![Build Status](https://github.com/goplus/llcppg/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/llcppg/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/llcppg)](https://goreportcard.com/report/github.com/goplus/llcppg)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/llcppg.svg?label=release)](https://github.com/goplus/llcppg/releases)
[![Coverage Status](https://codecov.io/gh/goplus/llcppg/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/llcppg)
[![Language](https://img.shields.io/badge/language-XGo-blue.svg)](https://github.com/goplus/gop)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/goplus/llcppg)

llcppg aims to be a tool for automatically generating LLGo bindings for C/C++ libraries, enhancing the experience of integrating LLGo with C! It is worth mentioning that several core components of llcppg are built using LLGo, fully leveraging its core capability of "better integrating with the C ecosystem" for development.

## How to install

This project depends on LLGo's C ecosystem integration capability, and some components of this tool must be compiled with LLGo. For LLGo installation, please refer to:
https://github.com/goplus/llgo?tab=readme-ov-file#how-to-install

```bash
brew install cjson # macos
apt-get install libcjson-dev # linux
llgo install ./_xtool/llcppsymg
llgo install ./_xtool/llcppsigfetch
go install ./cmd/llcppcfg
go install ./cmd/gogensig
go install .
```


## Usage

`llcppg.cfg` file is a configure file used by llcppg. Once `llcppg.cfg` is generated then you can run llcppg command to generate Go package for the C/C++ lib.

```sh
llcppg [config-file]
```

If `config-file` is not specified, a `llcppg.cfg` file is used in current directory.
Here's a demo configuration to generate LLGo bindings for cjson library:

```json
{
  "name": "cjson",
  "cflags": "$(pkg-config --cflags libcjson)",
  "include": ["cJSON.h","cJSON_Utils.h"],
  "libs": "$(pkg-config --libs libcjson libcjson_utils)",
  "trimPrefixes": ["cJSONUtils_","cJSON_"],
  "cplusplus": false,
  "deps": ["c"],
  "mix": false,
  "typeMap":{}
}
```

The configuration file supports the following options:

- `name`: The name of the generated package
- `cflags`: Compiler flags for the C/C++ library
- `include`: Header files to include in the binding generation
- `libs`: Library flags for linking
- `trimPrefixes`: Prefixes to remove from function names & type names
- `cplusplus`: Set to true for C++ libraries (not supported)
- `deps`: Dependencies (other packages & standard libraries)
- `mix`: Set to true when package header files are mixed with other header files in the same directory. In this mode, only files explicitly listed in `include` are processed as package files.
- `typeMap`: Custom name mapping from C types to Go types.
- `symMap`: Custom name mapping from C function names to Go function names.
- `staticLib`: Set to true to enable static library symbol reading instead of dynamic library linking. When enabled, llcppg will read symbols from static libraries (.a files) rather than dynamic libraries (.so/.dylib files).
- `headerOnly`: Set to true to enable header-only mode. In header-only processing mode, instead of matching library symbols with header declarations, it will generate the symbol table based solely on header files specified in cflags.

After creating the configuration file, run:

```bash
llcppg llcppg.cfg
```

If you're not in a Go module or want to create a separate module, you can use the `-mod` flag to create a new Go module for the generated package:

```bash
llcppg -mod github.com/author/cjson llcppg.cfg
```

After execution, LLGo Binding will be generated in a directory named after the config name (which is also the package name). For example, with the cjson configuration above, you'll see:

```bash
cjson/
├── cJSON.go
├── cJSON_Utils.go
├── cjson_autogen_link.go
├── llcppg.pub
├── go.mod  # Contains: module github.com/author/cjson (only when using -mod flag)
└── go.sum  # Contains dependency checksums (only when using -mod flag)
```

Import the generated cjson package and try this demo:

```go
package main

import (
    "github.com/author/cjson"
    "github.com/goplus/lib/c"
)

func main() {
	mod := cjson.CreateObject()
	mod.AddItemToObject(c.Str("hello"), cjson.CreateString(c.Str("llgo")))
	mod.AddItemToObject(c.Str("hello"), cjson.CreateString(c.Str("llcppg")))
	cstr := mod.PrintUnformatted()
	c.Printf(c.Str("%s\n"), cstr)
}
```
Run the demo with `llgo run .`, you will see the following output:
```
{"hello":"llgo","hello":"llcppg"}
```

### Generated Bindings

You can see that functions from the C header files have been automatically mapped and converted to corresponding LLGo binding functions.

Original C function:
```c
CJSON_PUBLIC(cJSON *) cJSON_CreateObject(void);
```

Converted to LLGo binding function:
```go
//go:linkname CreateObject C.cJSON_CreateObject
func CreateObject() *CJSON
```

The types defined in the header file are also converted, for example the cJSON struct:

Original C struct:
```c
typedef struct cJSON
{
    struct cJSON *next;
    struct cJSON *prev;
    struct cJSON *child;
    int type;
    char *valuestring;
    int valueint;
    double valuedouble;
    char *string;
} cJSON;
```

Converted to Go struct:
```go
type CJSON struct {
	Next        *CJSON
	Prev        *CJSON
	Child       *CJSON
	Type        c.Int
	Valuestring *int8
	Valueint    c.Int
	Valuedouble float64
	String      *int8
}
```

Notably, to make the API more idiomatic in Go, when a C function's first parameter is a converted type (like cJSON *), the function is automatically converted into a method of that type.

Original C function:
```c
CJSON_PUBLIC(cJSON_bool) cJSON_AddItemToArray(cJSON *array, cJSON *item);
```

Converted to Go method:
```go
// llgo:link (*CJSON).AddItemToObject C.cJSON_AddItemToObject
func (p *CJSON) AddItemToObject(string *int8, item *CJSON) CJSONBool {
	return 0
}
```

You can also observe the corresponding type name transformations. The generated `llcppg.pub` file contains a mapping table from C types to Go type names (which will be used for package dependency handling). For the example above, the `llcppg.pub` file looks like this, where the first field on the left is the C type and the first field on the right is the corresponding Go type name.
```
cJSON CJSON
cJSON_Hooks Hooks
cJSON_bool Bool
```
 You can customize these type mappings by editing this file (see [Customizing Bindings](#type-customization)).

### Customizing Bindings
#### Function Customization
After execution, you can directly see the mapped functions in the result.

```go
//go:linkname CreateObject C.cJSON_CreateObject
func CreateObject() *CJSON
// llgo:link (*CJSON).AddItemToObject C.cJSON_AddItemToObject
func (recv_ *CJSON) AddItemToObject(string *c.Char, item *CJSON) Bool {
	return 0
}
/* Render a cJSON entity to text for transfer/storage without any formatting. */
// llgo:link (*CJSON).PrintUnformatted C.cJSON_PrintUnformatted
func (recv_ *CJSON) PrintUnformatted() *c.Char {
	return nil
}
```
You can specify function mapping behavior in `llcppg.cfg` by config the `symMap` field:
```json
{
    "symMap":{
        "mangle":"<goFuncName> | <.goMethodName> | -"
    }
}
```
`mangle` is the symbol name of the function. For the value of `mangle`, you can customize it as:
  1. `goFuncName` - generates a regular function named `goFuncName`
  2. `.goMethodName` - generates a method named `goMethodName` (if it doesn't meet the rules for generating a method, it will be generated as a regular function)
  3. `-` - completely ignore this function

For example, to convert `(*CJSON).PrintUnformatted` from a method to a function, you can use follow config:

```json
{
  "symMap":{
    "cJSON_PrintUnformatted":"PrintUnformatted"
  }
}
```
This will generate a function instead of a method in the Go code:
```go
//go:linkname PrintUnformatted C.cJSON_PrintUnformatted
func PrintUnformatted(item *CJSON) *c.Char
```

You can also run llcppsymg separately. To do this, use the command:

```sh
llcppg -symbgen
```

If you only want to generate Go code using an already generated symbol table, execute:
```sh
llcppg -codegen
```
#### Type Customization
You can customize type name mappings to better suit your needs.

For instance, if you prefer to use `JSON` instead of `cJSON` as the Go type name, simply config the `typeMap` field of `llcppg.cfg` as follows:
```json
{
  "typeMap": {
    "cJSON": "JSON"
  }
}
```
After running llcppg again, all generated code will use the new type name. The struct definition and its methods will be automatically updated:
```go
type JSON struct {
  // .....
}
// llgo:link (*JSON).PrintBuffered C.cJSON_PrintBuffered
func (recv_ *JSON) PrintBuffered(prebuffer c.Int, fmt Bool) *int8 {
	return nil
}
```

More demo projects and configuration files can be found under `_llcppgtest` directory.

### Header File Concepts
In the `llcppg.cfg`, the `include` field specifies the list of interface header files to be converted. These header files are the primary source for generating Go code, and each listed header file will generate a corresponding .go file.

```json
{
  "name": "xslt",
  "cflags": "$(pkg-config --cflags libxslt)",
  "include": [
    "libxslt/xslt.h",
    "libxslt/security.h"
  ]
}
```
#### Package Header File Determination

llcppg determines whether a header file belongs to the current package based on the following rules:

1. **Interface header files**: Header files explicitly listed in the `include` field
2. **Implementation header files**: Other header files in the same root directory as interface header files

For example, if the configuration includes `libxslt/xslt.h`, and this file contains `#include "xsltexports.h"`, then:
- `xslt.h` is an interface header file, which will generate `xslt.go`
- `xsltexports.h` is an implementation header file, whose content will be generated into `xslt_autogen.go`

Header files that don't belong to the current package (such as standard libraries or third-party dependencies) won't be directly converted but are handled through dependency relationships.

#### Special Case: Mixed Header Files
For cases where package header files are mixed with other header files in the same directory (such as system headers or third-party libraries), you can handle this by setting `mix: true`:

```json
{
  "mix": true
}
```

In this case, only header files explicitly declared in the `include` field are considered package header files, and all others are treated as third-party header files. Note that in this mode, implementation header files of the package also need to be explicitly declared in `include`, otherwise they will be treated as third-party header files and won't be processed.

This is particularly useful in scenarios like Linux systems where library headers might be installed in common directories (e.g., `/usr/include/sqlite3.h` alongside system headers like `/usr/include/stdio.h`).

### Dependency
llcppg does not convert header files outside of the current package, including any referenced third-party or standard library headers. Instead, it manages cross-package type references and ensures conversion consistency through the `deps` declaration in `llcppg.cfg`, which must include standard library types as well.
```json
{
  "deps":["c/os","github.com/author/pkg"]
}
```

#### Dependency Package Structure
Each dependency package follows a unified file organization structure (using xml2 as an example):
* Converted Go source files
1. HTMLtree.go (generated from HTMLtree.h)
2. HTMLparser.go (generated from HTMLparser.h)
* Configuration files
1. llcppg.cfg (dependency information)
2. llcppg.pub (type mapping information)

#### Dependency Handling Logic
1. llcppg scans each dependency package's `llcppg.pub` file to obtain type mappings.
2. If the dependency package's `llcppg.cfg` also contains deps configuration, llcppg will recursively process these dependencies.
3. Type mappings from all dependency packages are loaded and registered into the conversion project.
When a header file in the current project references types from third-party packages, it directly searches within the current conversion project scope
 * If a mapped type is found, it is referenced;
 * Otherwise, the user is notified of the missing type and its source header file for conversion.

#### Special Dependency Aliases
In llcppg, there is a consistent pattern for naming aliases related to the standard library. Any alias that starts with `c/` corresponds to a remote repository in the github.com/goplus/llgo.

For example:
* The alias `c` → `github.com/goplus/lib/c`
* The alias `c/os` → `github.com/goplus/lib/c/os`
* The alias `c/time` → `github.com/goplus/lib/c/time`

> Note: Standard library type conversion in LLGo is not comprehensive. For standard library types that cannot be found in LLGo, you will need to supplement these types in the corresponding package at https://github.com/goplus/llgo.

#### Example
You can specify dependent package paths in the `deps` field of `llcppg.cfg`. For example, in the `_llcppgtest/libxslt` example, since libxslt depends on libxml2, its configuration file looks like this:
```json
{
  "name": "libxslt",
  "cflags": "$(pkg-config --cflags libxslt)",
  "libs": "$(pkg-config --libs libxslt)",
  "trimPrefixes": ["xslt"],
  "deps": ["c/os","github.com/goplus/llpkg/libxml2"],
  "include":["libxslt/xsltutils.h","libxslt/templates.h"]
}
```

In `libxslt/xsltutils.h`, there are dependencies on `libxml2`'s `xmlChar` and `xmlNodePtr`:
```c
#include <libxml/dict.h>
#include <libxml/xmlerror.h>
#include <libxml/xpath.h>
xmlChar * xsltGetNsProp(xmlNodePtr node, const xmlChar *name, const xmlChar *nameSpace);
```
If `xmlChar` and `xmlNodePtr` mappings are not found (not declare `llcppg-libxml` in `deps`), llcppg will notify the user of these missing types and indicate they are from `libxml2` header files.
The corresponding notification would be:
```bash
convert /path/to/include/libxml2/libxml/xmlstring.h first, declare its converted package in llcppg.cfg deps for load [xmlChar].
convert /path/to/libxml2/libxml/tree.h first, declare its converted package in llcppg.cfg deps for load [xmlNodePtr].
```

For this project, `llcppg` will automatically handle type references to libxml2. During the process, `llcppg` uses the `llcppg.pub` file from the generated libxml2 package to ensure type consistency.
You can see this in the generated code, where libxslt correctly references libxml2's types:
```go
package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llpkg/libxml2"
	"unsafe"
)

/*
 * Our own version of namespaced attributes lookup.
 */
//go:linkname GetNsProp C.xsltGetNsProp
func GetNsProp(node libxml2.NodePtr, name *libxml2.Char, nameSpace *libxml2.Char) *libxml2.Char
```

### Important Note on Header File Ordering

llcppg follows C language's dependency resolution order when processing header files. The order of files in the `include` configuration determines the processing sequence, and incorrect ordering can lead to type resolution failures. Here's an example from the LZMA library that demonstrates the dependency relationships:

`lzma/vli.h`:
```c
typedef uint64_t lzma_vli;
```

`lzma/filter.h`:
```c
/**
 * \file        lzma/filter.h
 * \brief       Common filter related types and functions
 * \note        Never include this file directly. Use <lzma.h> instead.
 */
typedef struct {
    lzma_vli id;      // Uses lzma_vli from vli.h
    void *options;
} lzma_filter;
```

`lzma.h`:
```c
/**
 * \file        lzma.h
 * \brief       Main LZMA library header
 */
#include "lzma/vli.h"
#include "lzma/filter.h"
```

Since llcppg processes headers in the order specified in the configuration, you need to list them in the correct dependency order:

✅ Correct ordering:
```json
{
   "include": ["lzma.h", "lzma/vli.h", "lzma/filter.h"]
}
```

❌ Incorrect ordering:
```json
{
   "include": ["lzma/filter.h", "lzma.h", "lzma/vli.h"]
}
```

The incorrect ordering will fail because `filter.h` uses the `lzma_vli` type before it's defined.

#### Translation Unit Processing

When processing headers with the correct ordering, llcppg employs a caching strategy:

1. First, it creates a translation unit starting from `lzma.h`. During this process:
   - It follows normal C language parsing order for includes
   - When processing `lzma/vli.h` (included by `lzma.h`), the `lzma_vli` type definition is cached in this translation unit
   - When it reaches `lzma/filter.h`, the type references work correctly because `lzma_vli` is already defined in this translation unit

2. For subsequent includes in the configuration:
   - When processing `lzma/vli.h`, llcppg creates a new translation unit
   - When reaching `lzma/filter.h`, if it was already processed in a previous translation unit (like the one starting from `lzma.h`), llcppg will use the cached version that has proper type resolution

This caching strategy ensures that type references are properly maintained when headers are processed in the correct order.

### Development Tools

### llcppcfg - Configuration Generator

```sh
llcppcfg [libname]
```

llcppcfg tool is used to generate `llcppg.cfg` file.

## Design

See [llcppg Design](design.md).
