// A global with a string-literal initializer must not panic the generator.
// Clang completes `text` to `char[3]` ("hi" plus the terminating NUL), and the
// initialized global itself is intentionally omitted from the bindings.
char text[] = "hi";

// A field whose extent is inferred from the string-completed array keeps the
// resolved bound: sizeof(text) is 3, so `a` is a fixed [3]c.Char array.
struct Owner {
    char a[sizeof(text)];
};
