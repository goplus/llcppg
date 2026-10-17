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
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"math/big"
	"strings"

	"github.com/goplus/gogen"
)

// -----------------------------------------------------------------------------

// Define describes a user-specified constant to inject into the generated Go
// package. Name is a Go identifier and Value is the raw, untrimmed right-hand
// side (everything after the first '=' of a "NAME=VALUE" entry). Source is a
// human-readable origin (e.g. "-d" or "llcppg.cfg") used in error messages.
//
// Group identifies the single source list an entry came from (one -d argument,
// or one Define config list). Duplicate names are an error only within the same
// Group; across groups a later entry overrides an earlier one.
type Define struct {
	Name   string
	Value  string
	Source string
	Group  int
}

// ParseDefineArgs parses the accumulated -d arguments into Define entries, in
// the order the flags were given. Entries within one argument are separated by
// ';'. Each -d argument is its own group, so two entries with the same NAME in
// one -d argument are a duplicate error, while a later -d may override an
// earlier one.
func ParseDefineArgs(args []string) (defs []Define, err error) {
	for i, arg := range args {
		entries, err := parseDefineArg(arg, i)
		if err != nil {
			return nil, err
		}
		defs = append(defs, entries...)
	}
	return
}

// parseDefineArg parses a single -d argument (group) into Define entries.
// Entries have the form "NAME=VALUE", are separated by ';', and are split at
// the first '='; leading/trailing whitespace around NAME and around the whole
// VALUE is trimmed. An entry with no '=' or whose NAME is not a valid Go
// identifier is an error that names the offending entry.
func parseDefineArg(arg string, group int) (defs []Define, err error) {
	for _, entry := range strings.Split(arg, ";") {
		if strings.TrimSpace(entry) == "" {
			continue // tolerate empty entries from trailing/duplicate ';'
		}
		def, err := parseDefineEntry(entry, "-d", group)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	return
}

func parseDefineEntry(entry, source string, group int) (def Define, err error) {
	pos := strings.IndexByte(entry, '=')
	if pos < 0 {
		return Define{}, fmt.Errorf("%s: invalid define %q: missing '='", source, strings.TrimSpace(entry))
	}
	name := strings.TrimSpace(entry[:pos])
	value := strings.TrimSpace(entry[pos+1:])
	if !token.IsIdentifier(name) || token.IsKeyword(name) {
		return Define{}, fmt.Errorf("%s: invalid define %q: %q is not a valid Go identifier", source, strings.TrimSpace(entry), name)
	}
	return Define{Name: name, Value: value, Source: source, Group: group}, nil
}

// -----------------------------------------------------------------------------

// applyDefines injects the configured Define constants into the generated
// package. It runs after the load phase, when every constant that can be
// referenced is known, and emits the new constants together in a single const
// block in the package's default (main) generated file, in the order the
// entries were given. Any error increments the context error count.
func (p *pkgCtx) applyDefines(defs []Define) {
	if len(defs) == 0 {
		return
	}

	// Within one source list (same Group) a repeated NAME is a typo and an
	// error. Across groups, a later entry overrides an earlier one: it keeps
	// the first occurrence's slot in declaration order but takes the later
	// value.
	type groupKey struct {
		name  string
		group int
	}
	inGroup := make(map[groupKey]bool, len(defs))
	byName := make(map[string]*Define, len(defs))
	order := make([]*Define, 0, len(defs))
	for i := range defs {
		d := &defs[i]
		key := groupKey{d.Name, d.Group}
		if inGroup[key] {
			p.defineErrorf(d, "%q is defined twice in the same %s source", d.Name, d.Source)
			return
		}
		inGroup[key] = true
		if _, ok := byName[d.Name]; !ok {
			order = append(order, d) // first occurrence fixes the slot
		}
		byName[d.Name] = d // later occurrence wins the value
	}

	pkg := p.pkg
	scope := pkg.Types.Scope()

	// A Define must not collide with a symbol the conversion already produced.
	for _, d := range order {
		if obj := scope.Lookup(byName[d.Name].Name); obj != nil {
			p.defineErrorf(byName[d.Name], "%q conflicts with an existing symbol in the package", d.Name)
			return
		}
	}

	// Resolve every entry's value, following $(name) references on demand.
	r := &defineResolver{ctx: p, scope: scope, byName: byName, resolved: make(map[string]constant.Value), inProgress: make(map[string]bool)}
	for _, d := range order {
		if _, ok := r.resolve(byName[d.Name]); !ok {
			return // error already reported
		}
	}

	// Emit the resolved constants as a single const block in the default file,
	// in declaration order.
	old := pkg.CurFile()
	defer pkg.RestoreCurFile(old)
	if p.defaultGoFile != "" {
		pkg.SetCurFile(p.defaultGoFile, true)
	}

	defs2 := pkg.NewConstDefs(scope)
	for _, d := range order {
		v := r.resolved[d.Name]
		name := byName[d.Name].Name
		defs2.New(func(cb *gogen.CodeBuilder) int {
			cb.Val(goValueOf(v))
			return 1
		}, 0, token.NoPos, nil, name)
	}
}

func (p *pkgCtx) defineErrorf(d *Define, format string, args ...any) {
	p.errCnt++
	p.checkFailFast()
	// Define entries have no clang node/position; report against the source.
	log.Printf("%s: define error: %s", d.Source, fmt.Sprintf(format, args...))
}

// -----------------------------------------------------------------------------

// defineResolver resolves Define values, expanding $(name) references against
// the generated package scope and other Define entries. Results are cached as
// go/constant.Value so a referenced Define is only resolved once.
type defineResolver struct {
	ctx        *pkgCtx
	scope      *types.Scope
	byName     map[string]*Define        // Define name => entry
	resolved   map[string]constant.Value // Define name => resolved value
	inProgress map[string]bool           // Define name => currently resolving (cycle guard)
}

func (r *defineResolver) resolve(d *Define) (val constant.Value, ok bool) {
	if v, done := r.resolved[d.Name]; done {
		return v, true
	}
	if r.inProgress[d.Name] {
		r.ctx.defineErrorf(d, "%q is part of a reference cycle", d.Name)
		return nil, false
	}
	r.inProgress[d.Name] = true
	val, ok = r.evalValue(d, d.Value)
	delete(r.inProgress, d.Name)
	if ok {
		r.resolved[d.Name] = val
	}
	return
}

// evalValue evaluates a define's raw value into a go/constant.Value, applying
// the typing and reference rules from the proposal.
func (r *defineResolver) evalValue(d *Define, raw string) (constant.Value, bool) {
	// Backtick form: always a string, references still expanded. The backticks
	// themselves are not part of the string.
	if strings.HasPrefix(raw, "`") {
		if len(raw) < 2 || !strings.HasSuffix(raw, "`") {
			r.ctx.defineErrorf(d, "%q has an unterminated backtick", d.Name)
			return nil, false
		}
		inner := raw[1 : len(raw)-1]
		s, ok := r.interpolate(d, inner)
		if !ok {
			return nil, false
		}
		return constant.MakeString(s), true
	}

	// A value that is exactly one reference inherits the referenced constant's
	// value and kind.
	if name, ok := soleRef(raw); ok {
		return r.lookupRef(d, name)
	}

	// Any value that still contains a reference interpolates to a string.
	// Otherwise the kind is decided from the literal syntax. In both cases we
	// run interpolate first so $$ -> $ is unescaped and unterminated $( is
	// reported.
	s, ok := r.interpolate(d, raw)
	if !ok {
		return nil, false
	}
	if strings.Contains(raw, "$(") {
		return constant.MakeString(s), true
	}
	return literalValue(s), true
}

// lookupRef resolves a $(name) reference to a go/constant.Value, consulting
// other Define entries first (resolved on demand) and then the package scope.
func (r *defineResolver) lookupRef(d *Define, name string) (constant.Value, bool) {
	if ref, ok := r.byName[name]; ok {
		return r.resolve(ref)
	}
	obj := r.scope.Lookup(name)
	if obj == nil {
		r.ctx.defineErrorf(d, "%q references undefined name %q", d.Name, name)
		return nil, false
	}
	c, ok := obj.(*types.Const)
	if !ok {
		r.ctx.defineErrorf(d, "%q references %q which is not a constant", d.Name, name)
		return nil, false
	}
	return c.Val(), true
}

// interpolate expands references and $-escapes in s, rendering referenced
// numeric constants in their decimal form, and returns the resulting string.
func (r *defineResolver) interpolate(d *Define, s string) (string, bool) {
	var b strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c != '$' {
			b.WriteByte(c)
			i++
			continue
		}
		// c == '$'
		if i+1 < len(s) && s[i+1] == '$' { // $$ -> literal $
			b.WriteByte('$')
			i += 2
			continue
		}
		if i+1 < len(s) && s[i+1] == '(' { // $(name)
			end := strings.IndexByte(s[i+2:], ')')
			if end < 0 {
				r.ctx.defineErrorf(d, "%q has an unterminated $(", d.Name)
				return "", false
			}
			name := s[i+2 : i+2+end]
			v, ok := r.lookupRef(d, name)
			if !ok {
				return "", false
			}
			b.WriteString(constantToString(v))
			i += 2 + end + 1
			continue
		}
		// A $ not followed by ( or $ is kept as is.
		b.WriteByte('$')
		i++
	}
	return b.String(), true
}

