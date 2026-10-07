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
	"strings"
)

// -----------------------------------------------------------------------------

type rule struct {
	prefix  string
	pattern []string
	result  string
}

func newRule(name, expr string) (*rule, error) {
	pos := strings.IndexByte(expr, ':')
	if pos < 0 {
		return nil, fmt.Errorf("invalid %s: %s", name, expr)
	}
	pattern := strings.TrimSpace(expr[:pos])
	parts := strings.Split(pattern, "*")
	result := strings.TrimSpace(expr[pos+1:])
	return &rule{
		prefix:  parts[0],
		pattern: parts[1:],
		result:  result,
	}, nil
}

func (p *rule) match(source string, matchFull bool) (ret string, matched bool) {
	source, ok := strings.CutPrefix(source, p.prefix)
	if !ok {
		return
	}
	n := 0
	np := len(p.pattern)
	match := make([]string, np)
	for i, p := range p.pattern {
		pos := strings.Index(source, p)
		if pos < 0 {
			return
		}
		if p == "" && i == np-1 {
			pos = len(source) // the last * matches the rest of the source
		}
		n += pos
		match[i] = source[:pos]
		source = source[pos+len(p):]
	}
	if matchFull && source != "" {
		return
	}
	return matchResult(p.result, match, n), true
}

func matchResult(result string, match []string, n int) string {
	b := make([]byte, 0, len(result)+n)
	for i := 0; i < len(result); i++ {
		if result[i] == '$' {
			if i+1 < len(result) {
				i++
				c := result[i]
				if c >= '1' && c <= '9' {
					if index := int(c - '1'); index < len(match) {
						b = append(b, match[index]...)
						continue
					}
				}
				b = append(b, '$', c)
				continue
			}
		}
		b = append(b, result[i])
	}
	return string(b)
}

// -----------------------------------------------------------------------------

type matcher []*rule

func newMatcher(name string, exprs []string) (matcher, error) {
	ret := make(matcher, len(exprs))
	for i, c := range exprs {
		m, err := newRule(name, c)
		if err != nil {
			return nil, err
		}
		ret[i] = m
	}
	return ret, nil
}

func (p matcher) match(source string, matchFull bool) (ret string, matched bool) {
	for _, m := range p {
		if r, ok := m.match(source, matchFull); ok {
			return r, true
		}
	}
	return "", false
}

func (p matcher) matchList(source string, matchFull bool) (ret []string, matched bool) {
	for _, m := range p {
		if r, ok := m.match(source, matchFull); ok {
			return parseResults(r), true
		}
	}
	return nil, false
}

func parseResults(result string) []string {
	items := strings.Split(result, ",")
	ret := items[:0]
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v != "" {
			ret = append(ret, v)
		}
	}
	return ret
}

// -----------------------------------------------------------------------------
