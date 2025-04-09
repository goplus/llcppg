llcppg Design
=====

## Usage

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
  "mix":false 
}
```

## Steps

1. llcppsymg: Generate symbol table for a C/C++ library
2. Manually modify the desired Go symbols in symbol table
3. llcppsigfetch: Fetch information of C/C++ symbols
4. gogensig: Generate a go package by information of symbols


### llcppsymg

```sh
llcppsymg config-file
llcppsymg -  # read config from stdin
```

It generates a symbol table file named `llcppg.symb.json`. Its file format is as follows:

```json
[
  {
    "mangle": "_ZN9INIReaderC1EPKcm",
    "c++": "INIReader::INIReader(char const*, unsigned long)",
    "go": "(*Reader).Init__0"
  }
]
```


### llcppsigfetch

```sh
llcppsigfetch config-file
llcppsigfetch -  # read config from stdin
```

It fetches information of C/C++ symbols and print to stdout. Its format is as follows:

```json
[
  {
    "path": "/path/to/file.h",
    "doc": {
      "decls": [],
      "macros": [],
      "includes": [
        {
          "path": "incfile.h"
        }
      ]
    }
  }
]
```

### gogensig

```sh
gogensig ast-file
gogensig -  # read AST from stdin
```

## Header File Processing Basic Rules

* llcppg only converts header files within the current package
* Does not convert third-party package and standard library header files
* For package header files, two types of mappings are generated during processing:
  * Function mapping table `llcppg.symb.json`
  * Type mapping table `llcppg.pub`

## Dependency Processing

1. The system scans each dependency package's llcppg.pub file to obtain type mappings.
2. If the dependency package's llcppg.cfg also contains deps configuration, the system will recursively process these dependencies.
3. Type mappings from all dependency packages are loaded and registered into the conversion project.
When a header file in the current project references types from third-party packages, it directly searches within the current conversion project scope
 * If a mapped type is found, it is referenced;
 * Otherwise, the user is notified of the missing type and its source header file for conversion.

### Dependency Package Structure
Each dependency package follows a unified file organization structure (using xml2 as an example):
* Converted Go source files
1. HTMLtree.go (generated from HTMLtree.h)
2. HTMLparser.go (generated from HTMLparser.h)
* Configuration files
1. llcppg.cfg (dependency information)
2. llcppg.pub (type mapping information)

### Example Project
```json
{
  "name":"xslt",
  "include": [
    "libxslt/xslt.h",
    "libxslt/security.h",
    "libexslt/exsltconfig.h"
  ],
  "deps": [
    "c",
    "github.com/goplus/..../xml2",
    "github.com/goplus/..../zlib"
  ]
}
```
In `libxslt/xslt.h`, there are dependencies on `libxml2`'s `xmlChar` and `xmlNodePtr`:
```c
#include <libxml/dict.h>
#include <libxml/xmlerror.h>
#include <libxml/xpath.h>
xmlChar * xsltGetNsProp(xmlNodePtr node, const xmlChar *name, const xmlChar *nameSpace);
```
If `xmlChar` and `xmlNodePtr` mappings are not found, llcppg will notify the user of these missing types and indicate they are from `libxml2` header files.
The corresponding notification would be:
```bash
xmlChar not found in `/path/to/libxml/dict.h`.
xmlNodePtr not found in `/path/to/libxml/dict.h`.
```
### Type Mapping Examples

Standard Library Type Mapping
NOTE: "c" is an alias for "github.com/goplus/lib/c"
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
## File Generation Rules
### Generated File Types
* Interface header files: Each header file generates a corresponding .go file
* Implementation files: All generated in a single libname_autogen.go file

### Header File Classification

1. Interface header files:
* Header files explicitly declared in include
2. Implementation header files:
* Other header files in the same root directory as interface header files

#### Example Explanation
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
## Special Case Handling
### System Path Processing
For third-party header files mixed in system paths (common in Linux apt-installed libraries, e.g., `/usr/include/sqlite3.h` and `/usr/include/stdio.h`):
1. Configuration method
```json
{
  "mix": true
}
```
2. Processing rules
* Header files explicitly declared in include are considered package header files; all others are treated as third-party header files.

### Cross-Platform Difference Handling
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