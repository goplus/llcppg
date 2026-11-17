struct AliasFirst {
    struct { int x; } a, b, c;
};

struct AliasMiddle {
    struct { int x; } a, b, c;
};

struct TaggedStruct { int value; };
union TaggedUnion { int integer; double decimal; };
enum TaggedEnum { TaggedEnumValue = 1 };
typedef struct { int value; } StructAlias;
typedef union { int integer; double decimal; } UnionAlias;
typedef enum { EnumAliasValue = 1 } EnumAlias;
typedef StructAlias StructAlias2;
typedef UnionAlias UnionAlias2;
typedef EnumAlias EnumAlias2;

struct StructAliasFields {
    struct TaggedStruct named_a, *named_b;
    StructAlias alias_a;
    StructAlias2 alias_b;
    StructAlias *pointer_a, *pointer_b;
    StructAlias array_a[2], array_b[3];
};

struct UnionAliasFields {
    union TaggedUnion named_a, *named_b;
    UnionAlias alias_a;
    UnionAlias2 alias_b;
    UnionAlias *pointer_a, *pointer_b;
    UnionAlias array_a[2], array_b[3];
};

struct EnumAliasFields {
    enum TaggedEnum named_a, *named_b;
    EnumAlias alias_a;
    EnumAlias2 alias_b;
    EnumAlias *pointer_a, *pointer_b;
    EnumAlias array_a[2], array_b[3];
};

union StructAliasMembers {
    struct TaggedStruct named_a, *named_b;
    StructAlias alias_a;
    StructAlias2 alias_b;
    StructAlias *pointer_a, *pointer_b;
    StructAlias array_a[2], array_b[3];
};

union UnionAliasMembers {
    union TaggedUnion named_a, *named_b;
    UnionAlias alias_a;
    UnionAlias2 alias_b;
    UnionAlias *pointer_a, *pointer_b;
    UnionAlias array_a[2], array_b[3];
};

union EnumAliasMembers {
    enum TaggedEnum named_a, *named_b;
    EnumAlias alias_a;
    EnumAlias2 alias_b;
    EnumAlias *pointer_a, *pointer_b;
    EnumAlias array_a[2], array_b[3];
};

struct InlineStructFields {
    struct InlineStructFieldsType { int value; } a, *b, c[2];
};

struct InlineUnionFields {
    union InlineUnionFieldsType { int integer; double decimal; } a, *b, c[2];
};

struct InlineEnumFields {
    enum InlineEnumFieldsType { InlineEnumFieldsTypeValue = 1 } a, *b, c[2];
};

union InlineStructMembers {
    struct InlineStructMembersType { int value; } a, *b, c[2];
};

union InlineUnionMembers {
    union InlineUnionMembersType { int integer; double decimal; } a, *b, c[2];
};

union InlineEnumMembers {
    enum InlineEnumMembersType { InlineEnumMembersTypeValue = 1 } a, *b, c[2];
};

typedef struct { int value; } StructValue, *StructPointer, StructArray[2];

struct StructTypedefFields {
    StructValue a;
    StructPointer b;
    StructArray c;
};

union StructTypedefMembers {
    StructValue a;
    StructPointer b;
    StructArray c;
};

typedef union { int integer; double decimal; } UnionValue, *UnionPointer, UnionArray[2];

struct UnionTypedefFields {
    UnionValue a;
    UnionPointer b;
    UnionArray c;
};

union UnionTypedefMembers {
    UnionValue a;
    UnionPointer b;
    UnionArray c;
};

typedef enum { EnumValueValue = 1 } EnumValue, *EnumPointer, EnumArray[2];

struct EnumTypedefFields {
    EnumValue a;
    EnumPointer b;
    EnumArray c;
};

union EnumTypedefMembers {
    EnumValue a;
    EnumPointer b;
    EnumArray c;
};

struct ScopedAliasFields {
    using Value = struct { int value; };
    using Pointer = Value *;
    Value a;
    Pointer b;
    Value c[2];
};

union ScopedAliasMembers {
    using Value = struct { int value; };
    using Pointer = Value *;
    Value a;
    Pointer b;
    Value c[2];
};

struct ScopedTypedefFields {
    typedef struct { int value; } Value, *Pointer;
    Value a;
    Pointer b;
};

union ScopedTypedefMembers {
    typedef struct { int value; } Value, *Pointer;
    Value a;
    Pointer b;
};

struct ScopedEnumFields {
    enum class Kind : unsigned { Value = 1 } a, b;
};

union ScopedEnumMembers {
    enum class Kind : unsigned { Value = 1 } a, b;
};

namespace left { struct Value { int value; }; }
namespace right { struct Value { int value; }; }
struct NamespaceFields {
    left::Value a;
    right::Value b;
};
