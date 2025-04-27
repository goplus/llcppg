package main

import (
	"fmt"
	"sort"

	"github.com/goplus/llcppg/_xtool/llcppsymg/names"
	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
)

func main() {
	TestNewSymbolProcessor()
	TestGenMethodName()
	TestAddSuffix()
	TestParseHeaderFile()
}

func TestNewSymbolProcessor() {
	fmt.Println("=== Test NewSymbolProcessor ===")
	process := parse.NewSymbolProcessor([]string{}, []string{"lua_", "luaL_"}, nil)
	fmt.Printf("Before: No prefixes After: Prefixes: %v\n", process.Prefixes)
	fmt.Println()
}

func TestGenMethodName() {
	fmt.Println("=== Test GenMethodName ===")
	process := &parse.SymbolProcessor{}

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
	process := parse.NewSymbolProcessor([]string{}, []string{"INI"}, nil)
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

		symbolMap, err := parse.ParseHeaderFile([]string{tc.content}, tc.prefixes, []string{}, nil, tc.isCpp, true)

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
