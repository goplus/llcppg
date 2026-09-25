package clang

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

// A parsed comment.
type Comment struct {
	ASTNode         unsafe.Pointer
	TranslationUnit TranslationUnit
}

// Describes the type of the comment AST node (\c CXComment).  A comment
// node can be considered block content (e. g., paragraph), inline content
// (plain text) or neither (the root AST node).
type CommentKind c.Int

const (
	// Null comment.  No AST node is constructed at the requested location
	// because there is no text or a syntax error.
	Comment_Null CommentKind = 0
	// Plain text.  Inline content.
	Comment_Text CommentKind = 1
	// A command with word-like arguments that is considered inline content.
	//
	// For example: \\c command.
	Comment_InlineCommand CommentKind = 2
	// HTML start tag with attributes (name-value pairs).  Considered
	// inline content.
	//
	// For example:
	// \verbatim
	// <br> <br /> <a href="http://example.org/">
	// \endverbatim
	Comment_HTMLStartTag CommentKind = 3
	// HTML end tag.  Considered inline content.
	//
	// For example:
	// \verbatim
	// </a>
	// \endverbatim
	Comment_HTMLEndTag CommentKind = 4
	// A paragraph, contains inline comment.  The paragraph itself is
	// block content.
	Comment_Paragraph CommentKind = 5
	// A command that has zero or more word-like arguments (number of
	// word-like arguments depends on command name) and a paragraph as an
	// argument.  Block command is block content.
	//
	// Paragraph argument is also a child of the block command.
	//
	// For example: \has 0 word-like arguments and a paragraph argument.
	//
	// AST nodes of special kinds that parser knows about (e. g., \\param
	// command) have their own node kinds.
	Comment_BlockCommand CommentKind = 6
	// A \\param or \\arg command that describes the function parameter
	// (name, passing direction, description).
	//
	// For example: \\param [in] ParamName description.
	Comment_ParamCommand CommentKind = 7
	// A \\tparam command that describes a template parameter (name and
	// description).
	//
	// For example: \\tparam T description.
	Comment_TParamCommand CommentKind = 8
	// A verbatim block command (e. g., preformatted code).  Verbatim
	// block has an opening and a closing command and contains multiple lines of
	// text (\c CXComment_VerbatimBlockLine child nodes).
	//
	// For example:
	// \\verbatim
	// aaa
	// \\endverbatim
	Comment_VerbatimBlockCommand CommentKind = 9
	// A line of text that is contained within a
	// CXComment_VerbatimBlockCommand node.
	Comment_VerbatimBlockLine CommentKind = 10
	// A verbatim line command.  Verbatim line has an opening command,
	// a single line of text (up to the newline after the opening command) and
	// has no closing command.
	Comment_VerbatimLine CommentKind = 11
	// A full comment attached to a declaration, contains block content.
	Comment_FullComment CommentKind = 12
)

// The most appropriate rendering mode for an inline command, chosen on
// command semantics in Doxygen.
type CommentInlineCommandRenderKind c.Int

const (
	// Command argument should be rendered in a normal font.
	CommentInlineCommandRenderKind_Normal CommentInlineCommandRenderKind = 0
	// Command argument should be rendered in a bold font.
	CommentInlineCommandRenderKind_Bold CommentInlineCommandRenderKind = 1
	// Command argument should be rendered in a monospaced font.
	CommentInlineCommandRenderKind_Monospaced CommentInlineCommandRenderKind = 2
	// Command argument should be rendered emphasized (typically italic
	// font).
	CommentInlineCommandRenderKind_Emphasized CommentInlineCommandRenderKind = 3
	// Command argument should not be rendered (since it only defines an anchor).
	CommentInlineCommandRenderKind_Anchor CommentInlineCommandRenderKind = 4
)

// Describes parameter passing direction for \\param or \\arg command.
type CommentParamPassDirection c.Int

const (
	// The parameter is an input parameter.
	CommentParamPassDirection_In CommentParamPassDirection = 0
	// The parameter is an output parameter.
	CommentParamPassDirection_Out CommentParamPassDirection = 1
	// The parameter is an input and output parameter.
	CommentParamPassDirection_InOut CommentParamPassDirection = 2
)

