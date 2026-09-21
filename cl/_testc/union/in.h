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

// Alignment 1: every member is byte-sized, so the storage element is uint8.
union Bytes {
	char c;
	unsigned char uc;
};

// Alignment 2: the widest member is a short, so the storage element is uint16.
union Shorts {
	short s;
	unsigned short us;
};

// All-float, width 4: every scalar leaf is a float of width 4, so the storage
// element is float32.
union Floats {
	float x;
	float y;
};

// A non-float leaf defeats the all-float rule even when the alignment is 8: the
// int member is not a float, so the storage falls back to an unsigned integer
// element ([4]uint64 for the 32-byte double[4]) rather than [4]float64. The
// array member also exercises an accessor whose type is a pointer to an array.
union BigFloat {
	int i;
	double arr[4];
};

union Empty {
}

// A union that is only forward-declared has no body and no known size, so it
// produces no Go declaration at all.
union Opaque;
