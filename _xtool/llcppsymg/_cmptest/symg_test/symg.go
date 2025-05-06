package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/_xtool/llcppsymg/names"
	"github.com/goplus/llcppg/_xtool/llcppsymg/symg"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llgo/xtool/nm"
)

func main() {
	// parse header file
	TestNewSymbolProcessor()
	TestGenMethodName()
	TestAddSuffix()
	TestParseHeaderFile()

	// syspath
	TestLdOutput()
	TestClangIncOutput()

	// test full process
	TestGen()
}

func TestLdOutput() {
	fmt.Println("=== TestLdOutput ===")
	res := symg.ParseLdOutput(
		`GNU ld (GNU Binutils for Ubuntu) 2.42
  Supported emulations:
   aarch64linux
   aarch64elf
   aarch64elf32
   aarch64elf32b
   aarch64elfb
   armelf
   armelfb
   aarch64linuxb
   aarch64linux32
   aarch64linux32b
   armelfb_linux_eabi
   armelf_linux_eabi
using internal linker script:
==================================================
/* Script for -z combreloc */
/* Copyright (C) 2014-2024 Free Software Foundation, Inc.
   Copying and distribution of this script, with or without modification,
   are permitted in any medium without royalty provided the copyright
   notice and this notice are preserved.  */
OUTPUT_FORMAT("elf64-littleaarch64", "elf64-bigaarch64",
              "elf64-littleaarch64")
OUTPUT_ARCH(aarch64)
ENTRY(_start)
SEARCH_DIR("=/usr/local/lib/aarch64-linux-gnu"); SEARCH_DIR("=/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/local/lib"); SEARCH_DIR("=/lib"); SEARCH_DIR("=/usr/lib"); SEARCH_DIR("=/usr/aarch64-linux-gnu/lib");
SECTIONS
{
  /* Read-only sections, merged into text segment: */
  PROVIDE (__executable_start = SEGMENT_START("text-segment", 0x400000)); . = SEGMENT_START("text-segment", 0x400000) + SIZEOF_HEADERS;
  .interp         : { *(.interp) }
  .note.gnu.build-id  : { *(.note.gnu.build-id) }
  .hash           : { *(.hash) }
  .gnu.hash       : { *(.gnu.hash) }
  .dynsym         : { *(.dynsym) }
  .dynstr         : { *(.dynstr) }
`)
	fmt.Println(res)
}

func TestClangIncOutput() {
	fmt.Println("=== TestClangIncOutput ===")
	res := symg.ParseClangIncOutput(
		`Ubuntu clang version 18.1.3 (1ubuntu1)
Target: aarch64-unknown-linux-gnu
Thread model: posix
InstalledDir: /usr/bin
Found candidate GCC installation: /usr/bin/../lib/gcc/aarch64-linux-gnu/13
Selected GCC installation: /usr/bin/../lib/gcc/aarch64-linux-gnu/13
Candidate multilib: .;@m64
Selected multilib: .;@m64
 (in-process)
 "/usr/lib/llvm-18/bin/clang" -cc1 -triple aarch64-unknown-linux-gnu -E -disable-free -clear-ast-before-backend -disable-llvm-verifier -discard-value-names -main-file-name null -mrelocation-model pic -pic-level 2 -pic-is-pie -mframe-pointer=non-leaf -fmath-errno -ffp-contract=on -fno-rounding-math -mconstructor-aliases -funwind-tables=2 -target-cpu generic -target-feature +v8a -target-feature +fp-armv8 -target-feature +neon -target-abi aapcs -debugger-tuning=gdb -fdebug-compilation-dir=/root/llcppg -v -fcoverage-compilation-dir=/root/llcppg -resource-dir /usr/lib/llvm-18/lib/clang/18 -internal-isystem /usr/lib/llvm-18/lib/clang/18/include -internal-isystem /usr/local/include -internal-isystem /usr/bin/../lib/gcc/aarch64-linux-gnu/13/../../../../aarch64-linux-gnu/include -internal-externc-isystem /usr/include/aarch64-linux-gnu -internal-externc-isystem /include -internal-externc-isystem /usr/include -ferror-limit 19 -fno-signed-char -fgnuc-version=4.2.1 -fskip-odr-check-in-gmf -fcolor-diagnostics -target-feature +outline-atomics -target-feature -fmv -faddrsig -D__GCC_HAVE_DWARF2_CFI_ASM=1 -o - -x c /dev/null
clang -cc1 version 18.1.3 based upon LLVM 18.1.3 default target aarch64-unknown-linux-gnu
ignoring nonexistent directory "/usr/bin/../lib/gcc/aarch64-linux-gnu/13/../../../../aarch64-linux-gnu/include"
ignoring nonexistent directory "/include"
#include "..." search starts here:
#include <...> search starts here:
 /usr/lib/llvm-18/lib/clang/18/include
 /usr/local/include
 /usr/include/aarch64-linux-gnu
 /usr/include
End of search list.
# 1 "/dev/null"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 399 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "/dev/null" 2
`)
	fmt.Println(res)
}

