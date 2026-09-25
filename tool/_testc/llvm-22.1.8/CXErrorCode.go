package clang

import "github.com/goplus/lib/c"

// Error codes returned by libclang routines.
//
// Zero (\c CXError_Success) is the only error code indicating success.  Other
// error codes, including not yet assigned non-zero values, indicate errors.
type ErrorCode c.Int

const (
	// No error.
	Error_Success ErrorCode = 0
	// A generic error code, no further details are available.
	//
	// Errors of this kind can get their own specific error codes in future
	// libclang versions.
	Error_Failure ErrorCode = 1
	// libclang crashed while performing the requested operation.
	Error_Crashed ErrorCode = 2
	// The function detected that the arguments violate the function
	// contract.
	Error_InvalidArguments ErrorCode = 3
	// An AST deserialization error has occurred.
	Error_ASTReadError ErrorCode = 4
)
