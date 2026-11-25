typedef int (*Callback)();

struct Holder {
    int (*callback)();
    Callback named;
    int (*empty)(void);
    int (*typed)(int, double);
    int (*variadic)(int, ...);
};
