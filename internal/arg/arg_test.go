package arg_test

import (
	"fmt"
	"reflect"
	"testing"

	llcppg "github.com/goplus/llcppg/config"
	"github.com/goplus/llcppg/internal/arg"
)

func TestParseArgs(t *testing.T) {
	fmt.Println("=== Test ParseArgs ===")

	swflags := map[string]bool{
		"--extract": true,
	}

	testCases := []struct {
		name   string
		input  []string
		expect struct {
			args   *arg.Args
			remain []string
		}
	}{
		{
			name:  "Basic flags",
			input: []string{"-h", "-v", "-"},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     true,
					Verbose:  true,
					UseStdin: true,
					CfgFile:  llcppg.LLCPPG_CFG,
				},
				remain: []string{},
			},
		},
		{
			name:  "Config file",
			input: []string{"lua.llcppg.cfg"},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     false,
					Verbose:  false,
					UseStdin: false,
					CfgFile:  "lua.llcppg.cfg",
				},
				remain: []string{},
			},
		},
		{
			name:  "Extract with multiple args",
			input: []string{"--extract", "file1.h", "file2.h", "-v"},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     false,
					Verbose:  true,
					UseStdin: false,
					CfgFile:  llcppg.LLCPPG_CFG,
				},
				remain: []string{"--extract", "file1.h", "file2.h"},
			},
		},
		{
			name:  "Non-skippable flags",
			input: []string{"--extract", "file1.h", "file2.h", "-out=true", "-cpp=true", "-v"},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     false,
					Verbose:  true,
					UseStdin: false,
					CfgFile:  llcppg.LLCPPG_CFG,
				},
				remain: []string{"--extract", "file1.h", "file2.h", "-out=true", "-cpp=true"},
			},
		},
		{
			name:  "Mixed flags",
			input: []string{"-v", "--extract", "file.h", "-out=true", "config.json"},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     false,
					Verbose:  true,
					UseStdin: false,
					CfgFile:  "config.json",
				},
				remain: []string{"--extract", "file.h", "-out=true"},
			},
		},
		{
			name:  "Empty input",
			input: []string{},
			expect: struct {
				args   *arg.Args
				remain []string
			}{
				args: &arg.Args{
					Help:     false,
					Verbose:  false,
					UseStdin: false,
					CfgFile:  llcppg.LLCPPG_CFG,
				},
				remain: []string{},
			},
		},
	}

	for _, tc := range testCases {
		result, filteredArgs := arg.ParseArgs(tc.input, llcppg.LLCPPG_CFG, swflags)
		if !reflect.DeepEqual(result, tc.expect.args) {
			t.Fatalf("Test case %s failed: expected %#v, got %#v", tc.name, tc.expect.args, result)
		}
		if !reflect.DeepEqual(filteredArgs, tc.expect.remain) {
			t.Fatalf("Test case %s failed: expected %#v, got %#v", tc.name, tc.expect.remain, filteredArgs)
		}
	}
}
