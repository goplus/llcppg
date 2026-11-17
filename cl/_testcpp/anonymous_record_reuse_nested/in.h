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
