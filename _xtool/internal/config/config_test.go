package config_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/config"
	llconfig "github.com/goplus/llcppg/config"
)

func TestGetConfByByte(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expect    *llconfig.Config
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
			expect: &llconfig.Config{
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
			expect: &llconfig.Config{
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
			result, err := config.GetConfByByte([]byte(tc.input))
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for test case %s, but got nil", tc.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error for test case %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(result.Config, tc.expect) {
				t.Fatalf("expected %#v, but got %#v", tc.expect, result.Config)
			}
		})
	}
}

func TestPkgHfileInfo(t *testing.T) {
	cases := []struct {
		conf *llconfig.Config
		want *config.PkgHfilesInfo
	}{
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
			},
			want: &config.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{"testdata/hfile/tempimpl.h"},
			},
		},
		{
			conf: &llconfig.Config{
				CFlags:  "-I./testdata/hfile -I ./testdata/thirdhfile",
				Include: []string{"temp1.h", "temp2.h"},
				Mix:     true,
			},
			want: &config.PkgHfilesInfo{
				Inters: []string{"testdata/hfile/temp1.h", "testdata/hfile/temp2.h"},
				Impls:  []string{},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			info := config.PkgHfileInfo(tc.conf, []string{})
			if !reflect.DeepEqual(info.Inters, tc.want.Inters) {
				t.Fatalf("inter expected %v, but got %v", tc.want.Inters, info.Inters)
			}
			if !reflect.DeepEqual(info.Impls, tc.want.Impls) {
				t.Fatalf("impl expected %v, but got %v", tc.want.Impls, info.Impls)
			}

			thirdhfile, err := filepath.Abs("./testdata/thirdhfile/third.h")
			if err != nil {
				t.Fatalf("failed to get abs path: %w", err)
			}
			tfileFound := false
			stdioFound := false
			for _, tfile := range info.Thirds {
				absTfile, err := filepath.Abs(tfile)
				if err != nil {
					t.Fatalf("failed to get abs path: %w", err)
				}
				if absTfile == thirdhfile {
					tfileFound = true
				}
				if strings.HasSuffix(absTfile, "stdio.h") {
					stdioFound = true
				}
			}
			if !tfileFound || !stdioFound {
				t.Fatalf("third hfile or std hfile not found")
			}
		})
	}
}
