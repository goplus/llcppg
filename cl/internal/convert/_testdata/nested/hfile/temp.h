#include <stddef.h>
struct struct1
{
    char *b;
    size_t n;
    union
    {
        long l;
        char b[60];
    } init;
};


// https://github.com/goplus/llcppg/issues/514
// named nested struct
struct struct_with_nested {
    struct inner_struct {
        long l;
    } init;
};

struct struct2
{
    char *b;
    size_t size;
    size_t n;
    struct
    {
        long l;
        char b[60];
        struct1 rec;
    } init;
};

union union1
{
    char *b;
    size_t size;
    size_t n;
    struct
    {
        long l;
        char b[60];
        struct2 rec;
    } init;
};

union union2
{
    char *b;
    size_t size;
    size_t n;
    union
    {
        long l;
        char b[60];
        struct2 rec;
    } init;
};


// https://github.com/goplus/llcppg/issues/514
struct a {
    struct b {
        struct c {
            int a;
        } c_field;
        struct d {
            int b;
        } d_field;
    } b_field;
    struct e {
        struct f {
            int b;
        } f_field;
    } e_field;
};
