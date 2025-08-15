llcppg Design
=====

### Type Mapping

#### Basic Type Mapping

All basic types are imported from `github.com/goplus/lib/c` and mapped accordingly.

| C Type | Go Type |
|--------|---------|
| void | c.Void |
| bool | bool |
| char | c.Char |
| wchar_t | int16 |
| char16_t | int16 |
| char32_t | int32 |
| short | int16 |
| unsigned short | uint16 |
| int | c.Int |
| unsigned int | c.Uint |
| long | c.Long |
| unsigned long | c.Ulong |
| long long | c.LongLong |
| unsigned long long | c.UlongLong |
| float | c.Float |
| double | c.Double |
| float complex | complex64 |
| double complex | complex128 |

#### Special Case

##### void *

The pointer type `void*` is mapped to `c.Pointer`.

```c
void *(luaL_testudata) (lua_State *L, int ud, const char *tname);
```
```go
//go:linkname Testudata C.luaL_testudata
func LuaLTestudata(L *State, ud c.Int, tname *c.Char) c.Pointer
```

##### Function pointer

C function pointer types are converted to Go function types with corresponding parameter and return type mappings, and LLGo needs to add `llgo:type C` tag to the function type.


```c
typedef int (*CallBack)(void *L);
```
```go
// llgo:type C
type CallBack func(c.Pointer) c.Int
```

For function pointer types referenced in function signatures & struct fields, the type is replaced with the converted Go function type.

```c
void exec(void *L, CallBack cb);
```
```go
// llgo:type C
func Exec(L c.Pointer, cb CallBack)
```

```c
typedef struct Stream {
    CallBack cb;
} Stream;
```
```go
type Stream struct {
	Cb CallBack
}
```

For cases where a parameter in a function signature is an anonymous function pointer (meaning it does not reference a pre-defined function pointer type), it is mapped to the corresponding Go function type.

```c
int sqlite3_exec(
  sqlite3*,                                  /* An open database */
  const char *sql,                           /* SQL to be evaluated */
  int (*callback)(void*,int,char**,char**),  /* Callback function */
  void *,                                    /* 1st argument to callback */
  char **errmsg                              /* Error msg written here */
);
```
```go
// llgo:link (*Sqlite3).Exec C.sqlite3_exec
func (recv_ *Sqlite3) Exec(sql *c.Char, callback func(c.Pointer, c.Int, **c.Char, **c.Char) c.Int, __llgo_arg_2 c.Pointer, errmsg **c.Char) c.Int {
	return 0
}
```

For struct fields that are anonymous function pointers, the field type is replaced with a `c.Pointer` for description.

```c
typedef struct Hooks {
    void *(*malloc_fn)(size_t sz);
    void (*free_fn)(void *ptr);
} Hooks;
```
```go
type Hooks struct {
	MallocFn c.Pointer
	FreeFn   c.Pointer
}
```

##### Array

Arrays in C are mapped differently depending on their context - function parameters versus struct fields.

###### As Function Param

Arrays in function parameters are converted to pointers.

```c
void foo(unsigned int a[], double b[3]);
```
```go
//go:linkname Foo C.foo
func Foo(a *c.Uint, b *c.Double)
```

###### As Struct Field

Arrays in struct fields maintain their fixed-length array form to preserve memory layout compatibility with the original C struct.

```c
typedef struct Foo {
    char a[4];
    int b[3][4];
} Foo;
```
```go
type Foo struct {
	A [4]c.Char
	B [3][4]c.Int
}
```
###### Multi-dimensional

Multi-dimensional arrays are supported in both contexts, with the same conversion rules applying:

```c
char matrix[3][4];  // In function parameter becomes **c.Char
char field[3][4];   // In struct field becomes [3][4]c.Char
```

##### Nested Struct

###### Anonymous Nested Struct

Anonymous nested structs/unions are converted to inline Go struct types within the parent struct.

```c
struct outer {
    struct {
        int x;
        int y;
    } inner;
};
```

```go
type Outer struct {
    Inner struct {
        X c.Int
        Y c.Int
    }
}
```

###### Named Nested Struct

Named nested structs in C are accessible in the global scope, not just as anonymous nested types. llcppg handles this by creating separate type declarations for both the outer and inner structs.

