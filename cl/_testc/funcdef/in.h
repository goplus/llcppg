#define PyAPI_FUNC(t) t

typedef struct {
    int v;
} PyObject;

PyAPI_FUNC(PyObject *) PyModule_New(
    const char *name
    );

PyAPI_FUNC(PyObject *) PyModule_GetDict(PyObject *);

PyAPI_FUNC(const char *) PyModule_GetName(PyObject *);

PyAPI_FUNC(int) PyObject_IsTrue(PyObject *o);