// -----------------------------------------------------------------------------

// soleRef reports whether raw is exactly a single $(name) reference with no
// other text, returning the referenced name. A value like "$(A)$(B)" or
// "x$(A)" is not a sole reference.
func soleRef(raw string) (name string, ok bool) {
	if !strings.HasPrefix(raw, "$(") || !strings.HasSuffix(raw, ")") {
		return "", false
	}
	inner := raw[2 : len(raw)-1]
	// Reject anything with a nested reference boundary, so only one $( ... )
	// spanning the whole value qualifies.
	if strings.Contains(inner, "$(") || strings.Contains(inner, ")") {
		return "", false
	}
	return inner, true
}

// literalValue decides the kind of a reference-free, non-backtick value from
// Go's integer and floating-point literal syntax: an integer literal yields an
// untyped integer constant, a floating-point literal an untyped float, and
// anything else an untyped string. Numeric values are kept exact (e.g. 3.14
// stays 3.14, not its binary64 approximation).
func literalValue(s string) constant.Value {
	if v, ok := parseGoNumber(s); ok {
		return v
	}
	return constant.MakeString(s)
}

// parseGoNumber parses s as a Go integer or floating-point literal (with an
// optional leading sign), preserving the exact value. It returns ok=false for
// anything that is not a numeric literal.
func parseGoNumber(s string) (constant.Value, bool) {
	neg := false
	body := s
	if len(body) > 0 && (body[0] == '+' || body[0] == '-') {
		neg = body[0] == '-'
		body = body[1:]
	}
	if body == "" {
		return nil, false
	}
	tok := token.INT
	if looksFloat(body) {
		tok = token.FLOAT
	}
	v := constant.MakeFromLiteral(body, tok, 0)
	if v.Kind() == constant.Unknown {
		return nil, false
	}
	if neg {
		v = constant.UnaryOp(token.SUB, v, 0)
	}
	return v, true
}

