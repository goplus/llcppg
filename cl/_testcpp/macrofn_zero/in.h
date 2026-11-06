typedef int (*zero_callback_t)(int value);

void set_options(int count, int *pointer, bool enabled);
#define default_options() set_options(0, 0, 0)
#define float_options() set_options(256, 0, 0.0)
#define pointer_options(pointer) set_options(256, pointer, 0)

int set_callback(int mode, zero_callback_t callback);
#define default_callback() set_callback(256, 0)

void set_context(int mode, void *context);
#define default_context() set_context(256, 0)

typedef void *handle_t;
int use_handle(int marker, handle_t value);
#define default_handle() use_handle(256, 0)

void set_array(int count, int values[3]);
#define default_array() set_array(1, 0)

int log_value(int count, ...);
#define log_zero() log_value(1, 0)
