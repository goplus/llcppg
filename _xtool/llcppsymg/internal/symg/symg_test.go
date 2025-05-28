package symg_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/config"
	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/symg"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/name"
	"github.com/goplus/llgo/xtool/nm"
)

func TestNewSymbolProcessor(t *testing.T) {
	process := symg.NewSymbolProcessor([]string{}, []string{"lua_", "luaL_"}, nil)
	expect := []string{"lua_", "luaL_"}
	if !reflect.DeepEqual(process.Prefixes, expect) {
		t.Fatalf("expect %v, but got %v", expect, process.Prefixes)
	}
}

func TestAddSuffix(t *testing.T) {
	process := symg.NewSymbolProcessor([]string{}, []string{"INI"}, nil)
	methods := []struct {
		method string
		expect string
	}{
		{"INIReader", "(*Reader).Init"},
		{"INIReader", "(*Reader).Init__1"},
		{"ParseError", "(*Reader).ParseError"},
		{"HasValue", "(*Reader).HasValue"},
	}
	for _, tc := range methods {
		t.Run(tc.method, func(t *testing.T) {
			goName := name.GoName(tc.method, process.Prefixes, true)
			className := name.GoName("INIReader", process.Prefixes, true)
			methodName := process.GenMethodName(className, goName, false, true)
			actual := process.AddSuffix(methodName)
			if actual != tc.expect {
				t.Fatalf("expect %s, but got %s", tc.expect, actual)
			}
		})
	}
}

func TestGenMethodName(t *testing.T) {
	process := &symg.SymbolProcessor{}

	testCases := []struct {
		class        string
		name         string
		isDestructor bool
		expect       string
	}{
		{"INIReader", "INIReader", false, "(*INIReader).Init"},
		{"INIReader", "INIReader", true, "(*INIReader).Dispose"},
		{"INIReader", "HasValue", false, "(*INIReader).HasValue"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			result := process.GenMethodName(tc.class, tc.name, tc.isDestructor, true)
			if result != tc.expect {
				t.Fatalf("expect %s, but got %s", tc.expect, result)
			}
		})
	}
}

