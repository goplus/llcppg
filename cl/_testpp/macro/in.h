#define FOO_H

#define f()  1
#define g(x) x

#define STRVAL  "abc"

#define TRUE  1
#define FALSE 0

#define THREE (TRUE + 2)
#define FVAL  (TRUE * 3.14 + 0.0)

#define MASK  (~FALSE)
#define MASK2 (~FALSE)
#define TWO   (THREE * TRUE + -TRUE)

#define THREE 3

#define QUO_ZERO (1 / 0)
#define REM_ZERO (1 % 0)
#define MATH_OP  (1 + 2 % 4 * 3 / 1 - 4.0 / 2)
#define BIT_OP   ((1 << 2 >> 2) ^ 2 | 0 & 1)

#define ERR_BIT_OP (2.0 & 1.0)
#define ERR_TILDE  ~ERR_BIT_OP

#define PANIC_BITOP (1 << -1)

#define ERROR    ('*' % 3.24)
#define QUESTION 1 ? 2
#define TILDE    1 ~ 2

typedef char BOOL, *PBOOL;
