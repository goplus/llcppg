package config_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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
			name: "Static library configuration",
			input: `{
		  "name": "mylib",
		  "cflags": "-I/opt/homebrew/include",
		  "include": ["mylib.h"],
		  "libs": "-L/opt/homebrew/lib -lmylib",
		  "libstatic": true
		}`,
			expect: llconfig.Config{
				Name:      "mylib",
				CFlags:    "-I/opt/homebrew/include",
				Include:   []string{"mylib.h"},
				Libs:      "-L/opt/homebrew/lib -lmylib",
				LibStatic: true,
			},
			mode: useFile,
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

func TestLookupSymbolOK(t *testing.T) {
	table, err := llconfig.GetSymTableFromFile(path.Join("./testdata", llconfig.LLCPPG_SYMB))
	if err != nil {
		t.Fatal(err)
	}
	entry, err := table.LookupSymbol("_ZNK9INIReader10GetBooleanERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_b")
	if err != nil {
		t.Fatal(err)
	}
	const expectCppName = "INIReader::GetBoolean(std::__1::basic_string<char, std::__1::char_traits<char>, std::__1::allocator<char>> const&, std::__1::basic_string<char, std::__1::char_traits<char>, std::__1::allocator<char>> const&, bool) const"
	const expectGoName = "(*Reader).GetBoolean"

	if entry.CPP != expectCppName {
		t.Fatalf("expect %s, got %s", expectCppName, entry.CPP)
	}
	if entry.Go != expectGoName {
		t.Fatalf("expect %s, got %s", expectGoName, entry.Go)
	}
}

func TestLookupSymbolError(t *testing.T) {
	_, err := llconfig.GetSymTableFromFile("./testdata/llcppg.symb.txt")
	if err == nil {
		t.Error("expect error")
	}
	table, err := llconfig.GetSymTableFromFile(path.Join("./testdata", llconfig.LLCPPG_SYMB))
	if err != nil {
		t.Fatal(err)
	}
	lookupSymbs := []string{
		"_ZNK9INIReader10GetBooleanERKNSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEES8_bXXX",
		"",
	}
	for _, lookupSymbol := range lookupSymbs {
		_, err := table.LookupSymbol(lookupSymbol)
		if err == nil {
			t.Error("expect error")
		}
	}
}

func TestReadPubFile(t *testing.T) {
	pub, err := llconfig.ReadPubFile(path.Join("./testdata", llconfig.LLCPPG_PUB))
	if err != nil {
		t.Fatal(err)
	}
	if len(pub) != 3 {
		t.Fatalf("expect 3 entries, got %d", len(pub))
	}
	if pub["file"] != "FILE" || pub["err"] != "Err" || pub["stdio"] != "" {
		t.Fatalf("expect file, err, stdio, got %v", pub)
	}
}

func TestReadPubFileError(t *testing.T) {
	pub, err := llconfig.ReadPubFile("./testdata/llcppg.txt")
	if !(pub != nil && len(pub) == 0 && err == nil) {
		t.Fatalf("expect empty map for llcppg.txt")
	}
	temp, err := os.CreateTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(temp.Name())
	content := `a b c`
	_, err = temp.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	_, err = llconfig.ReadPubFile(temp.Name())
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
}

func TestWritePubFile(t *testing.T) {
	pub := map[string]string{
		"file":  "FILE",
		"err":   "Err",
		"stdio": "",
	}
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	pubFile := filepath.Join(tempDir, llconfig.LLCPPG_PUB)
	err = llconfig.WritePubFile(pubFile, pub)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(pubFile)
	if err != nil {
		t.Fatal(err)
	}
	expect :=
		`err Err
file FILE
stdio`
	if string(content) != expect {
		t.Fatalf("expect %s, got %s", expect, string(content))
	}

	notExistFilePath := filepath.Join(tempDir, "not_exit_dir", "not_exist_file.pub")
	err = llconfig.WritePubFile(notExistFilePath, pub)
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expect os.IsNotExist error, got %v", err)
	}

	notExistFile := filepath.Join(tempDir, "not_exist_file.pub")
	err = llconfig.WritePubFile(notExistFile, make(map[string]string, 0))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	_, err = os.Stat(notExistFile)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("expect file %s, got error %v", notExistFile, err)
	}
}
