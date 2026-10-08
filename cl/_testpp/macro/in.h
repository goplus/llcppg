#define FOO_H

#define f()  1
#define g(x) x

#define STRVAL  "abc"

#define TRUE  1
#define FALSE 0

#define THREE (TRUE + 2)
#define FVAL  (TRUE * 3.14)

#define MASK  (~FALSE)
#define MASK2 (~FALSE)
#define TWO   (THREE * TRUE + -TRUE)

#define THREE 3

#define DIV_ZERO (1 / 0)
#define MATH_OP  (1 + 2 * 3 - 4.0 / 2)
#define BIT_OP   ((1 << 2 >> 2) ^ 2 | 0)

#define ERROR ('*' % 3.24)

typedef char BOOL, *PBOOL;
