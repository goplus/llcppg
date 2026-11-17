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
