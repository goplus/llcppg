struct {
    int a;
};

struct Foo1 {
    int a;
    int b;
};
struct Foo2 {
    int a, b;
};

struct Foo3 {
    int a;
    int (*Foo)(int, int);
};

struct Person {
    int age;
    struct {
        int year;
        int day;
        int month;
    } birthday;
};
