#include <assert.h>
#include <complex.h>
#include <ctype.h>
#include <errno.h>
#include <fenv.h>
#include <float.h>
#include <inttypes.h>
#include <iso646.h>
#include <limits.h>
#include <locale.h>
#include <math.h>
#include <setjmp.h>
#include <signal.h>
#include <stdalign.h>
#include <stdarg.h>
#include <stdatomic.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdnoreturn.h>
#include <string.h>
#include <tgmath.h>
// #include <threads.h> mac not supported
#include <time.h>
#include <wchar.h>
#include <wctype.h>

// https://en.cppreference.com/w/c/header
struct fenv {
    fenv_t t1;
    fexcept_t t2;
};

struct stddef {
    size_t t1;
    ptrdiff_t t2;
    nullptr_t t3;
    max_align_t t4;
};

struct stdint {
    int8_t t1;
    int16_t t2;
    int32_t t3;
    int64_t t4;
    int_fast8_t t5;
    int_fast16_t t6;
    int_fast32_t t7;
    int_fast64_t t8;
    int_least8_t t9;
    int_least16_t t10;
    int_least32_t t11;
    int_least64_t t12;
    intmax_t t13;
    intptr_t t14;
    uint8_t t15;
    uint16_t t16;
    uint32_t t17;
    uint64_t t18;
    uint_fast8_t t19;
    uint_fast16_t t20;
    uint_fast32_t t21;
    uint_fast64_t t22;
    uint_least8_t t23;
    uint_least16_t t24;
    uint_least32_t t25;
    uint_least64_t t26;
    uintmax_t t27;
    uintptr_t t28;
};

struct stdio {
    FILE *t1;
    fpos_t t2;
};

struct time {
    tm t1;
    time_t t2;
    clock_t t3;
    timespec t4;
};

struct wchar {
    wint_t t2;
    wctrans_t t3;
    wctype_t t4;
};