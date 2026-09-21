extern int foo;

extern const char* name;

extern unsigned _count;

// A const-qualified variable is an immutable var: it maps to a Go var linked
// to the same C symbol, exactly like a mutable var.
extern const int maxValue;

extern const double ratio;
