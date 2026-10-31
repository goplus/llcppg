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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// -----------------------------------------------------------------------------
// Platform-specific project layout.
//
// A platform-specific llcppg project keeps a self-contained header tree per
// target platform in a directory named exactly $GOOS-$GOARCH (following Go's
// own naming, e.g. "darwin-arm64"). Such a project opts in by starting its
// Dir with the literal "$GOOS-$GOARCH/" prefix:
//
//	"Dir": "$GOOS-$GOARCH/include"
//	"Dir": "$GOOS-$GOARCH/include/..."
//
// When the prefix is present, llcppg reads headers from <platform>/include and
// writes the generated package under <dest-dir>/<platform>. Projects without
// the prefix keep their existing behavior in full.
//
// See https://github.com/goplus/llcppg/issues/1125.

// platformPlaceholder is the only accepted placeholder form in Dir. It must
// appear as the leading path segment of Dir, followed by a path separator.
const platformPlaceholder = "$GOOS-$GOARCH"

// TargetPlatform returns the target GOOS and GOARCH for the current conversion.
// The GOOS and GOARCH environment variables select the target platform; when
// either is unset it falls back to the host value (runtime.GOOS / runtime.GOARCH),
// matching Go's own conventions.
func TargetPlatform() (goos, goarch string) {
	goos = os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch = os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return
}

// PlatformDir returns the $GOOS-$GOARCH directory name for the target platform,
// e.g. "darwin-arm64".
func PlatformDir() string {
	goos, goarch := TargetPlatform()
	return goos + "-" + goarch
}

// Platform reports whether the configuration describes a platform-specific
// project and, if so, the resolved $GOOS-$GOARCH directory name for the target
// platform (e.g. "darwin-arm64"). It also validates the Dir field:
//
//   - A Dir may reference $GOOS or $GOARCH only as the leading "$GOOS-$GOARCH/"
//     segment; any other usage is an error.
//   - A Dir without the prefix is an ordinary, platform-independent project
//     (platformDir == "", specific == false).
//
// Detection is driven purely by the Dir prefix.
func (cfg *Config) Platform() (platformDir string, specific bool, err error) {
	return platformOf(cfg.Dir)
}

// platformOf implements Config.Platform for a raw Dir value.
func platformOf(dir string) (platformDir string, specific bool, err error) {
	// The project is platform-specific iff Dir starts with the placeholder
	// segment "$GOOS-$GOARCH" followed by a path separator (or is exactly the
	// placeholder).
	rest, specific := cutPlatformPrefix(dir)
	if specific {
		// No further occurrence of the placeholder tokens is permitted. Scan
		// the remainder after the leading segment, but report the original dir
		// so the diagnostic matches the configured Dir value.
		if err = checkNoPlaceholder(dir, rest); err != nil {
			return
		}
		platformDir = PlatformDir()
		return
	}
	// Not prefixed: reject any stray $GOOS / $GOARCH elsewhere in the path so
	// that a misplaced placeholder is a clear error rather than a silent no-op.
	if err = checkNoPlaceholder(dir, dir); err != nil {
		return
	}
	return
}

// cutPlatformPrefix reports whether dir begins with the "$GOOS-$GOARCH" segment
// and, if so, returns the remainder of the path after that leading segment (and
// its following separator). The placeholder must be a whole leading path
// segment: "$GOOS-$GOARCH/include" matches, but "$GOOS-$GOARCHfoo" does not.
func cutPlatformPrefix(dir string) (rest string, ok bool) {
	if !strings.HasPrefix(dir, platformPlaceholder) {
		return dir, false
	}
	rest = dir[len(platformPlaceholder):]
	switch {
	case rest == "":
		// Dir is exactly "$GOOS-$GOARCH".
		return "", true
	case rest[0] == '/' || rest[0] == os.PathSeparator:
		return rest[1:], true
	default:
		// e.g. "$GOOS-$GOARCHfoo" — placeholder is not a whole segment.
		return dir, false
	}
}

