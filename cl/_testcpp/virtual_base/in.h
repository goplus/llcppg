// Virtual base classes, following the C++ Itanium ABI. A class that virtually
// derives from a base carries a vptr (needed to locate the shared base) and lays
// the virtual base subobject out once, after all of its own non-virtual data.
//
// This is the building block of C++'s iostream hierarchy: basic_istream and
// basic_ostream each *virtually* derive from basic_ios so that, when they are
// combined, a single shared basic_ios subobject results. The cases below model
// istream / ostream themselves (single virtual inheritance).

// Case 1: a plain (non-polymorphic) virtual base and a class that virtually
// derives from it. In has a vptr (to locate Ios) even though neither declares a
// virtual method; Ios is laid out at In's tail. Layout: vptr, igcount, Ios.
struct Ios
{
	int state;
};

struct In : virtual Ios
{
	int igcount;
};

// Case 2: a second class virtually deriving from the same base, mirroring the
// istream/ostream pair. Layout: vptr, oputcount, Ios.
struct Out : virtual Ios
{
	int oputcount;
};

// Case 3: an iostream-like polymorphic virtual base. IosBase has a virtual
// destructor, so it owns a vptr and a vtable (the destructor's two Itanium
// slots). IStream virtually derives from it and is itself polymorphic (a virtual
// base), so it carries its own vptr and lays IosBase out at its tail. IStream
// declares no virtual method of its own, so it gets no vtable struct. Layout:
// vptr, gpos, IosBase.
class IosBase
{
public:
	int flags;
	virtual ~IosBase();
};

class IStream : public virtual IosBase
{
public:
	int gpos;
};

// Case 4: a second polymorphic virtual-inheritance class (the ostream side).
// Layout: vptr, ppos, IosBase.
class OStream : public virtual IosBase
{
public:
	int ppos;
};
