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
The fix detects when a builtin `int` type comes from error recovery and marks it with a special signature to identify undefined types.

### Detection Method
1. Check if we're processing a builtin `int` type (`t.Kind == clang.TypeInt`)
2. Call `t.TypeDeclaration()` to get the type's declaration cursor
3. For legitimate builtin types, this should return a null cursor
4. For error-recovery types, this might return a non-null cursor
5. If non-null cursor detected, return a BuiltinType with `TypeKind: Void` and `TypeFlag: Signed` to mark as undefined type

### New Behavior
Functions with undefined types are still included in the AST but with a recognizable pattern:
- `Kind`: 0 (Void)
- `Flags`: 1 (Signed)

This allows downstream processing to:
- Identify potentially problematic functions with undefined types
- Handle undefined types appropriately without losing the function declaration
- Distinguish from legitimate void functions or missing functions

## Test Cases
- `testdata/undef_type/temp.h`: Contains `undef fn();` 
- `testdata/undef_type/expect.json`: Expected output includes the function with Void/Signed marking

## Expected Behavior
- Functions with undefined types: Processed with `TypeKind: Void` and `TypeFlag: Signed` marking
- Legitimate functions: Processed normally
- Related issue #109: Method conversion should work correctly with undefined types marked but not hidden

## Validation
Use clang AST dump to see the difference:
- Undefined: `FunctionDecl ... invalid fn 'int ()'` (marked invalid)
- Legitimate: `FunctionDecl ... fn 'int ()'` (normal)