void send_text(const char *text);
void send_pointer(void *pointer);

#define TEXT "_"
#define PREFIX _
#define EMPTY
#define TYPE_NAME int
#define PREFIX_ALIAS PREFIX
#define EMPTY_ALIAS EMPTY
#define BUILTIN_PREFIX __USER_LABEL_PREFIX__
#define NULL ((void *)0)

#define send_literal() send_text(TEXT)
#define send_value(text) send_text(text)
#define send_null() send_pointer(NULL)
#define send_identifier() send_text(PREFIX)
#define send_empty() send_text(EMPTY)
#define send_builtin() send_text(__USER_LABEL_PREFIX__)
