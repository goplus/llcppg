struct Base
{
	int a;
	unsigned b;
};

struct Derived : public Base
{
	int c;
};

struct Mid
{
	int m;
};

class Multi : public Base, public Mid
{
public:
	int n;
};
