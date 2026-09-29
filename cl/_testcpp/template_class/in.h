template <typename T>
struct Foo {
	T value;

	Foo(T val) : value(val) {}
	~Foo() {}

	void g(T val) {}
};

template <typename T1, typename T2>
class Bar {
private:
	T1  v1;
	T2  v2;
	int v3;

public:
	Bar(T1 val1, T2 val2, int val3) : v1(val1), v2(val2), v3(val3) {}
	virtual ~Bar() {}
	virtual T2 g(T1 val1, T2 val2, int val3) {
		return val2;
	}
};
