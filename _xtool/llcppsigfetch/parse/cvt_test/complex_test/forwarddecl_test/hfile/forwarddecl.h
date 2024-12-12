typedef struct sqlite3_pcache_page sqlite3_pcache_page;
struct sqlite3_pcache_page
{
    void *pExtra;
};

typedef struct sqlite3_pcache sqlite3_pcache;

typedef struct sqlite3_pcache_methods2 sqlite3_pcache_methods2;
struct sqlite3_pcache_methods2
{
    int iVersion;
    void (*xShutdown)(void *);
    sqlite3_pcache *(*xCreate)(int szPage, int szExtra, int bPurgeable);
};

#define LUA_IDSIZE 60

typedef struct lua_State lua_State;

typedef struct lua_Debug lua_Debug;

int(lua_getstack)(lua_State *L, int level, lua_Debug *ar);

struct lua_Debug
{
    char short_src[LUA_IDSIZE];
    /* private part */
    struct CallInfo *i_ci; /* active function */
};