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
	"go/token"
	"strconv"
	"strings"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/clang"
	lc "github.com/llarhub/clang-c"
)

// -----------------------------------------------------------------------------

func evalConstExpr(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token) (v any, ok bool) {
	defer func() {
		recover()
	}()
	kind, v, _, _, ok := parseExpr(ctx, tu, tokens, false)
	ok = ok && kind == fca_Const
	return
}

// -----------------------------------------------------------------------------

const (
	fca_Const = iota
	fca_Ident
	fca_Sizeof
)

type exprVal struct {
	kind int
	v    any
}

type funcCallExpr struct {
	fn   string
	args []exprVal
}

func (p *funcCallExpr) isMethodCall(recvName string) bool {
	if len(p.args) == 0 {
		return false
	}
	arg := p.args[0]
	return arg.kind == fca_Ident && arg.v.(string) == recvName
}

func evalFuncCallExpr(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token) (v funcCallExpr, ok bool) {
	v, left, ok := parseFuncCallExpr(ctx, tu, tokens)
	ok = ok && len(left) == 0
	return
}

func parseFuncCallExpr(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token) (v funcCallExpr, left []lc.Token, found bool) {
	tok, lit, left, ok := scanToken(tu, tokens)
	if !ok {
		return
	}
	switch tok {
	case token.IDENT:
		tok, _, left, ok = scanToken(tu, left)
		if !ok || tok != token.LPAREN {
			return
		}
		for !nextIsRParen(tu, left) {
			var kind int
			var arg any
			kind, arg, tok, left, ok = parseExpr(ctx, tu, left, true)
			if !ok {
				return
			}
			v.args = append(v.args, exprVal{kind, arg})
			if tok == token.RPAREN {
				break
			}
		}
		v.fn, found = lit, true
	case token.LPAREN:
		if v, left, ok = parseFuncCallExpr(ctx, tu, left); ok {
			tok, _, left, ok = scanToken(tu, left)
			found = ok && tok == token.RPAREN
		}
	}
	return
}

// -----------------------------------------------------------------------------

type operand struct {
	val  any
	tok  token.Token
	prec int
}

var opPrecs = map[token.Token]int{
	token.REM: 12,
	token.MUL: 12,
	token.QUO: 12,

	token.ADD: 10,
	token.SUB: 10,

	token.SHL: 8,
	token.SHR: 8,

	token.AND: 6,
	token.XOR: 4,
	token.OR:  2,
}

func parseExpr(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token, needRParen bool) (kind int, v any, tok token.Token, left []lc.Token, found bool) {
	kind, v, left, ok := parseOperand(ctx, tu, tokens)
	if !ok {
		return
	}
	if kind == fca_Ident {
		switch v.(string) {
		case "sizeof":
			tok, _, left, ok = scanToken(tu, left)
			if !ok || tok != token.LPAREN {
				return
			}
			tok, v, left, ok = scanToken(tu, left)
			if !ok || tok != token.IDENT {
				return
			}
			tok, _, left, ok = scanToken(tu, left)
			if !ok || tok != token.RPAREN {
				return
			}
			kind = fca_Sizeof
		}
		if needRParen {
			tok, _, left, ok = scanToken(tu, left)
			found = ok && (tok == token.RPAREN || tok == token.COMMA)
		} else {
			tok, found = token.ILLEGAL, len(left) == 0
		}
		return
	}
	var n, prec int
	var ops = []operand{{val: v, tok: token.ILLEGAL, prec: -1}}
	for len(left) > 0 {
		tok, _, left, ok = scanToken(tu, left)
		if !ok {
			return
		}
		prec, ok = opPrecs[tok]
		if !ok {
			if needRParen && (tok == token.RPAREN || tok == token.COMMA) {
				_, found = calc(ops, n, 0)
				v = ops[0].val
			}
			return
		}
		n, ok = calc(ops, n, prec)
		if !ok {
			return
		}
		kind, v, left, ok = parseOperand(ctx, tu, left)
		if !ok || kind != fca_Const {
			return
		}
		n++
		ops = append(ops[:n], operand{val: v, tok: tok, prec: prec})
	}
	tok = token.ILLEGAL
	if ok = !needRParen; ok {
		_, found = calc(ops, n, 0)
		v = ops[0].val
	}
	return
}

func calc(ops []operand, nlast, prec int) (n int, ok bool) {
	n, ok = nlast, true
	for ops[n].prec >= prec {
		switch op := ops[n].tok; op {
		case token.ADD, token.SUB, token.MUL, token.QUO:
			ops[n-1].val, ok = mathOp(op, ops[n-1].val, ops[n].val)
			if !ok {
				return
			}
		default:
			a, ok1 := ops[n-1].val.(int)
			b, ok2 := ops[n].val.(int)
			if ok = ok1 && ok2; !ok {
				return
			}
			switch op {
			case token.SHL:
				a <<= b
			case token.SHR:
				a >>= b
			case token.AND:
				a &= b
			case token.OR:
				a |= b
			case token.XOR:
				a ^= b
			case token.REM:
				if b == 0 {
					// integer divide by zero
					return 0, false
				}
				a %= b
			default:
				panic("calc: unknown op")
			}
			ops[n-1].val = a
		}
		n--
	}
	return
}

