int sort(void* a, void* b, int elementSize, int count, int (*cmp)(const void*, const void*));

void f(int (callback)(void));

void g(void);

void h(void (**callbackPtr)(void));
