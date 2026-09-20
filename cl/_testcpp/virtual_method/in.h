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
