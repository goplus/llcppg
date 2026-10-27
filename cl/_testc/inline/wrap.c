#include <foo.h>

int _llcppg_add(int a, int b) {
	return add(a, b);
}

int _llcppg_mul(int a, int b) {
	return mul(a, b);
}

int _llcppg_foo(int (*f)(int , int )) {
	return foo(f);
}
