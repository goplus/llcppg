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

package docserver

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/doc/comment"
	"go/printer"
	"go/scanner"
	"go/token"
	"html"
	"html/template"
	"sort"
	"strings"
)

// pageData is the data passed to the HTML template for one package page.
type pageData struct {
	ImportPath  string
	Name        string
	Dir         string
	ParseErrors []string

	Overview template.HTML // rendered package comment, or "" when none

	Consts []*valueDoc
	Vars   []*valueDoc
	Funcs  []*funcDoc
	Types  []*typeDoc

	Imports []string

	Summary summary
}

// summary is the small diagnostic box: counts of exported declarations and how
// many exported declarations lack a doc comment.
type summary struct {
	Types        int
	Funcs        int
	Methods      int
	Consts       int
	Vars         int
	Undocumented int
}

// valueDoc renders a const or var block.
type valueDoc struct {
	Anchor string
	Doc    template.HTML
	Decl   template.HTML
}

// funcDoc renders a function or method.
type funcDoc struct {
	Anchor string
	Name   string
	Recv   string // receiver type name for methods, "" for free functions
	Doc    template.HTML
	Decl   template.HTML
}

// typeDoc renders a type with its grouped declarations, constructors and
// methods.
type typeDoc struct {
	Anchor  string
	Name    string
	Doc     template.HTML
	Decl    template.HTML
	Consts  []*valueDoc
	Vars    []*valueDoc
	Funcs   []*funcDoc // constructor-style functions returning the type
	Methods []*funcDoc
}

// renderer renders a loaded package into pageData. It holds the per-render
// state (file set, comment printer, and the set of top-level type names used
// for in-page linking) so the small helper methods stay cohesive.
type renderer struct {
	p         *pkg
	fset      *token.FileSet
	cparser   *comment.Parser
	cprinter  *comment.Printer
	typeNames map[string]bool // top-level type names in the package
}

// render builds the template data for p.
func (p *pkg) render() *pageData {
	dp := p.Doc

	typeNames := make(map[string]bool, len(dp.Types))
	for _, t := range dp.Types {
		typeNames[t.Name] = true
	}

	r := &renderer{
		p:         p,
		fset:      p.FileSet,
		cparser:   dp.Parser(),
		cprinter:  dp.Printer(),
		typeNames: typeNames,
	}
	// Resolve [Name] and [Type.Method] doc links to in-page anchors.
	r.cprinter.DocLinkURL = func(link *comment.DocLink) string {
		if link.ImportPath != "" {
			// We do not resolve imported packages; leave such links as plain
			// text by returning no URL.
			return ""
		}
		return "#" + symAnchor(link.Name, link.Recv)
	}

	data := &pageData{
		ImportPath:  dp.ImportPath,
		Name:        dp.Name,
		Dir:         p.Dir,
		ParseErrors: p.ParseErrors,
		Imports:     append([]string(nil), dp.Imports...),
	}
	sort.Strings(data.Imports)

	if strings.TrimSpace(dp.Doc) != "" {
		data.Overview = r.comment(dp.Doc)
	}

	for _, c := range dp.Consts {
		data.Consts = append(data.Consts, r.value(c, "const"))
	}
	for _, v := range dp.Vars {
		data.Vars = append(data.Vars, r.value(v, "var"))
	}
	for _, f := range dp.Funcs {
		data.Funcs = append(data.Funcs, r.fn(f, ""))
	}
	for _, t := range dp.Types {
		data.Types = append(data.Types, r.typ(t))
	}

	data.Summary = r.summarize()
	return data
}

func (r *renderer) typ(t *doc.Type) *typeDoc {
	td := &typeDoc{
		Anchor: symAnchor(t.Name, ""),
		Name:   t.Name,
		Doc:    r.comment(t.Doc),
		Decl:   r.decl(t.Decl),
	}
	for _, c := range t.Consts {
		td.Consts = append(td.Consts, r.value(c, "const"))
	}
	for _, v := range t.Vars {
		td.Vars = append(td.Vars, r.value(v, "var"))
	}
	for _, f := range t.Funcs {
		td.Funcs = append(td.Funcs, r.fn(f, ""))
	}
	for _, m := range t.Methods {
		td.Methods = append(td.Methods, r.fn(m, t.Name))
	}
	return td
}

func (r *renderer) value(v *doc.Value, kind string) *valueDoc {
	anchor := ""
	if len(v.Names) > 0 {
		anchor = symAnchor(v.Names[0], "")
	}
	return &valueDoc{
		Anchor: anchor,
		Doc:    r.comment(v.Doc),
		Decl:   r.decl(v.Decl),
	}
}

func (r *renderer) fn(f *doc.Func, recv string) *funcDoc {
	return &funcDoc{
		Anchor: symAnchor(f.Name, recv),
		Name:   f.Name,
		Recv:   recv,
		Doc:    r.comment(f.Doc),
		Decl:   r.decl(f.Decl),
	}
}

