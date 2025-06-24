# Abstract
The LLCppg core configuration file, llcppg.cfg, can be complex and error-prone to configure. This is because it requires a deep understanding of the project structure and compilation details. We have designed llcppcfg to automatically generate the basic llcppg.cfg configuration file for users. It greatly simplifies the configuration process, allowing users to simply provide the target library's name as input. The tool then generates corresponding configuration content based on established rules or templates.

# Installation
`go install github.com/goplus/llcppg/cmd/llcppcfg`

# Basic Usage
`llcppcfg [options] <library actual PC name>`

## Example of Generating Configuration File
`llcppcfg cjson`

This command will generate the llcppg.cfg configuration file in the current directory.

## Command Line Option Details

| Option      | Default  | Description                                                                 |
|------------|----------|-----------------------------------------------------------------------------|
| `-cpp`     | false    | Specifies this is a C++ library (generates C++ related configuration when true) |
| `-tab`     | true     | Uses tab indentation to format the output configuration file                  |
| `-exts`    | ".h"     | List of included header file extensions (e.g., `-exts=".h .hpp .hh"`)        |
| `-excludes`| ""       | Excluded subdirectories (e.g., `-excludes="internal private"` to exclude these directories) |
| `-deps`    | ""       | Dependency library list (e.g., `-deps="zlib libssl"`)                        |
| `-help`    | false    | Displays help information                                                    |

### Advanced Usage Examples
Generate configuration file for C++ library:

`llcppcfg -cpp -exts=".h .hpp .hh" opencv`

Customize header file extensions and exclude specific directories:

`llcppcfg -exts=".h .hpp" -excludes="internal impl" curl`

Specify dependent libraries:

`llcppcfg -deps="github.com/goplus/llpkg/zlib@v1.0.2" openssl`

### Configuration File Generation Example
After executing the following command:

`llcppcfg -cpp -deps="github.com/goplus/llpkg/zlib@v1.0.2" -exts=".h .hpp" openssl`

The generated llcppg.cfg content will be similar to:

```json
{
	"name": "openssl",
	"cflags": "$(pkg-config --cflags openssl)",
	"libs": "$(pkg-config --libs openssl)",
	"cplusplus": true,
	"include": [
		"ssl.h",
		"crypto.h",
		// ...other header files
	],
	"deps": ["github.com/goplus/llpkg/zlib@v1.0.2"]
}
```

# Design

## Header File Acquisition

Based on the user-provided PC name, we can obtain its `cflags` information

Subsequently, based on the `cflags` information, we find the header file storage path and scan all files that meet the requirements

Then, by calling the `clang` preprocessor, we obtain its dependency graph and sort the header files based on the dependency graph information
