# Undefined Type Detection Fix

## Problem
When libclang encounters undefined types like `undef fn();`, it performs error recovery by defaulting the undefined type to `int`. This causes the parser to generate:

```json
"Ret": {
    "_Type": "BuiltinType",
    "Kind": 6,  // int
    "Flags": 0
}
```

This is misleading because `undef` is not actually an `int` type.

## Solution
The fix detects when a builtin `int` type comes from error recovery and skips the invalid function declaration.

### Detection Method
1. Check if we're processing a builtin `int` type (`t.Kind == clang.TypeInt`)
2. Call `t.TypeDeclaration()` to get the type's declaration cursor
3. For legitimate builtin types, this should return a null cursor
4. For error-recovery types, this might return a non-null cursor
5. If non-null cursor detected, skip the type (return nil)

### Error Handling
The fix includes proper nil handling throughout the type processing pipeline:
- `ProcessFunctionType`: Check if return type or parameter types are nil
- `ProcessType`: Check pointer, reference, and array element types for nil
- `createBaseField`: Check if field type is nil
- `ProcessFieldList`: Skip nil fields
- `visitTop`: Skip nil function declarations

## Test Cases
- `testdata/undef_type/temp.h`: Contains `undef fn();` 
- `testdata/undef_type/expect.json`: Expected empty output (no declarations)

## Expected Behavior
- Functions with undefined types: Skipped (no output)
- Legitimate functions: Processed normally
- Related issue #109: Method conversion should work correctly without interference from undefined types

## Validation
Use clang AST dump to see the difference:
- Undefined: `FunctionDecl ... invalid fn 'int ()'` (marked invalid)
- Legitimate: `FunctionDecl ... fn 'int ()'` (normal)