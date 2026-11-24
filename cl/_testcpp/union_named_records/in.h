typedef union {
    struct mutex_state {
        int lock;
    } data;
    long align;
} Mutex;

union NamedMembers {
    struct NamedStruct {
        int value;
        struct NamedStruct *next;
    } record;
    struct NamedStruct *pointer;
    struct NamedStruct array[2];
    union NamedUnion {
        int number;
        double real;
    } value;
    union NamedUnion *union_pointer;
};

union TypeDeclarations {
    static int count;
    struct Entry {
        int value;
    };
    class Detail {
    public:
        short code;
    };
    union Number {
        int number;
        double real;
    };
    Entry entry;
    Detail detail;
    Number number;
};
