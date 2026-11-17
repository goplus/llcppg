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

struct StructMixedPointerArrayFields {
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

struct UnionMixedPointerArrayFields {
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

union StructMixedPointerArrayMembers {
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

union UnionMixedPointerArrayMembers {
    union { int integer; double decimal; } *a[2], (*b)[2];
};

union UnionPointerConstMembers {
    union { int integer; double decimal; } *const a, *b;
};

union UnionAlternatingMembers {
    union { int integer; double decimal; } a, *b, c;
};

struct ReferenceFields {
    struct { int value; } &a, &b;
};

struct RvalueReferenceFields {
    struct { int value; } &&a, &&b;
};

struct MixedReferenceFields {
    struct { int value; } &a, *b;
};
