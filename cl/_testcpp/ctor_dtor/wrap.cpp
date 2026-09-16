#include <foo.h>

void _llcppg__ZN3barC1Ev(bar* this) {
	this->bar();
}

void _llcppg__ZN3barC1EPKc(bar* this, const char * a) {
	this->bar(a);
}

void _llcppg__ZN3barD1Ev(bar* this) {
	this->~bar();
}
