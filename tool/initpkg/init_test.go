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

package initpkg

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// Test helpers.

func mustGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := gitOutput(dir, args...)
	if err != nil {
		t.Fatalf("git %s (in %s): %v", strings.Join(args, " "), dir, err)
	}
	return out
}

// testEnv configures a deterministic git identity so commits work without a
// global git config.
func testEnv() []string {
	return append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_TERMINAL_PROMPT=0",
	)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = testEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s (in %s): %v\n%s", strings.Join(args, " "), dir, err, out)
	}
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// newTemplateRepo creates a bare template repository with a c and a main
// branch, each carrying the placeholder files the proposal expects. It returns
// the file:// URL callers pass as Options.TemplateURL.
func newTemplateRepo(t *testing.T) string {
	t.Helper()
	src := t.TempDir()
	runGit(t, src, "init", "-q", "-b", branchMain)

	// main branch content.
	writeTestFile(t, src, "go.mod", "module "+Placeholder+"\n\ngo 1.27.0\n")
	writeTestFile(t, src, "main.go", "package main\n")
	runGit(t, src, "add", "-A")
	runGit(t, src, "commit", "-q", "-m", "template main")

	// c branch content: an orphan branch so c and main files never mix. The
	// C-side go.mod and llcppg.cfg live under the c/ directory.
	runGit(t, src, "checkout", "-q", "--orphan", branchC)
	runGit(t, src, "rm", "-rf", "-q", ".")
	writeTestFile(t, src, "c/go.mod", "module "+Placeholder+"\n\ngo 1.27.0\n")
	writeTestFile(t, src, "c/llcppg.cfg", `{"Name":"`+Placeholder+`"}`+"\n")
	writeTestFile(t, src, "nested/header.h", "// header for "+Placeholder+"\n")
	runGit(t, src, "add", "-A")
	runGit(t, src, "commit", "-q", "-m", "template c")
	runGit(t, src, "checkout", "-q", branchMain)

	// A bare clone acts as the "remote" the cache fetches from.
	bare := filepath.Join(t.TempDir(), "template.git")
	cmd := exec.Command("git", "clone", "-q", "--bare", src, bare)
	cmd.Env = testEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("clone bare template: %v\n%s", err, out)
	}
	return "file://" + bare
}

// newEmptyTarget creates a freshly initialized (empty, unborn-HEAD) git repo,
// mirroring a just-cloned binding project.
func newEmptyTarget(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-q", "-b", branchMain)
	return dir
}

// runInit invokes Init with the test identity in the process environment so the
// git commits it makes succeed.
func runInit(t *testing.T, dir, module string, opts *Config) (string, string, error) {
	t.Helper()
	// gitOutput inherits os.Environ(); set the identity for this test run.
	for _, kv := range []string{
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	} {
		k, v, _ := strings.Cut(kv, "=")
		t.Setenv(k, v)
	}
	var stdout, stderr bytes.Buffer
	if opts == nil {
		opts = &Config{}
	}
	opts.Stdout = &stdout
	opts.Stderr = &stderr
	err := Init(dir, module, opts)
	return stdout.String(), stderr.String(), err
}

func fileOnBranch(t *testing.T, dir, branch, file string) (string, bool) {
	t.Helper()
	out, err := gitOutput(dir, "show", branch+":"+file)
	if err != nil {
		return "", false
	}
	return out, true
}

// -----------------------------------------------------------------------------
// Tests.

