# Abstract
The llcppg core configuration file, `llcppg.cfg`, can be complex and error-prone to configure. This is because it requires a deep understanding of the project structure and compilation details. We have designed llcppcfg to automatically generate the basic `llcppg.cfg` configuration file for users. It greatly simplifies the configuration process, allowing users to simply provide the target library's name as input. The tool then generates corresponding configuration content based on established rules or templates.

# Basic Usage
`llcppcfg [options] <library actual pkg-config name>`

## Example of Generating Configuration File
`llcppcfg cjson`

This command will generate the `llcppg.cfg` configuration file in the current directory.

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

The generated `llcppg.cfg` content will be similar to:

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

# Process Design

## Header File

1. **Extract Include Paths**
   Parse the `cflags` from pkg-config to identify include paths (`-I` flags)

2. **Discover Header Files**
   Scan each include path for files:
   - Matching specified extensions (`.h`, `.hpp`, etc.)
   - Excluding specified subdirectories (`internal`, `impl`, etc.)

3. **Analyze Dependencies**
   For each valid header file:
   - Execute `clang -MM [header_path]`
   - Capture dependency output in Makefile format:
     ```
     header.o: header.cpp \
         dependency1.h \
         dependency2.h
     ```

4. **Sort by Dependency Weight**
   Prioritize headers based on:
   - Dependency count (files with more dependencies rank higher).