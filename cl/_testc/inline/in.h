inline int add(int a, int b) {
    return a + b;
}

inline int mul(int a, int b) {
    return a * b;
}

inline int foo(int (*f)(int, int)) {
    return f(2, 3);
}

void g(void);
