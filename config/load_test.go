package config_test

import (
	"reflect"
	"strings"
	"testing"

	llconfig "github.com/goplus/llcppg/config"
)

func TestGetConfByByte(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expect    llconfig.Config
		expectErr bool
	}{
		{
			name: "SQLite configuration",
			input: `{
  "name": "sqlite",
  "cflags": "-I/opt/homebrew/opt/sqlite/include",
  "include": ["sqlite3.h"],
  "libs": "-L/opt/homebrew/opt/sqlite/lib -lsqlite3",
  "trimPrefixes": ["sqlite3_"],
  "cplusplus": false,
  "symMap": {
    "sqlite3_finalize":".Close"
  }
}`,
			expect: llconfig.Config{
				Name:         "sqlite",
				CFlags:       "-I/opt/homebrew/opt/sqlite/include",
				Include:      []string{"sqlite3.h"},
				Libs:         "-L/opt/homebrew/opt/sqlite/lib -lsqlite3",
				TrimPrefixes: []string{"sqlite3_"},
				Cplusplus:    false,
				SymMap: map[string]string{
					"sqlite3_finalize": ".Close",
				},
			},
		},

		{
			name: "Lua configuration",
			input: `{
		  "name": "lua",
		  "cflags": "-I/opt/homebrew/include/lua",
		  "include": ["lua.h"],
		  "libs": "-L/opt/homebrew/lib -llua -lm",
		  "trimPrefixes": ["lua_", "lua_"],
		  "cplusplus": false
		}`,
			expect: llconfig.Config{
				Name:         "lua",
				CFlags:       "-I/opt/homebrew/include/lua",
				Include:      []string{"lua.h"},
				Libs:         "-L/opt/homebrew/lib -llua -lm",
				TrimPrefixes: []string{"lua_", "lua_"},
			},
		},
		{
			name:      "Invalid JSON",
			input:     `{invalid json}`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := llconfig.ConfigFromReader(strings.NewReader(tc.input))
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for test case %s, but got nil", tc.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error for test case %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(result, tc.expect) {
				t.Fatalf("expected %#v, but got %#v", tc.expect, result)
			}
		})
	}
}
