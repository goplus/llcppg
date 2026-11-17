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

union StructPromotedUnionMembers {
    int before;
    struct {
        union { int integer; double decimal; } a, b;
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

union UnionPromotedUnionMembers {
    int before;
    union {
        union { int integer; double decimal; } a, b;
    };
    int after;
};

// C++ anonymous records cannot declare enums: use enum E { V }; union { E a, b; }.

enum StructPromotedEnumFieldsType { StructPromotedEnumFieldsValue = 1 };

struct StructPromotedEnumFields {
    int before;
    struct {
        StructPromotedEnumFieldsType a, b;
    };
    int after;
};

enum UnionPromotedEnumFieldsType { UnionPromotedEnumFieldsValue = 1 };

struct UnionPromotedEnumFields {
    int before;
    union {
        UnionPromotedEnumFieldsType a, b;
    };
    int after;
};

enum StructPromotedEnumMembersType { StructPromotedEnumMembersValue = 1 };

union StructPromotedEnumMembers {
    int before;
    struct {
        StructPromotedEnumMembersType a, b;
    };
    int after;
};

enum UnionPromotedEnumMembersType { UnionPromotedEnumMembersValue = 1 };

union UnionPromotedEnumMembers {
    int before;
    union {
        UnionPromotedEnumMembersType a, b;
    };
    int after;
};
