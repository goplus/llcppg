class base {
public:
	~base();
};

class bar : public base
{
public:
	bar(int a);
	bar();
	bar(const char* a) {}
	~bar() {}
};

inline bar::bar() {
}
