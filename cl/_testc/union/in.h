// A tagged C union (issue goplus/llcppg#764): bound as a Go struct that owns a
// storage array with the same size and alignment as the union, plus a typed
// XGo_union_<member> accessor for each member. A struct uses it by tag name.
union Value {
	int           i;
	double        d;
	char         *s;
	unsigned char raw[8];
};

struct Variant {
	int          kind;
	union Value  val;
};

// An all-double union: storage elements are float64 so by-value ABI
// classification is preserved (proposal §1 step 2).
union Vec2 {
	double xy[2];
	double v;
};
