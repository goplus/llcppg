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

struct _object {
    Py_ssize_t ob_refcnt;
};

typedef struct _object PyObject;

// Test if an object is the True singleton, the same as "x is True" in Python.
PyAPI_FUNC(int) Py_IsTrue(PyObject *x);

PyAPI_FUNC(int) PyObject_IsTrue(PyObject *x);

typedef struct {
    struct _object ob_base;
    Py_ssize_t ob_size;
} PyByteArrayObject;

PyAPI_FUNC(PyObject *) PyByteArray_New();
PyAPI_FUNC(PyObject *) PyByteArray_FromObject(PyObject *o);

typedef struct {
    struct _object ob_base;
    Py_ssize_t ob_size;
} PyListObject;

PyAPI_FUNC(PyObject *) PyList_GetItem(PyObject *x, PyObject *index);
PyAPI_FUNC(void) PyList_SetItem(PyObject *x, PyObject *index, PyObject *value);

typedef struct {
    struct _object ob_base;
    Py_ssize_t ma_used;
} PyDictObject;

PyAPI_FUNC(PyObject *) PyDict_GetItem(PyObject *x, PyObject *index);

PyAPI_FUNC(PyObject *) PyImport_Import(PyObject *name);

typedef enum {
    /* PyMem_RawMalloc(), PyMem_RawRealloc() and PyMem_RawFree() */
    PYMEM_DOMAIN_RAW,

    /* PyMem_Malloc(), PyMem_Realloc() and PyMem_Free() */
    PYMEM_DOMAIN_MEM,

    /* PyObject_Malloc(), PyObject_Realloc() and PyObject_Free() */
    PYMEM_DOMAIN_OBJ
} PyMemAllocatorDomain;

typedef struct {
    /* user context passed as the first argument to the 4 functions */
    void *ctx;

    /* allocate a memory block */
    void* (*malloc) (void *ctx, size_t size);

    /* allocate a memory block initialized by zeros */
    void* (*calloc) (void *ctx, size_t nelem, size_t elsize);

    /* allocate or resize a memory block */
    void* (*realloc) (void *ctx, void *ptr, size_t new_size);

    /* release a memory block */
    void (*free) (void *ctx, void *ptr);
} PyMemAllocatorEx;

/* Get the memory block allocator of the specified domain. */
PyAPI_FUNC(void) PyMem_GetAllocator(PyMemAllocatorDomain domain,
                                    PyMemAllocatorEx *allocator);

#ifndef Py_LIMITED_API
#  define Py_CPYTHON_PYTHREAD_H
#  include "cpython/pythread.h"
#  undef Py_CPYTHON_PYTHREAD_H
#endif

#ifdef __cplusplus
}
#endif
#endif /* !Py_PYTHREAD_H */