func TestNewSymbolProcessor() {
	fmt.Println("=== Test NewSymbolProcessor ===")
	process := symg.NewSymbolProcessor([]string{}, []string{"lua_", "luaL_"}, nil)
	fmt.Printf("Before: No prefixes After: Prefixes: %v\n", process.Prefixes)
	fmt.Println()
}

func TestGenMethodName() {
	fmt.Println("=== Test GenMethodName ===")
	process := &symg.SymbolProcessor{}

	testCases := []struct {
		class        string
		name         string
		isDestructor bool
	}{
		{"INIReader", "INIReader", false},
		{"INIReader", "INIReader", true},
		{"INIReader", "HasValue", false},
	}
	for _, tc := range testCases {
		input := fmt.Sprintf("Class: %s, Name: %s", tc.class, tc.name)
		result := process.GenMethodName(tc.class, tc.name, tc.isDestructor, true)
		fmt.Printf("Before: %s After: %s\n", input, result)
	}
	fmt.Println()
}

func TestAddSuffix() {
	fmt.Println("=== Test AddSuffix ===")
	process := symg.NewSymbolProcessor([]string{}, []string{"INI"}, nil)
	methods := []string{
		"INIReader",
		"INIReader",
		"ParseError",
		"HasValue",
	}
	for _, method := range methods {
		goName := names.GoName(method, process.Prefixes, true)
		className := names.GoName("INIReader", process.Prefixes, true)
		methodName := process.GenMethodName(className, goName, false, true)
		finalName := process.AddSuffix(methodName)
		input := fmt.Sprintf("Class: INIReader, Method: %s", method)
		fmt.Printf("Before: %s After: %s\n", input, finalName)
	}
	fmt.Println()
}

func TestParseHeaderFile() {
	testCases := []struct {
		name     string
		content  string
		isCpp    bool
		prefixes []string
	}{
		{
			name: "C++ Class with Methods",
			content: `
class INIReader {
  public:
    INIReader(const std::string &filename);
    INIReader(const char *buffer, size_t buffer_size);
    ~INIReader();
    int ParseError() const;
  private:
    static std::string MakeKey(const std::string &section, const std::string &name);
};
            `,
			isCpp:    true,
			prefixes: []string{"INI"},
		},
		{
			name: "C Functions",
			content: `
typedef struct lua_State lua_State;
int(lua_rawequal)(lua_State *L, int idx1, int idx2);
int(lua_compare)(lua_State *L, int idx1, int idx2, int op);
int(lua_sizecomp)(size_t s, int idx1, int idx2, int op);
            `,
			isCpp:    false,
			prefixes: []string{"lua_"},
		},
		{
			name: "InvalidReceiver",
			content: `
			typedef struct sqlite3 sqlite3;
			typedef const char *sqlite3_filename;
			SQLITE_API const char *sqlite3_uri_parameter(sqlite3_filename z, const char *zParam);
			SQLITE_API int sqlite3_errcode(sqlite3 *db);
			            `,
			isCpp:    false,
			prefixes: []string{"sqlite3_"},
		},
		{
			name: "InvalidReceiver PointerLevel > 1",
			content: `
			typedef struct asn1_node_st asn1_node_st;
			typedef asn1_node_st *asn1_node;
			extern ASN1_API int asn1_der_decoding (asn1_node * element, const void *ider, int ider_len, char *errorDescription);
						`,
			isCpp:    false,
			prefixes: []string{"asn1_"},
		},
		{
			name: "InvalidReceiver typ.NamedType.String is empty",
			content: `
			RLAPI void InitWindow(int width, int height, const char *title);
			`,
			isCpp:    false,
			prefixes: []string{""},
		},
		{
			name: "InvalidReceiver typ.canonicalType.Kind == clang.TypePointer",
			content: `
			typedef struct
			{
			int _mp_alloc;		/* Number of *limbs* allocated and pointed
							to by the _mp_d field.  */
			int _mp_size;			/* abs(_mp_size) is the number of limbs the
							last field points to.  If _mp_size is
							negative this is a negative number.  */
			} __mpz_struct;
			typedef __mpz_struct *mpz_ptr;
			inline void __mpz_set_ui_safe(mpz_ptr p, unsigned long l)
{
  p->_mp_size = (l != 0);
  p->_mp_d[0] = l & GMP_NUMB_MASK;
#if __GMPZ_ULI_LIMBS > 1
  l >>= GMP_NUMB_BITS;
  p->_mp_d[1] = l;
  p->_mp_size += (l != 0);
#endif
}
			`,
			isCpp:    false,
			prefixes: []string{""},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test Case: %s ===\n", tc.name)

		symbolMap, err := symg.ParseHeaderFile([]string{tc.content}, tc.prefixes, []string{}, nil, tc.isCpp, true)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("Parsed Symbols:")

		var keys []string
		for key := range symbolMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			info := symbolMap[key]
			fmt.Printf("Symbol Map GoName: %s, ProtoName In HeaderFile: %s, MangledName: %s\n", info.GoName, info.ProtoName, key)
		}
		fmt.Println()
	}
}

