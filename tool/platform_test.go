/*
 * Copyright (c) 2026 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tool

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// withPlatform pins GOOS/GOARCH for the duration of a test so results are
// deterministic regardless of the host.
func withPlatform(t *testing.T, goos, goarch string) {
	t.Helper()
	t.Setenv("GOOS", goos)
	t.Setenv("GOARCH", goarch)
}

func TestTargetPlatform(t *testing.T) {
	withPlatform(t, "darwin", "arm64")
	if goos, goarch := TargetPlatform(); goos != "darwin" || goarch != "arm64" {
		t.Fatalf("TargetPlatform() = %q/%q, want darwin/arm64", goos, goarch)
	}
	if got := PlatformDir(); got != "darwin-arm64" {
		t.Fatalf("PlatformDir() = %q, want darwin-arm64", got)
	}
}

func TestTargetPlatformFallsBackToHost(t *testing.T) {
	// Empty env values must fall back to the host runtime values.
	t.Setenv("GOOS", "")
	t.Setenv("GOARCH", "")
	goos, goarch := TargetPlatform()
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		t.Fatalf("TargetPlatform() = %q/%q, want host %q/%q",
			goos, goarch, runtime.GOOS, runtime.GOARCH)
	}
}

func TestPlatform(t *testing.T) {
	withPlatform(t, "linux", "amd64")

	for _, tt := range []struct {
		name         string
		dir          string
		wantSpecific bool
		wantDir      string // expected platformDir when specific
		wantErr      bool
	}{
		{"plain", "./include", false, "", false},
		{"plain recursive", "./include/...", false, "", false},
		{"bare dot", ".", false, "", false},
		{"prefixed", "$GOOS-$GOARCH/include", true, "linux-amd64", false},
		{"prefixed recursive", "$GOOS-$GOARCH/include/...", true, "linux-amd64", false},
		{"prefixed sep variant", "$GOOS-$GOARCH/include/clang-c", true, "linux-amd64", false},
		{"exact placeholder", "$GOOS-$GOARCH", true, "linux-amd64", false},
		// $GOOS/$GOARCH anywhere other than the leading segment is an error.
		{"trailing placeholder", "include/$GOOS-$GOARCH", false, "", true},
		{"goos only mid path", "include/$GOOS", false, "", true},
		{"goarch only mid path", "include/$GOARCH", false, "", true},
		{"placeholder after prefix", "$GOOS-$GOARCH/include/$GOOS", false, "", true},
		// The placeholder must be a whole leading segment, not a substring.
		{"not a whole segment", "$GOOS-$GOARCHfoo/include", false, "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Dir: tt.dir}
			platformDir, specific, err := cfg.Platform()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Platform(%q) = nil err, want error", tt.dir)
				}
				return
			}
			if err != nil {
				t.Fatalf("Platform(%q) unexpected err: %v", tt.dir, err)
			}
			if specific != tt.wantSpecific {
				t.Fatalf("Platform(%q) specific = %v, want %v", tt.dir, specific, tt.wantSpecific)
			}
			if specific && platformDir != tt.wantDir {
				t.Fatalf("Platform(%q) platformDir = %q, want %q", tt.dir, platformDir, tt.wantDir)
			}
			if !specific && platformDir != "" {
				t.Fatalf("Platform(%q) platformDir = %q, want empty", tt.dir, platformDir)
			}
		})
	}
}

func TestResolveDir(t *testing.T) {
	withPlatform(t, "darwin", "arm64")

	for _, tt := range []struct {
		name         string
		dir          string
		wantDir      string
		wantPlatform string
		wantErr      bool
	}{
		{"plain unchanged", "./include", "./include", "", false},
		{"plain recursive unchanged", "./include/...", "./include/...", "", false},
		{"prefixed", "$GOOS-$GOARCH/include", "darwin-arm64/include", "darwin-arm64", false},
		{"prefixed recursive", "$GOOS-$GOARCH/include/...", "darwin-arm64/include/...", "darwin-arm64", false},
		{"exact placeholder", "$GOOS-$GOARCH", "darwin-arm64", "darwin-arm64", false},
		{"invalid placeholder", "include/$GOOS", "", "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Dir: tt.dir}
			gotDir, gotPlatform, err := cfg.resolveDir()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("resolveDir(%q) = nil err, want error", tt.dir)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveDir(%q) unexpected err: %v", tt.dir, err)
			}
			if gotDir != tt.wantDir {
				t.Fatalf("resolveDir(%q) dir = %q, want %q", tt.dir, gotDir, tt.wantDir)
			}
			if gotPlatform != tt.wantPlatform {
				t.Fatalf("resolveDir(%q) platform = %q, want %q", tt.dir, gotPlatform, tt.wantPlatform)
			}
		})
	}
}

func TestCutPlatformPrefix(t *testing.T) {
	for _, tt := range []struct {
		dir      string
		wantRest string
		wantOK   bool
	}{
		{"$GOOS-$GOARCH/include", "include", true},
		{"$GOOS-$GOARCH/include/...", "include/...", true},
		{"$GOOS-$GOARCH", "", true},
		{"$GOOS-$GOARCHfoo", "$GOOS-$GOARCHfoo", false},
		{"./include", "./include", false},
		{"include/$GOOS-$GOARCH", "include/$GOOS-$GOARCH", false},
	} {
		rest, ok := cutPlatformPrefix(tt.dir)
		if ok != tt.wantOK || rest != tt.wantRest {
			t.Fatalf("cutPlatformPrefix(%q) = %q,%v; want %q,%v",
				tt.dir, rest, ok, tt.wantRest, tt.wantOK)
		}
	}
}

func TestCheckNoPlaceholderMessage(t *testing.T) {
	err := checkNoPlaceholder("include/$GOOS")
	if err == nil {
		t.Fatal("checkNoPlaceholder should reject stray $GOOS")
	}
	if !strings.Contains(err.Error(), platformPlaceholder) {
		t.Fatalf("error should mention %q, got %q", platformPlaceholder, err.Error())
	}
}

func TestSrcDir(t *testing.T) {
	withPlatform(t, "linux", "amd64")
	workDir := "/work"

	for _, tt := range []struct {
		name          string
		dir           string
		wantDir       string
		wantRecursive bool
		wantPlatform  string
	}{
		{"plain", "./include", filepath.Join(workDir, "include"), false, ""},
		{"plain recursive", "./include/...", filepath.Join(workDir, "include"), true, ""},
		{"prefixed", "$GOOS-$GOARCH/include", filepath.Join(workDir, "linux-amd64/include"), false, "linux-amd64"},
		{"prefixed recursive", "$GOOS-$GOARCH/include/...", filepath.Join(workDir, "linux-amd64/include"), true, "linux-amd64"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Dir: tt.dir}
			dir, recursive, platformDir, err := cfg.srcDir(workDir)
			if err != nil {
				t.Fatalf("srcDir(%q) err: %v", tt.dir, err)
			}
			if dir != tt.wantDir {
				t.Fatalf("srcDir(%q) dir = %q, want %q", tt.dir, dir, tt.wantDir)
			}
			if recursive != tt.wantRecursive {
				t.Fatalf("srcDir(%q) recursive = %v, want %v", tt.dir, recursive, tt.wantRecursive)
			}
			if platformDir != tt.wantPlatform {
				t.Fatalf("srcDir(%q) platformDir = %q, want %q", tt.dir, platformDir, tt.wantPlatform)
			}
		})
	}
}

func TestCheckPlatformDir(t *testing.T) {
	withPlatform(t, "linux", "amd64")

	// A non-platform-specific project never fails this check.
	plain := &Config{Dir: "./include"}
	if err := plain.CheckPlatformDir(t.TempDir()); err != nil {
		t.Fatalf("CheckPlatformDir for plain project err: %v", err)
	}

	// A platform-specific project whose directory is missing must fail with a
	// message naming the expected path.
	work := t.TempDir()
	cfg := &Config{Dir: "$GOOS-$GOARCH/include/..."}
	err := cfg.CheckPlatformDir(work)
	if err == nil {
		t.Fatal("CheckPlatformDir should fail when platform dir is missing")
	}
	want := filepath.Join(work, "linux-amd64", "include")
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error should name %q, got %q", want, err.Error())
	}

	// Once the directory exists the check passes.
	if e := os.MkdirAll(want, 0755); e != nil {
		t.Fatal(e)
	}
	if err := cfg.CheckPlatformDir(work); err != nil {
		t.Fatalf("CheckPlatformDir after mkdir err: %v", err)
	}
}
