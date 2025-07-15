#include <stddef.h>
typedef struct Foo
{
    struct Foo *next;
} Foo;
typedef struct Foo2
{
    struct Foo2 *next;
} Foo2;
// remove prefix Foo_ & can be a method of Foo (*Foo).Delete
char *Foo_Print(const Foo *item);
// config not be a method in llcppg.cfg/symMap
void Foo_Delete(Foo *item);
// normal function no be a method
Foo *Foo_ParseWithLength(const char *value, size_t buffer_length);
// only can be a normal function but config be a method,keep output as function
Foo *Foo_ParseWithSize(const char *value, size_t buffer_length);

Foo *Foo_ignoreFunc();

// config Foo_ForBar to Bar,so Foo_Bar to Bar__1
void Foo_Bar();
void Foo_ForBar();

void Foo_Prefix_BarMethod(Foo *item); // to BarMethod,but follow config the BarMethod,so it need add prefix
void Foo_BarMethod(Foo *item);        // config BarMethod
void Foo_ForBarMethod(Foo *item);     // config BarMethod,so it need add suffix

// first receiver Foo's method,with name 'Parse'
void Foo_ReceiverParse(Foo *item);
// first function with name 'Parse'
void Foo_FunctionParse();
// second receiver Foo's method,with name 'Parse', and the next method name is the same,so we need add suffix
void Foo_ReceiverParse2(Foo *item);
// not same receiver,but same function name,we don't need add suffix
void Foo_Receiver2Parse2(Foo2 *item);

void Foo_Valist(Foo2 *item, ...);