func TestGen() {
	testCases := []struct {
		name         string
		path         string
		dylibSymbols []string
	}{
		{
			name: "c",
			path: "./c",
			dylibSymbols: []string{
				"Foo_Print",
				"Foo_ParseWithLength",
				"Foo_Delete",
				"Foo_ParseWithSize",
				"Foo_ignoreFunc",
				"Foo_Bar",
				"Foo_ForBar",
				"Foo_Bar2",
				"Foo_ForBar2",
				"Foo_Prefix_BarMethod",
				"Foo_BarMethod",
				"Foo_ForBarMethod",
				"Foo_ReceiverParse",
				"Foo_FunctionParse",
				"Foo_ReceiverParse2",
				"Foo_Receiver2Parse2",
			},
		},
		{
			name: "cpp",
			path: "./cpp",
			dylibSymbols: []string{
				"ZN3FooC1EPKc",
				"ZN3FooC1EPKcl",
				"ZN3FooD1Ev",
				"ZNK3Foo8ParseBarEv",
				"ZNK3Foo3GetEPKcS1_S1_",
				"ZN3Foo6HasBarEv",
			},
		},
		{
			name: "inireader",
			path: "./inireader",
			dylibSymbols: []string{
				"ZN9INIReaderC1EPKc",
				"ZN9INIReaderC1EPKcl",
				"ZN9INIReaderD1Ev",
				"ZNK9INIReader10ParseErrorEv",
				"ZNK9INIReader3GetEPKcS1_S1_",
			},
		},
		{
			name: "lua",
			path: "./lua",
			dylibSymbols: []string{
				"lua_error",
				"lua_next",
				"lua_concat",
				"lua_stringtonumber",
			},
		},
		{
			name: "cjson",
			path: "./cjson",
			dylibSymbols: []string{
				"cJSON_Print",
				"cJSON_ParseWithLength",
				"cJSON_Delete",
				// mock multiple symbols
				"cJSON_Delete",
			},
		},
		{
			name: "isl",
			path: "./isl",
			dylibSymbols: []string{
				"isl_pw_qpolynomial_get_ctx",
			},
		},
		{
			name: "gpgerror",
			path: "./gpgerror",
			dylibSymbols: []string{
				"gpg_strsource",
				"gpg_strerror_r",
				"gpg_strerror",
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test Case: %s ===\n", tc.name)
		projPath, err := filepath.Abs(tc.path)
		if err != nil {
			fmt.Println("Get Abs Path Error:", err)
		}
		cfgdata, err := os.ReadFile(filepath.Join(projPath, llcppg.LLCPPG_CFG))
		if err != nil {
			fmt.Println("Read Cfg File Error:", err)
		}
		cfg, err := config.GetConf(cfgdata)
		if err != nil {
			fmt.Println("Get Conf Error:", err)
		}
		if err != nil {
			fmt.Println("Read Symb File Error:", err)
		}

		cfg.CFlags = "-I" + projPath
		pkgHfileInfo := config.PkgHfileInfo(cfg.Config, []string{})
		headerSymbolMap, err := symg.ParseHeaderFile(pkgHfileInfo.CurPkgFiles(), cfg.TrimPrefixes, strings.Fields(cfg.CFlags), cfg.SymMap, cfg.Cplusplus, false)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if err != nil {
			fmt.Printf("Failed to create temp file: %v\n", err)
			return
		}

		// trim to nm symbols
		var dylibsymbs []*nm.Symbol
		for _, symb := range tc.dylibSymbols {
			dylibsymbs = append(dylibsymbs, &nm.Symbol{Name: symg.AddSymbolPrefixUnder(symb, cfg.Cplusplus)})
		}
		symbolData, err := symg.GenerateSymTable(dylibsymbs, headerSymbolMap)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(string(symbolData))
	}
}
