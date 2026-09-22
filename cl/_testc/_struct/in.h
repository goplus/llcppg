struct Foo {
	struct Shorts {
		short s;
		unsigned short us;
	};
	struct {
		float x;
		float y;
	} u;
	struct Shorts shorts;
};

struct Bar {
	struct {
		short s;
		unsigned short us;
	};
};
