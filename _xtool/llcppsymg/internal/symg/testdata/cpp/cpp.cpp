#include "cpp.h"

Foo::Foo(const char *filename) {}
Foo::Foo(const char *buffer, long buffer_size) {}
Foo::Foo::~Foo() {}
int Foo::ParseBar() const {}
const char *Foo::Get(const char *section, const char *name,
                     const char *default_value) const {}

const char *Foo::MakeBar(const char *section, const char *name) {}
