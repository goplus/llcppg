package clang

import _ "unsafe"

// Installs error handler that prints error message to stderr and calls abort().
// Replaces currently installed error handler (if any).
//
//go:linkname InstallAbortingLlvmFatalErrorHandler C.clang_install_aborting_llvm_fatal_error_handler
func InstallAbortingLlvmFatalErrorHandler()

// Removes currently installed error handler (if any).
// If no error handler is intalled, the default strategy is to print error
// message to stderr and call exit(1).
//
//go:linkname UninstallLlvmFatalErrorHandler C.clang_uninstall_llvm_fatal_error_handler
func UninstallLlvmFatalErrorHandler()
