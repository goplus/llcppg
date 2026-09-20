enum Global {
	GA,
	GB
};

namespace bar {
	enum Color {
		Red,
		Green = 5,
		Blue
	};
}

class Shape {
public:
	enum Kind {
		Circle,
		Square,
		Triangle
	};
};
