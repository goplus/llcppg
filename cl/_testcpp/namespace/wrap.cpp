#include <foo.h>

extern "C" void _llcppg__ZN3bar6detail1fEv() {
	bar::detail::f();
}

extern "C" void _llcppg__ZN3bar4baseD1Ev(bar::base* this) {
	this->~base();
}

extern "C" void _llcppg__ZN3bar5printENS_4baseE(bar::base b) {
	bar::print(b);
}
