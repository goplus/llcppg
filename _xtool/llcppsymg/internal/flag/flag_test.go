package flag_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/goplus/llcppg/_xtool/llcppsymg/internal/flag"
)

func TestParseLibs(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect *flag.Libs
	}{
		{
			name:  "Lua library",
			input: "-L/opt/homebrew/lib -llua -lm",
			expect: &flag.Libs{
				Paths: []string{"/opt/homebrew/lib"},
				Names: []string{"lua", "m"},
			},
		},
		{
			name:  "SQLite library",
			input: "-L/opt/homebrew/opt/sqlite/lib -lsqlite3",
			expect: &flag.Libs{
				Paths: []string{"/opt/homebrew/opt/sqlite/lib"},
				Names: []string{"sqlite3"},
			},
		},
		{
			name:  "INIReader library",
			input: "-L/opt/homebrew/Cellar/inih/58/lib -lINIReader",
			expect: &flag.Libs{
				Paths: []string{"/opt/homebrew/Cellar/inih/58/lib"},
				Names: []string{"INIReader"},
			},
		},
		{
			name:  "Multiple library paths",
			input: "-L/opt/homebrew/lib -L/usr/lib -llua",
			expect: &flag.Libs{
				Paths: []string{"/opt/homebrew/lib", "/usr/lib"},
				Names: []string{"lua"},
			},
		},
		{
			name:  "No valid library",
			input: "-L/opt/homebrew/lib",
			expect: &flag.Libs{
				Paths: []string{"/opt/homebrew/lib"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := flag.ParseLibs(tc.input)
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
		conf         *flag.Libs
		defaultPaths []string
		expectErr    bool
		want         []string
		wantNotFound []string
	}{
		{
			name: "existing dylib",
			conf: &flag.Libs{
				Names: []string{"symb1"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{},
			want:         []string{dylib1},
			expectErr:    false,
		},
		{
			name: "existing dylibs",
			conf: &flag.Libs{
				Names: []string{"symb1", "symb2"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{},
			want:         []string{dylib1, dylib2},
			expectErr:    false,
		},
		{
			name: "existint default paths",
			conf: &flag.Libs{
				Names: []string{"symb1", "symb3"},
				Paths: []string{tempDir},
			},
			defaultPaths: []string{tempDefaultPath},
			want:         []string{dylib1, defaultDylib3},
			expectErr:    false,
		},
		{
			name: "existint default paths & not found",
			conf: &flag.Libs{
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
			conf: &flag.Libs{
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
			paths, notFounds, err := tc.conf.GenDylibPaths(tc.defaultPaths)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %w", err)
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

func TestParseCFlags(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "Single include path",
			input:  "-I/usr/include",
			expect: []string{"/usr/include"},
		},
		{
			name:   "Multiple include paths",
			input:  "-I/usr/include -I/opt/homebrew/include",
			expect: []string{"/usr/include", "/opt/homebrew/include"},
		},
		{
			name:   "Include paths mixed with other flags",
			input:  "-I/usr/include -DDEBUG -I/opt/local/include -Wall",
			expect: []string{"/usr/include", "/opt/local/include"},
		},
		{
			name:  "Empty input",
			input: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := flag.ParseCFlags(tc.input)
			if !reflect.DeepEqual(conf.Paths, tc.expect) {
				t.Fatalf("expected paths %v, got %v", tc.expect, conf.Paths)
			}
		})
	}
}

func TestGenHeaderFilePath(t *testing.T) {
	tempDir := os.TempDir()
	temDir2 := filepath.Join(tempDir, "include")
	tempFile1 := filepath.Join(tempDir, "test1.h")
	tempFile2 := filepath.Join(tempDir, "test2.h")
	tempFile3 := filepath.Join(temDir2, "test3.h")
	os.MkdirAll(temDir2, 0755)
	os.Create(tempFile1)
	os.Create(tempFile2)
	os.Create(tempFile3)
	defer os.Remove(tempFile1)
	defer os.Remove(tempFile2)
	defer os.Remove(tempFile3)
	defer os.Remove(temDir2)

	testCases := []struct {
		name      string
		cflags    string
		files     []string
		notFounds []string
		expect    []string
		expectErr bool
	}{
		{
			name:   "Valid files",
			cflags: "-I" + tempDir,
			files:  []string{"test1.h", "test2.h"},
			expect: []string{"test1.h", "test2.h"},
		},
		{
			name:      "Mixed existing and non-existing files",
			cflags:    "-I" + tempDir,
			files:     []string{"test1.h", "nonexistent.h"},
			notFounds: []string{"nonexistent.h"},
			expect:    []string{"test1.h"},
		},
		{
			name:   "Multiple include paths",
			cflags: "-I" + tempDir + " -I" + temDir2,
			files:  []string{"test1.h", "test2.h", "test3.h"},
			expect: []string{"test1.h", "test2.h", "test3.h"},
		},
		{
			name:      "No existing files",
			cflags:    "-I" + tempDir,
			files:     []string{"nonexistent1.h", "nonexistent2.h"},
			notFounds: []string{"nonexistent1.h", "nonexistent2.h"},
			expectErr: true,
		},
		{
			name:      "Empty file list",
			cflags:    "-I/usr/include",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cflag := flag.ParseCFlags(tc.cflags)
			result, notFounds, err := cflag.GenHeaderFilePaths(tc.files, []string{})
			if !reflect.DeepEqual(notFounds, tc.notFounds) {
				t.Fatalf("expected notFounds %v, got %v", tc.notFounds, notFounds)
			}
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %w", err)
			}

			relativeResult := make([]string, len(result))
			for i, path := range result {
				relativeResult[i] = filepath.Base(path)
			}
			if !reflect.DeepEqual(relativeResult, tc.expect) {
				t.Fatalf("expected files %v, got %v", tc.expect, relativeResult)
			}
		})
	}
}