func TestInit_HappyPath(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	target := newEmptyTarget(t)

	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// main is checked out.
	if got := strings.TrimSpace(mustGit(t, target, "rev-parse", "--abbrev-ref", "HEAD")); got != branchMain {
		t.Fatalf("HEAD = %q, want %q", got, branchMain)
	}

	// Both branches exist with exactly one commit each.
	for _, b := range []string{branchC, branchMain} {
		count := strings.TrimSpace(mustGit(t, target, "rev-list", "--count", b))
		if count != "1" {
			t.Errorf("branch %q has %s commits, want 1", b, count)
		}
	}

	// go.mod has the module name and no placeholder on each branch: at the root
	// on main, and under c/ on the c branch.
	goModOnBranch := map[string]string{branchC: "c/go.mod", branchMain: "go.mod"}
	for _, b := range []string{branchC, branchMain} {
		content, ok := fileOnBranch(t, target, b, goModOnBranch[b])
		if !ok {
			t.Fatalf("%s missing on branch %q", goModOnBranch[b], b)
		}
		if strings.Contains(content, Placeholder) {
			t.Errorf("branch %q go.mod still contains %q:\n%s", b, Placeholder, content)
		}
		if !strings.Contains(content, "module cjson") {
			t.Errorf("branch %q go.mod missing module name:\n%s", b, content)
		}
	}

	// c/llcppg.cfg on c has the module name and no placeholder.
	cfg, ok := fileOnBranch(t, target, branchC, "c/llcppg.cfg")
	if !ok {
		t.Fatal("c/llcppg.cfg missing on c branch")
	}
	if strings.Contains(cfg, Placeholder) || !strings.Contains(cfg, "cjson") {
		t.Errorf("c/llcppg.cfg not substituted: %q", cfg)
	}

	// A nested template file came through on c.
	if _, ok := fileOnBranch(t, target, branchC, "nested/header.h"); !ok {
		t.Error("nested/header.h missing on c branch")
	}
}

func TestInit_RefusesExistingCBranch(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	target := newEmptyTarget(t)

	if _, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("first Init: %v", err)
	}

	// Record state before the refused second run.
	headBefore := strings.TrimSpace(mustGit(t, target, "rev-parse", "HEAD"))

	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err == nil {
		t.Fatal("second Init: expected error on existing c branch, got nil")
	}
	if !strings.Contains(err.Error(), `branch "c" already exists`) {
		t.Errorf("unexpected error: %v", err)
	}
	if headAfter := strings.TrimSpace(mustGit(t, target, "rev-parse", "HEAD")); headAfter != headBefore {
		t.Errorf("HEAD changed on refused run: %s -> %s", headBefore, headAfter)
	}
}

func TestInit_RefusesExistingGoMod(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	target := newEmptyTarget(t)
	writeTestFile(t, target, "go.mod", "module existing\n")

	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err == nil {
		t.Fatal("expected error on existing go.mod, got nil")
	}
	if !strings.Contains(err.Error(), "go.mod already exists") {
		t.Errorf("unexpected error: %v", err)
	}
	// No branch was created on the refused run.
	if exists, _ := branchCExists(target); exists {
		t.Error("c branch created despite refused run")
	}
}

func TestInit_RefusedRunDoesNotTouchCache(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	target := newEmptyTarget(t)
	writeTestFile(t, target, "go.mod", "module existing\n") // forces a refusal

	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err == nil {
		t.Fatal("expected refusal")
	}
	// The cache must remain empty: safety checks run before the cache is touched.
	entries, _ := os.ReadDir(cache)
	if len(entries) != 0 {
		t.Errorf("cache was touched on refused run: %v", entries)
	}
}

