#include <stddef.h>
typedef unsigned int gpg_error_t;
const char *gpg_strerror (gpg_error_t err);
int gpg_strerror_r (gpg_error_t err, char *buf, size_t buflen);
const char *gpg_strsource (gpg_error_t err);