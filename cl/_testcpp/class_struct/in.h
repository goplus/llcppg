class Foo {
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

class Bar {
	struct {
		short s;
		unsigned short us;
	};
public:
    struct {
        char c;
        unsigned char uc;
    };
	class {
    public:
		float x;
		float y;
	};
};