func TestInit_NotAGitRepo(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	dir := t.TempDir() // plain directory, no git init

	_, _, err := runInit(t, dir, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err == nil {
		t.Fatal("expected error outside a git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not the root of a git repository") {
		t.Errorf("unexpected error: %v", err)
	}
	// No files or branches left behind.
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("directory not left clean: %v", entries)
	}
}

func TestInit_DerivesModuleNameFromDir(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()

	// A target whose directory base name is a valid module name. Passing an
	// empty (here, all-whitespace) module name must infer it from that name.
	target := filepath.Join(t.TempDir(), "cjson")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, target, "init", "-q", "-b", branchMain)

	stdout, _, err := runInit(t, target, "   ", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err != nil {
		t.Fatalf("Init with inferred module name: %v", err)
	}
	if !strings.Contains(stdout, `using module name "cjson"`) {
		t.Errorf("expected the inferred module name to be reported, got: %q", stdout)
	}
	// go.mod on each branch uses the inferred name: root on main, c/ on c.
	goModOnBranch := map[string]string{branchC: "c/go.mod", branchMain: "go.mod"}
	for _, b := range []string{branchC, branchMain} {
		content, ok := fileOnBranch(t, target, b, goModOnBranch[b])
		if !ok {
			t.Fatalf("%s missing on branch %q", goModOnBranch[b], b)
		}
		if strings.Contains(content, Placeholder) || !strings.Contains(content, "module cjson") {
			t.Errorf("branch %q go.mod not substituted with the inferred name:\n%s", b, content)
		}
	}
}

func TestInit_InferredNameInvalidReportsError(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()

	// A directory name that cannot serve as a module path (contains a space).
	target := filepath.Join(t.TempDir(), "bad name")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, target, "init", "-q", "-b", branchMain)

	_, _, err := runInit(t, target, "", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err == nil {
		t.Fatal("expected an error for an un-usable inferred module name, got nil")
	}
	if !strings.Contains(err.Error(), "pass the module name explicitly") {
		t.Errorf("unexpected error: %v", err)
	}
	// The refusal happens before the cache is touched.
	entries, _ := os.ReadDir(cache)
	if len(entries) != 0 {
		t.Errorf("cache was touched on refused run: %v", entries)
	}
}

func TestInit_InvalidExplicitModuleName(t *testing.T) {
	tmpl := newTemplateRepo(t)
	target := newEmptyTarget(t)

	_, _, err := runInit(t, target, "bad name", &Config{TemplateURL: tmpl, CacheDir: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "invalid module name") {
		t.Fatalf("expected an invalid-module-name error, got %v", err)
	}
}

func TestInit_CacheClonedThenUpdatedInPlace(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	tmplCacheDir := templateCacheDir(cache, tmpl)

	// First run clones into the cache.
	target1 := newEmptyTarget(t)
	if _, _, err := runInit(t, target1, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("first Init: %v", err)
	}
	if !isGitRepoRoot(tmplCacheDir) {
		t.Fatal("template was not cloned into the cache")
	}
	headAfterClone := strings.TrimSpace(mustGit(t, tmplCacheDir, "rev-parse", "origin/"+branchMain))

	// Push a new commit to the template remote so the next run must update.
	bare := strings.TrimPrefix(tmpl, "file://")
	work := t.TempDir()
	cmd := exec.Command("git", "clone", "-q", bare, work)
	cmd.Env = testEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("clone template remote: %v\n%s", err, out)
	}
	runGit(t, work, "checkout", "-q", branchMain)
	writeTestFile(t, work, "NEW.md", "new template file\n")
	runGit(t, work, "add", "-A")
	runGit(t, work, "commit", "-q", "-m", "add NEW.md")
	runGit(t, work, "push", "-q", "origin", branchMain)

	// Second run (different empty repo) must update the existing clone in place
	// (not re-clone) and pick up the new commit.
	target2 := newEmptyTarget(t)
	if _, _, err := runInit(t, target2, "zlib", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("second Init: %v", err)
	}
	headAfterUpdate := strings.TrimSpace(mustGit(t, tmplCacheDir, "rev-parse", "origin/"+branchMain))
	if headAfterUpdate == headAfterClone {
		t.Error("cache was not updated with the new template commit")
	}
	if _, ok := fileOnBranch(t, target2, branchMain, "NEW.md"); !ok {
		t.Error("new template file did not reach the generated project")
	}
}

func TestInit_LocalCacheModificationsDoNotLeak(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	tmplCacheDir := templateCacheDir(cache, tmpl)

	// Populate the cache via a first run.
	target1 := newEmptyTarget(t)
	if _, _, err := runInit(t, target1, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("first Init: %v", err)
	}

	// Make a local, uncommitted modification inside the cached checkout.
	writeTestFile(t, tmplCacheDir, "go.mod", "module LOCALLY_HACKED\n")

	// A fresh run must ignore the dirty working tree and use branch contents.
	target2 := newEmptyTarget(t)
	if _, _, err := runInit(t, target2, "zlib", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("second Init: %v", err)
	}
	content, _ := fileOnBranch(t, target2, branchMain, "go.mod")
	if strings.Contains(content, "LOCALLY_HACKED") {
		t.Errorf("local cache modification leaked into project: %q", content)
	}
}

func TestInit_OfflineUsesCachedCopy(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()

	// Populate the cache.
	target1 := newEmptyTarget(t)
	if _, _, err := runInit(t, target1, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("first Init: %v", err)
	}

	// Delete the remote so a fetch fails, simulating being offline.
	if err := os.RemoveAll(strings.TrimPrefix(tmpl, "file://")); err != nil {
		t.Fatal(err)
	}

	target2 := newEmptyTarget(t)
	_, stderr, err := runInit(t, target2, "zlib", &Config{TemplateURL: tmpl, CacheDir: cache})
	if err != nil {
		t.Fatalf("offline Init with populated cache should succeed, got: %v", err)
	}
	if !strings.Contains(stderr, "warning") {
		t.Errorf("expected an offline warning on stderr, got: %q", stderr)
	}
	if _, ok := fileOnBranch(t, target2, branchMain, "go.mod"); !ok {
		t.Error("offline run did not produce a go.mod from the cached copy")
	}
}

func TestInit_OfflineEmptyCacheFails(t *testing.T) {
	cache := t.TempDir()
	target := newEmptyTarget(t)
	// A template URL that cannot be cloned, with an empty cache.
	bogus := "file://" + filepath.Join(t.TempDir(), "does-not-exist.git")

	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: bogus, CacheDir: cache})
	if err == nil {
		t.Fatal("expected failure with empty cache and no network")
	}
	if !strings.Contains(err.Error(), "unavailable") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInit_BrokenCacheIsReclone(t *testing.T) {
	tmpl := newTemplateRepo(t)
	cache := t.TempDir()
	tmplCacheDir := templateCacheDir(cache, tmpl)

	// Simulate an interrupted clone: a directory that exists but is not a valid
	// git repo.
	if err := os.MkdirAll(tmplCacheDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, tmplCacheDir, "partial.tmp", "junk\n")

	target := newEmptyTarget(t)
	if _, _, err := runInit(t, target, "cjson", &Config{TemplateURL: tmpl, CacheDir: cache}); err != nil {
		t.Fatalf("Init should recover from a broken cache, got: %v", err)
	}
	if !isGitRepoRoot(tmplCacheDir) {
		t.Error("broken cache was not re-cloned into a valid repo")
	}
	if _, ok := fileOnBranch(t, target, branchMain, "go.mod"); !ok {
		t.Error("recovered run did not produce a go.mod")
	}
}

