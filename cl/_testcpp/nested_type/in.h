class Outer {
public:
	typedef int value_type;

	struct Inner {
		int x;
	};

	class Detail {
		int secret;
	};

	Inner inner;
};