**Reason**: In C, named nested structs are declared in the global scope and can be used independently. This means `struct inner_struct` can be used anywhere in the code, not just within the context of the outer struct.

```c
typedef struct struct2 {
    char *b;
    struct inner_struct {
        long l;
    } init;
} struct2;

// This is valid C - inner_struct is in global scope
struct inner_struct inner = {123};
```

**Generated Go code**:
```go
type InnerStruct struct {
    L c.Long
}

type Struct2 struct {
    B    *c.Char
    Init InnerStruct
}
```

This is equivalent to:
```c
struct inner_struct {
    long l;
};
struct struct2 {
    char *b;
    struct inner_struct init;
};
```

##### Nested Enum

Similar to nested structs, nested enums can also be accessed in the global scope. llcppg handles named nested enums by creating separate type declarations that are accessible globally.

###### Anonymous Nested Enum

Anonymous nested enums are converted to inline enum constants within the parent struct context, with the enum values defaulting to `c.Int` type.

```c
typedef struct Config {
    int version;
    enum {
        MODE_DEBUG = 0,
        MODE_RELEASE = 1
    } mode;
} Config;
```

**Generated Go code**:
```go
type Config struct {
    Version c.Int
    Mode    c.Int
}

const (
    MODE_DEBUG  c.Int = 0
    MODE_RELEASE c.Int = 1
)
```

###### Named Nested Enum

Named nested enums in C are accessible in the global scope, not just within the context of the outer struct. llcppg handles this by creating separate enum type declarations.

**Reason**: In C, named nested enums are declared in the global scope and can be used independently. This means the enum type can be used anywhere in the code, not just within the context of the outer struct.