func TestInit_MissingTemplateFileReportsTemplateError(t *testing.T) {
	// Build a template whose c branch lacks the required c/llcppg.cfg.
	src := t.TempDir()
	runGit(t, src, "init", "-q", "-b", branchMain)
	writeTestFile(t, src, "go.mod", "module "+Placeholder+"\n")
	runGit(t, src, "add", "-A")
	runGit(t, src, "commit", "-q", "-m", "main")
	runGit(t, src, "checkout", "-q", "--orphan", branchC)
	runGit(t, src, "rm", "-rf", "-q", ".")
	writeTestFile(t, src, "c/go.mod", "module "+Placeholder+"\n") // no c/llcppg.cfg
	runGit(t, src, "add", "-A")
	runGit(t, src, "commit", "-q", "-m", "c")
	runGit(t, src, "checkout", "-q", branchMain)
	bare := filepath.Join(t.TempDir(), "t.git")
	cmd := exec.Command("git", "clone", "-q", "--bare", src, bare)
	cmd.Env = testEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bare clone: %v\n%s", err, out)
	}

	target := newEmptyTarget(t)
	_, _, err := runInit(t, target, "cjson", &Config{TemplateURL: "file://" + bare, CacheDir: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "template error") {
		t.Fatalf("expected a template error for missing c/llcppg.cfg, got: %v", err)
	}
}

// -----------------------------------------------------------------------------
