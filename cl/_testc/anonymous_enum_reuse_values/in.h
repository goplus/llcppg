struct EnumValueFields {
    enum { EnumValueFieldsValue = 1 } a, b, c;
};

struct EnumConstFields {
    const enum { EnumConstFieldsValue = 1 } a, b;
};

union EnumValueMembers {
    enum { EnumValueMembersValue = 1 } a, b, c;
};

union EnumConstMembers {
    const enum { EnumConstMembersValue = 1 } a, b;
};
