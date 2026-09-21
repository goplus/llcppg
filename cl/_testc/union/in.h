// A simple C union: three members sharing storage. The alignment is 8 (double),
// the size is 8, so the storage is [1]uint64 (mixed int/float/pointer, so the
// element type is an unsigned integer, not a float). Each accessible member
// gets a XGof_ref_<member> accessor returning a typed pointer into the storage.
union Value {
	int i;
	double d;
	char *p;
};

// A member named like a Go keyword is kept verbatim: the accessor is
// XGof_ref_type, not renamed.
union Keyword {
	int type;
	unsigned func;
};

// An all-float union: every scalar leaf is a double of width 8, so the storage
// element is float64 to keep by-value ABI classification correct.
union Doubles {
	double x;
	double y;
};
