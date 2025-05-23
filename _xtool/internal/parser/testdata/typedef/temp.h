typedef int INT;

typedef INT STANDARD_INT;

typedef int NewINT, *NewIntPtr, NewIntArr[];

typedef int (*Foo1)(int, int, ...);

typedef int (*Foo2)(int, int), (*Bar)(void *, void *);

namespace A {
typedef class Foo {
    int x;
} MyClass, *MyClassPtr, MyClassArray[];
} // namespace A

typedef struct {
    int x;
} MyStruct1;

typedef union {
    int x;
} MyUnion1;
typedef enum { MyEnum1RED, MyEnum1GREEN, MyEnum1BLUE } MyEnum1;

typedef struct {
    int x;
} MyStruct2, MyStruct3, *StructPtr, StructArr[];

typedef enum { MyEnum2RED, MyEnum2GREEN, MyEnum2BLUE } MyEnum2, MyEnum3, *EnumPtr, EnumArr[];

namespace A {
namespace B {
typedef struct {
    int x;
} MyStruct, MyStruct2, *StructPtr, StructArr[];
} // namespace B
} // namespace A

typedef enum algorithm { AlgorithmA, AlgorithmB } algorithm_t;
typedef algorithm_t algorithm;
