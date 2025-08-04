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
struct struct_with_nested
{
    struct inner_struct
    {
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
struct a
{
    struct b
    {
        struct c
        {
            int a;
        } c_field;
        struct d
        {
            int b;
        } d_field;
    } b_field;
    struct e
    {
        struct f
        {
            int b;
        } f_field;
    } e_field;
};
struct NestedEnum
{
    enum
    {
        APR_BUCKET_DATA1 = 0,
        APR_BUCKET_METADATA2 = 1
    } is_metadata1;

    struct a
    {
        int b;
    } a_t;
};

struct NestedEnum2
{
    enum
    {
        APR_BUCKET_DATA3 = 0,
        APR_BUCKET_METADATA4 = 1
    };
};

struct NestedEnum3
{
    enum is_metadata3
    {
        APR_BUCKET_DATA5 = 0,
        APR_BUCKET_METADATA6 = 1
    };
};

struct NestedEnum4
{
    enum is_metadata4
    {
        APR_BUCKET_DATA7 = 0,
        APR_BUCKET_METADATA8 = 1
    } key;
};

enum OuterEnum
{
    APR_BUCKET_DATA9 = 0,
    APR_BUCKET_METADATA10 = 1
};
struct Enum
{
    enum OuterEnum k;
};
