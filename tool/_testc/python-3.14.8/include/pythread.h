#ifndef Py_PYTHREAD_H
#define Py_PYTHREAD_H

typedef void *PyThread_type_lock;

#ifdef __cplusplus
extern "C" {
#endif

/* Return status codes for Python lock acquisition.  Chosen for maximum
 * backwards compatibility, ie failure -> 0, success -> 1.  */
typedef enum PyLockStatus {
    PY_LOCK_FAILURE = 0,
    PY_LOCK_ACQUIRED = 1,
    PY_LOCK_INTR
} PyLockStatus;

typedef struct _Py_tss_t Py_tss_t;  /* opaque */

#ifndef Py_LIMITED_API
#  define Py_CPYTHON_PYTHREAD_H
#  include "cpython/pythread.h"
#  undef Py_CPYTHON_PYTHREAD_H
#endif

#ifdef __cplusplus
}
#endif
#endif /* !Py_PYTHREAD_H */
