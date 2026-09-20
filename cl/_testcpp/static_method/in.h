class Bar
{
public:
	unsigned f(int a);
	static int create();
	static int create(int a);
};

inline static int Bar::create() {
	reurn 0;
}