func mathOp(op token.Token, a, b any) (any, bool) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			switch op {
			case token.ADD:
				return a + b, true
			case token.SUB:
				return a - b, true
			case token.MUL:
				return a * b, true
			case token.QUO:
				if b == 0 {
					// integer divide by zero
					return nil, false
				}
				return a / b, true
			}
		case float64:
			return floatMathOp(op, float64(a), b), true
		}
	case float64:
		switch b := b.(type) {
		case int:
			return floatMathOp(op, a, float64(b)), true
		case float64:
			return floatMathOp(op, a, b), true
		}
	}
	return nil, false
}

func floatMathOp(op token.Token, a, b float64) float64 {
	switch op {
	case token.ADD:
		return a + b
	case token.SUB:
		return a - b
	case token.MUL:
		return a * b
	case token.QUO:
		return a / b
	}
	panic("floatMathOp: unknown op")
}

func parseOperand(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token) (kind int, v any, left []lc.Token, ok bool) {
	tok, lit, left, ok := scanToken(tu, tokens)
	if !ok {
		return
	}
	switch tok {
	case token.FLOAT:
		if strings.IndexByte(lit, '.') >= 0 {
			val, e := strconv.ParseFloat(lit, 64)
			v, ok = val, e == nil
		} else {
			val, e := strconv.ParseInt(lit, 0, 64)
			v, ok = int(val), e == nil
		}
	case token.LPAREN:
		kind, v, _, left, ok = parseExpr(ctx, tu, left, true)
	case token.IDENT:
		if v, ok = ctx.macroVals[lit]; !ok {
			kind, v, ok = fca_Ident, lit, true
		}
	case token.STRING:
		val, e := strconv.Unquote(lit)
		v, ok = val, e == nil
	case token.CHAR:
		ok = false // not supported
	case token.SUB: // -
		_, v, left, ok = parseOperand(ctx, tu, left)
		if !ok {
			return
		}
		switch x := v.(type) {
		case int:
			v = -x
		case float64:
			v = -x
		default:
			ok = false
		}
	case token.TILDE: // ~
		_, v, left, ok = parseOperand(ctx, tu, left)
		if !ok {
			return
		}
		switch x := v.(type) {
		case int:
			v = ^x
		default:
			ok = false
		}
	}
	return
}

// -----------------------------------------------------------------------------

var c2goOps = map[string]token.Token{
	"(": token.LPAREN,
	")": token.RPAREN,
	",": token.COMMA,

	"%": token.REM,
	"*": token.MUL,
	"/": token.QUO,

	"+": token.ADD,
	"-": token.SUB,

	"<<": token.SHL,
	">>": token.SHR,

	"&": token.AND,
	"~": token.TILDE,
	"^": token.XOR,
	"|": token.OR,
}

func scanToken(tu clang.TranslationUnit, tokens []lc.Token) (ret token.Token, lit string, left []lc.Token, ok bool) {
	for len(tokens) > 0 {
		tok := tokens[0]
		kind := tok.Kind()
		switch kind {
		case lc.Token_Punctuation:
			op := tu.TokenSpelling(tok)
			if op == "\\" {
				tokens = tokens[1:]
				continue
			}
			left = tokens[1:]
			ret, ok = c2goOps[op]
		case lc.Token_Literal:
			lit, ok = tu.TokenSpelling(tok), true
			switch lit[0] {
			case '"':
				ret = token.STRING
			case '\'':
				ret = token.CHAR
			default:
				ret = token.FLOAT
			}
			left = tokens[1:]
		case lc.Token_Identifier, lc.Token_Keyword:
			ret = token.IDENT
			lit, ok = tu.TokenSpelling(tok), true
			left = tokens[1:]
		case lc.Token_Comment:
			tokens = tokens[1:]
			continue
		}
		break
	}
	return
}

func nextIsRParen(tu clang.TranslationUnit, tokens []lc.Token) bool {
	for len(tokens) > 0 {
		tok := tokens[0]
		kind := tok.Kind()
		switch kind {
		case lc.Token_Punctuation:
			op := tu.TranslationUnit.TokenSpelling(tok).CStr()
			if c.Strcmp(op, c.Str(")")) == 0 {
				return true
			}
		case lc.Token_Comment:
			tokens = tokens[1:]
			continue
		}
		break
	}
	return false
}

// -----------------------------------------------------------------------------