func TestGetCommonSymbols(t *testing.T) {
	testCases := []struct {
		name          string
		dylibSymbols  []*nm.Symbol
		headerSymbols map[string]*symg.SymbolInfo
		expect        []*llcppg.SymbolInfo
	}{
		{
			name: "Lua symbols",
			dylibSymbols: []*nm.Symbol{
				{Name: symg.AddSymbolPrefixUnder("lua_absindex", false)},
				{Name: symg.AddSymbolPrefixUnder("lua_arith", false)},
				{Name: symg.AddSymbolPrefixUnder("lua_atpanic", false)},
				{Name: symg.AddSymbolPrefixUnder("lua_callk", false)},
				{Name: symg.AddSymbolPrefixUnder("lua_lib_nonexistent", false)},
			},
			headerSymbols: map[string]*symg.SymbolInfo{
				"lua_absindex":           {ProtoName: "lua_absindex(lua_State *, int)", GoName: "Absindex"},
				"lua_arith":              {ProtoName: "lua_arith(lua_State *, int)", GoName: "Arith"},
				"lua_atpanic":            {ProtoName: "lua_atpanic(lua_State *, lua_CFunction)", GoName: "Atpanic"},
				"lua_callk":              {ProtoName: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", GoName: "Callk"},
				"lua_header_nonexistent": {ProtoName: "lua_header_nonexistent()", GoName: "HeaderNonexistent"},
			},
			expect: []*llcppg.SymbolInfo{
				{Mangle: "lua_absindex", CPP: "lua_absindex(lua_State *, int)", Go: "Absindex"},
				{Mangle: "lua_arith", CPP: "lua_arith(lua_State *, int)", Go: "Arith"},
				{Mangle: "lua_atpanic", CPP: "lua_atpanic(lua_State *, lua_CFunction)", Go: "Atpanic"},
				{Mangle: "lua_callk", CPP: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", Go: "Callk"},
			},
		},
		{
			name: "INIReader and Std library symbols",
			dylibSymbols: []*nm.Symbol{
				{Name: symg.AddSymbolPrefixUnder("ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x", true)},
				{Name: symg.AddSymbolPrefixUnder("ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d", true)},
				{Name: symg.AddSymbolPrefixUnder("ZNK9INIReader10ParseErrorEv", true)},
			},
			headerSymbols: map[string]*symg.SymbolInfo{
				"_ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x":  {GoName: "(*Reader).GetInteger64", ProtoName: "INIReader::GetInteger64(const std::string &, const std::string &, int64_t)"},
				"_ZNK9INIReader13GetUnsigned64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_y": {GoName: "(*Reader).GetUnsigned64", ProtoName: "INIReader::GetUnsigned64(const std::string &, const std::string &, uint64_t)"},
				"_ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d":        {GoName: "(*Reader).GetReal", ProtoName: "INIReader::GetReal(const std::string &, const std::string &, double)"},
				"_ZNK9INIReader10ParseErrorEv": {GoName: "(*Reader).ParseError", ProtoName: "INIReader::ParseError()"},
				"_ZNK9INIReader10GetBooleanERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_b": {GoName: "(*Reader).GetBoolean", ProtoName: "INIReader::GetBoolean(const std::string &, const std::string &, bool)"},
			},
			expect: []*llcppg.SymbolInfo{
				{Mangle: "_ZNK9INIReader10ParseErrorEv", CPP: "INIReader::ParseError()", Go: "(*Reader).ParseError"},
				{Mangle: "_ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x", CPP: "INIReader::GetInteger64(const std::string &, const std::string &, int64_t)", Go: "(*Reader).GetInteger64"},
				{Mangle: "_ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d", CPP: "INIReader::GetReal(const std::string &, const std::string &, double)", Go: "(*Reader).GetReal"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commonSymbols := symg.GetCommonSymbols(tc.dylibSymbols, tc.headerSymbols)
			if !reflect.DeepEqual(commonSymbols, tc.expect) {
				t.Fatalf("expect %v, but got %v", tc.expect, commonSymbols)
			}
		})
	}
}

func TestGenSymbolTableData(t *testing.T) {
	commonSymbols := []*llcppg.SymbolInfo{
		{Mangle: "lua_absindex", CPP: "lua_absindex(lua_State *, int)", Go: "Absindex"},
		{Mangle: "lua_arith", CPP: "lua_arith(lua_State *, int)", Go: "Arith"},
		{Mangle: "lua_atpanic", CPP: "lua_atpanic(lua_State *, lua_CFunction)", Go: "Atpanic"},
		{Mangle: "lua_callk", CPP: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", Go: "Callk"},
	}

	data, err := symg.GenSymbolTableData(commonSymbols)
	if err != nil {
		t.Fatal(err)
	}
	expect := strings.TrimSpace(`
[{
		"mangle":	"lua_absindex",
		"c++":	"lua_absindex(lua_State *, int)",
		"go":	"Absindex"
	}, {
		"mangle":	"lua_arith",
		"c++":	"lua_arith(lua_State *, int)",
		"go":	"Arith"
	}, {
		"mangle":	"lua_atpanic",
		"c++":	"lua_atpanic(lua_State *, lua_CFunction)",
		"go":	"Atpanic"
	}, {
		"mangle":	"lua_callk",
		"c++":	"lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)",
		"go":	"Callk"
	}]
`)

	if res := strings.TrimSpace(string(data)); res != expect {
		t.Fatalf("expect %s, but got %s", expect, res)
	}
}

func TestParseHeaderFile(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		isCpp    bool
		prefixes []string
		expect   []*llcppg.SymbolInfo
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
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "(*Reader).Init__1",
					CPP:    "INIReader::INIReader(const char *, int)",
					Mangle: "_ZN9INIReaderC1EPKci",
				},
				{
					Go:     "(*Reader).Init",
					CPP:    "INIReader::INIReader(const int &)",
					Mangle: "_ZN9INIReaderC1ERKi",
				},
				{
					Go:     "(*Reader).Dispose",
					CPP:    "INIReader::~INIReader()",
					Mangle: "_ZN9INIReaderD1Ev",
				},
				{
					Go:     "(*Reader).ParseError",
					CPP:    "INIReader::ParseError()",
					Mangle: "_ZNK9INIReader10ParseErrorEv",
				},
			},
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
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "(*State).Compare",
					CPP:    "lua_compare(lua_State *, int, int, int)",
					Mangle: "lua_compare",
				},
				{
					Go:     "(*State).Rawequal",
					CPP:    "lua_rawequal(lua_State *, int, int)",
					Mangle: "lua_rawequal",
				},
				{
					Go:     "Sizecomp",
					CPP:    "lua_sizecomp(int, int, int, int)",
					Mangle: "lua_sizecomp",
				},
			},
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
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "(*Sqlite3).Errcode",
					CPP:    "sqlite3_errcode(sqlite3 *)",
					Mangle: "sqlite3_errcode",
				},
				{
					Go:     "UriParameter",
					CPP:    "sqlite3_uri_parameter(sqlite3_filename, const char *)",
					Mangle: "sqlite3_uri_parameter",
				},
			},
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
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "DerDecoding",
					CPP:    "asn1_der_decoding(asn1_node *, const void *, int, char *)",
					Mangle: "asn1_der_decoding",
				},
			},
		},

		{
			name: "InvalidReceiver typ.NamedType.String is empty",
			content: `
					RLAPI void InitWindow(int width, int height, const char *title);
					`,
			isCpp:    false,
			prefixes: []string{""},
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "InitWindow",
					CPP:    "InitWindow(int, int, const char *)",
					Mangle: "InitWindow",
				},
			},
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
			expect: []*llcppg.SymbolInfo{
				{
					Go:     "X__mpzSetUiSafe",
					CPP:    "__mpz_set_ui_safe(mpz_ptr, unsigned long)",
					Mangle: "__mpz_set_ui_safe",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			symbolMap, err := symg.ParseHeaderFile([]string{tc.content}, tc.prefixes, []string{}, nil, tc.isCpp, true)
			if err != nil {
				log.Fatal(err)
			}

			var keys []string
			for key := range symbolMap {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			result := make([]*llcppg.SymbolInfo, 0, len(keys))
			for _, key := range keys {
				info := symbolMap[key]
				result = append(result, &llcppg.SymbolInfo{
					Go:     info.GoName,
					CPP:    info.ProtoName,
					Mangle: key,
				})
			}
			if !reflect.DeepEqual(result, tc.expect) {
				t.Fatalf("expect %#v, but got %#v", tc.expect, result)
			}
		})
	}
}

