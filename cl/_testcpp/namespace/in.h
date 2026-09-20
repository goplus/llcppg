namespace bar {
	namespace detail {
		unsigned f(int a);
		inline void f() {}
	}
}

namespace bar {
	class base
	{
	public:
		~base() {}
	};

	inline void print(base b) {}

	typedef char boolean;
}
