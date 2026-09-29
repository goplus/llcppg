template <typename T>
class Foo {
private:
	T value;

public:
	Foo(T val) : value(val) {}
	~Foo() {}

	void g(T val) {}
};