**NOTE:** Should we add type alias here in the future? See discussion: [#530](https://github.com/goplus/llcppg/pull/530)

```c
typedef struct Config {
    int version;
    enum LogLevel {
        LOG_DEBUG = 0,
        LOG_INFO = 1,
        LOG_ERROR = 2
    } level;
} Config;

// This is valid C - LogLevel is in global scope
struct Config cfg;
cfg.level = LOG_INFO;
```

**Generated Go code**:
```go
type LogLevel c.Int

const (
    LOG_DEBUG LogLevel = 0
    LOG_INFO  LogLevel = 1
    LOG_ERROR LogLevel = 2
)

type Config struct {
    Version c.Int
    Level   LogLevel
}
```

This is equivalent to:
```c
enum LogLevel {
    LOG_DEBUG = 0,
    LOG_INFO = 1,
    LOG_ERROR = 2
};
struct Config {
    int version;
    enum LogLevel level;
};
```

##### Function

###### To Normal Function

LLGo needs to add `//go:linkname <funcName> C.<mangleName>` tag to the function type.

```c
void foo(int a, int b);
```
```go
//go:linkname Foo C.foo
func Foo(a c.Int, b c.Int)
```

###### To Method

When a C function could be a Go method, llcppg automatically converts the function to a Go method, moving the first parameter to the receiver position, and using recv_ as the receiver name.

Since Go's `//go:linkname` directive doesn't support methods, LLGo uses `// llgo:link` to mark the connection between methods and C symbols. And generated methods return zero values of their return types as placeholders.

And LLGo should not treat C functions with variable parameters as methods. Variadic functions (those using ... in their parameter list) will be generated as regular Go functions rather than methods, even if they otherwise meet the criteria for method conversion.

* value receiver

```c
typedef struct Vector3 {
    int x;
    int y;
    int z;
} Vector3;

Vector3 Vector3Barycenter(Vector3 p, Vector3 a, Vector3 b, Vector3 c);
```
```go
// llgo:link Vector3.Vector3Barycenter C.Vector3Barycenter
func (recv_ Vector3) Vector3Barycenter(a Vector3, b Vector3, c Vector3) Vector3 {
	return Vector3{}
}
```

* pointer receiver
```c
typedef struct sqlite3 sqlite3;
SQLITE_API int sqlite3_close(sqlite3*);
```
```go
// llgo:link (*Sqlite3).Close C.sqlite3_close
func (recv_ *Sqlite3) Close() c.Int {
	return 0
}
```

#### Name Mapping Rules

The llcppg system converts C/C++ type names to Go-compatible identifiers following specific transformation rules. These rules ensure generated Go code follows Go naming conventions while maintaining clarity and avoiding conflicts.

##### Public Name Processing
Names starting with underscore or digit are prefixed with "X" to create valid Go identifiers.

##### Type Name Conversion (struct, union, typedef, enum)

1. Remove configured prefixes from `trimPrefixes`
2. Convert to PascalCase if the name starts with a letter
3. If the name starts with an underscore, apply Public Name Processing and preserve the original case after underscores, then convert to PascalCase format

Examples without `trimPrefixes` sconfiguration:
* C: `cJSON_Hooks` → Go: `CJSONHooks`
* C: `xmlAttrHashBucket` → Go: `XmlAttrHashBucket`
* C: `sqlite3_destructor_type` → Go: `Sqlite3DestructorType`

Examples with `trimPrefixes: ["cJSON_", "sqlite3_", "xml"]`:

* C: `cJSON_Hooks` → Go: `Hooks`
* C: `sqlite3_destructor_type` → Go: `DestructorType`
* C: `xmlAttrHashBucket` → Go: `AttrHashBucket`

Examples which is start with underscore:

* C: `_gmp_err` → Go: `X_gmpErr`


##### Macro and Enum Special Rules
For macros and enums after prefix removal:

Letter-starting names: Capitalize first letter only, preserve original format
Underscore/digit-starting names: Apply public name processing, preserve original format

##### Custom Type Mappings

Types with explicit mappings in typeMap configuration bypass all other processing rules:
```json
{
  "typeMap": {
    "cJSON": "JSON"
  }
}
```
Example: C: `cJSON` → Go: `JSON`

##### Field Name Conversion

Field names must be exportable (public) in Go to allow external access. The conversion rules:

1. Letter-starting fields: Convert to PascalCase
2. Underscore/digit-starting fields: Apply public processing, then convert to PascalCase while preserving case after underscores

##### Param Name Conversion

Parameter names are preserved in their original style without conversion, with only the following special cases being handled:

1. When a parameter conflicts with a keyword, a `_` suffix is added to the parameter name

```c
void(lua_sethook)(lua_State *L, lua_Hook func, int mask, int count);
```
```go
//go:linkname Sethook C.lua_sethook
func LuaSethook(L *LuaState, func_ LuaHook, mask c.Int, count c.Int)
```

2. For variadic parameters, the parameter name is `__llgo_va_list`

```c
LUA_API int(lua_gc)(lua_State *L, int what, ...);
```
```go
//go:linkname Gc C.lua_gc
func LuaGc(L *State, what c.Int, __llgo_va_list ...interface{}) c.Int
```

3. For function signatures where all parameters have no names, the corresponding function signature will not generate parameter names.

4. Once there are named parameters in the function signature, according to Go's rules, all parameter names must be generated in the corresponding Go signature.

C allows mixing named and unnamed parameters in function signatures. For this case, the rule is to generate parameter names like `__llgo_arg_N` for unnamed parameters based on their index in the parameter list.

```c
int OSSL_PROVIDER_add_builtin(OSSL_LIB_CTX *, const char *name);
```
```go
//go:linkname ProviderAddBuiltin C.OSSL_PROVIDER_add_builtin
func OSSLProviderAddBuiltin(__llgo_arg_0 *OSSLLIBCTX, name *c.Char) c.Int
```

And for cases where only variadic parameters appear, LLGo requires ` __llgo_va_list ...interface{}` to describe variadic parameters, and the same placeholder name generation processing is needed for this case.

```c
char *mprintf(const char*,...);
```
```go
//go:linkname Mprintf C.mprintf
func Mprintf(__llgo_arg_0 *c.Char, __llgo_va_list ...interface{}) *c.Char
```


### File Generation Rules

#### Generated File Types
* Interface header files: Each header file generates a corresponding .go file
* Implementation files: All generated in a single libname_autogen.go file
* Third-party header files: Skip generation, only as a dependency

#### Header File Concepts
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
##### Package Header File Determination

llcppg determines whether a header file belongs to the current package based on the following rules:

1. **Interface header files**: Header files explicitly listed in the `include` field
2. **Implementation header files**: Other header files in the same root directory as interface header files
3. **Third-party header files**: Header files that don't belong to the current package (such as standard libraries or third-party dependencies) won't be directly converted but are handled through dependency relationships.

For example, if the configuration includes `libxslt/xslt.h`, and this file contains `#include "xsltexports.h"`, then:
- `xslt.h` is an interface header file, which will generate `xslt.go`
- `xsltexports.h` is an implementation header file, whose content will be generated into `xslt_autogen.go`

###### Example Explanation

For example, the header file paths obtained after linking with Clang in the above example:
```
/opt/homebrew/Cellar/libxslt/1.1.42_1/include/libxslt/xslt.h
/opt/homebrew/Cellar/libxslt/1.1.42_1/include/libxslt/security.h
/opt/homebrew/Cellar/libxslt/1.1.42_1/include/libexslt/exsltconfig.h
```
The calculated common root directory is:
```
/opt/homebrew/Cellar/libxslt/1.1.42_1/include/
```
In `libxslt/xslt.h`, the following header files are referenced:
```c
#include <libxml/tree.h>
#include "xsltexports.h"
```
The corresponding paths are:
`libxml/tree.h` -> `/opt/homebrew/Cellar/libxml2/2.13.5/include/libxml2/libxml/tree.h` (third-party dependency)
`xsltexports.h` -> `/opt/homebrew/Cellar/libxslt/1.1.42_1/include/libxslt/xsltexports.h` (package implementation file)
Since `xsltexports.h` is in the same directory as `libxslt/xslt.h`, it's considered a package implementation file, and its content is generated in `xslt_autogen.go`. Meanwhile, `libxml/tree.h` is not in the same directory and is considered a third-party dependency.

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
1. `llcppg.cfg` (dependency information)
2. `llcppg.pub` (type mapping information)

##### TypeMapping Examples (`llcppg.pub`)

* C types on the left and corresponding Go type names on the right
* If the Go Name is same with C type name, only need keep one column

Standard Library Type Mapping
`github.com/goplus/lib/c/llcppg.pub`
```
size_t SizeT
intptr_t IntptrT
FILE
```
XML2 Type Mapping From Expamle
`github.com/goplus/..../xml2/llcppg.pub`
```
xml_doc XmlDoc
```

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
convert /path/to/include/libxml2/libxml/xmlstring.h first, declare its converted package in `llcppg.cfg` deps for load [xmlChar].
convert /path/to/libxml2/libxml/tree.h first, declare its converted package in `llcppg.cfg` deps for load [xmlNodePtr].
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

### Cross-Platform Difference Handling

C/C++ libraries often use conditional compilation to provide different implementations for different platforms. This creates platform-specific differences in header file content that must be handled correctly during binding generation.

Consider a simple C library header file that adapts to different platforms:

```c
// mock_platform.h - Simple cross-platform example
typedef struct PlatformData {
    int common_field;
#ifdef __APPLE__
    int mac_field;
#elif defined(__linux__)
    int linux_field;
#endif
} PlatformData;

#ifdef __APPLE__
void mac_function(int x);
#else
void other_function(int x);
#endif
```

When llcppg processes this header file, the actual content varies between platforms:
- **macOS**: The struct contains `mac_field` and `mac_function` is available
- **Linux**: The struct contains `linux_field` and `other_function` is available

#### Configuration Solution

Use impl.files configuration to handle platform differences:
```json
{
    "impl": [
        {
            "files": ["t1.h", "t2.h"],
            "cond": {
                "os": ["macos", "linux"],
                "arch": ["arm64", "amd64"]
            }
        }
    ]
}
```

This configuration tells llcppg to generate platform-specific versions of the specified header files, ensuring that the correct platform-specific definitions are used for each target platform.

#### Generated Platform-Specific Files

The generated t1.go & t2.go files will have platform-specific build tags at the beginning:
macos arm64 t1_macos_arm64.go  t2_macos_arm64.go
```go
// +build macos,arm64
package xxx
```
linux arm64 `t1_linux_arm64.go`  `t2_linux_arm64.go`
```go
// +build linux,arm64
package xxx
```
macos amd64  `t1_macos_amd64.go`  `t2_macos_amd64.go`
```go
// +build macos,amd64
package xxx
```
linux amd64 `t1_linux_amd64.go`  `t2_linux_amd64.go`
```go
// +build linux,amd64
package xxx
```

## Input

```sh
llcppg [config-file]
```

If `config-file` is not specified, a `llcppg.cfg` file is used in current directory. The configuration file format is as follows:

```json
{
  "name": "inih",
  "cflags": "$(pkg-config --cflags inireader)",
  "include": [
    "INIReader.h",
    "AnotherHeaderFile.h"
  ],
  "libs": "$(pkg-config --libs inireader)",
  "trimPrefixes": ["Ini", "INI"],
  "cplusplus":true,
  "deps":["c","github.com/..../third"],
  "mix":false,
  "staticLib": false,
  "headerOnly": false
}
```

The configuration file supports the following options:

- `name`: The name of the generated package
- `cflags`: Compiler flags for the C/C++ library
- `include`: Header files to include in the binding generation
- `libs`: Library flags for linking
- `trimPrefixes`: Prefixes to remove from function names & type names
- `cplusplus`: Set to true for C++ libraries(not support)
- `deps`: Dependencies (other packages & standard libraries)
- `mix`: Set to true when package header files are mixed with other header files in the same directory. In this mode, only files explicitly listed in `include` are processed as package files.
- `typeMap`: Custom name mapping from C types to Go types.
- `symMap`: Custom name mapping from C function names to Go function names.
- `staticLib`: Set to true to enable static library symbol reading instead of dynamic library linking. When enabled, llcppg will read symbols from static libraries (.a files) rather than dynamic libraries (.so/.dylib files).
- `headerOnly`: Set to true to ​​skip the symbol intersection process​​ described in [step 3](#llcppsymg).

## Output

After running llcppg, LLGo bindings will be generated in a directory named by `name` field in `llcppg.cfg`, and the `name` field is also the package name of all Go files. The generated file structure is as follows:

### Go Source Files

#### C Header File Corresponding Go File

* A corresponding .go file is generated for each header file listed in the `include` field in `llcppg.cfg`.
* File names are based on header file names, e.g., cJSON.h generates cJSON.go, cJSON_Utils.h generates cJSON_Utils.go
* Implementation files are all generated at `{name}_autogen.go` file, determined file type by [Package Header File Determination](#Package-Header-File-Determination)

#### Auto generated Link File

* Generates a `{name}_autogen_link.go` file containing linking information and necessary imports
* This file includes the `LLGoPackage` constant to specify the lib link flags from `libs` field in `llcppg.cfg`. for example: `"libs": "$(pkg-config --libs libxslt)"`, will generate:

```json
{
  "libs": "$(pkg-config --libs libxslt)"
}
```
```go
const LLGoPackage string = "link: $(pkg-config --libs libxslt);"
```

* blank import for every dependency package in `deps` field in `llcppg.cfg`, for example:

```json
{
  "deps": ["c/os","github.com/goplus/llpkg/libxml2@v1.0.1"]
}
```
```go
import (
	_ "github.com/goplus/lib/c"
	_ "github.com/goplus/lib/c/os"
	_ "github.com/goplus/llpkg/libxml2"
)
```

### Type Mapping File

* Generates an `llcppg.pub` file containing a mapping table from C types to Go type names, is used for package dependency handling, example and concept see [Dependency](#Dependency)

### Go Module Files (Optional)

* When using the `-mod` flag, go.mod and go.sum files are generated with the specified module name.

for example: `llcppg -mod github.com/author/cjson` will generate a go.mod file with the module name `github.com/author/cjson` and a go.sum file in the result directory.

### Example

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

Using the cjson configuration as an example, the generated directory structure would be:

```bash
cjson/
├── cJSON.go                    # Bindings generated from cJSON.h
├── cJSON_Utils.go             # Bindings generated from cJSON_Utils.h
├── cjson_autogen_link.go      # Auto-generated link file
├── llcppg.pub                 # Type mapping information
├── go.mod                     # Go module file (when using -mod flag)
└── go.sum                     # Dependency checksums (when using -mod flag)
```

## Process Steps

The llcppg tool orchestrates a three-stage pipeline that automatically generates Go bindings for C/C++ libraries by coordinating symbol table generation, signature extraction, and Go code generation components.

1. llcppsymg: Generate symbol table for a C/C++ library
2. llcppsigfetch: Fetch information of C/C++ symbols
3. gogensig: Generate a Go package by information of symbols

### llcppsymg

```sh
llcppsymg config-file
llcppsymg -  # read config from stdin
```

llcppsymg is the symbol table generator in the llcppg toolchain, responsible for analyzing C/C++ libraries (dynamic or static) and header files to generate symbol mapping tables. Its main functions are:

1. Parse library symbols: Extract exported symbols from libraries using the nm tool
2. Parse header file declarations: Analyze C/C++ header files using libclang for function declarations
3. Find intersection: Match library symbols with header declarations and then generate symbol table named `llcppg.symb.json`.

#### Static Library Support

When `staticLib: true` is configured in `llcppg.cfg`, llcppsymg switches to static library mode:

- **Dynamic Library Mode (default)**: Uses nm tool to extract symbols from .so/.dylib files
- **Static Library Mode**: Uses nm tool to extract symbols from .a files

#### Header-Only Mode

When `headerOnly: true` is configured in `llcppg.cfg`, llcppg operates in header-only processing mode.

In header-only processing mode, instead of matching library symbols with header declarations, it will generate the symbol table based solely on header files specified in cflags.

#### Symbol Table

This symbol table determines whether the function appears in the generated Go code、its actual name and if it is a method. Its file format is as follows:

```json
[
  {
    "mangle": "cJSON_Delete",
    "c++": "cJSON_Delete(cJSON *)",
    "go": "(*CJSON).Delete"
  },
]
```

* mangle: mangled name of function
* c++: C/C++ function prototype declaration string
* go: corresponding Go function or method name, during the process, llcppg will automatically check if the current function can be a method
  1. When go is "-", the function is ignored (not generated)
  2. When go is a valid function name, the function name will be named as the mangle
  3. When go is `(*Type).MethodName` or `Type.MethodName`, the function will be generated as a method with Receiver as Type/*Type, and Name as MethodName

#### Custom Symbol Table generation

Specify function mapping behavior in `llcppg.cfg` by config the `symMap` field:
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
and the `llcppg.symb.json` will be:
```json
[
  {
    "mangle": "cJSON_PrintUnformatted",
    "c++": "cJSON_PrintUnformatted(cJSON *)",
    "go": "PrintUnformatted"
  }
]
```

### llcppsigfetch

llcppsigfetch is a tool that extracts type information and function signatures from C/C++ header files. It uses Clang & Libclang to parse C/C++ header files and outputs a JSON-formatted package information structure.

```sh
llcppsigfetch config-file
llcppsigfetch -  # read config from stdin
```

* Preprocesses C/C++ header files
* Creates translation units using libclang and traverses the preprocessed header file to extract ast info.

#### Output:

The output is a `pkg-info` structure that contains comprehensive package information needed for Go code generation. This `pkg-info` consists of two main components:

* File: Contains the AST with decls, includes, and macros.
* FileMap: Maps file paths to file types, where FileType indicates file classification (interface, implementation, or third-party files)

```json
{
    "File": {
        "decls": [],
        "include": [],
        "macros": []
    },
    "FileMap": {
        "usr/include/sys/_types/_rsize_t.h": {
            "FileType": 3
        },
        "/opt/homebrew/include/lua/lua.h": {
            "FileType": 1
        },
        "/opt/homebrew/include/lua/luaconf.h": {
            "FileType": 2
        }
    }
}
```

### gogensig

gogensig is the final component in the pipeline, responsible for converting C/C++ type declarations and function signatures into Go code. It reads the `pkg-info` structure generated by llcppsigfetch.

```sh
gogensig pkg-info-file
gogensig -  # read pkg-info-file from stdin
```

#### Function Generation
During execution, gogensig only generates functions whose corresponding mangle exists in `llcppg.symb.json`, determining whether to generate functions/methods with specified Go names by parsing the go field corresponding to the mangle.

1. Regular function format: "FunctionName"
  * Generates regular functions, using `//go:linkname` annotation
2. Pointer receiver method format: "(*TypeName).MethodName"
  * Generates methods with pointer receivers, using `// llgo:link` annotation
3. Value receiver method format: "TypeName.MethodName"
  * Generates methods with value receivers, using `// llgo:link` annotation
4. Ignore function format: "-"
  * Completely ignores the function, generates no code
