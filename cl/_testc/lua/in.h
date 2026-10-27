#include <stddef.h>

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

#define LUA_INTEGER     int
#define LUA_UNSIGNED    unsigned LUA_INTEGER

/* type for integer functions */
typedef LUA_INTEGER lua_Integer;

/* unsigned integer type */
typedef LUA_UNSIGNED lua_Unsigned;

typedef struct lua_State lua_State;

LUALIB_API lua_Integer (luaL_checkinteger) (lua_State *L, int arg);
LUALIB_API lua_Integer (luaL_optinteger) (lua_State *L, int arg,
                                          lua_Integer def);

LUALIB_API int (luaL_loadfilex) (lua_State *L, const char *filename,
                                               const char *mode);

#define luaL_loadfile(L,f)  luaL_loadfilex(L,f,NULL)

#define luaL_checkunsigned(L,a) ((lua_Unsigned)luaL_checkinteger(L,a))
#define luaL_optunsigned(L,a,d) \
  ((lua_Unsigned)luaL_optinteger(L,a,(lua_Integer)(d)))

#define luaL_checkint(L,n)  ((int)luaL_checkinteger(L, (n)))
#define luaL_optint(L,n,d)  ((int)luaL_optinteger(L, (n), (d)))

#define luaL_checklong(L,n) ((long)luaL_checkinteger(L, (n)))
#define luaL_optlong(L,n,d) ((long)luaL_optinteger(L, (n), (d)))

#define foo(p) ((int)(*p))
#define bar(p) (*p)
#define baz(p) *p