type APISetImpl struct {
}

// CXAPISet is an opaque type that represents a data structure containing all
// the API information for a given translation unit. This can be used for a
// single symbol symbol graph for a given symbol.
type APISet = *APISetImpl

// Given a cursor that represents a documentable entity (e.g.,
// declaration), return the associated parsed comment as a
// \c CXComment_FullComment AST node.
//
// llgo:link Cursor.ParsedComment C.clang_Cursor_getParsedComment
func (C Cursor) ParsedComment() Comment {
	return Comment{}
}

// \param Comment AST node of any kind.
//
// \returns the type of the AST node.
//
// llgo:link Comment.Kind C.clang_Comment_getKind
func (Comment Comment) Kind() CommentKind {
	return 0
}

// \param Comment AST node of any kind.
//
// \returns number of children of the AST node.
//
// llgo:link Comment.NumChildren C.clang_Comment_getNumChildren
func (Comment Comment) NumChildren() c.Uint {
	return 0
}

// \param Comment AST node of any kind.
//
// \param ChildIdx child index (zero-based).
//
// \returns the specified child of the AST node.
//
// llgo:link Comment.Child C.clang_Comment_getChild
func (Comment Comment) Child(ChildIdx c.Uint) Comment {
	return Comment
}

// A \c CXComment_Paragraph node is considered whitespace if it contains
// only \c CXComment_Text nodes that are empty or whitespace.
//
// Other AST nodes (except \c CXComment_Paragraph and \c CXComment_Text) are
// never considered whitespace.
//
// \returns non-zero if \c Comment is whitespace.
//
// llgo:link Comment.IsWhitespace C.clang_Comment_isWhitespace
func (Comment Comment) IsWhitespace() c.Uint {
	return 0
}

// \returns non-zero if \c Comment is inline content and has a newline
// immediately following it in the comment text.  Newlines between paragraphs
// do not count.
//
// llgo:link Comment.InlineContentCommentHasTrailingNewline C.clang_InlineContentComment_hasTrailingNewline
func (Comment Comment) InlineContentCommentHasTrailingNewline() c.Uint {
	return 0
}

// \param Comment a \c CXComment_Text AST node.
//
// \returns text contained in the AST node.
//
// llgo:link Comment.TextCommentGetText C.clang_TextComment_getText
func (Comment Comment) TextCommentGetText() String {
	return String{}
}

// \param Comment a \c CXComment_InlineCommand AST node.
//
// \returns name of the inline command.
//
// llgo:link Comment.InlineCommandCommentGetCommandName C.clang_InlineCommandComment_getCommandName
func (Comment Comment) InlineCommandCommentGetCommandName() String {
	return String{}
}

// \param Comment a \c CXComment_InlineCommand AST node.
//
// \returns the most appropriate rendering mode, chosen on command
// semantics in Doxygen.
//
// llgo:link Comment.InlineCommandCommentGetRenderKind C.clang_InlineCommandComment_getRenderKind
func (Comment Comment) InlineCommandCommentGetRenderKind() CommentInlineCommandRenderKind {
	return 0
}

// \param Comment a \c CXComment_InlineCommand AST node.
//
// \returns number of command arguments.
//
// llgo:link Comment.InlineCommandCommentGetNumArgs C.clang_InlineCommandComment_getNumArgs
func (Comment Comment) InlineCommandCommentGetNumArgs() c.Uint {
	return 0
}

// \param Comment a \c CXComment_InlineCommand AST node.
//
// \param ArgIdx argument index (zero-based).
//
// \returns text of the specified argument.
//
// llgo:link Comment.InlineCommandCommentGetArgText C.clang_InlineCommandComment_getArgText
func (Comment Comment) InlineCommandCommentGetArgText(ArgIdx c.Uint) String {
	return String{}
}

