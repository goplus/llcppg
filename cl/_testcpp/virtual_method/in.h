// Case 1: no base class; the class itself has virtual methods.
class Shape
{
public:
	int id;
	virtual int area();
};

// A plain base class without any virtual methods.
class Tag
{
public:
	int tag;
};

// Case 2: has a base class, but the base has no virtual methods; the class
// itself has virtual methods.
class Widget : public Tag
{
public:
	int w;
	virtual int paint();
};

// Case 3: single base class, and the base class has virtual methods. The
// derived class reuses the base's vptr, so it introduces no new vptr.
class Circle : public Shape
{
public:
	int r;
	int area();
};

// Case 4: multiple base classes; an earlier base has no virtual methods, while
// a later one does. The polymorphic base (Shape) is the primary base, so the
// class reuses its vptr and introduces no new one; the primary base is laid out
// first (before the non-polymorphic Tag) to keep the shared vptr at offset 0.
class Button : public Tag, public Shape
{
public:
	int b;
	virtual int click();
};

// Case 5: a virtual destructor occupies two consecutive vtable slots (the
// complete-object and deleting destructors, per the Itanium ABI). They are kept
// as reserved placeholders so the following virtual method (Draw) lands at its
// real slot index (2) rather than 0. Stream, whose only virtual member is the
// destructor, keeps its vptr field but gets no typed vtable.
class Stream
{
public:
	int fd;
	virtual ~Stream();
};

class Canvas
{
public:
	virtual ~Canvas();
	virtual int draw();
};

// Case 6: a non-public virtual method still occupies a vtable slot, but it is
// reserved as an unexported placeholder rather than exposed. The private step()
// keeps stop() at its real slot index (2) without leaking a private API.
class Machine
{
public:
	virtual int start();
private:
	virtual int step();
public:
	virtual int stop();
};
