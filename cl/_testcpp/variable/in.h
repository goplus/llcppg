extern int total;

// A const-qualified variable is an immutable var: it maps to a Go var linked
// to the same C++ symbol, exactly like a mutable var.
extern const int maxValue;

namespace bar {
	extern int value;

	extern const double ratio;
}

class Config
{
public:
	static int count;
	static const char* name;
};
