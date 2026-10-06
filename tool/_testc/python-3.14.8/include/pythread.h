#ifndef Py_PYTHREAD_H
#define Py_PYTHREAD_H

typedef void *PyThread_type_lock;

#ifdef __cplusplus
extern "C" {
#endif

/* Return status codes for Python lock acquisition.  Chosen for maximum
 * backwards compatibility, ie failure -> 0, success -> 1.  */
typedef enum Py_Lock_Status {
    PY_LOCK_FAILURE = 0,
    PY_LOCK_ACQUIRED = 1,
    PY_LOCK_INTR
} Py_Lock_Status;

typedef struct _Py_tss_t Py_tss_t;  /* opaque */

struct PyObject {};

// Test if an object is the True singleton, the same as "x is True" in Python.
PyAPI_FUNC(int) Py_IsTrue(PyObject *x);

PyAPI_FUNC(int) PyObject_IsTrue(PyObject *x);

PyAPI_FUNC(PyObject *) PyList_GetItem(PyObject *x, PyObject *index);

PyAPI_FUNC(PyObject *) PyDict_GetItem(PyObject *x, PyObject *index);

#ifndef Py_LIMITED_API
#  define Py_CPYTHON_PYTHREAD_H
#  include "cpython/pythread.h"
#  undef Py_CPYTHON_PYTHREAD_H
#endif

#ifdef __cplusplus
}
#endif
#endif /* !Py_PYTHREAD_H */
