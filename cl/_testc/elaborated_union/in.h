// Regression cases from v0.7.8 _xtool/internal/parser/testdata/union/temp.h.
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
