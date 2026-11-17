struct StructStruct {
    struct { int x; } a, b, c;
};

union UnionStruct {
    struct { int x; } a, b;
};

struct StructUnion {
    union { int x; double y; } a, b;
};

union UnionUnion {
    union { int x; double y; } a, b;
};

struct WrappedFields {
    struct { int x; } a[2], b[2];
    union { int x; double y; } *c, *d;
};

union WrappedMembers {
    struct { int x; } a[2], b[2];
    union { int x; double y; } *c, *d;
};

struct DistinctDeclarations {
    struct { int x; } a;
    struct { int x; } b;
};

union PromotedFields {
    struct {
        struct { int x; } a, b;
    };
};

struct AliasFirst {
    struct { int x; } a, b, c;
};

struct AliasMiddle {
    struct { int x; } a, b, c;
};

union PromotedUnionFields {
    int outer;
    union {
        int x;
        struct { int x; } a, b;
    };
    int tail;
};

union NestedPromotedUnionFields {
    int outer;
    union {
        int x;
        union {
            int y;
            struct { int value; } a, b;
        };
        int afterW;
    };
    int tail;
};
