struct StructValuePointerFields {
    struct { int value; } a, *b;
};

struct StructPointerValueFields {
    struct { int value; } *a, b;
};

struct StructPointerDepthFields {
    struct { int value; } *a, **b;
};

struct StructValueArrayFields {
    struct { int value; } a, b[2];
};

struct StructArrayExtentFields {
    struct { int value; } a[2], b[3];
};

struct StructPointerArrayFields {
    struct { int value; } *a[2], (*b)[2];
};

struct StructPointerConstFields {
    struct { int value; } *const a, *b;
};

struct StructAlternatingFields {
    struct { int value; } a, *b, c;
};

struct UnionValuePointerFields {
    union { int integer; double decimal; } a, *b;
};

struct UnionPointerValueFields {
    union { int integer; double decimal; } *a, b;
};

struct UnionPointerDepthFields {
    union { int integer; double decimal; } *a, **b;
};

struct UnionValueArrayFields {
    union { int integer; double decimal; } a, b[2];
};

struct UnionArrayExtentFields {
    union { int integer; double decimal; } a[2], b[3];
};

struct UnionPointerArrayFields {
    union { int integer; double decimal; } *a[2], (*b)[2];
};

struct UnionPointerConstFields {
    union { int integer; double decimal; } *const a, *b;
};

struct UnionAlternatingFields {
    union { int integer; double decimal; } a, *b, c;
};

union StructValuePointerMembers {
    struct { int value; } a, *b;
};

union StructPointerValueMembers {
    struct { int value; } *a, b;
};

union StructPointerDepthMembers {
    struct { int value; } *a, **b;
};

union StructValueArrayMembers {
    struct { int value; } a, b[2];
};

union StructArrayExtentMembers {
    struct { int value; } a[2], b[3];
};

union StructPointerArrayMembers {
    struct { int value; } *a[2], (*b)[2];
};

union StructPointerConstMembers {
    struct { int value; } *const a, *b;
};

union StructAlternatingMembers {
    struct { int value; } a, *b, c;
};

union UnionValuePointerMembers {
    union { int integer; double decimal; } a, *b;
};

union UnionPointerValueMembers {
    union { int integer; double decimal; } *a, b;
};

union UnionPointerDepthMembers {
    union { int integer; double decimal; } *a, **b;
};

union UnionValueArrayMembers {
    union { int integer; double decimal; } a, b[2];
};

union UnionArrayExtentMembers {
    union { int integer; double decimal; } a[2], b[3];
};

union UnionPointerArrayMembers {
    union { int integer; double decimal; } *a[2], (*b)[2];
};

union UnionPointerConstMembers {
    union { int integer; double decimal; } *const a, *b;
};

union UnionAlternatingMembers {
    union { int integer; double decimal; } a, *b, c;
};
