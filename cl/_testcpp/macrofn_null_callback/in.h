typedef int (*callback_t)(int value);

int invoke_callback(int mode, callback_t callback);

#define invoke_empty(mode) invoke_callback(mode, 0)
#define invoke_zero(callback) invoke_callback(0, callback)

void set_options(int count, int *pointer, bool enabled);
#define set_default_options() set_options(0, 0, 0)

int open_callbacks(const void *cookie,
    int (*read)(void *, char *, int),
    int (*write)(void *, const char *, int),
    long (*seek)(void *, long, int),
    int (*close)(void *));

#define open_read(cookie, fn) open_callbacks(cookie, fn, 0, 0, 0)
#define open_write(cookie, fn) open_callbacks(cookie, 0, fn, 0, 0)

#ifdef __cplusplus
#define NULL 0
#else
#define NULL ((void *)0)
#endif
#define open_empty(cookie) open_callbacks(cookie, NULL, NULL, NULL, NULL)

typedef struct owner owner_t;
int owner_callback(owner_t *owner, int mode, callback_t callback);
#define owner_empty(owner) owner_callback(owner, 0, 0)
