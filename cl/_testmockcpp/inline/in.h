class bar
{
public:
	unsigned f(int a);
	inline void _g() {}
	signed int xprintf(const char* fmt, ...);

private:
	void privateMethod();
};

inline unsigned bar::f(int a) {
	return 0;
}
