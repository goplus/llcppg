package clang

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// A compilation database holds all information used to compile files in a
// project. For each file in the database, it can be queried for the working
// directory or the command line used for the compiler invocation.
//
// Must be freed by \c clang_CompilationDatabase_dispose
type CompilationDatabase uintptr

// Contains the results of a search in the compilation database
//
// When searching for the compile command for a file, the compilation db can
// return several commands, as the file may have been compiled with
// different options in different places of the project. This choice of compile
// commands is wrapped in this opaque data structure. It must be freed by
// \c clang_CompileCommands_dispose.
type CompileCommands uintptr

// Represents the command line invocation to compile a specific file.
type CompileCommand uintptr

// Error codes for Compilation Database
type CompilationDatabase_Error c.Int

const (
	CompilationDatabase_NoError            CompilationDatabase_Error = 0
	CompilationDatabase_CanNotLoadDatabase CompilationDatabase_Error = 1
)

// Creates a compilation database from the database found in directory
// buildDir. For example, CMake can output a compile_commands.json which can
// be used to build the database.
//
// It must be freed by \c clang_CompilationDatabase_dispose.
//
//go:linkname CompilationDatabaseFromDirectory C.clang_CompilationDatabase_fromDirectory
func CompilationDatabaseFromDirectory(BuildDir *c.Char, ErrorCode *CompilationDatabase_Error) CompilationDatabase

// Free the given compilation database
//
// llgo:link CompilationDatabase.Dispose C.clang_CompilationDatabase_dispose
func (_llcppg_param1 CompilationDatabase) Dispose() {
}

// Find the compile commands used for a file. The compile commands
// must be freed by \c clang_CompileCommands_dispose.
//
// llgo:link CompilationDatabase.CompileCommands C.clang_CompilationDatabase_getCompileCommands
func (_llcppg_param1 CompilationDatabase) CompileCommands(CompleteFileName *c.Char) CompileCommands {
	return 0
}

// Get all the compile commands in the given compilation database.
//
// llgo:link CompilationDatabase.AllCompileCommands C.clang_CompilationDatabase_getAllCompileCommands
func (_llcppg_param1 CompilationDatabase) AllCompileCommands() CompileCommands {
	return 0
}

// Free the given CompileCommands
//
// llgo:link CompileCommands.Dispose C.clang_CompileCommands_dispose
func (_llcppg_param1 CompileCommands) Dispose() {
}

// Get the number of CompileCommand we have for a file
//
// llgo:link CompileCommands.Size C.clang_CompileCommands_getSize
func (_llcppg_param1 CompileCommands) Size() c.Uint {
	return 0
}

// Get the I'th CompileCommand for a file
//
// Note : 0 <= i < clang_CompileCommands_getSize(CXCompileCommands)
//
// llgo:link CompileCommands.Command C.clang_CompileCommands_getCommand
func (_llcppg_param1 CompileCommands) Command(I c.Uint) CompileCommand {
	return 0
}

// Get the working directory where the CompileCommand was executed from
//
// llgo:link CompileCommand.Directory C.clang_CompileCommand_getDirectory
func (_llcppg_param1 CompileCommand) Directory() String {
	return String{}
}

// Get the filename associated with the CompileCommand.
//
// llgo:link CompileCommand.Filename C.clang_CompileCommand_getFilename
func (_llcppg_param1 CompileCommand) Filename() String {
	return String{}
}

// Get the number of arguments in the compiler invocation.
//
// llgo:link CompileCommand.NumArgs C.clang_CompileCommand_getNumArgs
func (_llcppg_param1 CompileCommand) NumArgs() c.Uint {
	return 0
}

// Get the I'th argument value in the compiler invocations
//
// Invariant :
//  - argument 0 is the compiler executable
//
// llgo:link CompileCommand.Arg C.clang_CompileCommand_getArg
func (_llcppg_param1 CompileCommand) Arg(I c.Uint) String {
	return String{}
}

// Get the number of source mappings for the compiler invocation.
//
// llgo:link CompileCommand.NumMappedSources C.clang_CompileCommand_getNumMappedSources
func (_llcppg_param1 CompileCommand) NumMappedSources() c.Uint {
	return 0
}

// Get the I'th mapped source path for the compiler invocation.
//
// llgo:link CompileCommand.MappedSourcePath C.clang_CompileCommand_getMappedSourcePath
func (_llcppg_param1 CompileCommand) MappedSourcePath(I c.Uint) String {
	return String{}
}

// Get the I'th mapped source content for the compiler invocation.
//
// llgo:link CompileCommand.MappedSourceContent C.clang_CompileCommand_getMappedSourceContent
func (_llcppg_param1 CompileCommand) MappedSourceContent(I c.Uint) String {
	return String{}
}
