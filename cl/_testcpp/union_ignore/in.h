union UnsupportedMember {
    int value;
    int vector __attribute__((vector_size(8)));
};

struct UsesUnsupportedMember {
    union UnsupportedMember value;
};

struct AnonymousUnionField {
    union {
        int vector __attribute__((vector_size(8)));
    } value;
};

struct PromotedUnion {
    union {
        int vector __attribute__((vector_size(8)));
    };
};

union NestedUnionField {
    union {
        int vector __attribute__((vector_size(8)));
    } value;
};

union Supported {
    int value;
};
