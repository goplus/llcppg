// A callback with unspecified parameters.
typedef int (*Callback)();
typedef Callback CallbackAlias;
typedef int Function();
typedef Function *FunctionPointer;
typedef void (*Visitor)(int (*callback)());

// The holder's documentation stays on its type.
struct Holder {
    int before;
    int (*callback)();
    // Keep both comments on this field.
    int (*documented)(); // The original inline comment.
    Callback named;
    CallbackAlias alias;
    Function *function;
    FunctionPointer function_pointer;
    int (**pointer)();
    int (*callbacks[2])();
    int (*const fixed)();
    int (*first)(), (*second)();
    int after;
};

struct Nested {
    struct {
        int (*callback)();
        int value;
    } entries[2]; // The record array's own comment.
    int after;
};

union Value {
    int (*callback)();
    int number;
    struct {
        int (*nested)(); // The nested callback's comment.
    };
};

typedef int CustomFunction();
typedef int (*CustomCallback)();
struct Overrides {
    CustomFunction *function;
    CustomCallback callback;
    int (*renamed)();
};

struct Prototypes {
    int (*empty)(void);
    int (*typed)(int, double);
    int (*variadic)(int, ...);
};

// The active callback.
extern int (*global_callback)();
void accept_callback(int (*callback)());
int (*get_callback(void))();
void accept_visitor(void (*visitor)(int (*callback)()));
