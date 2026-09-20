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

// Case 4: multiple base classes; the first base has no virtual methods, while
// another base does. The class introduces its own vptr (the first base is the
// primary base and is non-polymorphic).
class Button : public Tag, public Shape
{
public:
	int b;
	virtual int click();
};
