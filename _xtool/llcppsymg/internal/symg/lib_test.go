package symg_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/symbol"
	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/symg"
)

func TestParseLibs(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect *symg.Libs
	}{
		{
			name:  "Lua library",
			input: "-L/opt/homebrew/lib -llua -lm",
			expect: &symg.Libs{
				Paths: []string{"/opt/homebrew/lib"},
				Names: []string{"lua", "m"},
			},
		},
		{
			name:  "SQLite library",
			input: "-L/opt/homebrew/opt/sqlite/lib -lsqlite3",
			expect: &symg.Libs{
				Paths: []string{"/opt/homebrew/opt/sqlite/lib"},
				Names: []string{"sqlite3"},
			},
		},
		{
			name:  "INIReader library",
			input: "-L/opt/homebrew/Cellar/inih/58/lib -lINIReader",
			expect: &symg.Libs{
				Paths: []string{"/opt/homebrew/Cellar/inih/58/lib"},
				Names: []string{"INIReader"},
			},
		},
		{
			name:  "Multiple library paths",
			input: "-L/opt/homebrew/lib -L/usr/lib -llua",
			expect: &symg.Libs{
				Paths: []string{"/opt/homebrew/lib", "/usr/lib"},
				Names: []string{"lua"},
			},
		},
		{
			name:  "No valid library",
			input: "-L/opt/homebrew/lib",
			expect: &symg.Libs{
				Paths: []string{"/opt/homebrew/lib"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := symg.ParseLibs(tc.input)
			if !reflect.DeepEqual(conf, tc.expect) {
				t.Errorf("expected %#v, got %#v", tc.expect, conf)
			}
		})
	}
}

func TestGenDylibPaths(t *testing.T) {
	tempDir := os.TempDir()
	tempDefaultPath := filepath.Join(tempDir, "symblib")
	affix := ".dylib"
	if runtime.GOOS == "linux" {
		affix = ".so"
	}
	err := os.MkdirAll(tempDefaultPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create temp default path: %v\n", err)
		return
	}

	dylib1 := filepath.Join(tempDir, "libsymb1"+affix)
	dylib2 := filepath.Join(tempDir, "libsymb2"+affix)
	defaultDylib3 := filepath.Join(tempDefaultPath, "libsymb3"+affix)

	os.Create(dylib1)
	os.Create(dylib2)
	os.Create(defaultDylib3)
	defer os.Remove(dylib1)
	defer os.Remove(dylib2)
	defer os.Remove(defaultDylib3)
	defer os.Remove(tempDefaultPath)

	testCase := []struct {
		name         string
		conf         *symg.Libs
		defaultPaths []string
		expectErr    bool
		want         []string
		wantNotFound []string
	}{
		{
			name: "existing dylib",
			conf: &symg.Libs{
				Names: []string{"symb1"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{},
			want:         []string{dylib1},
			expectErr:    false,
		},
		{
			name: "existing dylibs",
			conf: &symg.Libs{
				Names: []string{"symb1", "symb2"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{},
			want:         []string{dylib1, dylib2},
			expectErr:    false,
		},
		{
			name: "existint default paths",
			conf: &symg.Libs{
				Names: []string{"symb1", "symb3"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{tempDefaultPath},
			want:         []string{dylib1, defaultDylib3},
			expectErr:    false,
		},
		{
			name: "existint default paths & not found",
			conf: &symg.Libs{
				Names: []string{"symb1", "symb3", "math"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{tempDefaultPath},
			want:         []string{dylib1, defaultDylib3},
			wantNotFound: []string{"math"},
			expectErr:    false,
		},
		{
			name: "no existing dylib",
			conf: &symg.Libs{
				Names: []string{"notexist"},
				Paths: []string{tempDir},
			},
			want:         []string{},
			wantNotFound: []string{"notexist"},
			expectErr:    true,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			paths, notFounds, err := tc.conf.Files(tc.defaultPaths, symbol.ModeDynamic)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(notFounds, tc.wantNotFound) {
				t.Fatalf("expected notFounds %v, got %v", tc.wantNotFound, notFounds)
			}

			for _, wantPath := range tc.want {
				found := false
				for _, path := range paths {
					if path == wantPath {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected path %s, but not found", wantPath)
				}
			}
		})

	}
}
