class Bar
{
public:
	unsigned f(int a);
	static int create();
	static int create(int a);
};

inline int Bar::create() {
	return 0;
}
