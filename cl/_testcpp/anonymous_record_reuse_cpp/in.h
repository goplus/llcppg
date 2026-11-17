struct ReferenceFields {
    struct { int value; } &a, &b;
};

struct RvalueReferenceFields {
    struct { int value; } &&a, &&b;
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