// comment renders a doc comment to HTML using the package's go/doc comment
// parser and printer, so headings, lists, code blocks and doc links follow the
// current Go doc comment syntax. Empty input renders to "".
func (r *renderer) comment(text string) template.HTML {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	parsed := r.cparser.Parse(text)
	return template.HTML(r.cprinter.HTML(parsed))
}

// decl prints a declaration node with go/printer, with its doc comment cleared
// so it is not shown twice, HTML-escapes the result, and wraps same-package type
// identifiers in anchor links. The returned value is safe HTML.
func (r *renderer) decl(node ast.Node) template.HTML {
	src := r.printDecl(node)
	return template.HTML(r.linkTypes(src))
}

// printDecl renders a declaration to Go source text, clearing any attached doc
// comment first so the rendered comment above the declaration is not duplicated.
func (r *renderer) printDecl(node ast.Node) string {
	// Clear the doc comment on a shallow copy so the original AST (shared with
	// go/doc) is left untouched across reloads.
	switch d := node.(type) {
	case *ast.GenDecl:
		cp := *d
		cp.Doc = nil
		node = &cp
	case *ast.FuncDecl:
		cp := *d
		cp.Doc = nil
		node = &cp
	}
	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 4}
	if err := cfg.Fprint(&buf, r.fset, node); err != nil {
		return ""
	}
	return buf.String()
}

// linkTypes HTML-escapes src and wraps every identifier that names a top-level
// type in the package in an anchor to that type. It uses go/scanner so it never
// needs to type-check, which would require resolving the dependencies we have
// deliberately decided not to resolve; this covers the common case of
// same-package type references in signatures.
func (r *renderer) linkTypes(src string) string {
	var sf token.FileSet
	file := sf.AddFile("", sf.Base(), len(src))
	var s scanner.Scanner
	s.Init(file, []byte(src), nil, 0)

	var b strings.Builder
	prev := 0 // offset in src up to which we have already written output
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.IDENT || !r.typeNames[lit] {
			continue
		}
		off := file.Offset(pos)
		// Emit the (escaped) gap before this identifier, then the linked
		// identifier itself.
		b.WriteString(html.EscapeString(src[prev:off]))
		b.WriteString(`<a href="#`)
		b.WriteString(symAnchor(lit, ""))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(lit))
		b.WriteString(`</a>`)
		prev = off + len(lit)
	}
	b.WriteString(html.EscapeString(src[prev:]))
	return b.String()
}

// summarize counts exported declarations and how many lack a doc comment, for
// the summary box. doc.AllDecls may add unexported entries; the summary only
// counts exported names so it stays a consumer-facing view.
func (r *renderer) summarize() summary {
	dp := r.p.Doc
	var s summary
	for _, c := range dp.Consts {
		s.Consts += len(exported(c.Names))
		countUndoc(&s, c.Doc, exported(c.Names))
	}
	for _, v := range dp.Vars {
		s.Vars += len(exported(v.Names))
		countUndoc(&s, v.Doc, exported(v.Names))
	}
	for _, f := range dp.Funcs {
		if !ast.IsExported(f.Name) {
			continue
		}
		s.Funcs++
		countUndoc(&s, f.Doc, []string{f.Name})
	}
	for _, t := range dp.Types {
		if !ast.IsExported(t.Name) {
			continue
		}
		s.Types++
		countUndoc(&s, t.Doc, []string{t.Name})
		for _, c := range t.Consts {
			s.Consts += len(exported(c.Names))
			countUndoc(&s, c.Doc, exported(c.Names))
		}
		for _, v := range t.Vars {
			s.Vars += len(exported(v.Names))
			countUndoc(&s, v.Doc, exported(v.Names))
		}
		for _, f := range t.Funcs {
			if !ast.IsExported(f.Name) {
				continue
			}
			s.Funcs++
			countUndoc(&s, f.Doc, []string{f.Name})
		}
		for _, m := range t.Methods {
			if !ast.IsExported(m.Name) {
				continue
			}
			s.Methods++
			countUndoc(&s, m.Doc, []string{m.Name})
		}
	}
	return s
}

// countUndoc adds the number of exported names in a declaration to the
// undocumented count when the declaration has no doc comment.
func countUndoc(s *summary, docText string, exportedNames []string) {
	if len(exportedNames) == 0 {
		return
	}
	if strings.TrimSpace(docText) == "" {
		s.Undocumented += len(exportedNames)
	}
}

// exported returns the exported subset of names.
func exported(names []string) []string {
	var out []string
	for _, n := range names {
		if ast.IsExported(n) {
			out = append(out, n)
		}
	}
	return out
}

// symAnchor returns the HTML id/anchor for a symbol. A method uses
// "Recv.Method"; everything else uses its own name. The names come from Go
// identifiers so they are already safe as anchors.
func symAnchor(name, recv string) string {
	if recv != "" {
		return recv + "." + name
	}
	return name
}
