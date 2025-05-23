package clangtool_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/clangtool"
)

func TestComposeIncludes(t *testing.T) {
	fmt.Println("=== Test ComposeIncludes ===")
	testCases := []struct {
		name   string
		files  []string
		expect string
	}{
		{
			name:  "One file",
			files: []string{"file1.h"},
			expect: `#include <file1.h>
`,
		},
		{
			name:  "Two files",
			files: []string{"file1.h", "file2.h"},
			expect: `#include <file1.h>
#include <file2.h>
`,
		},
		{
			name:   "Empty files",
			files:  []string{},
			expect: "",
		},
	}
	for _, tc := range testCases {
		outfile, err := os.CreateTemp("", "compose_*.h")
		if err != nil {
			t.Fatal(err)
		}

		err = clangtool.ComposeIncludes(tc.files, outfile.Name())
		if err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(outfile.Name())
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != tc.expect {
			t.Fatalf("expect %s, but got %s", tc.expect, string(content))
		}
		outfile.Close()
		os.Remove(outfile.Name())
	}
}

func TestSysRoot(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name: "macos-sysroot",
			input: `Homebrew clang version 19.1.7
Target: arm64-apple-darwin23.6.0
Thread model: posix
InstalledDir: /opt/homebrew/Cellar/llvm@19/19.1.7/bin
Configuration file: /opt/homebrew/etc/clang/arm64-apple-darwin23.cfg
System configuration file directory: /opt/homebrew/etc/clang
 (in-process)
 "/opt/homebrew/Cellar/llvm@19/19.1.7/bin/clang-19" -cc1 -triple arm64-apple-macosx14.0.0 -Wundef-prefix=TARGET_OS_ -Werror=undef-prefix -Wdeprecated-objc-isa-usage -Werror=deprecated-objc-isa-usage -E -disable-free -clear-ast-before-backend -disable-llvm-verifier -discard-value-names -main-file-name null -mrelocation-model pic -pic-level 2 -mframe-pointer=non-leaf -ffp-contract=on -fno-rounding-math -funwind-tables=1 -target-sdk-version=14.4 -fcompatibility-qualified-id-block-type-checking -fvisibility-inlines-hidden-static-local-var -fbuiltin-headers-in-system-modules -fdefine-target-os-macros -target-cpu apple-m1 -target-feature +zcm -target-feature +zcz -target-feature +v8.4a -target-feature +aes -target-feature +altnzcv -target-feature +ccdp -target-feature +complxnum -target-feature +crc -target-feature +dotprod -target-feature +fp-armv8 -target-feature +fp16fml -target-feature +fptoint -target-feature +fullfp16 -target-feature +jsconv -target-feature +lse -target-feature +neon -target-feature +pauth -target-feature +perfmon -target-feature +predres -target-feature +ras -target-feature +rcpc -target-feature +rdm -target-feature +sb -target-feature +sha2 -target-feature +sha3 -target-feature +specrestrict -target-feature +ssbs -target-abi darwinpcs -debugger-tuning=lldb  -target-linker-version 1115.7.3 -v  -resource-dir /opt/homebrew/Cellar/llvm@19/19.1.7/lib/clang/19 -isysroot /Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk -internal-isystem /Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/local/include -internal-isystem /opt/homebrew/Cellar/llvm@19/19.1.7/lib/clang/19/include -internal-externc-isystem /Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include -ferror-limit 19 -stack-protector 1 -fblocks -fencode-extended-block-signature -fregister-global-dtors-with-atexit -fgnuc-version=4.2.1 -fskip-odr-check-in-gmf -fmax-type-align=16 -fcolor-diagnostics -D__GCC_HAVE_DWARF2_CFI_ASM=1 -o - -x c /dev/null
clang -cc1 version 19.1.7 based upon LLVM 19.1.7 default target arm64-apple-darwin23.6.0
ignoring nonexistent directory "/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/local/include"
ignoring nonexistent directory "/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/Library/Frameworks"
#include "..." search starts here:
#include <...> search starts here:
 /opt/homebrew/Cellar/llvm@19/19.1.7/lib/clang/19/include
 /Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include
 /Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/System/Library/Frameworks (framework directory)
End of search list.
# 1 "/dev/null"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 455 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "/dev/null" 2
`,
			expect: []string{
				"-resource-dir/opt/homebrew/Cellar/llvm@19/19.1.7/lib/clang/19",
				"-isysroot/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk",
				"-internal-isystem/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/local/include",
				"-internal-isystem/opt/homebrew/Cellar/llvm@19/19.1.7/lib/clang/19/include",
				"-internal-externc-isystem/Library/Developer/CommandLineTools/SDKs/MacOSX14.sdk/usr/include",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := clangtool.ParseSystemPath(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(output, tc.expect) {
				t.Fatalf("parse sysroot failed: want: %v got %v", tc.expect, output)
			}
		})
	}
}