// checkNoPlaceholder returns an error if scan still mentions $GOOS or $GOARCH.
// dir is the original Dir value used only in the error message, so the
// diagnostic names what the user configured rather than an internal substring.
func checkNoPlaceholder(dir, scan string) error {
	if strings.Contains(scan, "$GOOS") || strings.Contains(scan, "$GOARCH") {
		return fmt.Errorf(
			"invalid Dir %q: $GOOS/$GOARCH may only appear as the leading %q segment",
			dir, platformPlaceholder)
	}
	return nil
}

// resolveDir expands the leading $GOOS-$GOARCH placeholder of Dir (if any) to
// the concrete platform directory for the target platform and returns the
// resolved path together with the platform directory name. For ordinary
// projects it returns Dir unchanged and an empty platform name.
func (cfg *Config) resolveDir() (resolvedDir, platformDir string, err error) {
	platformDir, specific, err := cfg.Platform()
	if err != nil {
		return
	}
	if !specific {
		return cfg.Dir, "", nil
	}
	rest, _ := cutPlatformPrefix(cfg.Dir)
	if rest == "" {
		resolvedDir = platformDir
	} else {
		resolvedDir = platformDir + "/" + rest
	}
	return
}

// srcDir resolves Dir to the absolute source directory that holds the headers
// for the target platform. It expands the $GOOS-$GOARCH placeholder, strips a
// trailing "/..." recursive marker (reported via recursive), and joins a
// relative path onto workDir. platformDir is the $GOOS-$GOARCH directory name
// for a platform-specific project, or empty otherwise.
func (cfg *Config) srcDir(workDir string) (dir string, recursive bool, platformDir string, err error) {
	dir, platformDir, err = cfg.resolveDir()
	if err != nil {
		return
	}
	recursive = strings.HasSuffix(dir, "/...")
	if recursive {
		dir = dir[:len(dir)-4]
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(workDir, dir)
	}
	return
}

// CheckPlatformDir validates, for a platform-specific project, that the source
// directory selected by Dir for the target platform exists under workDir. It is
// a no-op for ordinary (non-platform-specific) projects. It lets a caller fail
// fast — before creating any output — when the headers for the requested
// platform are not present.
func (cfg *Config) CheckPlatformDir(workDir string) error {
	dir, _, platformDir, err := cfg.srcDir(workDir)
	if err != nil {
		return err
	}
	if platformDir == "" {
		return nil
	}
	return checkPlatformDir(dir, platformDir)
}

// checkPlatformDir reports an error, naming the expected path, when the
// resolved platform source directory does not exist.
func checkPlatformDir(dir, platformDir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("platform directory for target %q not found: %s",
				platformDir, dir)
		}
		return err
	}
	return nil
}

// -----------------------------------------------------------------------------
// Mapping a $GOOS-$GOARCH platform to the Clang target triple used when parsing
// headers, so conditional compilation, macros, and type layout are interpreted
// for the requested platform rather than the host.
//
// See https://github.com/goplus/llcppg/issues/1124.

// clangTriples maps a $GOOS-$GOARCH directory name to the Clang target triple.
// Only the platforms recognized by the initial implementation are listed; an
// unmapped platform yields no triple (TargetTriple reports ok == false).
var clangTriples = map[string]string{
	"linux-amd64":  "x86_64-unknown-linux-gnu",
	"linux-arm64":  "aarch64-unknown-linux-gnu",
	"darwin-amd64": "x86_64-apple-darwin",
	"darwin-arm64": "arm64-apple-darwin",
}

// targetTriple returns the Clang target triple for the target platform (the one
// selected by TargetPlatform) and whether a mapping exists. It is used to pass
// "-target <triple>" to libclang for platform-specific projects; ok is false
// for platforms without a known triple, in which case the caller parses with
// the host default as before.
func targetTriple(platformDir string) (triple string, ok bool) {
	triple, ok = clangTriples[platformDir]
	return
}

// -----------------------------------------------------------------------------
