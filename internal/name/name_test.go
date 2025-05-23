package name_test

import (
	"testing"

	"github.com/goplus/llcppg/internal/name"
)

func TestToGoName(t *testing.T) {
	testCases := []struct {
		prefixes []string
		input    string
		expect   string
	}{
		{[]string{"lua_", "luaL_"}, "lua_closethread", "Closethread"},
		{[]string{"lua_", "luaL_"}, "luaL_checknumber", "Checknumber"},
		{[]string{"sqlite3_", "sqlite3_"}, "sqlite3_close_v2", "CloseV2"},
		{[]string{"sqlite3_", "sqlite3_"}, "sqlite3_callback", "Callback"},
		{[]string{"INI"}, "GetReal", "GetReal"},
		{[]string{"INI"}, "GetBoolean", "GetBoolean"},
		{[]string{"INI"}, "INIReader", "Reader"},
	}

	for _, tc := range testCases {
		result := name.GoName(tc.input, tc.prefixes, true)
		if result != tc.expect {
			t.Fatalf("TestToGoName failed, input: %s, expected: %s, got: %s", tc.input, tc.expect, result)
		}
	}
}

func TestPubName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"sqlite_file", "SqliteFile"},
		{"_gmp_err", "X_gmpErr"},
		{"123illegal", "X123illegal"},
		{"alreadyCamel", "AlreadyCamel"},
		{"_ab__cd", "X_abCd"},
		{"_ab_cd", "X_abCd"},
		{"_ab___cd", "X_abCd"},
		{"ab_cd", "AbCd"},
		{"ab__cd", "AbCd"},
		{"ab_cd_", "AbCd_"},
		{"ab__cd_", "AbCd_"},
		{"ab__cd__", "AbCd__"},
		{"_", "X_"},
		{"__", "X__"},
		{"___", "X___"},
	}

	for _, tc := range testCases {
		result := name.PubName(tc.input)
		if result != tc.expected {
			t.Fatalf("TestPubName failed, input: %s, expected: %s, got: %s", tc.input, tc.expected, result)
		}
	}
}

func TestExportName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"sqlite_file", "Sqlite_file"},
		{"_sqlite_file", "X_sqlite_file"},
		{"123illegal", "X123illegal"},
		{"CODE_MASK", "CODE_MASK"},
		{"_CODE_MASK", "X_CODE_MASK"},
		{"_ab__cd", "X_ab__cd"},
		{"_ab_cd", "X_ab_cd"},
		{"_ab___cd", "X_ab___cd"},
		{"ab_cd", "Ab_cd"},
		{"ab__cd", "Ab__cd"},
		{"ab_cd_", "Ab_cd_"},
		{"ab__cd_", "Ab__cd_"},
		{"_", "X_"},
		{"__", "X__"},
		{"___", "X___"},
	}
	for _, tc := range testCases {
		result := name.ExportName(tc.input)
		if result != tc.expected {
			t.Fatalf("TestExportName failed, input: %s, expected: %s, got: %s", tc.input, tc.expected, result)
		}
	}
}

func TestHeaderFileToGo(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"/path/to/foo.h", "foo.go"},
		{"/path/to/_intptr.h", "X_intptr.go"},
		{"header.h", "header.go"},
		{"_impl.h", "X_impl.go"},
	}

	for _, tc := range testCases {
		result := name.HeaderFileToGo(tc.input)
		if result != tc.expected {
			t.Fatalf("TestHeaderFileToGo failed, input: %s, expected: %s, got: %s", tc.input, tc.expected, result)
		}
	}
}
