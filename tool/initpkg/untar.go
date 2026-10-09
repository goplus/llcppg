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
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

// untar extracts the tar stream r into destDir. Only regular files and
// directories are written; the archive produced by `git archive` for a tree
// contains exactly those. Entries are confined to destDir to guard against a
// malformed path escaping the destination.
func untar(r io.Reader, destDir string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target, err := safeJoin(destDir, hdr.Name)
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := writeFile(target, tr, os.FileMode(hdr.Mode).Perm()); err != nil {
				return err
			}
		default:
			// Symlinks and other special entries are not part of a plain tree
			// archive; skip anything unexpected rather than acting on it.
		}
	}
}

func writeFile(target string, r io.Reader, mode os.FileMode) error {
	if mode == 0 {
		mode = 0644
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// safeJoin joins name onto base and verifies the result stays within base,
// rejecting absolute paths and ".." traversal.
func safeJoin(base, name string) (string, error) {
	clean := filepath.Clean("/" + strings.ReplaceAll(name, "\\", "/"))
	target := filepath.Join(base, clean)
	rel, err := filepath.Rel(base, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("llcppg: template error: unsafe path %q in template archive", name)
	}
	return target, nil
}

// -----------------------------------------------------------------------------
