enum Global : int;

enum Global : int {
	GA,
	GB
};

enum class Local : int;

enum class Local : int {
	LA,
	LB
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
