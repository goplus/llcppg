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

package cl

import (
	"go/ast"
	"strings"

	"github.com/goplus/llcppg/clang"
)

// -----------------------------------------------------------------------------

// docComments returns the doc comment lines of the given declaration, converted
// into Go "//" line comments, or nil when keepDoc is off or the declaration has
// no associated comment.
//
// The C/C++ comment markers ("/** ... */", "///", "//!", etc.) are stripped and
// the content is re-emitted as a block of "//" line comments so it renders as a
// leading Go doc comment. Callers place these comments before any generated
// directive (e.g. //go:linkname), which must stay adjacent to its declaration.
func (p *pkgCtx) docComments(decl clang.Cursor) []*ast.Comment {
	if !p.keepDoc {
		return nil
	}
	raw := clang.RawComment(decl)
	if raw == "" {
		return nil
	}
	return toLineComments(raw)
}

// toLineComments converts a raw C/C++ documentation comment into a slice of Go
// "//" line comments. It handles both block comments ("/* ... */", "/** ... */")
// and consecutive line comments ("//", "///", "//!"), stripping the markers and
// any common leading " * " decoration from block comments. It returns nil when
// there is no meaningful content.
func toLineComments(raw string) []*ast.Comment {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var lines []string
	switch {
	case strings.HasPrefix(raw, "/*"):
		body := strings.TrimPrefix(raw, "/*")
		body = strings.TrimSuffix(body, "*/")
		for _, line := range strings.Split(body, "\n") {
			line = strings.TrimSpace(line)
			// Drop a Doxygen/Javadoc-style leading "*" decoration.
			if line == "*" {
				line = ""
			} else if s, ok := strings.CutPrefix(line, "* "); ok {
				line = s
			} else if s, ok := strings.CutPrefix(line, "*"); ok {
				line = s
			}
			lines = append(lines, line)
		}
	default:
		for _, line := range strings.Split(raw, "\n") {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "//")
			// Doxygen line-comment markers: "///" and "//!".
			line = strings.TrimPrefix(line, "/")
			line = strings.TrimPrefix(line, "!")
			line = strings.TrimPrefix(line, " ")
			lines = append(lines, line)
		}
	}

	// Trim leading/trailing blank lines that come from the markers being on
	// their own line (e.g. "/**\n * foo\n */").
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil
	}

	comments := make([]*ast.Comment, len(lines))
	for i, line := range lines {
		text := "//"
		if line != "" {
			text += " " + line
		}
		comments[i] = &ast.Comment{Text: text}
	}
	return comments
}

// docCommentGroup returns a *ast.CommentGroup holding the declaration's doc
// comment, or nil when there is nothing to attach. It is used for declarations
// whose only leading comment is the doc (types, consts), where no generated
// directive needs to be preserved alongside it.
func (p *pkgCtx) docCommentGroup(decl clang.Cursor) *ast.CommentGroup {
	comments := p.docComments(decl)
	if len(comments) == 0 {
		return nil
	}
	comments[0].Text = "\n" + comments[0].Text
	return &ast.CommentGroup{List: comments}
}

// directiveComments builds the comment group for a declaration that carries a
// generated directive (e.g. "//go:linkname" or "// llgo:link"): the C/C++ doc
// comment (when keepDoc is on) followed by the directive. The directive is
// always emitted last so it stays adjacent to the declaration it applies to.
//
// directive is expected to start with a leading "\n" (matching the original
// directive-only form, which relied on it to leave a blank line before the
// declaration). When a doc comment is present, the "\n" is moved to the front
// of the whole group so the blank line separates the doc block from the
// preceding code, and the directive follows the doc with no blank line.
func (p *pkgCtx) directiveComments(decl clang.Cursor, directive string) *ast.CommentGroup {
	doc := p.docComments(decl)
	if len(doc) == 0 {
		return &ast.CommentGroup{List: []*ast.Comment{{Text: directive}}}
	}
	list := make([]*ast.Comment, 0, len(doc)+2)
	list = append(list, &ast.Comment{Text: "\n" + doc[0].Text})
	list = append(list, doc[1:]...)
	list = append(list, &ast.Comment{Text: "//"}, &ast.Comment{Text: strings.TrimPrefix(directive, "\n")})
	return &ast.CommentGroup{List: list}
}

// -----------------------------------------------------------------------------

func (p *pkgCtx) directiveTypeC(decl clang.Cursor, hasCallback bool) *ast.CommentGroup {
	if hasCallback {
		return p.directiveComments(decl, "\n// llgo:type C")
	}
	return p.docCommentGroup(decl)
}

// -----------------------------------------------------------------------------