// \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST
// node.
//
// \returns HTML tag name.
//
// llgo:link Comment.HTMLTagCommentGetTagName C.clang_HTMLTagComment_getTagName
func (Comment Comment) HTMLTagCommentGetTagName() String {
	return String{}
}

// \param Comment a \c CXComment_HTMLStartTag AST node.
//
// \returns non-zero if tag is self-closing (for example, &lt;br /&gt;).
//
// llgo:link Comment.HTMLStartTagCommentIsSelfClosing C.clang_HTMLStartTagComment_isSelfClosing
func (Comment Comment) HTMLStartTagCommentIsSelfClosing() c.Uint {
	return 0
}

// \param Comment a \c CXComment_HTMLStartTag AST node.
//
// \returns number of attributes (name-value pairs) attached to the start tag.
//
// llgo:link Comment.HTMLStartTagGetNumAttrs C.clang_HTMLStartTag_getNumAttrs
func (Comment Comment) HTMLStartTagGetNumAttrs() c.Uint {
	return 0
}

// \param Comment a \c CXComment_HTMLStartTag AST node.
//
// \param AttrIdx attribute index (zero-based).
//
// \returns name of the specified attribute.
//
// llgo:link Comment.HTMLStartTagGetAttrName C.clang_HTMLStartTag_getAttrName
func (Comment Comment) HTMLStartTagGetAttrName(AttrIdx c.Uint) String {
	return String{}
}

// \param Comment a \c CXComment_HTMLStartTag AST node.
//
// \param AttrIdx attribute index (zero-based).
//
// \returns value of the specified attribute.
//
// llgo:link Comment.HTMLStartTagGetAttrValue C.clang_HTMLStartTag_getAttrValue
func (Comment Comment) HTMLStartTagGetAttrValue(AttrIdx c.Uint) String {
	return String{}
}

// \param Comment a \c CXComment_BlockCommand AST node.
//
// \returns name of the block command.
//
// llgo:link Comment.BlockCommandCommentGetCommandName C.clang_BlockCommandComment_getCommandName
func (Comment Comment) BlockCommandCommentGetCommandName() String {
	return String{}
}

// \param Comment a \c CXComment_BlockCommand AST node.
//
// \returns number of word-like arguments.
//
// llgo:link Comment.BlockCommandCommentGetNumArgs C.clang_BlockCommandComment_getNumArgs
func (Comment Comment) BlockCommandCommentGetNumArgs() c.Uint {
	return 0
}

// \param Comment a \c CXComment_BlockCommand AST node.
//
// \param ArgIdx argument index (zero-based).
//
// \returns text of the specified word-like argument.
//
// llgo:link Comment.BlockCommandCommentGetArgText C.clang_BlockCommandComment_getArgText
func (Comment Comment) BlockCommandCommentGetArgText(ArgIdx c.Uint) String {
	return String{}
}

// \param Comment a \c CXComment_BlockCommand or
// \c CXComment_VerbatimBlockCommand AST node.
//
// \returns paragraph argument of the block command.
//
// llgo:link Comment.BlockCommandCommentGetParagraph C.clang_BlockCommandComment_getParagraph
func (Comment Comment) BlockCommandCommentGetParagraph() Comment {
	return Comment
}

// \param Comment a \c CXComment_ParamCommand AST node.
//
// \returns parameter name.
//
// llgo:link Comment.ParamCommandCommentGetParamName C.clang_ParamCommandComment_getParamName
func (Comment Comment) ParamCommandCommentGetParamName() String {
	return String{}
}

// \param Comment a \c CXComment_ParamCommand AST node.
//
// \returns non-zero if the parameter that this AST node represents was found
// in the function prototype and \c clang_ParamCommandComment_getParamIndex
// function will return a meaningful value.
//
// llgo:link Comment.ParamCommandCommentIsParamIndexValid C.clang_ParamCommandComment_isParamIndexValid
func (Comment Comment) ParamCommandCommentIsParamIndexValid() c.Uint {
	return 0
}

// \param Comment a \c CXComment_ParamCommand AST node.
//
// \returns zero-based parameter index in function prototype.
//
// llgo:link Comment.ParamCommandCommentGetParamIndex C.clang_ParamCommandComment_getParamIndex
func (Comment Comment) ParamCommandCommentGetParamIndex() c.Uint {
	return 0
}

