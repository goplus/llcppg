#include <foo.h>

extern "C" void _llcppg__ZN3barC1Ev(bar* this) {
	this->bar();
}

extern "C" void _llcppg__ZN3barC1EPKc(bar* this, const char *a) {
	this->bar(a);
}

extern "C" void _llcppg__ZN3barD1Ev(bar* this) {
	this->~bar();
}
