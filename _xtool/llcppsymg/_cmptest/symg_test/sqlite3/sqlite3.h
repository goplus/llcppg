
#ifndef __SQlITE3_H__
#define __SQlITE3_H__

#ifndef SQLITE_API
# define SQLITE_API
#endif

typedef struct sqlite3_stmt sqlite3_stmt;

SQLITE_API int sqlite3_finalize(sqlite3_stmt *pStmt);

#endif