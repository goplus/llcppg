// An enumeration type used as a value type.
enum Color {
	Red,
	Green,
	Blue
};

// A struct: fixed-size arrays declared as fields are true arrays (T[N]).
struct Matrix {
	int rows;
	int cols;
	double data[16];
	char name[32];
};

// Various primitive C types as function parameters and result.
double primitives(char ch, unsigned char uc, signed char sc, short s,
	unsigned short us, int i, unsigned int ui, long l, unsigned long ul,
	long long ll, unsigned long long ull, float flt);

// An enum value passed and returned by value.
enum Color pick(enum Color c);

// Function parameters: an array formally written as T[N] or T[] is a
// pseudo-array that decays to a pointer T*.
int sum(int values[10], int count);

void fill(double buf[], int n);

// Multidimensional and pointer-to-array parameters: only the outermost array
// of a parameter decays. int matrix[3][4] decays its outer [3] to a pointer
// but keeps the inner [4] (=> *[4]int); int (*p)[10] keeps the pointee [10]
// (=> *[10]int).
void transform(int matrix[3][4], int (*p)[10]);

// Function-pointer parameters and results.
int sort(void* a, void* b, int elementSize, int count, int (*cmp)(const void*, const void*));

void f(int (callback)(void));

void g(void);

void h(void (**callbackPtr)(void));
