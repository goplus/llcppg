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
