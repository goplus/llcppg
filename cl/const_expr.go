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

	"github.com/goplus/llcppg/clang"
	lc "github.com/goplus/llcppg/lib/clang"
)

// -----------------------------------------------------------------------------

func evalConstExpr(ctx *pkgCtx, tokens []lc.Token) (v any, ok bool) {
	v, _, ok = parseExpr(ctx, ctx.tu, tokens, false)
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

func parseExpr(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token, needRParen bool) (v any, left []lc.Token, ok bool) {
	v, left, ok = parseOperand(ctx, tu, tokens)
	if !ok {
		return
	}
	var tok token.Token
	var n, prec int
	var ops = []operand{{val: v, tok: token.ILLEGAL, prec: -1}}
	for len(left) > 0 {
		tok, _, left, ok = scanToken(tu, left)
		if !ok {
			return
		}
		prec, ok = opPrecs[tok]
		if !ok {
			if tok == token.RPAREN && needRParen {
				_, ok = calc(ops, n, 0)
				v = ops[0].val
			}
			return
		}
		n, ok = calc(ops, n, prec)
		if !ok {
			return
		}
		v, left, ok = parseOperand(ctx, tu, left)
		if !ok {
			return
		}
		n++
		ops = append(ops[:n], operand{val: v, tok: tok, prec: prec})
	}
	if ok = !needRParen; ok {
		_, ok = calc(ops, n, 0)
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
				a %= b
			default:
				panic("parseExpr: unknown op")
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

func parseOperand(ctx *pkgCtx, tu clang.TranslationUnit, tokens []lc.Token) (v any, left []lc.Token, ok bool) {
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
		return parseExpr(ctx, tu, left, true)
	case token.IDENT:
		v, ok = ctx.macroVals[lit]
	case token.CHAR, token.STRING:
		ok = false // not supported
	case token.SUB: // -
		v, left, ok = parseOperand(ctx, tu, left)
		if !ok {
			return
		}
		switch v := v.(type) {
		case int:
			v = -v
		case float64:
			v = -v
		default:
			ok = false
		}
	case token.XOR: // ~
		v, left, ok = parseOperand(ctx, tu, left)
		if !ok {
			return
		}
		switch v := v.(type) {
		case int:
			v = ^v
		default:
			ok = false
		}
	}
	return
}

// -----------------------------------------------------------------------------

var goOps = map[string]token.Token{
	"(": token.LPAREN,
	")": token.RPAREN,

	"%": token.REM,
	"*": token.MUL,
	"/": token.QUO,

	"+": token.ADD,
	"-": token.SUB,

	"<<": token.SHL,
	">>": token.SHR,

	"&": token.AND,
	"^": token.XOR,
	"|": token.OR,
}

func scanToken(tu clang.TranslationUnit, tokens []lc.Token) (ret token.Token, lit string, left []lc.Token, ok bool) {
	for len(tokens) > 0 {
		tok := tokens[0]
		kind := tok.Kind()
		switch kind {
		case lc.Punctuation:
			op := tu.Token(tok)
			left = tokens[1:]
			ret, ok = goOps[op]
		case lc.Literal:
			lit, ok = tu.Token(tok), true
			switch lit[0] {
			case '"':
				ret = token.STRING
			case '\'':
				ret = token.CHAR
			default:
				ret = token.FLOAT
			}
			left = tokens[1:]
		case lc.Identifier:
			ret = token.IDENT
			lit, ok = tu.Token(tok), true
			left = tokens[1:]
		case lc.Comment:
			tokens = tokens[1:]
			continue
		}
		break
	}
	return
}

// -----------------------------------------------------------------------------
