#define LUALIB_API

typedef unsigned long size_t;

typedef struct {
  char *b;  /* buffer address */
  size_t size;  /* buffer size */
  size_t n;  /* number of characters in buffer */
} luaL_Buffer;

LUALIB_API char *(luaL_prepbuffsize) (luaL_Buffer *B, size_t sz);
LUALIB_API void (luaL_addlstring) (luaL_Buffer *B, const char *s, size_t l);
LUALIB_API void (luaM_addlstring) (luaL_Buffer *B, const char *s, size_t l);
