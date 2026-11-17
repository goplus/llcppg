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
