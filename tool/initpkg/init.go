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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gomodule "golang.org/x/mod/module"
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
// occurrences are rewritten with the module name, as repository-relative paths.
// On the c branch these live under the c/ directory (c/go.mod, c/llcppg.cfg);
// on main the go.mod is at the root. Every listed file must exist in the
// template branch; a missing one is reported as a template error.
var substituteFiles = map[string][]string{
	branchC:    {"c/go.mod", "c/llcppg.cfg"},
	branchMain: {"go.mod"},
}

// Config configures Init. The zero value uses the public template and the
// per-user cache directory; tests and advanced callers can override both.
type Config struct {
	// TemplateURL is the git URL of the template repository. Empty means
	// DefaultTemplateURL.
	TemplateURL string

	// CacheDir is the directory that holds the cached template clone. required.
	CacheDir string

	// Stdout receives the human-readable progress and completion report. Empty
	// means os.Stdout.
	Stdout io.Writer

	// Stderr receives warnings (for example, the offline fallback notice).
	// Empty means os.Stderr.
	Stderr io.Writer
}

func (o *Config) templateURL() string {
	if o.TemplateURL != "" {
		return o.TemplateURL
	}
	return DefaultTemplateURL
}

func (o *Config) stdout() io.Writer {
	if o.Stdout != nil {
		return o.Stdout
	}
	return os.Stdout
}

func (o *Config) stderr() io.Writer {
	if o.Stderr != nil {
		return o.Stderr
	}
	return os.Stderr
}

// templateCacheDir returns the directory under cacheDir that holds the cached
// clone of templateURL. The subpath is derived from the URL (host, then path
// segments) so that different template URLs sharing a cache dir never collide
// and the location is not hardcoded to a single owner/repo. All templates live
// under a shared "templates" parent.
func templateCacheDir(cacheDir, templateURL string) string {
	segs := []string{cacheDir, "templates"}
	segs = append(segs, urlCacheSegments(templateURL)...)
	return filepath.Join(segs...)
}

// urlCacheSegments splits a git URL into sanitized path segments (host and path
// components) suitable for use as a cache subpath. It handles the common forms
// (https://host/owner/repo, scp-like git@host:owner/repo, and file:// URLs)
// without pulling in a full URL parser, and falls back to a single sanitized
// segment for anything it does not recognize so the result is always usable.
func urlCacheSegments(templateURL string) []string {
	s := templateURL
	// Strip a scheme (scheme://...) or an scp-like user@host: prefix.
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+len("://"):]
		if at := strings.Index(s, "@"); at >= 0 && at < strings.IndexAny(s+"/", "/") {
			s = s[at+1:]
		}
	} else if at := strings.Index(s, "@"); at >= 0 {
		// scp-like syntax git@host:owner/repo -> host/owner/repo.
		s = s[at+1:]
		s = strings.Replace(s, ":", "/", 1)
	}
	s = strings.TrimSuffix(s, ".git")
	var segs []string
	for _, part := range strings.FieldsFunc(s, func(r rune) bool { return r == '/' }) {
		if clean := sanitizeSegment(part); clean != "" {
			segs = append(segs, clean)
		}
	}
	if len(segs) == 0 {
		return []string{"template"}
	}
	return segs
}

// sanitizeSegment keeps a path segment safe to use as a directory name: it drops
// any path separators or traversal and replaces characters that are awkward on
// common filesystems, so a crafted URL can never escape the cache directory.
func sanitizeSegment(seg string) string {
	seg = strings.TrimSpace(seg)
	if seg == "" || seg == "." || seg == ".." {
		return ""
	}
	var b strings.Builder
	for _, r := range seg {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := b.String()
	// A leading dot is legal and intentional (e.g. the ".llcppg" repo).
	if out == "." || out == ".." {
		return ""
	}
	return out
}

// Init bootstraps a new binding repository from the template, in the current
// repository rooted at dir. It:
//
//  1. determines the module name (the given module, or the current directory's
//     name when module is empty) and validates the environment (dir is a git
//     repo, no existing c branch, no existing go.mod);
//  2. prepares the template in the cache (clone on first use, update otherwise,
//     falling back to the cached copy when offline);
//  3. creates and commits a c branch with the template's c-branch files, with
//     MODULE_NAME substituted in c/go.mod and c/llcppg.cfg;
//  4. switches to the main branch (creating it only if it does not already
//     exist) and commits the template's main-branch files, with MODULE_NAME
//     substituted in go.mod;
//  5. leaves main checked out and reports what was done.
//
// It never pushes, adds or modifies remotes, or changes git configuration.
func Init(dir, module string, opts *Config) error {
	if opts == nil {
		opts = &Config{}
	}

	// Step 1: determine the module name and validate the environment up front,
	// before touching the cache or writing any file, so a refused run leaves
	// everything exactly as it was.
	module, err := moduleName(dir, module)
	if err != nil {
		return err
	}
	if err := checkEnvironment(dir); err != nil {
		return err
	}

	out := opts.stdout()
	// Report the module name before making any change, so the user sees which
	// name is in effect whether they passed it or it was inferred.
	fmt.Fprintf(out, "llcppg: using module name %q\n", module)

	// Step 2: prepare the cached template.
	cacheDir := opts.CacheDir
	tmplDir := templateCacheDir(cacheDir, opts.templateURL())
	if err := prepareTemplate(tmplDir, opts.templateURL(), opts.stderr()); err != nil {
		return err
	}

	// Step 3: create, populate, and commit the c branch. c is the first branch;
	// on a freshly cloned, empty repository it is created on the unborn HEAD.
	if err := setupBranch(dir, tmplDir, branchC, module, false, out); err != nil {
		return err
	}

	// Step 4: switch to main (a freshly cloned repository already has it as its
	// default branch) and populate and commit it. If main does not exist it is
	// created as an orphan so the C-side and Go-side files never share history
	// or mix (see the proposal's open question #1).
	if err := setupBranch(dir, tmplDir, branchMain, module, true, out); err != nil {
		return err
	}

	// Step 5: report. main is left checked out by setupBranch above.
	fmt.Fprintf(out, "\nllcppg: initialized %q with the %q and %q branches.\n", module, branchC, branchMain)
	fmt.Fprintln(out, "Nothing has been pushed. Review the result, then push with:")
	fmt.Fprintf(out, "\n    git push -u origin %s %s\n", branchC, branchMain)
	return nil
}

// moduleName resolves the module name to use. A non-empty module (after
// trimming surrounding spaces) is used as given; an empty module is inferred
// from the last element of dir's absolute path, matching the directory a
// freshly cloned project sits in. Either way the result must be a valid module
// path; an inferred name that is not (for example a directory name with spaces)
// is rejected with a message asking the user to pass the name explicitly.
func moduleName(dir, module string) (string, error) {
	module = strings.TrimSpace(module)
	inferred := false
	if module == "" {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		module = filepath.Base(abs)
		inferred = true
	}
	if err := gomodule.CheckImportPath(module); err != nil {
		if inferred {
			return "", fmt.Errorf("llcppg: cannot init: cannot use the current directory name %q as a module name (%v); pass the module name explicitly: llcppg -init <module-name>", module, err)
		}
		return "", fmt.Errorf("llcppg: cannot init: invalid module name %q: %v", module, err)
	}
	return module, nil
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