// looksFloat reports whether s has the shape of a Go floating-point literal
// (contains a '.', an exponent, or is a hex float). This keeps values like
// "1e9" and "3.14" as floats while never misclassifying a plain integer.
func looksFloat(s string) bool {
	ls := strings.ToLower(s)
	if strings.HasPrefix(ls, "0x") || strings.HasPrefix(ls, "-0x") || strings.HasPrefix(ls, "+0x") {
		return strings.ContainsAny(ls, "pP") || strings.Contains(ls, ".")
	}
	return strings.ContainsAny(ls, ".eE")
}

// -----------------------------------------------------------------------------

// constantToString renders a go/constant.Value for interpolation: numeric
// constants in their decimal form, strings as their raw content.
func constantToString(v constant.Value) string {
	switch v.Kind() {
	case constant.String:
		return constant.StringVal(v)
	default:
		return v.ExactString()
	}
}

// goValueOf converts a resolved go/constant.Value into a native Go value
// suitable for gogen's CodeBuilder.Val, producing an untyped constant whose
// kind follows the value (int, float64, or string).
func goValueOf(v constant.Value) any {
	switch v.Kind() {
	case constant.Int:
		if iv, ok := constant.Int64Val(v); ok {
			return int(iv)
		}
		// Too large for int64: keep it an integer constant via *big.Int.
		if bi, ok := constant.Val(v).(*big.Int); ok {
			return bi
		}
		return v.ExactString()
	case constant.Float:
		fv, _ := constant.Float64Val(v)
		return fv
	case constant.Bool:
		return constant.BoolVal(v)
	default: // String and anything else
		return constant.StringVal(v)
	}
}

// -----------------------------------------------------------------------------
