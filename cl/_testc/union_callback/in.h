typedef int (*Callback)(int);

union DirectCallback {
    int (*callback)(int);
    void (*notify)(int);
    Callback named;
    int value;
};

struct EmbeddedCallback {
    union {
        int (*callback)(int);
        int value;
    };
};

union PromotedStructCallback {
    struct {
        int tag;
        int (*callback)(int);
    };
};

union NestedUnionCallback {
    struct {
        int tag;
        union {
            int (*callback)(int);
            int value;
        };
    };
};
