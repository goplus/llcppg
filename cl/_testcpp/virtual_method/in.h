// Case 1: no base class; the class itself has virtual methods.
//
// It also declares a virtual destructor and a private virtual method to
// exercise the two extra vtable slot rules:
//   - a virtual destructor takes two consecutive Itanium slots, XGo_dtor
//     (complete object) then XGo_dtor_deleting, both func(this *Shape);
//   - a private virtual method still occupies a slot, emitted as an unexported
//     placeholder so later slots keep their index.
// Declaration order fixes the vtable order: area, ~Shape (two slots), reset.
class Shape
{
public:
	int id;
	virtual int area();
	virtual ~Shape();

private:
	virtual void reset();
};

// A plain base class without any virtual methods.
class Tag
{
public:
	int tag;
};

// Case 2: has a base class, but the base has no virtual methods; the class
// itself has virtual methods. It also declares a virtual destructor, so its own
// vtable holds paint plus the two destructor slots.
class Widget : public Tag
{
public:
	int w;
	virtual int paint();
	virtual ~Widget();
};

// Case 3: single base class, and the base class has virtual methods. The
// derived class reuses the base's vptr, so it introduces no new vptr. It does
// not declare its own destructor, so Shape's destructor and private-method slots
// are inherited unchanged (with this re-typed to *Circle).
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
// Shape's inherited slots (area, the two destructor slots, the private-method
// placeholder) precede the new click slot.
class Button : public Tag, public Shape
{
public:
	int b;
	virtual int click();
};
