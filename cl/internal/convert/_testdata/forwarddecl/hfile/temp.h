#include "impl.h"
typedef struct foo foo;
struct bar {
    foo *a;
};

typedef struct sqlite3_file sqlite3_file;
struct sqlite3_file {
    const struct sqlite3_io_methods *pMethods; /* Methods for an open file */
};

typedef struct sqlite3_io_methods sqlite3_io_methods;
struct sqlite3_io_methods {
    int (*xUnfetch)(sqlite3_file *, int iOfst, void *p);
};

typedef struct sqlite3_pcache_page sqlite3_pcache_page;
struct sqlite3_pcache_page {
    void *pBuf;
    void *pExtra;
};

typedef struct sqlite3_pcache sqlite3_pcache;

typedef struct sqlite3_pcache_methods2 sqlite3_pcache_methods2;
struct sqlite3_pcache_methods2 {
    int iVersion;
    void *pArg;
    int (*xInit)(void *);
    void (*xShutdown)(void *);
    sqlite3_pcache *(*xCreate)(int szPage, int szExtra, int bPurgeable);
    void (*xCachesize)(sqlite3_pcache *, int nCachesize);
    int (*xPagecount)(sqlite3_pcache *);
    sqlite3_pcache_page *(*xFetch)(sqlite3_pcache *, unsigned key, int createFlag);
    void (*xUnpin)(sqlite3_pcache *, sqlite3_pcache_page *, int discard);
    void (*xRekey)(sqlite3_pcache *, sqlite3_pcache_page *, unsigned oldKey, unsigned newKey);
    void (*xTruncate)(sqlite3_pcache *, unsigned iLimit);
    void (*xDestroy)(sqlite3_pcache *);
    void (*xShrink)(sqlite3_pcache *);
};

#define LUA_IDSIZE 60

typedef struct lua_State lua_State;

typedef struct lua_Debug lua_Debug;

int(lua_getstack)(lua_State *L, int level, lua_Debug *ar);

struct lua_Debug {
    int event;
    const char *name;
    const char *namewhat;
    const char *what;
    const char *source;
    int currentline;
    int linedefined;
    int lastlinedefined;
    unsigned char nups;
    unsigned char nparams;
    char isvararg;
    char istailcall;
    unsigned short ftransfer;
    unsigned short ntransfer;
    char short_src[LUA_IDSIZE];
    /* private part */
    struct CallInfo *i_ci; /* active function */
};

struct testReferImplicateForward {
    struct CallInfo *i_ci; /* active function */
};

typedef struct Fts5ExtensionApi Fts5ExtensionApi;
typedef struct Fts5Context Fts5Context;
typedef struct Fts5PhraseIter Fts5PhraseIter;

typedef struct sqlite3_value sqlite3_value;
typedef struct sqlite3_context sqlite3_context;
typedef void (*fts5_extension_function)(const Fts5ExtensionApi *pApi, /* API offered by current FTS version */
                                        Fts5Context *pFts,            /* First arg to pass to pApi functions */
                                        sqlite3_context *pCtx,        /* Context for returning result/error */
                                        int nVal,                     /* Number of values in apVal[] array */
                                        sqlite3_value **apVal         /* Array of trailing arguments */
);

struct Fts5PhraseIter {
    const unsigned char *a;
    const unsigned char *b;
};

struct _xmlParserCtxt;

typedef struct _xmlParserCtxt xmlParserCtxt;

typedef xmlParserCtxt htmlParserCtxt;

// https://github.com/goplus/llcppg/issues/526
struct ForwardOnly;
struct EmptyStruct {};