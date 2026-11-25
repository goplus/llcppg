struct Tag {
    int tag; // The non-polymorphic base's field.
};

struct Shape {
    /// The shape identifier.
    int id; // After the implicit vptr.
    virtual int size();
};

struct Derived : Tag, Shape {
    int value; // After the reordered base classes.
};

template <class T>
struct Box {
    /// The template's stored value.
    T value; // An inline template field comment.
};

class Access {
    int hidden; // A private field.
public:
    int visible; // A public field.
};

union Value {
    int number; // A C++ union member.
    struct {
        int first; // A promoted C++ struct member.
        int second; // A promoted C++ struct member at a nonzero offset.
    };
};
