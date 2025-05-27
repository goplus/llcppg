package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	llconfig "github.com/goplus/llcppg/config"
)

type testMode int

const (
	useStdin testMode = 1 << iota
	useFile
)

func TestGetConfByByte(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		mode      testMode
		expect    llconfig.Config
		expectErr bool
	}{
		{
			name: "SQLite configuration(File)",
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
			mode: useFile,
		},

		{
			name: "SQLite configuration(Stdin)",
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
			mode: useStdin,
		},

		{
			name: "Lua configuration(file)",
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
			mode: useFile,
		},

		{
			name: "Lua configuration(stdin)",
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
			mode: useStdin,
		},

		{
			name:      "Invalid JSON",
			input:     `{invalid json}`,
			expectErr: true,
			mode:      useStdin,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result llconfig.Config

			file, err := os.CreateTemp("", "configtest")
			if err != nil {
				t.Fatal(err)
				return
			}
			defer os.Remove(file.Name())

			_, err = file.Write([]byte(tc.input))
			if err != nil {
				t.Fatal(err)
				return
			}
			err = fmt.Errorf("config: no mode is specified")

			if tc.mode&useStdin != 0 {
				var fileR *os.File
				stdin := os.Stdin
				defer func() { os.Stdin = stdin }()

				//	Seek isn't implemented in llgo, so we have to open the file again to reset the cursor
				fileR, err = os.Open(file.Name())
				if err != nil {
					t.Fatal(err)
					return
				}
				defer fileR.Close()

				// swap for testing
				os.Stdin = fileR

				result, err = llconfig.GetConfFromStdin()
			}

			if tc.mode&useFile != 0 {
				result, err = llconfig.GetConfFromFile(file.Name())
			}

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
