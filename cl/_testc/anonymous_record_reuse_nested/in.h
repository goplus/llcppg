union PromotedFields {
    struct {
        struct { int x; } a, b;
    };
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

struct StructNamedStructFields {
    int before;
    struct {
        struct { int value; } a, b;
    } nested;
    int after;
};

struct StructNamedUnionFields {
    int before;
    struct {
        union { int integer; double decimal; } a, b;
    } nested;
    int after;
};

struct StructNamedEnumFields {
    int before;
    struct {
        enum { StructNamedEnumFieldsValue = 1 } a, b;
    } nested;
    int after;
};

struct StructPromotedStructFields {
    int before;
    struct {
        struct { int value; } a, b;
    };
    int after;
};

struct StructPromotedUnionFields {
    int before;
    struct {
        union { int integer; double decimal; } a, b;
    };
    int after;
};

struct StructPromotedEnumFields {
    int before;
    struct {
        enum { StructPromotedEnumFieldsValue = 1 } a, b;
    };
    int after;
};

struct UnionNamedStructFields {
    int before;
    union {
        struct { int value; } a, b;
    } nested;
    int after;
};

struct UnionNamedUnionFields {
    int before;
    union {
        union { int integer; double decimal; } a, b;
    } nested;
    int after;
};

struct UnionNamedEnumFields {
    int before;
    union {
        enum { UnionNamedEnumFieldsValue = 1 } a, b;
    } nested;
    int after;
};

struct UnionPromotedStructFields {
    int before;
    union {
        struct { int value; } a, b;
    };
    int after;
};

struct UnionPromotedUnionFields {
    int before;
    union {
        union { int integer; double decimal; } a, b;
    };
    int after;
};

struct UnionPromotedEnumFields {
    int before;
    union {
        enum { UnionPromotedEnumFieldsValue = 1 } a, b;
    };
    int after;
};

union StructNamedStructMembers {
    int before;
    struct {
        struct { int value; } a, b;
    } nested;
    int after;
};

union StructNamedUnionMembers {
    int before;
    struct {
        union { int integer; double decimal; } a, b;
    } nested;
    int after;
};

union StructNamedEnumMembers {
    int before;
    struct {
        enum { StructNamedEnumMembersValue = 1 } a, b;
    } nested;
    int after;
};

union StructPromotedUnionMembers {
    int before;
    struct {
        union { int integer; double decimal; } a, b;
    };
    int after;
};

union StructPromotedEnumMembers {
    int before;
    struct {
        enum { StructPromotedEnumMembersValue = 1 } a, b;
    };
    int after;
};

union UnionNamedStructMembers {
    int before;
    union {
        struct { int value; } a, b;
    } nested;
    int after;
};

union UnionNamedUnionMembers {
    int before;
    union {
        union { int integer; double decimal; } a, b;
    } nested;
    int after;
};

union UnionNamedEnumMembers {
    int before;
    union {
        enum { UnionNamedEnumMembersValue = 1 } a, b;
    } nested;
    int after;
};

union UnionPromotedUnionMembers {
    int before;
    union {
        union { int integer; double decimal; } a, b;
    };
    int after;
};

union UnionPromotedEnumMembers {
    int before;
    union {
        enum { UnionPromotedEnumMembersValue = 1 } a, b;
    };
    int after;
};
