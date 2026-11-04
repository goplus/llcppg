typedef int (*callback_t)(int value);

int invoke_callback(int mode, callback_t callback);

#define invoke_empty(mode) invoke_callback(mode, 0)
#define invoke_zero(callback) invoke_callback(0, callback)

void set_options(int count, int *pointer, _Bool enabled);
#define set_default_options() set_options(0, 0, 0)
#define set_parenthesized_options() set_options((0), ((/* null */ 0)), -0)
#define set_default_pointer(pointer) set_options(256, pointer, 0)
#define set_computed_options() set_options(256, 1 - 1, 2 - 2)

void set_flag(int mode, _Bool enabled);
#define set_float_zero() set_flag(0, 0.0)
#define set_float_negative_zero() set_flag(0, -0.0)

void set_context(int mode, void *context);
#define set_empty_context() set_context(256, 0x0)
#define set_computed_context() set_context(256, 1 - 1)

typedef void *handle_t;
int use_handle(int marker, handle_t value);
#define use_empty_handle() use_handle(256, 00)
#define use_computed_handle() use_handle(256, 1 - 1)

int log_value(int count, ...);
#define log_zero() log_value(1, 0)
#define log_computed() log_value(1, 1 - 1)

void set_fixed(int count, int values[3]);
void set_open(int count, int values[]);
void set_variable(int count, int values[count]);
#define set_fixed_empty() set_fixed(1, 0)
#define set_fixed_computed() set_fixed(1, 1 - 1)
#define set_open_empty() set_open(1, 0)
#define set_variable_empty() set_variable(1, 0)

int invoke_direct(int mode, int callback(int));
#define invoke_direct_empty() invoke_direct(256, 0)
#define invoke_direct_computed() invoke_direct(256, 1 - 1)

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
void owner_configure(owner_t *owner, int *value, _Bool enabled);
#define owner_reset(owner) owner_configure(owner, 0, 0)
