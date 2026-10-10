// A C name typedef'd more than once with the same underlying type, as happens
// when a forward typedef and its defining typedef (or two headers) describe the
// same type. Before the fix a second "type MyInt" declaration was emitted,
// panicking in gogen with "MyInt redeclared in this block". The generator must
// instead reuse the single definition. See issue goplus/llcppg#1001.
typedef int MyInt;

typedef int MyInt;
