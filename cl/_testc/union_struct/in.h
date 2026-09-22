struct Foo {
    union Shorts {
        short s;
        unsigned short us;
    };
    union {
        float x;
        float y;
    } u;
};

struct Bar {
    union {
        short s;
        unsigned short us;
    };
};
