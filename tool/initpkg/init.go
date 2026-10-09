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

// Package initpkg implements `llcppg -init`, which bootstraps a new binding
// repository from the llcppg template. See the package-level Init function for
// the full behavior.
package initpkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

const (
	// DefaultTemplateURL is the git URL of the llcppg project template.
	DefaultTemplateURL = "https://github.com/llarhub/.llcppg"

	// Placeholder is the token substituted with the module name in the files
	// listed by substituteFiles.
	Placeholder = "MODULE_NAME"

	// branchC and branchMain are the two branches a binding project uses.
	branchC    = "c"
	branchMain = "main"
)

// substituteFiles maps each branch to the files on it whose Placeholder
// occurrences are rewritten with the module name. Every listed file must exist
// in the template branch; a missing one is reported as a template error.
var substituteFiles = map[string][]string{
	branchC:    {"go.mod", "llcppg.cfg"},
	branchMain: {"go.mod"},
}

// Options configures Init. The zero value uses the public template and the
// per-user cache directory; tests and advanced callers can override both.
type Options struct {
	// TemplateURL is the git URL of the template repository. Empty means
	// DefaultTemplateURL.
	TemplateURL string

	// CacheDir is the directory that holds the cached template clone. Empty
	// means the location returned by CacheDir().
	CacheDir string

	// Stdout receives the human-readable progress and completion report. Empty
	// means os.Stdout.
	Stdout io.Writer

	// Stderr receives warnings (for example, the offline fallback notice).
	// Empty means os.Stderr.
	Stderr io.Writer
}

func (o *Options) templateURL() string {
	if o.TemplateURL != "" {
		return o.TemplateURL
	}
	return DefaultTemplateURL
}

func (o *Options) cacheDir() (string, error) {
	if o.CacheDir != "" {
		return o.CacheDir, nil
	}
	return CacheDir()
}

func (o *Options) stdout() io.Writer {
	if o.Stdout != nil {
		return o.Stdout
	}
	return os.Stdout
}

func (o *Options) stderr() io.Writer {
	if o.Stderr != nil {
		return o.Stderr
	}
	return os.Stderr
}

// CacheDir returns the directory llcppg uses to cache downloadable resources
// such as the project template. The LLCPPG_CACHE environment variable overrides
// the default, which is <os.UserCacheDir>/llcppg.
func CacheDir() (string, error) {
	if v := os.Getenv("LLCPPG_CACHE"); v != "" {
		return v, nil
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate user cache dir: %w", err)
	}
	return filepath.Join(base, "llcppg"), nil
}

// Init bootstraps a new binding repository from the template, in the current
// repository rooted at dir. It:
//
//  1. validates the environment (dir is a git repo, no existing c branch, no
//     existing go.mod);
//  2. prepares the template in the cache (clone on first use, update otherwise,
//     falling back to the cached copy when offline);
//  3. creates and commits a c branch with the template's c-branch files, with
//     MODULE_NAME substituted in go.mod and llcppg.cfg;
//  4. creates and commits a main branch with the template's main-branch files,
//     with MODULE_NAME substituted in go.mod;
//  5. leaves main checked out and reports what was done.
//
// It never pushes, adds or modifies remotes, or changes git configuration.
func Init(dir, module string, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}
	module = strings.TrimSpace(module)
	if module == "" {
		return errors.New("llcppg: cannot init: module name must not be empty")
	}

	// Step 1: validate the environment up front, before touching the cache or
	// writing any file, so a refused run leaves everything exactly as it was.
	if err := checkEnvironment(dir); err != nil {
		return err
	}

	// Step 2: prepare the cached template.
	cacheDir, err := opts.cacheDir()
	if err != nil {
		return err
	}
	tmplDir := filepath.Join(cacheDir, "templates", "llarhub", ".llcppg")
	if err := prepareTemplate(tmplDir, opts.templateURL(), opts.stderr()); err != nil {
		return err
	}

	out := opts.stdout()

	// Steps 3-6: create, populate, and commit the c branch. c is the first
	// branch; on a freshly cloned, empty repository it is created on the unborn
	// HEAD.
	if err := setupBranch(dir, tmplDir, branchC, module, false, out); err != nil {
		return err
	}

	// Steps 7-10: create, populate, and commit the main branch. main is created
	// as an orphan so the C-side and Go-side files never share history or mix
	// (see the proposal's open question #1).
	if err := setupBranch(dir, tmplDir, branchMain, module, true, out); err != nil {
		return err
	}

	// Step 11: report. main is left checked out by setupBranch above.
	fmt.Fprintf(out, "\nllcppg: initialized %q with the %q and %q branches.\n", module, branchC, branchMain)
	fmt.Fprintln(out, "Nothing has been pushed. Review the result, then push with:")
	fmt.Fprintf(out, "\n    git push -u origin %s %s\n", branchC, branchMain)
	return nil
}

// setupBranch creates branch (checking it out), copies the template branch's
// files into dir's working tree, substitutes MODULE_NAME in the configured
// files, and commits the result. When orphan is true the branch is created as a
// fresh root with no parent and no inherited files, so the two branches never
// mix content.
func setupBranch(dir, tmplDir, branch, module string, orphan bool, out io.Writer) error {
	fmt.Fprintf(out, "llcppg: setting up the %q branch...\n", branch)
	if err := checkoutNewBranch(dir, branch, orphan); err != nil {
		return err
	}
	if err := extractBranch(tmplDir, branch, dir); err != nil {
		return err
	}
	if err := substitute(dir, branch, module); err != nil {
		return err
	}
	msg := fmt.Sprintf("init: add llcppg template (%s branch) for %s", branch, module)
	return commitAll(dir, msg)
}

// substitute rewrites every Placeholder occurrence with module in each file
// configured for branch. A configured file missing from the extracted template
// is reported as a template error.
func substitute(dir, branch, module string) error {
	for _, name := range substituteFiles[branch] {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("llcppg: template error: %q is missing from the template %q branch", name, branch)
			}
			return err
		}
		replaced := strings.ReplaceAll(string(data), Placeholder, module)
		if replaced == string(data) {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(replaced), info.Mode().Perm()); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