// \param Comment a \c CXComment_ParamCommand AST node.
//
// \returns non-zero if parameter passing direction was specified explicitly in
// the comment.
//
// llgo:link Comment.ParamCommandCommentIsDirectionExplicit C.clang_ParamCommandComment_isDirectionExplicit
func (Comment Comment) ParamCommandCommentIsDirectionExplicit() c.Uint {
	return 0
}

// \param Comment a \c CXComment_ParamCommand AST node.
//
// \returns parameter passing direction.
//
// llgo:link Comment.ParamCommandCommentGetDirection C.clang_ParamCommandComment_getDirection
func (Comment Comment) ParamCommandCommentGetDirection() CommentParamPassDirection {
	return 0
}

// \param Comment a \c CXComment_TParamCommand AST node.
//
// \returns template parameter name.
//
// llgo:link Comment.TParamCommandCommentGetParamName C.clang_TParamCommandComment_getParamName
func (Comment Comment) TParamCommandCommentGetParamName() String {
	return String{}
}

// \param Comment a \c CXComment_TParamCommand AST node.
//
// \returns non-zero if the parameter that this AST node represents was found
// in the template parameter list and
// \c clang_TParamCommandComment_getDepth and
// \c clang_TParamCommandComment_getIndex functions will return a meaningful
// value.
//
// llgo:link Comment.TParamCommandCommentIsParamPositionValid C.clang_TParamCommandComment_isParamPositionValid
func (Comment Comment) TParamCommandCommentIsParamPositionValid() c.Uint {
	return 0
}

// \param Comment a \c CXComment_TParamCommand AST node.
//
// \returns zero-based nesting depth of this parameter in the template parameter list.
//
// For example,
// \verbatim
//     template<typename C, template<typename T> class TT>
//     void test(TT<int> aaa);
// \endverbatim
// for C and TT nesting depth is 0,
// for T nesting depth is 1.
//
// llgo:link Comment.TParamCommandCommentGetDepth C.clang_TParamCommandComment_getDepth
func (Comment Comment) TParamCommandCommentGetDepth() c.Uint {
	return 0
}

// \param Comment a \c CXComment_TParamCommand AST node.
//
// \returns zero-based parameter index in the template parameter list at a
// given nesting depth.
//
// For example,
// \verbatim
//     template<typename C, template<typename T> class TT>
//     void test(TT<int> aaa);
// \endverbatim
// for C and TT nesting depth is 0, so we can ask for index at depth 0:
// at depth 0 C's index is 0, TT's index is 1.
//
// For T nesting depth is 1, so we can ask for index at depth 0 and 1:
// at depth 0 T's index is 1 (same as TT's),
// at depth 1 T's index is 0.
//
// llgo:link Comment.TParamCommandCommentGetIndex C.clang_TParamCommandComment_getIndex
func (Comment Comment) TParamCommandCommentGetIndex(Depth c.Uint) c.Uint {
	return 0
}

// \param Comment a \c CXComment_VerbatimBlockLine AST node.
//
// \returns text contained in the AST node.
//
// llgo:link Comment.VerbatimBlockLineCommentGetText C.clang_VerbatimBlockLineComment_getText
func (Comment Comment) VerbatimBlockLineCommentGetText() String {
	return String{}
}

// \param Comment a \c CXComment_VerbatimLine AST node.
//
// \returns text contained in the AST node.
//
// llgo:link Comment.VerbatimLineCommentGetText C.clang_VerbatimLineComment_getText
func (Comment Comment) VerbatimLineCommentGetText() String {
	return String{}
}

// Convert an HTML tag AST node to string.
//
// \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST
// node.
//
// \returns string containing an HTML tag.
//
// llgo:link Comment.HTMLTagCommentGetAsString C.clang_HTMLTagComment_getAsString
func (Comment Comment) HTMLTagCommentGetAsString() String {
	return String{}
}

