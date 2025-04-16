// origin field is float,but if use float,we will ref the c.Float
// we only can get the c.Float's underlying type,not the c.Float
// so we use int to replace float, https://github.com/goplus/llcppg/issues/249
typedef struct Vector3 {
    int x;
    int y;
    int z;
} Vector3;

// in the case we want to check the return type is a zero named struct type,not a anonymous struct type
Vector3 Vector3Barycenter(Vector3 p, Vector3 a, Vector3 b, Vector3 c);
