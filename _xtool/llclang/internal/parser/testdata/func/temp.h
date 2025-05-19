void foo1();
void foo2(int a);
void foo3(int a, ...);
float *foo4(int a, double b);
static inline int foo5(int a, int b);

typedef void(fntype1)();
fntype1 bar1;

typedef long(fntype2)(long a);
typedef fntype2 fntype3;
fntype3 bar2;

typedef struct OSSL_CORE_HANDLE OSSL_CORE_HANDLE;
typedef struct OSSL_DISPATCH OSSL_DISPATCH;
typedef int(OSSL_provider_init_fn)(const OSSL_CORE_HANDLE *handle, const OSSL_DISPATCH *in, const OSSL_DISPATCH **out,
                                   void **provctx);
OSSL_provider_init_fn OSSL_provider_init;
