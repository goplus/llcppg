package symg_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/clangtool"
	"github.com/goplus/llcppg/_xtool/internal/header"
	"github.com/goplus/llcppg/_xtool/internal/symbol"
	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/symg"
	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/name"
	"github.com/goplus/llgo/xtool/nm"
)

func TestAddSuffix(t *testing.T) {
	prefix := []string{"INI"}
	process := symg.NewSymbolProcessor([]string{}, prefix, nil)
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
			goName := name.GoName(tc.method, prefix, true)
			className := name.GoName("INIReader", prefix, true)
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
				{Name: addSymbolPrefixUnder("lua_absindex", false)},
				{Name: addSymbolPrefixUnder("lua_arith", false)},
				{Name: addSymbolPrefixUnder("lua_atpanic", false)},
				{Name: addSymbolPrefixUnder("lua_callk", false)},
				{Name: addSymbolPrefixUnder("lua_lib_nonexistent", false)},
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
				{Name: addSymbolPrefixUnder("ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x", true)},
				{Name: addSymbolPrefixUnder("ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d", true)},
				{Name: addSymbolPrefixUnder("ZNK9INIReader10ParseErrorEv", true)},
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

	data, err := json.MarshalIndent(commonSymbols, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expect := strings.TrimSpace(`
[
  {
    "mangle": "lua_absindex",
    "c++": "lua_absindex(lua_State *, int)",
    "go": "Absindex"
  },
  {
    "mangle": "lua_arith",
    "c++": "lua_arith(lua_State *, int)",
    "go": "Arith"
  },
  {
    "mangle": "lua_atpanic",
    "c++": "lua_atpanic(lua_State *, lua_CFunction)",
    "go": "Atpanic"
  },
  {
    "mangle": "lua_callk",
    "c++": "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)",
    "go": "Callk"
  }
]
`)

	if res := strings.TrimSpace(string(data)); res != expect {
		t.Fatalf("expect \n%s, but got \n%s", expect, res)
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
			f, err := os.CreateTemp("", "temp*.h")
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.Write([]byte(tc.content))
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())
			symbolMap, err := symg.ParseHeaderFile(f.Name(), []string{f.Name()}, tc.prefixes, []string{}, nil, tc.isCpp)
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

			tempFile, err := os.CreateTemp("", "combine*.h")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFile.Name())
			clangtool.ComposeIncludes(cfg.Include, tempFile.Name())

			pkgHfileInfo := header.PkgHfileInfo(&header.Config{
				Includes: cfg.Include,
				Args:     strings.Fields(cfg.CFlags),
				Mix:      false,
			})
			headerSymbolMap, err := symg.ParseHeaderFile(tempFile.Name(), pkgHfileInfo.CurPkgFiles(), cfg.TrimPrefixes, strings.Fields(cfg.CFlags), cfg.SymMap, cfg.Cplusplus)
			if err != nil {
				t.Fatal(err)
			}

			// trim to nm symbols
			var dylibsymbs []*nm.Symbol
			for _, symb := range tc.dylibSymbols {
				dylibsymbs = append(dylibsymbs, &nm.Symbol{Name: addSymbolPrefixUnder(symb, cfg.Cplusplus)})
			}
			symbols := symg.GetCommonSymbols(dylibsymbs, headerSymbolMap)
			if err != nil {
				t.Fatal(err)
			}
			symbolData, err := json.MarshalIndent(symbols, "", "  ")
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

// For mutiple os test,the nm output's symbol name is different.
func addSymbolPrefixUnder(name string, isCpp bool) string {
	prefix := ""
	if runtime.GOOS == "darwin" {
		prefix = prefix + "_"
	}
	if isCpp {
		prefix = prefix + "_"
	}
	return prefix + name
}

func TestFetchSymbols(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_fetch_symbols_*")
	if err != nil {
		t.Fatal(err)
	}
	// todo(zzy): remove this after test,need llgo support
	// defer os.RemoveAll(tempDir)

	cSource := `
#include <stdio.h>

void test_function_1(void) {
    printf("test function 1\n");
}

int test_function_2(int x) {
    return x * 2;
}

const char* test_function_3(void) {
    return "hello world";
}
`

	cSourcePath := filepath.Join(tempDir, "test.c")
	err = os.WriteFile(cSourcePath, []byte(cSource), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(cSourcePath)

	var dylibPath string
	var compileCmd []string
	if runtime.GOOS == "darwin" {
		dylibPath = filepath.Join(tempDir, "libtest.dylib")
		compileCmd = []string{"clang", "-shared", "-fPIC", "-o", dylibPath, cSourcePath}
	} else if runtime.GOOS == "linux" {
		dylibPath = filepath.Join(tempDir, "libtest.so")
		compileCmd = []string{"gcc", "-shared", "-fPIC", "-o", dylibPath, cSourcePath}
	} else {
		t.Skip("Unsupported platform for this test")
	}

	cmd := exec.Command(compileCmd[0], compileCmd[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile test library: %v\nOutput: %s", err, output)
	}

	if _, err := os.Stat(dylibPath); os.IsNotExist(err) {
		t.Fatal("Dynamic library was not created")
	}

	libDir := tempDir
	libFlag := fmt.Sprintf("-L%s -ltest", libDir)

	symbols, err := symg.FetchSymbols(libFlag, symbol.ModeDynamic)
	if err != nil {
		t.Fatalf("FetchSymbols failed: %v", err)
	}

	if len(symbols) == 0 {
		t.Fatal("No symbols found")
	}

	expectedSymbols := []string{"test_function_1", "test_function_2", "test_function_3"}
	foundSymbols := make(map[string]bool)

	for _, sym := range symbols {
		// On Darwin, symbols have '_' prefix, so trim it
		symName := sym.Name
		if runtime.GOOS == "darwin" {
			symName = strings.TrimPrefix(symName, "_")
		}
		foundSymbols[symName] = true
	}

	for _, expected := range expectedSymbols {
		if !foundSymbols[expected] {
			t.Errorf("Expected symbol %s not found in library symbols", expected)
		}
	}

	t.Logf("Successfully found %d symbols including expected test functions", len(symbols))
}
