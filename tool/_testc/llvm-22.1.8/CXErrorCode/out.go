package CXErrorCode

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L$(llvm-config --libdir) -lclang; -lclang"

// Error codes returned by libclang routines.
//
// Zero (\c CXError_Success) is the only error code indicating success.  Other
// error codes, including not yet assigned non-zero values, indicate errors.
type ErrorCode c.Int

const (
	Error_Success          ErrorCode = 0
	Error_Failure          ErrorCode = 1
	Error_Crashed          ErrorCode = 2
	Error_InvalidArguments ErrorCode = 3
	Error_ASTReadError     ErrorCode = 4
)
