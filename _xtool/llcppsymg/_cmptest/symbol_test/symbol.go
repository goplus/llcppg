package main

import (
	"fmt"

	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/symbol"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/xtool/nm"
)

func main() {
	TestGetCommonSymbols()
	TestGenSymbolTableData()
}

func TestGetCommonSymbols() {
	fmt.Println("=== Test GetCommonSymbols ===")
	testCases := []struct {
		name          string
		dylibSymbols  []*nm.Symbol
		headerSymbols map[string]*parse.SymbolInfo
	}{
		{
			name: "Lua symbols",
			dylibSymbols: []*nm.Symbol{
				{Name: symbol.AddSymbolPrefixUnder("lua_absindex", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_arith", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_atpanic", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_callk", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_lib_nonexistent", false)},
			},
			headerSymbols: map[string]*parse.SymbolInfo{
				"lua_absindex":           {ProtoName: "lua_absindex(lua_State *, int)", GoName: "Absindex"},
				"lua_arith":              {ProtoName: "lua_arith(lua_State *, int)", GoName: "Arith"},
				"lua_atpanic":            {ProtoName: "lua_atpanic(lua_State *, lua_CFunction)", GoName: "Atpanic"},
				"lua_callk":              {ProtoName: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", GoName: "Callk"},
				"lua_header_nonexistent": {ProtoName: "lua_header_nonexistent()", GoName: "HeaderNonexistent"},
			},
		},
		{
			name: "INIReader and Std library symbols",
			dylibSymbols: []*nm.Symbol{
				{Name: symbol.AddSymbolPrefixUnder("ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZNK9INIReader10ParseErrorEv", true)},
			},
			headerSymbols: map[string]*parse.SymbolInfo{
				"_ZNK9INIReader12GetInteger64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_x":  {GoName: "(*Reader).GetInteger64", ProtoName: "INIReader::GetInteger64(const std::string &, const std::string &, int64_t)"},
				"_ZNK9INIReader13GetUnsigned64ERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_y": {GoName: "(*Reader).GetUnsigned64", ProtoName: "INIReader::GetUnsigned64(const std::string &, const std::string &, uint64_t)"},
				"_ZNK9INIReader7GetRealERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_d":        {GoName: "(*Reader).GetReal", ProtoName: "INIReader::GetReal(const std::string &, const std::string &, double)"},
				"_ZNK9INIReader10ParseErrorEv": {GoName: "(*Reader).ParseError", ProtoName: "INIReader::ParseError()"},
				"_ZNK9INIReader10GetBooleanERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_b": {GoName: "(*Reader).GetBoolean", ProtoName: "INIReader::GetBoolean(const std::string &, const std::string &, bool)"},
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\nTest Case: %s\n", tc.name)
		commonSymbols := symbol.GetCommonSymbols(tc.dylibSymbols, tc.headerSymbols)
		fmt.Printf("Common Symbols (%d):\n", len(commonSymbols))
		for _, sym := range commonSymbols {
			fmt.Printf("Mangle: %s, CPP: %s, Go: %s\n", sym.Mangle, sym.CPP, sym.Go)
		}
	}
	fmt.Println()
}

func TestGenSymbolTableData() {
	fmt.Println("=== Test GenSymbolTableData ===")

	commonSymbols := []*llcppg.SymbolInfo{
		{Mangle: "lua_absindex", CPP: "lua_absindex(lua_State *, int)", Go: "Absindex"},
		{Mangle: "lua_arith", CPP: "lua_arith(lua_State *, int)", Go: "Arith"},
		{Mangle: "lua_atpanic", CPP: "lua_atpanic(lua_State *, lua_CFunction)", Go: "Atpanic"},
		{Mangle: "lua_callk", CPP: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", Go: "Callk"},
	}

	// existingSymbols := map[string]llcppg.SymbolInfo{
	// 	"lua_absindex": {Mangle: "lua_absindex", CPP: "lua_absindex(lua_State *, int)", Go: "Absindex"},
	// 	"lua_arith":    {Mangle: "lua_arith", CPP: "lua_arith(lua_State *, int)", Go: "Arith"},
	// 	"lua_callk":    {Mangle: "lua_callk", CPP: "lua_callk(lua_State *, int, int, lua_KContext, lua_KFunction)", Go: "ModifiedCallk"},
	// }

	data, err := symbol.GenSymbolTableData(commonSymbols)
	if err != nil {
		fmt.Printf("Error generating symbol table data: %v\n", err)
		return
	}
	fmt.Println(string(data))
	fmt.Println()
}
