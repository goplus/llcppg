union A {
    int a;
    int b;
};

union OuterUnion {
    int i;
    float f;
    union {
        int c;
        short s;
    } inner;
};
