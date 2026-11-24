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

union PromotedStruct {
    struct {
        int tag;
        struct ViewState {
            short first;
            short second;
        } view;
        union ViewValue {
            int code;
            double real;
        } item;
    };
    long align;
};

union PromotedUnion {
    union {
        struct PromotedState {
            int code;
        } value;
    };
    long align;
};
