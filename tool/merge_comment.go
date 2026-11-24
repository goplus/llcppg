package tool

import (
	"go/ast"
	"go/doc/comment"
	"reflect"
	"slices"
	"strings"
	"unicode"
)

type commentgroup struct {
	doc        []*ast.Comment
	line       *ast.CommentGroup
	directives []string
}

func (g *commentgroup) load(doc, lineComment *ast.CommentGroup) {
	*g = commentgroup{line: lineComment}
	if doc == nil {
		return
	}
	for _, comment := range doc.List {
		// Generated comments can start with a newline, as in "\n//go:linkname …".
		text := strings.TrimLeft(comment.Text, "\r\n")
		if !strings.HasPrefix(text, "//") {
			g.doc = append(g.doc, comment)
			continue
		}
		for line := range strings.SplitSeq(text, "\n") {
			if directive := commentDirective(line); directive != "" {
				if directive != "go:build" && !strings.HasPrefix(directive, "go:build ") {
					g.directives = append(g.directives, directive)
				}
				continue
			}
			if line == text {
				g.doc = append(g.doc, comment)
			} else {
				g.doc = append(g.doc, &ast.Comment{Slash: comment.Slash, Text: line})
			}
		}
	}
	// A blank // before //go:linkname separates the documentation from the directive.
	for len(g.doc) > 0 && strings.TrimSpace(g.doc[0].Text) == "//" {
		g.doc = g.doc[1:]
	}
	for len(g.doc) > 0 && strings.TrimSpace(g.doc[len(g.doc)-1].Text) == "//" {
		g.doc = g.doc[:len(g.doc)-1]
	}
}

func (g commentgroup) equalDirectives(other commentgroup) bool {
	return slices.Equal(g.directives, other.directives)
}

func (g commentgroup) equalDoc(other commentgroup) bool {
	return slices.EqualFunc(g.doc, other.doc, func(a, b *ast.Comment) bool {
		return strings.TrimLeft(a.Text, "\r\n") == strings.TrimLeft(b.Text, "\r\n")
	})
}

func equalLineComments(a, b *ast.CommentGroup) bool {
	if a == nil || b == nil {
		return a == b
	}
	return slices.EqualFunc(a.List, b.List, func(a, b *ast.Comment) bool {
		return a.Text == b.Text
	})
}

func (g *commentgroup) intersect(other commentgroup) {
	if g.line != nil && !equalLineComments(g.line, other.line) {
		g.line = nil
	}
	if len(g.doc) == 0 || g.equalDoc(other) {
		return
	}
	left, right := g.paragraphs(), other.paragraphs()
	var common int
	// A matching tail can depend on different introductions, as in [A, X] and [B, X].
	for common < min(len(left), len(right)) && reflect.DeepEqual(left[common], right[common]) {
		common++
	}
	if common == 0 {
		g.doc = nil
		return
	}
	if common == len(left) {
		return
	}

	var printer comment.Printer
	text := strings.TrimSuffix(string(printer.Comment(&comment.Doc{Content: left[:common]})), "\n")
	g.doc = nil
	for line := range strings.SplitSeq(text, "\n") {
		if line != "" {
			line = " " + line
		}
		g.doc = append(g.doc, &ast.Comment{Text: "//" + line})
	}
}

func (g commentgroup) paragraphs() []comment.Block {
	group := ast.CommentGroup{List: slices.Clone(g.doc)}
	for i, node := range group.List {
		text := strings.TrimLeft(node.Text, "\r\n")
		if !strings.HasPrefix(text, "//") && !(len(text) >= 4 && strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/")) {
			return nil
		}
		if text != node.Text {
			group.List[i] = &ast.Comment{Text: text}
		}
	}
	var parser comment.Parser
	doc := parser.Parse(group.Text())
	// Keep a link definition with its document, such as [Flags]: https://example.org/flags.
	if len(doc.Links) != 0 {
		return nil
	}
	for _, block := range doc.Content {
		if _, ok := block.(*comment.Paragraph); !ok {
			return nil
		}
	}
	return doc.Content
}

func (g commentgroup) withDirectives(original *ast.CommentGroup) *ast.CommentGroup {
	// Anchor shortened documentation at the tail so //go:linkname still abuts func F.
	anchor := original.List[len(original.List)-1].Slash
	list := make([]*ast.Comment, 0, len(g.doc)+len(g.directives)+1)
	for _, node := range g.doc {
		list = append(list, &ast.Comment{Slash: anchor, Text: strings.TrimLeft(node.Text, "\r\n")})
	}
	for _, node := range original.List {
		text := strings.TrimLeft(node.Text, "\r\n")
		if !strings.HasPrefix(text, "//") {
			continue
		}
		for line := range strings.SplitSeq(text, "\n") {
			if commentDirective(line) == "" {
				continue
			}
			if len(g.doc) > 0 && len(list) == len(g.doc) {
				list = append(list, &ast.Comment{Slash: anchor, Text: "//"})
			}
			list = append(list, &ast.Comment{Slash: anchor, Text: line})
		}
	}
	if len(list) == 0 {
		return nil
	}
	first := original.List[0].Text
	list[0].Text = first[:len(first)-len(strings.TrimLeft(first, "\r\n"))] + list[0].Text
	return &ast.CommentGroup{List: list}
}

func commentDirective(line string) string {
	text, ok := strings.CutPrefix(line, "//")
	if !ok {
		return ""
	}
	word, _, ok := strings.Cut(text, ":")
	if !ok || word == "" || strings.IndexFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	}) >= 0 {
		return ""
	}
	return normalizeDirectiveSpace(text)
}

// Keep whitespace inside quoted arguments, as in //go:embed "a  b.txt".
func normalizeDirectiveSpace(text string) string {
	var result strings.Builder
	var quote rune
	var escaped, space bool
	for _, r := range text {
		if quote != 0 {
			result.WriteRune(r)
			switch {
			case escaped:
				escaped = false
			case r == '\\' && quote != '`':
				escaped = true
			case r == quote:
				quote = 0
			}
			continue
		}
		if unicode.IsSpace(r) {
			space = result.Len() != 0
			continue
		}
		if space {
			result.WriteByte(' ')
			space = false
		}
		result.WriteRune(r)
		switch r {
		case '"', '\'', '`':
			quote = r
		}
	}
	return result.String()
}
