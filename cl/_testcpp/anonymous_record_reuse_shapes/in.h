struct StructPointerFields {
    struct { int value; } *a, *b;
};

struct StructDoublePointerFields {
    struct { int value; } **a, **b;
};

struct StructArrayFields {
    struct { int value; } a[2], b[1 + 1];
};

struct StructMatrixFields {
    struct { int value; } a[2][3], b[2][3];
};

struct StructPointerArrayFields {
    struct { int value; } *a[2], *b[2];
};

struct StructArrayPointerFields {
    struct { int value; } (*a)[2], (*b)[2];
};

struct StructConstFields {
    const struct { int value; } a, b;
};

struct UnionDoublePointerFields {
    union { int integer; double decimal; } **a, **b;
};

struct UnionArrayFields {
    union { int integer; double decimal; } a[2], b[1 + 1];
};

struct UnionMatrixFields {
    union { int integer; double decimal; } a[2][3], b[2][3];
};

struct UnionPointerArrayFields {
    union { int integer; double decimal; } *a[2], *b[2];
};

struct UnionArrayPointerFields {
    union { int integer; double decimal; } (*a)[2], (*b)[2];
};

struct UnionConstFields {
    const union { int integer; double decimal; } a, b;
};

struct CallbackRecordFields {
    struct { void (*call)(int); } a, b;
};

union StructPointerMembers {
    struct { int value; } *a, *b;
};

union StructDoublePointerMembers {
    struct { int value; } **a, **b;
};

union StructArrayMembers {
    struct { int value; } a[2], b[1 + 1];
};

union StructMatrixMembers {
    struct { int value; } a[2][3], b[2][3];
};

union StructPointerArrayMembers {
    struct { int value; } *a[2], *b[2];
};

union StructArrayPointerMembers {
    struct { int value; } (*a)[2], (*b)[2];
};

union StructConstMembers {
    const struct { int value; } a, b;
};

union UnionDoublePointerMembers {
    union { int integer; double decimal; } **a, **b;
};

union UnionArrayMembers {
    union { int integer; double decimal; } a[2], b[1 + 1];
};

union UnionMatrixMembers {
    union { int integer; double decimal; } a[2][3], b[2][3];
};

union UnionPointerArrayMembers {
    union { int integer; double decimal; } *a[2], *b[2];
};

union UnionArrayPointerMembers {
    union { int integer; double decimal; } (*a)[2], (*b)[2];
};

union UnionConstMembers {
    const union { int integer; double decimal; } a, b;
};

union CallbackRecordMembers {
    struct { void (*call)(int); } a, b;
};
