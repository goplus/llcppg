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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

// checkEnvironment runs the up-front safety checks: dir must be the root of a
// git repository, must not already have a c branch (local or remote-tracking),
// and must not already have a go.mod. All checks run before anything is written
// or the cache is touched, so a refused run makes no changes. The c-branch
// check runs before the go.mod check because an existing c branch is the
// stronger signal that the project has already been initialized (a prior run
// leaves a go.mod behind too).
func checkEnvironment(dir string) error {
	if !isGitRepoRoot(dir) {
		return fmt.Errorf("llcppg: cannot init: %q is not the root of a git repository", dir)
	}
	exists, err := branchCExists(dir)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("llcppg: cannot init: branch %q already exists in this repository", branchC)
	}
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return fmt.Errorf("llcppg: cannot init: go.mod already exists in the current directory")
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// isGitRepoRoot reports whether dir is the top level of a git working tree.
func isGitRepoRoot(dir string) bool {
	out, err := gitOutput(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return false
	}
	top, err := filepath.Abs(strings.TrimSpace(out))
	if err != nil {
		return false
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	// Resolve symlinks so that, e.g., macOS /var vs /private/var compares equal.
	if resolved, err := filepath.EvalSymlinks(top); err == nil {
		top = resolved
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	return top == abs
}

// branchCExists reports whether a local "c" branch or a remote-tracking "c"
// branch (such as origin/c) exists; either indicates the project has already
// been initialized.
func branchCExists(dir string) (bool, error) {
	out, err := gitOutput(dir, "for-each-ref", "--format=%(refname)",
		"refs/heads/"+branchC, "refs/remotes/*/"+branchC)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// localBranchExists reports whether a local branch with the given name already
// exists in dir.
func localBranchExists(dir, branch string) bool {
	_, err := gitOutput(dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// -----------------------------------------------------------------------------

// prepareTemplate ensures a usable clone of templateURL exists at tmplDir and
// its c and main branches match the remote. On first use it makes a blobless
// partial clone (file contents are fetched on demand by the later `git
// archive`, so full history is never downloaded); on later use it fetches and
// resets the branches. If the directory exists but is not a valid clone it is
// removed and re-cloned. When the network is unavailable but a cached copy
// exists, it warns and continues with the cached copy; with no cached copy and
// a failed clone, it returns an error.
func prepareTemplate(tmplDir, templateURL string, stderr io.Writer) error {
	if strings.HasPrefix(templateURL, "-") {
		// Defense in depth: a URL that looks like an option could be
		// interpreted by git as a flag rather than a repository.
		return fmt.Errorf("llcppg: cannot init: invalid template URL %q: must not start with %q", templateURL, "-")
	}
	if !isGitRepoRoot(tmplDir) {
		// A leftover directory that is not a valid clone (e.g. an interrupted
		// clone) is discarded so the fresh clone can proceed.
		if _, err := os.Stat(tmplDir); err == nil {
			if err := os.RemoveAll(tmplDir); err != nil {
				return err
			}
		}
		if err := os.MkdirAll(filepath.Dir(tmplDir), 0755); err != nil {
			return err
		}
		// A blobless partial clone of all branches: only the two branch tips are
		// ever read (via `git archive`), so file blobs are fetched on demand and
		// full history is skipped, keeping the first-run cost small even as the
		// template grows. `--` guards against a URL that begins with a dash being
		// parsed as an option.
		if _, err := gitOutput("", "clone", "--quiet", "--filter=blob:none", "--no-single-branch", "--", templateURL, tmplDir); err != nil {
			return fmt.Errorf("llcppg: cannot init: template %q is unavailable: %w", templateURL, err)
		}
		return verifyTemplateBranches(tmplDir)
	}

	// Existing clone: update it in place. If the fetch fails (e.g. offline),
	// fall back to the cached copy with a warning rather than failing. The
	// fetched origin/<branch> refs are what later steps read from, so a
	// successful fetch alone brings the cache to the remote state and any local
	// modification inside the cache is ignored.
	if _, err := gitOutput(tmplDir, "fetch", "--quiet", "--prune", "origin"); err != nil {
		fmt.Fprintf(stderr, "llcppg: warning: could not update the template cache (%v); using the cached copy\n", err)
	}
	return verifyTemplateBranches(tmplDir)
}

// templateRef returns the ref extractBranch reads for the given branch. It
// prefers the remote-tracking ref (origin/<branch>), which reflects the latest
// fetched remote state and ignores the cache's local checkout and any local
// modification; it falls back to the local branch when no remote-tracking ref
// exists.
func templateRef(tmplDir, branch string) (string, bool) {
	for _, ref := range []string{"refs/remotes/origin/" + branch, "refs/heads/" + branch} {
		if _, err := gitOutput(tmplDir, "rev-parse", "--verify", "--quiet", ref); err == nil {
			return ref, true
		}
	}
	return "", false
}

// verifyTemplateBranches confirms the cached template has both required
// branches, so a later extract fails early with a clear message.
func verifyTemplateBranches(tmplDir string) error {
	for _, branch := range []string{branchC, branchMain} {
		if _, ok := templateRef(tmplDir, branch); !ok {
			return fmt.Errorf("llcppg: template error: branch %q not found in the cached template", branch)
		}
	}
	return nil
}

// extractBranch writes the contents of the template's branch into destDir's
// working tree, without copying the template's .git directory. It streams a tar
// archive produced by `git archive` and unpacks it with a self-contained tar
// reader, so the cache's checked-out state never affects the result.
func extractBranch(tmplDir, branch, destDir string) error {
	ref, ok := templateRef(tmplDir, branch)
	if !ok {
		return fmt.Errorf("llcppg: template error: branch %q not found in the cached template", branch)
	}
	cmd := exec.Command("git", "-C", tmplDir, "archive", "--format=tar", ref)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := untar(stdout, destDir); err != nil {
		_ = cmd.Wait()
		return err
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("git archive %s: %v: %s", branch, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// -----------------------------------------------------------------------------

// checkoutNewBranch checks branch out in dir, creating it when it does not yet
// exist. On an unborn HEAD (a freshly cloned, empty repository) `git checkout
// -b` makes branch the first branch to receive a commit.
//
// When orphan is true and the branch does not exist, it is created with `git
// checkout --orphan`, i.e. as a fresh root with no parent commit, and the index
// and working tree are then cleared so none of the previous branch's files
// carry over; the caller repopulates them from the template. This keeps the c
// and main branches from ever sharing history or mixing content.
//
// When orphan is true and the branch already exists — the common case for main,
// which a freshly cloned repository already has as its default branch — init
// switches to it with a plain `git checkout` rather than failing to recreate
// it. The caller then lays the template files on top and commits.
func checkoutNewBranch(dir, branch string, orphan bool) error {
	if orphan {
		if localBranchExists(dir, branch) {
			if _, err := gitOutput(dir, "checkout", branch); err != nil {
				return fmt.Errorf("llcppg: cannot init: switch to branch %q: %w", branch, err)
			}
			return nil
		}
		if _, err := gitOutput(dir, "checkout", "--orphan", branch); err != nil {
			return fmt.Errorf("llcppg: cannot init: create branch %q: %w", branch, err)
		}
		// `git checkout --orphan` keeps the previous branch's files staged and
		// in the working tree; remove them so main starts from a clean slate.
		if _, err := gitOutput(dir, "rm", "-rf", "--quiet", "--ignore-unmatch", "."); err != nil {
			return fmt.Errorf("llcppg: cannot init: clear working tree for %q: %w", branch, err)
		}
		return nil
	}
	if _, err := gitOutput(dir, "checkout", "-b", branch); err != nil {
		return fmt.Errorf("llcppg: cannot init: create branch %q: %w", branch, err)
	}
	return nil
}

// commitAll stages every change (including deletions) in dir and commits it
// with the given message.
func commitAll(dir, message string) error {
	if _, err := gitOutput(dir, "add", "-A"); err != nil {
		return err
	}
	if _, err := gitOutput(dir, "commit", "--quiet", "-m", message); err != nil {
		return fmt.Errorf("llcppg: cannot init: commit %q: %w", message, err)
	}
	return nil
}

// -----------------------------------------------------------------------------

// gitOutput runs git with the given args. When dir is non-empty the command is
// run with `-C dir`. It returns git's raw (untrimmed) stdout, or an error that
// includes git's stderr. Callers that compare the output trim it themselves.
func gitOutput(dir string, args ...string) (string, error) {
	if dir != "" {
		args = append([]string{"-C", dir}, args...)
	}
	cmd := exec.Command("git", args...)
	// A deterministic, non-interactive environment: never prompt for
	// credentials and never let user hooks or config surprise the run.
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return stdout.String(), nil
}

// -----------------------------------------------------------------------------
