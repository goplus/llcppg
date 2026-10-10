// Nested union regression adapted from v0.7.8 _xtool/internal/parser/testdata/union/temp.h.
union OuterUnion {
    int i;
    float f;
    // Exercise different outer and inner sizes and alignments.
    double d[2];
    union {
        int c;
        short s;
    } inner;
};

// The inner union determines the outer size and alignment, not the int member.
union OuterWithLargeInner {
    int i;
    union {
        double d[2];
        short s;
    } inner;
};
