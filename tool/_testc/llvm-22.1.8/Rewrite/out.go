package Rewrite

import (
	"clang/CXSourceLocation"
	"clang/Index"
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

type Rewriter uintptr

// Create CXRewriter.
//
//go:linkname RewriterCreate C.clang_CXRewriter_create
func RewriterCreate(TU Index.TranslationUnit) Rewriter

// Insert the specified string at the specified location in the original buffer.
//
// llgo:link Rewriter.InsertTextBefore C.clang_CXRewriter_insertTextBefore
func (Rew Rewriter) InsertTextBefore(Loc CXSourceLocation.SourceLocation, Insert *c.Char) {
}

// Replace the specified range of characters in the input with the specified
// replacement.
//
// llgo:link Rewriter.ReplaceText C.clang_CXRewriter_replaceText
func (Rew Rewriter) ReplaceText(ToBeReplaced CXSourceLocation.SourceRange, Replacement *c.Char) {
}

// Remove the specified range.
//
// llgo:link Rewriter.RemoveText C.clang_CXRewriter_removeText
func (Rew Rewriter) RemoveText(ToBeRemoved CXSourceLocation.SourceRange) {
}

// Save all changed files to disk.
// Returns 1 if any files were not saved successfully, returns 0 otherwise.
//
// llgo:link Rewriter.OverwriteChangedFiles C.clang_CXRewriter_overwriteChangedFiles
func (Rew Rewriter) OverwriteChangedFiles() c.Int {
	return 0
}

// Write out rewritten version of the main file to stdout.
//
// llgo:link Rewriter.WriteMainFileToStdOut C.clang_CXRewriter_writeMainFileToStdOut
func (Rew Rewriter) WriteMainFileToStdOut() {
}

// Free the given CXRewriter.
//
// llgo:link Rewriter.Dispose C.clang_CXRewriter_dispose
func (Rew Rewriter) Dispose() {
}