// Convert a given full parsed comment to an HTML fragment.
//
// Specific details of HTML layout are subject to change.  Don't try to parse
// this HTML back into an AST, use other APIs instead.
//
// Currently the following CSS classes are used:
// \li "para-brief" for \paragraph and equivalent commands;
// \li "para-returns" for \\returns paragraph and equivalent commands;
// \li "word-returns" for the "Returns" word in \\returns paragraph.
//
// Function argument documentation is rendered as a \<dl\> list with arguments
// sorted in function prototype order.  CSS classes used:
// \li "param-name-index-NUMBER" for parameter name (\<dt\>);
// \li "param-descr-index-NUMBER" for parameter description (\<dd\>);
// \li "param-name-index-invalid" and "param-descr-index-invalid" are used if
// parameter index is invalid.
//
// Template parameter documentation is rendered as a \<dl\> list with
// parameters sorted in template parameter list order.  CSS classes used:
// \li "tparam-name-index-NUMBER" for parameter name (\<dt\>);
// \li "tparam-descr-index-NUMBER" for parameter description (\<dd\>);
// \li "tparam-name-index-other" and "tparam-descr-index-other" are used for
// names inside template template parameters;
// \li "tparam-name-index-invalid" and "tparam-descr-index-invalid" are used if
// parameter position is invalid.
//
// \param Comment a \c CXComment_FullComment AST node.
//
// \returns string containing an HTML fragment.
//
// llgo:link Comment.FullCommentGetAsHTML C.clang_FullComment_getAsHTML
func (Comment Comment) FullCommentGetAsHTML() String {
	return String{}
}

// Convert a given full parsed comment to an XML document.
//
// A Relax NG schema for the XML can be found in comment-xml-schema.rng file
// inside clang source tree.
//
// \param Comment a \c CXComment_FullComment AST node.
//
// \returns string containing an XML document.
//
// llgo:link Comment.FullCommentGetAsXML C.clang_FullComment_getAsXML
func (Comment Comment) FullCommentGetAsXML() String {
	return String{}
}

// Traverses the translation unit to create a \c CXAPISet.
//
// \param tu is the \c CXTranslationUnit to build the \c CXAPISet for.
//
// \param out_api is a pointer to the output of this function. It is needs to be
// disposed of by calling clang_disposeAPISet.
//
// \returns Error code indicating success or failure of the APISet creation.
//
// llgo:link (*TranslationUnitImpl).CreateAPISet C.clang_createAPISet
func (tu *TranslationUnitImpl) CreateAPISet(out_api *APISet) ErrorCode {
	return 0
}

// Dispose of an APISet.
//
// The provided \c CXAPISet can not be used after this function is called.
//
// llgo:link (*APISetImpl).Dispose C.clang_disposeAPISet
func (api *APISetImpl) Dispose() {
}

// Generate a single symbol symbol graph for the given USR. Returns a null
// string if the associated symbol can not be found in the provided \c CXAPISet.
//
// The output contains the symbol graph as well as some additional information
// about related symbols.
//
// \param usr is a string containing the USR of the symbol to generate the
// symbol graph for.
//
// \param api the \c CXAPISet to look for the symbol in.
//
// \returns a string containing the serialized symbol graph representation for
// the symbol being queried or a null string if it can not be found in the
// APISet.
//
//go:linkname GetSymbolGraphForUSR C.clang_getSymbolGraphForUSR
func GetSymbolGraphForUSR(usr *c.Char, api APISet) String

// Generate a single symbol symbol graph for the declaration at the given
// cursor. Returns a null string if the AST node for the cursor isn't a
// declaration.
//
// The output contains the symbol graph as well as some additional information
// about related symbols.
//
// \param cursor the declaration for which to generate the single symbol symbol
// graph.
//
// \returns a string containing the serialized symbol graph representation for
// the symbol being queried or a null string if it can not be found in the
// APISet.
//
// llgo:link Cursor.SymbolGraphFor C.clang_getSymbolGraphForCursor
func (cursor Cursor) SymbolGraphFor() String {
	return String{}
}
