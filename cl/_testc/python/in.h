#define PyAPI_FUNC(t) t

typedef struct {
    int v;
} PyObject;

PyAPI_FUNC(PyObject *) PyModule_NewObject(
    PyObject *name
    );

PyAPI_FUNC(PyObject *) PyModule_New(
    const char *name            /* UTF-8 encoded string */
    );

PyAPI_FUNC(PyObject *) PyModule_GetDict(PyObject *);

PyAPI_FUNC(PyObject *) PyModule_GetNameObject(PyObject *);

PyAPI_FUNC(const char *) PyModule_GetName(PyObject *);
