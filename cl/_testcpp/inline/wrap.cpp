#include <foo.h>

extern "C" unsigned int _llcppg__ZN3bar1fEi(bar* this, int a) {
	return this->f(a);
}

extern "C" void _llcppg__ZN3bar2_gEv(bar* this) {
	this->_g();
}
