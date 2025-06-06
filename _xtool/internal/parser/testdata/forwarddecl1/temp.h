struct Foo {
    struct bar *b;
};
typedef struct bar bar;
struct bar {
    int a;
};

typedef struct sqlite3_pcache_page sqlite3_pcache_page;
struct sqlite3_pcache_page {
    void *pExtra;
};

typedef struct sqlite3_pcache sqlite3_pcache;

typedef struct sqlite3_pcache_methods2 sqlite3_pcache_methods2;
struct sqlite3_pcache_methods2 {
    int iVersion;
    void (*xShutdown)(void *);
    sqlite3_pcache *(*xCreate)(int szPage, int szExtra, int bPurgeable);
};

typedef struct sqlite3_file sqlite3_file;
struct sqlite3_file {
    const struct sqlite3_io_methods *pMethods; /* Methods for an open file */
};

typedef struct sqlite3_io_methods sqlite3_io_methods;
struct sqlite3_io_methods {
    int (*xUnfetch)(sqlite3_file *, int iOfst, void *p);
};

#define LUA_IDSIZE 60

typedef struct lua_State lua_State;

typedef struct lua_Debug lua_Debug;

int(lua_getstack)(lua_State *L, int level, lua_Debug *ar);

struct lua_Debug {
    // char in ast will got unsigned char & signed char and they are same in go
    // but in ast,will have different,but with compare test,we need avoid these senario
    int short_src[LUA_IDSIZE];
    /* private part */
    struct CallInfo *i_ci; /* active function */
};