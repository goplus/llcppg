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
The fix detects when a builtin `int` type comes from error recovery and marks it appropriately instead of skipping the function entirely.

### Detection Method
1. Check if we're processing a builtin `int` type (`t.Kind == clang.TypeInt`)
2. Call `t.TypeDeclaration()` to get the type's declaration cursor
3. For legitimate builtin types, this should return a null cursor
4. For error-recovery types, this might return a non-null cursor
5. If non-null cursor detected, return a BuiltinType with TypeFlag 0 to mark as undefined type

### New Behavior
Instead of completely skipping functions with undefined types (which was too aggressive), the function is still included in the AST but with a recognizable pattern:
- `Kind`: 6 (int)
- `Flags`: 0 (no type modifiers)

This allows downstream processing to:
- Identify potentially problematic functions
- Handle undefined types appropriately without losing the function declaration
- Distinguish from completely missing functions

## Test Cases
- `testdata/undef_type/temp.h`: Contains `undef fn();` 
- `testdata/undef_type/expect.json`: Expected output includes the function with int/flags=0

## Expected Behavior
- Functions with undefined types: Processed with TypeFlag 0 marking
- Legitimate functions: Processed normally
- Related issue #109: Method conversion should work correctly with undefined types marked but not hidden

## Validation
Use clang AST dump to see the difference:
- Undefined: `FunctionDecl ... invalid fn 'int ()'` (marked invalid)
- Legitimate: `FunctionDecl ... fn 'int ()'` (normal)