func TestLdOutput(t *testing.T) {
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
	expect := []string{
		"/usr/local/lib/aarch64-linux-gnu",
		"/lib/aarch64-linux-gnu",
		"/usr/lib/aarch64-linux-gnu",
		"/usr/local/lib",
		"/lib",
		"/usr/lib",
		"/usr/aarch64-linux-gnu/lib",
	}
	if !reflect.DeepEqual(res, expect) {
		t.Fatalf("expect %v, but got %v", expect, res)
	}
}

func TestGen(t *testing.T) {
	gen := false
	testCases := []struct {
		name         string
		path         string
		dylibSymbols []string
	}{
		{
			name: "c",
			path: "./testdata/c",
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
			path: "./testdata/cpp",
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
			path: "./testdata/inireader",
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
			path: "./testdata/lua",
			dylibSymbols: []string{
				"lua_error",
				"lua_next",
				"lua_concat",
				"lua_stringtonumber",
			},
		},
		{
			name: "cjson",
			path: "./testdata/cjson",
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
			path: "./testdata/isl",
			dylibSymbols: []string{
				"isl_pw_qpolynomial_get_ctx",
			},
		},
		{
			name: "gpgerror",
			path: "./testdata/gpgerror",
			dylibSymbols: []string{
				"gpg_strsource",
				"gpg_strerror_r",
				"gpg_strerror",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projPath, err := filepath.Abs(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			cfg, err := llcppg.GetConfFromFile(filepath.Join(projPath, llcppg.LLCPPG_CFG))
			if err != nil {
				t.Fatal(err)
			}

			cfg.CFlags = "-I" + projPath
			pkgHfileInfo := config.PkgHfileInfo(cfg.Include, strings.Fields(cfg.CFlags), false)
			headerSymbolMap, err := symg.ParseHeaderFile(pkgHfileInfo.CurPkgFiles(), cfg.TrimPrefixes, strings.Fields(cfg.CFlags), cfg.SymMap, cfg.Cplusplus, false)
			if err != nil {
				t.Fatal(err)
			}

			// trim to nm symbols
			var dylibsymbs []*nm.Symbol
			for _, symb := range tc.dylibSymbols {
				dylibsymbs = append(dylibsymbs, &nm.Symbol{Name: symg.AddSymbolPrefixUnder(symb, cfg.Cplusplus)})
			}
			symbolData, err := symg.GenerateSymTable(dylibsymbs, headerSymbolMap)
			if err != nil {
				t.Fatal(err)
			}
			expectFile := filepath.Join(projPath, "expect.json")
			if gen {
				os.WriteFile(expectFile, symbolData, 0644)
			} else {
				expectData, err := os.ReadFile(expectFile)
				if err != nil {
					t.Fatal(err)
				}
				if string(symbolData) != string(expectData) {
					t.Fatalf("expect %s, but got %s", expectData, symbolData)
				}
			}
		})
	}
}
