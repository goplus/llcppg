// Bit-field layout and accessor generation. See issue goplus/llcppg#770.

// A plain run of unsigned bit-fields packed into one storage word, plus a
// regular member that keeps its natural conversion.
struct Flags {
	unsigned int a : 1;
	unsigned int b : 3;
	unsigned int c : 12;
	int value;
};

// Signed and bool bit-fields, with an unnamed zero-width separator that forces
// the next field into a fresh storage unit (no accessor for the separator).
struct Mixed {
	int sx : 5;
	unsigned int : 0;
	_Bool flag : 1;
	unsigned int rest : 20;
};
