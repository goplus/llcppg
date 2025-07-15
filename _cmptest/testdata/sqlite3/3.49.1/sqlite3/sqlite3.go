package sqlite3

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const VERSION = "3.49.1"
const VERSION_NUMBER = 3049001
const SOURCE_ID = "2025-02-18 13:38:58 873d4e274b4988d260ba8354a9718324a1c26187a4ab4c1cc0227c03d0f10e70"
const OK = 0
const ERROR = 1
const INTERNAL = 2
const PERM = 3
const ABORT = 4
const BUSY = 5
const LOCKED = 6
const NOMEM = 7
const READONLY = 8
const INTERRUPT = 9
const IOERR = 10
const CORRUPT = 11
const NOTFOUND = 12
const FULL = 13
const CANTOPEN = 14
const PROTOCOL = 15
const EMPTY = 16
const SCHEMA = 17
const TOOBIG = 18
const CONSTRAINT = 19
const MISMATCH = 20
const MISUSE = 21
const NOLFS = 22
const AUTH = 23
const FORMAT = 24
const RANGE = 25
const NOTADB = 26
const NOTICE = 27
const WARNING = 28
const ROW = 100
const DONE = 101
const OPEN_READONLY = 0x00000001
const OPEN_READWRITE = 0x00000002
const OPEN_CREATE = 0x00000004
const OPEN_DELETEONCLOSE = 0x00000008
const OPEN_EXCLUSIVE = 0x00000010
const OPEN_AUTOPROXY = 0x00000020
const OPEN_URI = 0x00000040
const OPEN_MEMORY = 0x00000080
const OPEN_MAIN_DB = 0x00000100
const OPEN_TEMP_DB = 0x00000200
const OPEN_TRANSIENT_DB = 0x00000400
const OPEN_MAIN_JOURNAL = 0x00000800
const OPEN_TEMP_JOURNAL = 0x00001000
const OPEN_SUBJOURNAL = 0x00002000
const OPEN_SUPER_JOURNAL = 0x00004000
const OPEN_NOMUTEX = 0x00008000
const OPEN_FULLMUTEX = 0x00010000
const OPEN_SHAREDCACHE = 0x00020000
const OPEN_PRIVATECACHE = 0x00040000
const OPEN_WAL = 0x00080000
const OPEN_NOFOLLOW = 0x01000000
const OPEN_EXRESCODE = 0x02000000
const OPEN_MASTER_JOURNAL = 0x00004000
const IOCAP_ATOMIC = 0x00000001
const IOCAP_ATOMIC512 = 0x00000002
const IOCAP_ATOMIC1K = 0x00000004
const IOCAP_ATOMIC2K = 0x00000008
const IOCAP_ATOMIC4K = 0x00000010
const IOCAP_ATOMIC8K = 0x00000020
const IOCAP_ATOMIC16K = 0x00000040
const IOCAP_ATOMIC32K = 0x00000080
const IOCAP_ATOMIC64K = 0x00000100
const IOCAP_SAFE_APPEND = 0x00000200
const IOCAP_SEQUENTIAL = 0x00000400
const IOCAP_UNDELETABLE_WHEN_OPEN = 0x00000800
const IOCAP_POWERSAFE_OVERWRITE = 0x00001000
const IOCAP_IMMUTABLE = 0x00002000
const IOCAP_BATCH_ATOMIC = 0x00004000
const IOCAP_SUBPAGE_READ = 0x00008000
const LOCK_NONE = 0
const LOCK_SHARED = 1
const LOCK_RESERVED = 2
const LOCK_PENDING = 3
const LOCK_EXCLUSIVE = 4
const SYNC_NORMAL = 0x00002
const SYNC_FULL = 0x00003
const SYNC_DATAONLY = 0x00010
const FCNTL_LOCKSTATE = 1
const FCNTL_GET_LOCKPROXYFILE = 2
const FCNTL_SET_LOCKPROXYFILE = 3
const FCNTL_LAST_ERRNO = 4
const FCNTL_SIZE_HINT = 5
const FCNTL_CHUNK_SIZE = 6
const FCNTL_FILE_POINTER = 7
const FCNTL_SYNC_OMITTED = 8
const FCNTL_WIN32_AV_RETRY = 9
const FCNTL_PERSIST_WAL = 10
const FCNTL_OVERWRITE = 11
const FCNTL_VFSNAME = 12
const FCNTL_POWERSAFE_OVERWRITE = 13
const FCNTL_PRAGMA = 14
const FCNTL_BUSYHANDLER = 15
const FCNTL_TEMPFILENAME = 16
const FCNTL_MMAP_SIZE = 18
const FCNTL_TRACE = 19
const FCNTL_HAS_MOVED = 20
const FCNTL_SYNC = 21
const FCNTL_COMMIT_PHASETWO = 22
const FCNTL_WIN32_SET_HANDLE = 23
const FCNTL_WAL_BLOCK = 24
const FCNTL_ZIPVFS = 25
const FCNTL_RBU = 26
const FCNTL_VFS_POINTER = 27
const FCNTL_JOURNAL_POINTER = 28
const FCNTL_WIN32_GET_HANDLE = 29
const FCNTL_PDB = 30
const FCNTL_BEGIN_ATOMIC_WRITE = 31
const FCNTL_COMMIT_ATOMIC_WRITE = 32
const FCNTL_ROLLBACK_ATOMIC_WRITE = 33
const FCNTL_LOCK_TIMEOUT = 34
const FCNTL_DATA_VERSION = 35
const FCNTL_SIZE_LIMIT = 36
const FCNTL_CKPT_DONE = 37
const FCNTL_RESERVE_BYTES = 38
const FCNTL_CKPT_START = 39
const FCNTL_EXTERNAL_READER = 40
const FCNTL_CKSM_FILE = 41
const FCNTL_RESET_CACHE = 42
const FCNTL_NULL_IO = 43
const ACCESS_EXISTS = 0
const ACCESS_READWRITE = 1
const ACCESS_READ = 2
const SHM_UNLOCK = 1
const SHM_LOCK = 2
const SHM_SHARED = 4
const SHM_EXCLUSIVE = 8
const SHM_NLOCK = 8
const CONFIG_SINGLETHREAD = 1
const CONFIG_MULTITHREAD = 2
const CONFIG_SERIALIZED = 3
const CONFIG_MALLOC = 4
const CONFIG_GETMALLOC = 5
const CONFIG_SCRATCH = 6
const CONFIG_PAGECACHE = 7
const CONFIG_HEAP = 8
const CONFIG_MEMSTATUS = 9
const CONFIG_MUTEX = 10
const CONFIG_GETMUTEX = 11
const CONFIG_LOOKASIDE = 13
const CONFIG_PCACHE = 14
const CONFIG_GETPCACHE = 15
const CONFIG_LOG = 16
const CONFIG_URI = 17
const CONFIG_PCACHE2 = 18
const CONFIG_GETPCACHE2 = 19
const CONFIG_COVERING_INDEX_SCAN = 20
const CONFIG_SQLLOG = 21
const CONFIG_MMAP_SIZE = 22
const CONFIG_WIN32_HEAPSIZE = 23
const CONFIG_PCACHE_HDRSZ = 24
const CONFIG_PMASZ = 25
const CONFIG_STMTJRNL_SPILL = 26
const CONFIG_SMALL_MALLOC = 27
const CONFIG_SORTERREF_SIZE = 28
const CONFIG_MEMDB_MAXSIZE = 29
const CONFIG_ROWID_IN_VIEW = 30
const DBCONFIG_MAINDBNAME = 1000
const DBCONFIG_LOOKASIDE = 1001
const DBCONFIG_ENABLE_FKEY = 1002
const DBCONFIG_ENABLE_TRIGGER = 1003
const DBCONFIG_ENABLE_FTS3_TOKENIZER = 1004
const DBCONFIG_ENABLE_LOAD_EXTENSION = 1005
const DBCONFIG_NO_CKPT_ON_CLOSE = 1006
const DBCONFIG_ENABLE_QPSG = 1007
const DBCONFIG_TRIGGER_EQP = 1008
const DBCONFIG_RESET_DATABASE = 1009
const DBCONFIG_DEFENSIVE = 1010
const DBCONFIG_WRITABLE_SCHEMA = 1011
const DBCONFIG_LEGACY_ALTER_TABLE = 1012
const DBCONFIG_DQS_DML = 1013
const DBCONFIG_DQS_DDL = 1014
const DBCONFIG_ENABLE_VIEW = 1015
const DBCONFIG_LEGACY_FILE_FORMAT = 1016
const DBCONFIG_TRUSTED_SCHEMA = 1017
const DBCONFIG_STMT_SCANSTATUS = 1018
const DBCONFIG_REVERSE_SCANORDER = 1019
const DBCONFIG_ENABLE_ATTACH_CREATE = 1020
const DBCONFIG_ENABLE_ATTACH_WRITE = 1021
const DBCONFIG_ENABLE_COMMENTS = 1022
const DBCONFIG_MAX = 1022
const DENY = 1
const IGNORE = 2
const CREATE_INDEX = 1
const CREATE_TABLE = 2
const CREATE_TEMP_INDEX = 3
const CREATE_TEMP_TABLE = 4
const CREATE_TEMP_TRIGGER = 5
const CREATE_TEMP_VIEW = 6
const CREATE_TRIGGER = 7
const CREATE_VIEW = 8
const DELETE = 9
const DROP_INDEX = 10
const DROP_TABLE = 11
const DROP_TEMP_INDEX = 12
const DROP_TEMP_TABLE = 13
const DROP_TEMP_TRIGGER = 14
const DROP_TEMP_VIEW = 15
const DROP_TRIGGER = 16
const DROP_VIEW = 17
const INSERT = 18
const PRAGMA = 19
const READ = 20
const SELECT = 21
const TRANSACTION = 22
const UPDATE = 23
const ATTACH = 24
const DETACH = 25
const ALTER_TABLE = 26
const REINDEX = 27
const ANALYZE = 28
const CREATE_VTABLE = 29
const DROP_VTABLE = 30
const FUNCTION = 31
const SAVEPOINT = 32
const COPY = 0
const RECURSIVE = 33
const TRACE_STMT = 0x01
const TRACE_PROFILE = 0x02
const TRACE_ROW = 0x04
const TRACE_CLOSE = 0x08
const LIMIT_LENGTH = 0
const LIMIT_SQL_LENGTH = 1
const LIMIT_COLUMN = 2
const LIMIT_EXPR_DEPTH = 3
const LIMIT_COMPOUND_SELECT = 4
const LIMIT_VDBE_OP = 5
const LIMIT_FUNCTION_ARG = 6
const LIMIT_ATTACHED = 7
const LIMIT_LIKE_PATTERN_LENGTH = 8
const LIMIT_VARIABLE_NUMBER = 9
const LIMIT_TRIGGER_DEPTH = 10
const LIMIT_WORKER_THREADS = 11
const PREPARE_PERSISTENT = 0x01
const PREPARE_NORMALIZE = 0x02
const PREPARE_NO_VTAB = 0x04
const PREPARE_DONT_LOG = 0x10
const INTEGER = 1
const FLOAT = 2
const BLOB = 4
const NULL = 5
const TEXT = 3
const SQLITE3_TEXT = 3
const UTF8 = 1
const UTF16LE = 2
const UTF16BE = 3
const UTF16 = 4
const ANY = 5
const UTF16_ALIGNED = 8
const DETERMINISTIC = 0x000000800
const DIRECTONLY = 0x000080000
const SUBTYPE = 0x000100000
const INNOCUOUS = 0x000200000
const RESULT_SUBTYPE = 0x001000000
const SELFORDER1 = 0x002000000
const WIN32_DATA_DIRECTORY_TYPE = 1
const WIN32_TEMP_DIRECTORY_TYPE = 2
const TXN_NONE = 0
const TXN_READ = 1
const TXN_WRITE = 2
const INDEX_SCAN_UNIQUE = 0x00000001
const INDEX_SCAN_HEX = 0x00000002
const INDEX_CONSTRAINT_EQ = 2
const INDEX_CONSTRAINT_GT = 4
const INDEX_CONSTRAINT_LE = 8
const INDEX_CONSTRAINT_LT = 16
const INDEX_CONSTRAINT_GE = 32
const INDEX_CONSTRAINT_MATCH = 64
const INDEX_CONSTRAINT_LIKE = 65
const INDEX_CONSTRAINT_GLOB = 66
const INDEX_CONSTRAINT_REGEXP = 67
const INDEX_CONSTRAINT_NE = 68
const INDEX_CONSTRAINT_ISNOT = 69
const INDEX_CONSTRAINT_ISNOTNULL = 70
const INDEX_CONSTRAINT_ISNULL = 71
const INDEX_CONSTRAINT_IS = 72
const INDEX_CONSTRAINT_LIMIT = 73
const INDEX_CONSTRAINT_OFFSET = 74
const INDEX_CONSTRAINT_FUNCTION = 150
const MUTEX_FAST = 0
const MUTEX_RECURSIVE = 1
const MUTEX_STATIC_MAIN = 2
const MUTEX_STATIC_MEM = 3
const MUTEX_STATIC_MEM2 = 4
const MUTEX_STATIC_OPEN = 4
const MUTEX_STATIC_PRNG = 5
const MUTEX_STATIC_LRU = 6
const MUTEX_STATIC_LRU2 = 7
const MUTEX_STATIC_PMEM = 7
const MUTEX_STATIC_APP1 = 8
const MUTEX_STATIC_APP2 = 9
const MUTEX_STATIC_APP3 = 10
const MUTEX_STATIC_VFS1 = 11
const MUTEX_STATIC_VFS2 = 12
const MUTEX_STATIC_VFS3 = 13
const MUTEX_STATIC_MASTER = 2
const TESTCTRL_FIRST = 5
const TESTCTRL_PRNG_SAVE = 5
const TESTCTRL_PRNG_RESTORE = 6
const TESTCTRL_PRNG_RESET = 7
const TESTCTRL_FK_NO_ACTION = 7
const TESTCTRL_BITVEC_TEST = 8
const TESTCTRL_FAULT_INSTALL = 9
const TESTCTRL_BENIGN_MALLOC_HOOKS = 10
const TESTCTRL_PENDING_BYTE = 11
const TESTCTRL_ASSERT = 12
const TESTCTRL_ALWAYS = 13
const TESTCTRL_RESERVE = 14
const TESTCTRL_JSON_SELFCHECK = 14
const TESTCTRL_OPTIMIZATIONS = 15
const TESTCTRL_ISKEYWORD = 16
const TESTCTRL_GETOPT = 16
const TESTCTRL_SCRATCHMALLOC = 17
const TESTCTRL_INTERNAL_FUNCTIONS = 17
const TESTCTRL_LOCALTIME_FAULT = 18
const TESTCTRL_EXPLAIN_STMT = 19
const TESTCTRL_ONCE_RESET_THRESHOLD = 19
const TESTCTRL_NEVER_CORRUPT = 20
const TESTCTRL_VDBE_COVERAGE = 21
const TESTCTRL_BYTEORDER = 22
const TESTCTRL_ISINIT = 23
const TESTCTRL_SORTER_MMAP = 24
const TESTCTRL_IMPOSTER = 25
const TESTCTRL_PARSER_COVERAGE = 26
const TESTCTRL_RESULT_INTREAL = 27
const TESTCTRL_PRNG_SEED = 28
const TESTCTRL_EXTRA_SCHEMA_CHECKS = 29
const TESTCTRL_SEEK_COUNT = 30
const TESTCTRL_TRACEFLAGS = 31
const TESTCTRL_TUNE = 32
const TESTCTRL_LOGEST = 33
const TESTCTRL_USELONGDOUBLE = 34
const TESTCTRL_LAST = 34
const STATUS_MEMORY_USED = 0
const STATUS_PAGECACHE_USED = 1
const STATUS_PAGECACHE_OVERFLOW = 2
const STATUS_SCRATCH_USED = 3
const STATUS_SCRATCH_OVERFLOW = 4
const STATUS_MALLOC_SIZE = 5
const STATUS_PARSER_STACK = 6
const STATUS_PAGECACHE_SIZE = 7
const STATUS_SCRATCH_SIZE = 8
const STATUS_MALLOC_COUNT = 9
const DBSTATUS_LOOKASIDE_USED = 0
const DBSTATUS_CACHE_USED = 1
const DBSTATUS_SCHEMA_USED = 2
const DBSTATUS_STMT_USED = 3
const DBSTATUS_LOOKASIDE_HIT = 4
const DBSTATUS_LOOKASIDE_MISS_SIZE = 5
const DBSTATUS_LOOKASIDE_MISS_FULL = 6
const DBSTATUS_CACHE_HIT = 7
const DBSTATUS_CACHE_MISS = 8
const DBSTATUS_CACHE_WRITE = 9
const DBSTATUS_DEFERRED_FKS = 10
const DBSTATUS_CACHE_USED_SHARED = 11
const DBSTATUS_CACHE_SPILL = 12
const DBSTATUS_MAX = 12
const STMTSTATUS_FULLSCAN_STEP = 1
const STMTSTATUS_SORT = 2
const STMTSTATUS_AUTOINDEX = 3
const STMTSTATUS_VM_STEP = 4
const STMTSTATUS_REPREPARE = 5
const STMTSTATUS_RUN = 6
const STMTSTATUS_FILTER_MISS = 7
const STMTSTATUS_FILTER_HIT = 8
const STMTSTATUS_MEMUSED = 99
const CHECKPOINT_PASSIVE = 0
const CHECKPOINT_FULL = 1
const CHECKPOINT_RESTART = 2
const CHECKPOINT_TRUNCATE = 3
const VTAB_CONSTRAINT_SUPPORT = 1
const VTAB_INNOCUOUS = 2
const VTAB_DIRECTONLY = 3
const VTAB_USES_ALL_SCHEMAS = 4
const ROLLBACK = 1
const FAIL = 3
const REPLACE = 5
const SCANSTAT_NLOOP = 0
const SCANSTAT_NVISIT = 1
const SCANSTAT_EST = 2
const SCANSTAT_NAME = 3
const SCANSTAT_EXPLAIN = 4
const SCANSTAT_SELECTID = 5
const SCANSTAT_PARENTID = 6
const SCANSTAT_NCYCLE = 7
const SCANSTAT_COMPLEX = 0x0001
const SERIALIZE_NOCOPY = 0x001
const DESERIALIZE_FREEONCLOSE = 1
const DESERIALIZE_RESIZEABLE = 2
const DESERIALIZE_READONLY = 4
const NOT_WITHIN = 0
const PARTLY_WITHIN = 1
const FULLY_WITHIN = 2
const FTS5_TOKENIZE_QUERY = 0x0001
const FTS5_TOKENIZE_PREFIX = 0x0002
const FTS5_TOKENIZE_DOCUMENT = 0x0004
const FTS5_TOKENIZE_AUX = 0x0008
const FTS5_TOKEN_COLOCATED = 0x0001

//go:linkname Libversion C.sqlite3_libversion
func Libversion() *c.Char

//go:linkname Sourceid C.sqlite3_sourceid
func Sourceid() *c.Char

//go:linkname LibversionNumber C.sqlite3_libversion_number
func LibversionNumber() c.Int

/*
** CAPI3REF: Run-Time Library Compilation Options Diagnostics
**
** ^The sqlite3_compileoption_used() function returns 0 or 1
** indicating whether the specified option was defined at
** compile time.  ^The SQLITE_ prefix may be omitted from the
** option name passed to sqlite3_compileoption_used().
**
** ^The sqlite3_compileoption_get() function allows iterating
** over the list of options that were defined at compile time by
** returning the N-th compile time option string.  ^If N is out of range,
** sqlite3_compileoption_get() returns a NULL pointer.  ^The SQLITE_
** prefix is omitted from any strings returned by
** sqlite3_compileoption_get().
**
** ^Support for the diagnostic functions sqlite3_compileoption_used()
** and sqlite3_compileoption_get() may be omitted by specifying the
** [SQLITE_OMIT_COMPILEOPTION_DIAGS] option at compile time.
**
** See also: SQL functions [sqlite_compileoption_used()] and
** [sqlite_compileoption_get()] and the [compile_options pragma].
 */
//go:linkname CompileoptionUsed C.sqlite3_compileoption_used
func CompileoptionUsed(zOptName *c.Char) c.Int

//go:linkname CompileoptionGet C.sqlite3_compileoption_get
func CompileoptionGet(N c.Int) *c.Char

/*
** CAPI3REF: Test To See If The Library Is Threadsafe
**
** ^The sqlite3_threadsafe() function returns zero if and only if
** SQLite was compiled with mutexing code omitted due to the
** [SQLITE_THREADSAFE] compile-time option being set to 0.
**
** SQLite can be compiled with or without mutexes.  When
** the [SQLITE_THREADSAFE] C preprocessor macro is 1 or 2, mutexes
** are enabled and SQLite is threadsafe.  When the
** [SQLITE_THREADSAFE] macro is 0,
** the mutexes are omitted.  Without the mutexes, it is not safe
** to use SQLite concurrently from more than one thread.
**
** Enabling mutexes incurs a measurable performance penalty.
** So if speed is of utmost importance, it makes sense to disable
** the mutexes.  But for maximum safety, mutexes should be enabled.
** ^The default behavior is for mutexes to be enabled.
**
** This interface can be used by an application to make sure that the
** version of SQLite that it is linking against was compiled with
** the desired setting of the [SQLITE_THREADSAFE] macro.
**
** This interface only reports on the compile-time mutex setting
** of the [SQLITE_THREADSAFE] flag.  If SQLite is compiled with
** SQLITE_THREADSAFE=1 or =2 then mutexes are enabled by default but
** can be fully or partially disabled using a call to [sqlite3_config()]
** with the verbs [SQLITE_CONFIG_SINGLETHREAD], [SQLITE_CONFIG_MULTITHREAD],
** or [SQLITE_CONFIG_SERIALIZED].  ^(The return value of the
** sqlite3_threadsafe() function shows only the compile-time setting of
** thread safety, not any run-time changes to that setting made by
** sqlite3_config(). In other words, the return value from sqlite3_threadsafe()
** is unchanged by calls to sqlite3_config().)^
**
** See the [threading mode] documentation for additional information.
 */
//go:linkname Threadsafe C.sqlite3_threadsafe
func Threadsafe() c.Int

type Sqlite3 struct {
	Unused [8]uint8
}
type SqliteInt64 c.LongLong
type SqliteUint64 c.UlongLong
type Int64 SqliteInt64
type Uint64 SqliteUint64

/*
** CAPI3REF: Closing A Database Connection
** DESTRUCTOR: sqlite3
**
** ^The sqlite3_close() and sqlite3_close_v2() routines are destructors
** for the [sqlite3] object.
** ^Calls to sqlite3_close() and sqlite3_close_v2() return [SQLITE_OK] if
** the [sqlite3] object is successfully destroyed and all associated
** resources are deallocated.
**
** Ideally, applications should [sqlite3_finalize | finalize] all
** [prepared statements], [sqlite3_blob_close | close] all [BLOB handles], and
** [sqlite3_backup_finish | finish] all [sqlite3_backup] objects associated
** with the [sqlite3] object prior to attempting to close the object.
** ^If the database connection is associated with unfinalized prepared
** statements, BLOB handlers, and/or unfinished sqlite3_backup objects then
** sqlite3_close() will leave the database connection open and return
** [SQLITE_BUSY]. ^If sqlite3_close_v2() is called with unfinalized prepared
** statements, unclosed BLOB handlers, and/or unfinished sqlite3_backups,
** it returns [SQLITE_OK] regardless, but instead of deallocating the database
** connection immediately, it marks the database connection as an unusable
** "zombie" and makes arrangements to automatically deallocate the database
** connection after all prepared statements are finalized, all BLOB handles
** are closed, and all backups have finished. The sqlite3_close_v2() interface
** is intended for use with host languages that are garbage collected, and
** where the order in which destructors are called is arbitrary.
**
** ^If an [sqlite3] object is destroyed while a transaction is open,
** the transaction is automatically rolled back.
**
** The C parameter to [sqlite3_close(C)] and [sqlite3_close_v2(C)]
** must be either a NULL
** pointer or an [sqlite3] object pointer obtained
** from [sqlite3_open()], [sqlite3_open16()], or
** [sqlite3_open_v2()], and not previously closed.
** ^Calling sqlite3_close() or sqlite3_close_v2() with a NULL pointer
** argument is a harmless no-op.
 */
// llgo:link (*Sqlite3).Close C.sqlite3_close
func (recv_ *Sqlite3) Close() c.Int {
	return 0
}

// llgo:link (*Sqlite3).CloseV2 C.sqlite3_close_v2
func (recv_ *Sqlite3) CloseV2() c.Int {
	return 0
}

// llgo:type C
type Callback func(c.Pointer, c.Int, **c.Char, **c.Char) c.Int

/*
** CAPI3REF: One-Step Query Execution Interface
** METHOD: sqlite3
**
** The sqlite3_exec() interface is a convenience wrapper around
** [sqlite3_prepare_v2()], [sqlite3_step()], and [sqlite3_finalize()],
** that allows an application to run multiple statements of SQL
** without having to use a lot of C code.
**
** ^The sqlite3_exec() interface runs zero or more UTF-8 encoded,
** semicolon-separate SQL statements passed into its 2nd argument,
** in the context of the [database connection] passed in as its 1st
** argument.  ^If the callback function of the 3rd argument to
** sqlite3_exec() is not NULL, then it is invoked for each result row
** coming out of the evaluated SQL statements.  ^The 4th argument to
** sqlite3_exec() is relayed through to the 1st argument of each
** callback invocation.  ^If the callback pointer to sqlite3_exec()
** is NULL, then no callback is ever invoked and result rows are
** ignored.
**
** ^If an error occurs while evaluating the SQL statements passed into
** sqlite3_exec(), then execution of the current statement stops and
** subsequent statements are skipped.  ^If the 5th parameter to sqlite3_exec()
** is not NULL then any error message is written into memory obtained
** from [sqlite3_malloc()] and passed back through the 5th parameter.
** To avoid memory leaks, the application should invoke [sqlite3_free()]
** on error message strings returned through the 5th parameter of
** sqlite3_exec() after the error message string is no longer needed.
** ^If the 5th parameter to sqlite3_exec() is not NULL and no errors
** occur, then sqlite3_exec() sets the pointer in its 5th parameter to
** NULL before returning.
**
** ^If an sqlite3_exec() callback returns non-zero, the sqlite3_exec()
** routine returns SQLITE_ABORT without invoking the callback again and
** without running any subsequent SQL statements.
**
** ^The 2nd argument to the sqlite3_exec() callback function is the
** number of columns in the result.  ^The 3rd argument to the sqlite3_exec()
** callback is an array of pointers to strings obtained as if from
** [sqlite3_column_text()], one for each column.  ^If an element of a
** result row is NULL then the corresponding string pointer for the
** sqlite3_exec() callback is a NULL pointer.  ^The 4th argument to the
** sqlite3_exec() callback is an array of pointers to strings where each
** entry represents the name of corresponding result column as obtained
** from [sqlite3_column_name()].
**
** ^If the 2nd parameter to sqlite3_exec() is a NULL pointer, a pointer
** to an empty string, or a pointer that contains only whitespace and/or
** SQL comments, then no SQL statements are evaluated and the database
** is not changed.
**
** Restrictions:
**
** <ul>
** <li> The application must ensure that the 1st parameter to sqlite3_exec()
**      is a valid and open [database connection].
** <li> The application must not close the [database connection] specified by
**      the 1st parameter to sqlite3_exec() while sqlite3_exec() is running.
** <li> The application must not modify the SQL statement text passed into
**      the 2nd parameter of sqlite3_exec() while sqlite3_exec() is running.
** <li> The application must not dereference the arrays or string pointers
**       passed as the 3rd and 4th callback parameters after it returns.
** </ul>
 */
// llgo:link (*Sqlite3).Exec C.sqlite3_exec
func (recv_ *Sqlite3) Exec(sql *c.Char, callback func(c.Pointer, c.Int, **c.Char, **c.Char) c.Int, __llgo_arg_2 c.Pointer, errmsg **c.Char) c.Int {
	return 0
}

type File struct {
	PMethods *IoMethods
}

type IoMethods struct {
	IVersion               c.Int
	XClose                 c.Pointer
	XRead                  c.Pointer
	XWrite                 c.Pointer
	XTruncate              c.Pointer
	XSync                  c.Pointer
	XFileSize              c.Pointer
	XLock                  c.Pointer
	XUnlock                c.Pointer
	XCheckReservedLock     c.Pointer
	XFileControl           c.Pointer
	XSectorSize            c.Pointer
	XDeviceCharacteristics c.Pointer
	XShmMap                c.Pointer
	XShmLock               c.Pointer
	XShmBarrier            c.Pointer
	XShmUnmap              c.Pointer
	XFetch                 c.Pointer
	XUnfetch               c.Pointer
}

type Mutex struct {
	Unused [8]uint8
}

type ApiRoutines struct {
	AggregateContext     c.Pointer
	AggregateCount       c.Pointer
	BindBlob             c.Pointer
	BindDouble           c.Pointer
	BindInt              c.Pointer
	BindInt64            c.Pointer
	BindNull             c.Pointer
	BindParameterCount   c.Pointer
	BindParameterIndex   c.Pointer
	BindParameterName    c.Pointer
	BindText             c.Pointer
	BindText16           c.Pointer
	BindValue            c.Pointer
	BusyHandler          c.Pointer
	BusyTimeout          c.Pointer
	Changes              c.Pointer
	Close                c.Pointer
	CollationNeeded      c.Pointer
	CollationNeeded16    c.Pointer
	ColumnBlob           c.Pointer
	ColumnBytes          c.Pointer
	ColumnBytes16        c.Pointer
	ColumnCount          c.Pointer
	ColumnDatabaseName   c.Pointer
	ColumnDatabaseName16 c.Pointer
	ColumnDecltype       c.Pointer
	ColumnDecltype16     c.Pointer
	ColumnDouble         c.Pointer
	ColumnInt            c.Pointer
	ColumnInt64          c.Pointer
	ColumnName           c.Pointer
	ColumnName16         c.Pointer
	ColumnOriginName     c.Pointer
	ColumnOriginName16   c.Pointer
	ColumnTableName      c.Pointer
	ColumnTableName16    c.Pointer
	ColumnText           c.Pointer
	ColumnText16         c.Pointer
	ColumnType           c.Pointer
	ColumnValue          c.Pointer
	CommitHook           c.Pointer
	Complete             c.Pointer
	Complete16           c.Pointer
	CreateCollation      c.Pointer
	CreateCollation16    c.Pointer
	CreateFunction       c.Pointer
	CreateFunction16     c.Pointer
	CreateModule         c.Pointer
	DataCount            c.Pointer
	DbHandle             c.Pointer
	DeclareVtab          c.Pointer
	EnableSharedCache    c.Pointer
	Errcode              c.Pointer
	Errmsg               c.Pointer
	Errmsg16             c.Pointer
	Exec                 c.Pointer
	Expired              c.Pointer
	Finalize             c.Pointer
	Free                 c.Pointer
	FreeTable            c.Pointer
	GetAutocommit        c.Pointer
	GetAuxdata           c.Pointer
	GetTable             c.Pointer
	GlobalRecover        c.Pointer
	Interruptx           c.Pointer
	LastInsertRowid      c.Pointer
	Libversion           c.Pointer
	LibversionNumber     c.Pointer
	Malloc               c.Pointer
	Mprintf              c.Pointer
	Open                 c.Pointer
	Open16               c.Pointer
	Prepare              c.Pointer
	Prepare16            c.Pointer
	Profile              c.Pointer
	ProgressHandler      c.Pointer
	Realloc              c.Pointer
	Reset                c.Pointer
	ResultBlob           c.Pointer
	ResultDouble         c.Pointer
	ResultError          c.Pointer
	ResultError16        c.Pointer
	ResultInt            c.Pointer
	ResultInt64          c.Pointer
	ResultNull           c.Pointer
	ResultText           c.Pointer
	ResultText16         c.Pointer
	ResultText16be       c.Pointer
	ResultText16le       c.Pointer
	ResultValue          c.Pointer
	RollbackHook         c.Pointer
	SetAuthorizer        c.Pointer
	SetAuxdata           c.Pointer
	Xsnprintf            c.Pointer
	Step                 c.Pointer
	TableColumnMetadata  c.Pointer
	ThreadCleanup        c.Pointer
	TotalChanges         c.Pointer
	Trace                c.Pointer
	TransferBindings     c.Pointer
	UpdateHook           c.Pointer
	UserData             c.Pointer
	ValueBlob            c.Pointer
	ValueBytes           c.Pointer
	ValueBytes16         c.Pointer
	ValueDouble          c.Pointer
	ValueInt             c.Pointer
	ValueInt64           c.Pointer
	ValueNumericType     c.Pointer
	ValueText            c.Pointer
	ValueText16          c.Pointer
	ValueText16be        c.Pointer
	ValueText16le        c.Pointer
	ValueType            c.Pointer
	Vmprintf             c.Pointer
	OverloadFunction     c.Pointer
	PrepareV2            c.Pointer
	Prepare16V2          c.Pointer
	ClearBindings        c.Pointer
	CreateModuleV2       c.Pointer
	BindZeroblob         c.Pointer
	BlobBytes            c.Pointer
	BlobClose            c.Pointer
	BlobOpen             c.Pointer
	BlobRead             c.Pointer
	BlobWrite            c.Pointer
	CreateCollationV2    c.Pointer
	FileControl          c.Pointer
	MemoryHighwater      c.Pointer
	MemoryUsed           c.Pointer
	MutexAlloc           c.Pointer
	MutexEnter           c.Pointer
	MutexFree            c.Pointer
	MutexLeave           c.Pointer
	MutexTry             c.Pointer
	OpenV2               c.Pointer
	ReleaseMemory        c.Pointer
	ResultErrorNomem     c.Pointer
	ResultErrorToobig    c.Pointer
	Sleep                c.Pointer
	SoftHeapLimit        c.Pointer
	VfsFind              c.Pointer
	VfsRegister          c.Pointer
	VfsUnregister        c.Pointer
	Xthreadsafe          c.Pointer
	ResultZeroblob       c.Pointer
	ResultErrorCode      c.Pointer
	TestControl          c.Pointer
	Randomness           c.Pointer
	ContextDbHandle      c.Pointer
	ExtendedResultCodes  c.Pointer
	Limit                c.Pointer
	NextStmt             c.Pointer
	Sql                  c.Pointer
	Status               c.Pointer
	BackupFinish         c.Pointer
	BackupInit           c.Pointer
	BackupPagecount      c.Pointer
	BackupRemaining      c.Pointer
	BackupStep           c.Pointer
	CompileoptionGet     c.Pointer
	CompileoptionUsed    c.Pointer
	CreateFunctionV2     c.Pointer
	DbConfig             c.Pointer
	DbMutex              c.Pointer
	DbStatus             c.Pointer
	ExtendedErrcode      c.Pointer
	Log                  c.Pointer
	SoftHeapLimit64      c.Pointer
	Sourceid             c.Pointer
	StmtStatus           c.Pointer
	Strnicmp             c.Pointer
	UnlockNotify         c.Pointer
	WalAutocheckpoint    c.Pointer
	WalCheckpoint        c.Pointer
	WalHook              c.Pointer
	BlobReopen           c.Pointer
	VtabConfig           c.Pointer
	VtabOnConflict       c.Pointer
	CloseV2              c.Pointer
	DbFilename           c.Pointer
	DbReadonly           c.Pointer
	DbReleaseMemory      c.Pointer
	Errstr               c.Pointer
	StmtBusy             c.Pointer
	StmtReadonly         c.Pointer
	Stricmp              c.Pointer
	UriBoolean           c.Pointer
	UriInt64             c.Pointer
	UriParameter         c.Pointer
	Xvsnprintf           c.Pointer
	WalCheckpointV2      c.Pointer
	AutoExtension        c.Pointer
	BindBlob64           c.Pointer
	BindText64           c.Pointer
	CancelAutoExtension  c.Pointer
	LoadExtension        c.Pointer
	Malloc64             c.Pointer
	Msize                c.Pointer
	Realloc64            c.Pointer
	ResetAutoExtension   c.Pointer
	ResultBlob64         c.Pointer
	ResultText64         c.Pointer
	Strglob              c.Pointer
	ValueDup             c.Pointer
	ValueFree            c.Pointer
	ResultZeroblob64     c.Pointer
	BindZeroblob64       c.Pointer
	ValueSubtype         c.Pointer
	ResultSubtype        c.Pointer
	Status64             c.Pointer
	Strlike              c.Pointer
	DbCacheflush         c.Pointer
	SystemErrno          c.Pointer
	TraceV2              c.Pointer
	ExpandedSql          c.Pointer
	SetLastInsertRowid   c.Pointer
	PrepareV3            c.Pointer
	Prepare16V3          c.Pointer
	BindPointer          c.Pointer
	ResultPointer        c.Pointer
	ValuePointer         c.Pointer
	VtabNochange         c.Pointer
	ValueNochange        c.Pointer
	VtabCollation        c.Pointer
	KeywordCount         c.Pointer
	KeywordName          c.Pointer
	KeywordCheck         c.Pointer
	StrNew               c.Pointer
	StrFinish            c.Pointer
	StrAppendf           c.Pointer
	StrVappendf          c.Pointer
	StrAppend            c.Pointer
	StrAppendall         c.Pointer
	StrAppendchar        c.Pointer
	StrReset             c.Pointer
	StrErrcode           c.Pointer
	StrLength            c.Pointer
	StrValue             c.Pointer
	CreateWindowFunction c.Pointer
	NormalizedSql        c.Pointer
	StmtIsexplain        c.Pointer
	ValueFrombind        c.Pointer
	DropModules          c.Pointer
	HardHeapLimit64      c.Pointer
	UriKey               c.Pointer
	FilenameDatabase     c.Pointer
	FilenameJournal      c.Pointer
	FilenameWal          c.Pointer
	CreateFilename       c.Pointer
	FreeFilename         c.Pointer
	DatabaseFileObject   c.Pointer
	TxnState             c.Pointer
	Changes64            c.Pointer
	TotalChanges64       c.Pointer
	AutovacuumPages      c.Pointer
	ErrorOffset          c.Pointer
	VtabRhsValue         c.Pointer
	VtabDistinct         c.Pointer
	VtabIn               c.Pointer
	VtabInFirst          c.Pointer
	VtabInNext           c.Pointer
	Deserialize          c.Pointer
	Serialize            c.Pointer
	DbName               c.Pointer
	ValueEncoding        c.Pointer
	IsInterrupted        c.Pointer
	StmtExplain          c.Pointer
	GetClientdata        c.Pointer
	SetClientdata        c.Pointer
}
type Filename *c.Char

type Vfs struct {
	IVersion          c.Int
	SzOsFile          c.Int
	MxPathname        c.Int
	PNext             *Vfs
	ZName             *c.Char
	PAppData          c.Pointer
	XOpen             c.Pointer
	XDelete           c.Pointer
	XAccess           c.Pointer
	XFullPathname     c.Pointer
	XDlOpen           c.Pointer
	XDlError          c.Pointer
	XDlSym            c.Pointer
	XDlClose          c.Pointer
	XRandomness       c.Pointer
	XSleep            c.Pointer
	XCurrentTime      c.Pointer
	XGetLastError     c.Pointer
	XCurrentTimeInt64 c.Pointer
	XSetSystemCall    c.Pointer
	XGetSystemCall    c.Pointer
	XNextSystemCall   c.Pointer
}

// llgo:type C
type SyscallPtr func()

/*
** CAPI3REF: Initialize The SQLite Library
**
** ^The sqlite3_initialize() routine initializes the
** SQLite library.  ^The sqlite3_shutdown() routine
** deallocates any resources that were allocated by sqlite3_initialize().
** These routines are designed to aid in process initialization and
** shutdown on embedded systems.  Workstation applications using
** SQLite normally do not need to invoke either of these routines.
**
** A call to sqlite3_initialize() is an "effective" call if it is
** the first time sqlite3_initialize() is invoked during the lifetime of
** the process, or if it is the first time sqlite3_initialize() is invoked
** following a call to sqlite3_shutdown().  ^(Only an effective call
** of sqlite3_initialize() does any initialization.  All other calls
** are harmless no-ops.)^
**
** A call to sqlite3_shutdown() is an "effective" call if it is the first
** call to sqlite3_shutdown() since the last sqlite3_initialize().  ^(Only
** an effective call to sqlite3_shutdown() does any deinitialization.
** All other valid calls to sqlite3_shutdown() are harmless no-ops.)^
**
** The sqlite3_initialize() interface is threadsafe, but sqlite3_shutdown()
** is not.  The sqlite3_shutdown() interface must only be called from a
** single thread.  All open [database connections] must be closed and all
** other SQLite resources must be deallocated prior to invoking
** sqlite3_shutdown().
**
** Among other things, ^sqlite3_initialize() will invoke
** sqlite3_os_init().  Similarly, ^sqlite3_shutdown()
** will invoke sqlite3_os_end().
**
** ^The sqlite3_initialize() routine returns [SQLITE_OK] on success.
** ^If for some reason, sqlite3_initialize() is unable to initialize
** the library (perhaps it is unable to allocate a needed resource such
** as a mutex) it returns an [error code] other than [SQLITE_OK].
**
** ^The sqlite3_initialize() routine is called internally by many other
** SQLite interfaces so that an application usually does not need to
** invoke sqlite3_initialize() directly.  For example, [sqlite3_open()]
** calls sqlite3_initialize() so the SQLite library will be automatically
** initialized when [sqlite3_open()] is called if it has not be initialized
** already.  ^However, if SQLite is compiled with the [SQLITE_OMIT_AUTOINIT]
** compile-time option, then the automatic calls to sqlite3_initialize()
** are omitted and the application must call sqlite3_initialize() directly
** prior to using any other SQLite interface.  For maximum portability,
** it is recommended that applications always invoke sqlite3_initialize()
** directly prior to using any other SQLite interface.  Future releases
** of SQLite may require this.  In other words, the behavior exhibited
** when SQLite is compiled with [SQLITE_OMIT_AUTOINIT] might become the
** default behavior in some future release of SQLite.
**
** The sqlite3_os_init() routine does operating-system specific
** initialization of the SQLite library.  The sqlite3_os_end()
** routine undoes the effect of sqlite3_os_init().  Typical tasks
** performed by these routines include allocation or deallocation
** of static resources, initialization of global variables,
** setting up a default [sqlite3_vfs] module, or setting up
** a default configuration using [sqlite3_config()].
**
** The application should never invoke either sqlite3_os_init()
** or sqlite3_os_end() directly.  The application should only invoke
** sqlite3_initialize() and sqlite3_shutdown().  The sqlite3_os_init()
** interface is called automatically by sqlite3_initialize() and
** sqlite3_os_end() is called by sqlite3_shutdown().  Appropriate
** implementations for sqlite3_os_init() and sqlite3_os_end()
** are built into SQLite when it is compiled for Unix, Windows, or OS/2.
** When [custom builds | built for other platforms]
** (using the [SQLITE_OS_OTHER=1] compile-time
** option) the application must supply a suitable implementation for
** sqlite3_os_init() and sqlite3_os_end().  An application-supplied
** implementation of sqlite3_os_init() or sqlite3_os_end()
** must return [SQLITE_OK] on success and some other [error code] upon
** failure.
 */
//go:linkname Initialize C.sqlite3_initialize
func Initialize() c.Int

//go:linkname Shutdown C.sqlite3_shutdown
func Shutdown() c.Int

//go:linkname OsInit C.sqlite3_os_init
func OsInit() c.Int

//go:linkname OsEnd C.sqlite3_os_end
func OsEnd() c.Int

/*
** CAPI3REF: Configuring The SQLite Library
**
** The sqlite3_config() interface is used to make global configuration
** changes to SQLite in order to tune SQLite to the specific needs of
** the application.  The default configuration is recommended for most
** applications and so this routine is usually not necessary.  It is
** provided to support rare applications with unusual needs.
**
** <b>The sqlite3_config() interface is not threadsafe. The application
** must ensure that no other SQLite interfaces are invoked by other
** threads while sqlite3_config() is running.</b>
**
** The first argument to sqlite3_config() is an integer
** [configuration option] that determines
** what property of SQLite is to be configured.  Subsequent arguments
** vary depending on the [configuration option]
** in the first argument.
**
** For most configuration options, the sqlite3_config() interface
** may only be invoked prior to library initialization using
** [sqlite3_initialize()] or after shutdown by [sqlite3_shutdown()].
** The exceptional configuration options that may be invoked at any time
** are called "anytime configuration options".
** ^If sqlite3_config() is called after [sqlite3_initialize()] and before
** [sqlite3_shutdown()] with a first argument that is not an anytime
** configuration option, then the sqlite3_config() call will return SQLITE_MISUSE.
** Note, however, that ^sqlite3_config() can be called as part of the
** implementation of an application-defined [sqlite3_os_init()].
**
** ^When a configuration option is set, sqlite3_config() returns [SQLITE_OK].
** ^If the option is unknown or SQLite is unable to set the option
** then this routine returns a non-zero [error code].
 */
//go:linkname Config C.sqlite3_config
func Config(__llgo_arg_0 c.Int, __llgo_va_list ...interface{}) c.Int

/*
** CAPI3REF: Configure database connections
** METHOD: sqlite3
**
** The sqlite3_db_config() interface is used to make configuration
** changes to a [database connection].  The interface is similar to
** [sqlite3_config()] except that the changes apply to a single
** [database connection] (specified in the first argument).
**
** The second argument to sqlite3_db_config(D,V,...)  is the
** [SQLITE_DBCONFIG_LOOKASIDE | configuration verb] - an integer code
** that indicates what aspect of the [database connection] is being configured.
** Subsequent arguments vary depending on the configuration verb.
**
** ^Calls to sqlite3_db_config() return SQLITE_OK if and only if
** the call is considered successful.
 */
//go:linkname DbConfig C.sqlite3_db_config
func DbConfig(__llgo_arg_0 *Sqlite3, op c.Int, __llgo_va_list ...interface{}) c.Int

type MemMethods struct {
	XMalloc   c.Pointer
	XFree     c.Pointer
	XRealloc  c.Pointer
	XSize     c.Pointer
	XRoundup  c.Pointer
	XInit     c.Pointer
	XShutdown c.Pointer
	PAppData  c.Pointer
}

/*
** CAPI3REF: Enable Or Disable Extended Result Codes
** METHOD: sqlite3
**
** ^The sqlite3_extended_result_codes() routine enables or disables the
** [extended result codes] feature of SQLite. ^The extended result
** codes are disabled by default for historical compatibility.
 */
// llgo:link (*Sqlite3).ExtendedResultCodes C.sqlite3_extended_result_codes
func (recv_ *Sqlite3) ExtendedResultCodes(onoff c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Last Insert Rowid
** METHOD: sqlite3
**
** ^Each entry in most SQLite tables (except for [WITHOUT ROWID] tables)
** has a unique 64-bit signed
** integer key called the [ROWID | "rowid"]. ^The rowid is always available
** as an undeclared column named ROWID, OID, or _ROWID_ as long as those
** names are not also used by explicitly declared columns. ^If
** the table has a column of type [INTEGER PRIMARY KEY] then that column
** is another alias for the rowid.
**
** ^The sqlite3_last_insert_rowid(D) interface usually returns the [rowid] of
** the most recent successful [INSERT] into a rowid table or [virtual table]
** on database connection D. ^Inserts into [WITHOUT ROWID] tables are not
** recorded. ^If no successful [INSERT]s into rowid tables have ever occurred
** on the database connection D, then sqlite3_last_insert_rowid(D) returns
** zero.
**
** As well as being set automatically as rows are inserted into database
** tables, the value returned by this function may be set explicitly by
** [sqlite3_set_last_insert_rowid()]
**
** Some virtual table implementations may INSERT rows into rowid tables as
** part of committing a transaction (e.g. to flush data accumulated in memory
** to disk). In this case subsequent calls to this function return the rowid
** associated with these internal INSERT operations, which leads to
** unintuitive results. Virtual table implementations that do write to rowid
** tables in this way can avoid this problem by restoring the original
** rowid value using [sqlite3_set_last_insert_rowid()] before returning
** control to the user.
**
** ^(If an [INSERT] occurs within a trigger then this routine will
** return the [rowid] of the inserted row as long as the trigger is
** running. Once the trigger program ends, the value returned
** by this routine reverts to what it was before the trigger was fired.)^
**
** ^An [INSERT] that fails due to a constraint violation is not a
** successful [INSERT] and does not change the value returned by this
** routine.  ^Thus INSERT OR FAIL, INSERT OR IGNORE, INSERT OR ROLLBACK,
** and INSERT OR ABORT make no changes to the return value of this
** routine when their insertion fails.  ^(When INSERT OR REPLACE
** encounters a constraint violation, it does not fail.  The
** INSERT continues to completion after deleting rows that caused
** the constraint problem so INSERT OR REPLACE will always change
** the return value of this interface.)^
**
** ^For the purposes of this routine, an [INSERT] is considered to
** be successful even if it is subsequently rolled back.
**
** This function is accessible to SQL statements via the
** [last_insert_rowid() SQL function].
**
** If a separate thread performs a new [INSERT] on the same
** database connection while the [sqlite3_last_insert_rowid()]
** function is running and thus changes the last insert [rowid],
** then the value returned by [sqlite3_last_insert_rowid()] is
** unpredictable and might not equal either the old or the new
** last insert [rowid].
 */
// llgo:link (*Sqlite3).LastInsertRowid C.sqlite3_last_insert_rowid
func (recv_ *Sqlite3) LastInsertRowid() Int64 {
	return 0
}

/*
** CAPI3REF: Set the Last Insert Rowid value.
** METHOD: sqlite3
**
** The sqlite3_set_last_insert_rowid(D, R) method allows the application to
** set the value returned by calling sqlite3_last_insert_rowid(D) to R
** without inserting a row into the database.
 */
// llgo:link (*Sqlite3).SetLastInsertRowid C.sqlite3_set_last_insert_rowid
func (recv_ *Sqlite3) SetLastInsertRowid(Int64) {
}

/*
** CAPI3REF: Count The Number Of Rows Modified
** METHOD: sqlite3
**
** ^These functions return the number of rows modified, inserted or
** deleted by the most recently completed INSERT, UPDATE or DELETE
** statement on the database connection specified by the only parameter.
** The two functions are identical except for the type of the return value
** and that if the number of rows modified by the most recent INSERT, UPDATE,
** or DELETE is greater than the maximum value supported by type "int", then
** the return value of sqlite3_changes() is undefined. ^Executing any other
** type of SQL statement does not modify the value returned by these functions.
** For the purposes of this interface, a CREATE TABLE AS SELECT statement
** does not count as an INSERT, UPDATE or DELETE statement and hence the rows
** added to the new table by the CREATE TABLE AS SELECT statement are not
** counted.
**
** ^Only changes made directly by the INSERT, UPDATE or DELETE statement are
** considered - auxiliary changes caused by [CREATE TRIGGER | triggers],
** [foreign key actions] or [REPLACE] constraint resolution are not counted.
**
** Changes to a view that are intercepted by
** [INSTEAD OF trigger | INSTEAD OF triggers] are not counted. ^The value
** returned by sqlite3_changes() immediately after an INSERT, UPDATE or
** DELETE statement run on a view is always zero. Only changes made to real
** tables are counted.
**
** Things are more complicated if the sqlite3_changes() function is
** executed while a trigger program is running. This may happen if the
** program uses the [changes() SQL function], or if some other callback
** function invokes sqlite3_changes() directly. Essentially:
**
** <ul>
**   <li> ^(Before entering a trigger program the value returned by
**        sqlite3_changes() function is saved. After the trigger program
**        has finished, the original value is restored.)^
**
**   <li> ^(Within a trigger program each INSERT, UPDATE and DELETE
**        statement sets the value returned by sqlite3_changes()
**        upon completion as normal. Of course, this value will not include
**        any changes performed by sub-triggers, as the sqlite3_changes()
**        value will be saved and restored after each sub-trigger has run.)^
** </ul>
**
** ^This means that if the changes() SQL function (or similar) is used
** by the first INSERT, UPDATE or DELETE statement within a trigger, it
** returns the value as set when the calling statement began executing.
** ^If it is used by the second or subsequent such statement within a trigger
** program, the value returned reflects the number of rows modified by the
** previous INSERT, UPDATE or DELETE statement within the same trigger.
**
** If a separate thread makes changes on the same database connection
** while [sqlite3_changes()] is running then the value returned
** is unpredictable and not meaningful.
**
** See also:
** <ul>
** <li> the [sqlite3_total_changes()] interface
** <li> the [count_changes pragma]
** <li> the [changes() SQL function]
** <li> the [data_version pragma]
** </ul>
 */
// llgo:link (*Sqlite3).Changes C.sqlite3_changes
func (recv_ *Sqlite3) Changes() c.Int {
	return 0
}

// llgo:link (*Sqlite3).Changes64 C.sqlite3_changes64
func (recv_ *Sqlite3) Changes64() Int64 {
	return 0
}

/*
** CAPI3REF: Total Number Of Rows Modified
** METHOD: sqlite3
**
** ^These functions return the total number of rows inserted, modified or
** deleted by all [INSERT], [UPDATE] or [DELETE] statements completed
** since the database connection was opened, including those executed as
** part of trigger programs. The two functions are identical except for the
** type of the return value and that if the number of rows modified by the
** connection exceeds the maximum value supported by type "int", then
** the return value of sqlite3_total_changes() is undefined. ^Executing
** any other type of SQL statement does not affect the value returned by
** sqlite3_total_changes().
**
** ^Changes made as part of [foreign key actions] are included in the
** count, but those made as part of REPLACE constraint resolution are
** not. ^Changes to a view that are intercepted by INSTEAD OF triggers
** are not counted.
**
** The [sqlite3_total_changes(D)] interface only reports the number
** of rows that changed due to SQL statement run against database
** connection D.  Any changes by other database connections are ignored.
** To detect changes against a database file from other database
** connections use the [PRAGMA data_version] command or the
** [SQLITE_FCNTL_DATA_VERSION] [file control].
**
** If a separate thread makes changes on the same database connection
** while [sqlite3_total_changes()] is running then the value
** returned is unpredictable and not meaningful.
**
** See also:
** <ul>
** <li> the [sqlite3_changes()] interface
** <li> the [count_changes pragma]
** <li> the [changes() SQL function]
** <li> the [data_version pragma]
** <li> the [SQLITE_FCNTL_DATA_VERSION] [file control]
** </ul>
 */
// llgo:link (*Sqlite3).TotalChanges C.sqlite3_total_changes
func (recv_ *Sqlite3) TotalChanges() c.Int {
	return 0
}

// llgo:link (*Sqlite3).TotalChanges64 C.sqlite3_total_changes64
func (recv_ *Sqlite3) TotalChanges64() Int64 {
	return 0
}

/*
** CAPI3REF: Interrupt A Long-Running Query
** METHOD: sqlite3
**
** ^This function causes any pending database operation to abort and
** return at its earliest opportunity. This routine is typically
** called in response to a user action such as pressing "Cancel"
** or Ctrl-C where the user wants a long query operation to halt
** immediately.
**
** ^It is safe to call this routine from a thread different from the
** thread that is currently running the database operation.  But it
** is not safe to call this routine with a [database connection] that
** is closed or might close before sqlite3_interrupt() returns.
**
** ^If an SQL operation is very nearly finished at the time when
** sqlite3_interrupt() is called, then it might not have an opportunity
** to be interrupted and might continue to completion.
**
** ^An SQL operation that is interrupted will return [SQLITE_INTERRUPT].
** ^If the interrupted SQL operation is an INSERT, UPDATE, or DELETE
** that is inside an explicit transaction, then the entire transaction
** will be rolled back automatically.
**
** ^The sqlite3_interrupt(D) call is in effect until all currently running
** SQL statements on [database connection] D complete.  ^Any new SQL statements
** that are started after the sqlite3_interrupt() call and before the
** running statement count reaches zero are interrupted as if they had been
** running prior to the sqlite3_interrupt() call.  ^New SQL statements
** that are started after the running statement count reaches zero are
** not effected by the sqlite3_interrupt().
** ^A call to sqlite3_interrupt(D) that occurs when there are no running
** SQL statements is a no-op and has no effect on SQL statements
** that are started after the sqlite3_interrupt() call returns.
**
** ^The [sqlite3_is_interrupted(D)] interface can be used to determine whether
** or not an interrupt is currently in effect for [database connection] D.
** It returns 1 if an interrupt is currently in effect, or 0 otherwise.
 */
// llgo:link (*Sqlite3).Interrupt C.sqlite3_interrupt
func (recv_ *Sqlite3) Interrupt() {
}

// llgo:link (*Sqlite3).IsInterrupted C.sqlite3_is_interrupted
func (recv_ *Sqlite3) IsInterrupted() c.Int {
	return 0
}

/*
** CAPI3REF: Determine If An SQL Statement Is Complete
**
** These routines are useful during command-line input to determine if the
** currently entered text seems to form a complete SQL statement or
** if additional input is needed before sending the text into
** SQLite for parsing.  ^These routines return 1 if the input string
** appears to be a complete SQL statement.  ^A statement is judged to be
** complete if it ends with a semicolon token and is not a prefix of a
** well-formed CREATE TRIGGER statement.  ^Semicolons that are embedded within
** string literals or quoted identifier names or comments are not
** independent tokens (they are part of the token in which they are
** embedded) and thus do not count as a statement terminator.  ^Whitespace
** and comments that follow the final semicolon are ignored.
**
** ^These routines return 0 if the statement is incomplete.  ^If a
** memory allocation fails, then SQLITE_NOMEM is returned.
**
** ^These routines do not parse the SQL statements thus
** will not detect syntactically incorrect SQL.
**
** ^(If SQLite has not been initialized using [sqlite3_initialize()] prior
** to invoking sqlite3_complete16() then sqlite3_initialize() is invoked
** automatically by sqlite3_complete16().  If that initialization fails,
** then the return value from sqlite3_complete16() will be non-zero
** regardless of whether or not the input SQL is complete.)^
**
** The input to [sqlite3_complete()] must be a zero-terminated
** UTF-8 string.
**
** The input to [sqlite3_complete16()] must be a zero-terminated
** UTF-16 string in native byte order.
 */
//go:linkname Complete C.sqlite3_complete
func Complete(sql *c.Char) c.Int

//go:linkname Complete16 C.sqlite3_complete16
func Complete16(sql c.Pointer) c.Int

/*
** CAPI3REF: Register A Callback To Handle SQLITE_BUSY Errors
** KEYWORDS: {busy-handler callback} {busy handler}
** METHOD: sqlite3
**
** ^The sqlite3_busy_handler(D,X,P) routine sets a callback function X
** that might be invoked with argument P whenever
** an attempt is made to access a database table associated with
** [database connection] D when another thread
** or process has the table locked.
** The sqlite3_busy_handler() interface is used to implement
** [sqlite3_busy_timeout()] and [PRAGMA busy_timeout].
**
** ^If the busy callback is NULL, then [SQLITE_BUSY]
** is returned immediately upon encountering the lock.  ^If the busy callback
** is not NULL, then the callback might be invoked with two arguments.
**
** ^The first argument to the busy handler is a copy of the void* pointer which
** is the third argument to sqlite3_busy_handler().  ^The second argument to
** the busy handler callback is the number of times that the busy handler has
** been invoked previously for the same locking event.  ^If the
** busy callback returns 0, then no additional attempts are made to
** access the database and [SQLITE_BUSY] is returned
** to the application.
** ^If the callback returns non-zero, then another attempt
** is made to access the database and the cycle repeats.
**
** The presence of a busy handler does not guarantee that it will be invoked
** when there is lock contention. ^If SQLite determines that invoking the busy
** handler could result in a deadlock, it will go ahead and return [SQLITE_BUSY]
** to the application instead of invoking the
** busy handler.
** Consider a scenario where one process is holding a read lock that
** it is trying to promote to a reserved lock and
** a second process is holding a reserved lock that it is trying
** to promote to an exclusive lock.  The first process cannot proceed
** because it is blocked by the second and the second process cannot
** proceed because it is blocked by the first.  If both processes
** invoke the busy handlers, neither will make any progress.  Therefore,
** SQLite returns [SQLITE_BUSY] for the first process, hoping that this
** will induce the first process to release its read lock and allow
** the second process to proceed.
**
** ^The default busy callback is NULL.
**
** ^(There can only be a single busy handler defined for each
** [database connection].  Setting a new busy handler clears any
** previously set handler.)^  ^Note that calling [sqlite3_busy_timeout()]
** or evaluating [PRAGMA busy_timeout=N] will change the
** busy handler and thus clear any previously set busy handler.
**
** The busy callback should not take any actions which modify the
** database connection that invoked the busy handler.  In other words,
** the busy handler is not reentrant.  Any such actions
** result in undefined behavior.
**
** A busy handler must not close the database connection
** or [prepared statement] that invoked the busy handler.
 */
// llgo:link (*Sqlite3).BusyHandler C.sqlite3_busy_handler
func (recv_ *Sqlite3) BusyHandler(func(c.Pointer, c.Int) c.Int, c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: Set A Busy Timeout
** METHOD: sqlite3
**
** ^This routine sets a [sqlite3_busy_handler | busy handler] that sleeps
** for a specified amount of time when a table is locked.  ^The handler
** will sleep multiple times until at least "ms" milliseconds of sleeping
** have accumulated.  ^After at least "ms" milliseconds of sleeping,
** the handler returns 0 which causes [sqlite3_step()] to return
** [SQLITE_BUSY].
**
** ^Calling this routine with an argument less than or equal to zero
** turns off all busy handlers.
**
** ^(There can only be a single busy handler for a particular
** [database connection] at any given moment.  If another busy handler
** was defined  (using [sqlite3_busy_handler()]) prior to calling
** this routine, that other busy handler is cleared.)^
**
** See also:  [PRAGMA busy_timeout]
 */
// llgo:link (*Sqlite3).BusyTimeout C.sqlite3_busy_timeout
func (recv_ *Sqlite3) BusyTimeout(ms c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Convenience Routines For Running Queries
** METHOD: sqlite3
**
** This is a legacy interface that is preserved for backwards compatibility.
** Use of this interface is not recommended.
**
** Definition: A <b>result table</b> is memory data structure created by the
** [sqlite3_get_table()] interface.  A result table records the
** complete query results from one or more queries.
**
** The table conceptually has a number of rows and columns.  But
** these numbers are not part of the result table itself.  These
** numbers are obtained separately.  Let N be the number of rows
** and M be the number of columns.
**
** A result table is an array of pointers to zero-terminated UTF-8 strings.
** There are (N+1)*M elements in the array.  The first M pointers point
** to zero-terminated strings that  contain the names of the columns.
** The remaining entries all point to query results.  NULL values result
** in NULL pointers.  All other values are in their UTF-8 zero-terminated
** string representation as returned by [sqlite3_column_text()].
**
** A result table might consist of one or more memory allocations.
** It is not safe to pass a result table directly to [sqlite3_free()].
** A result table should be deallocated using [sqlite3_free_table()].
**
** ^(As an example of the result table format, suppose a query result
** is as follows:
**
** <blockquote><pre>
**        Name        | Age
**        -----------------------
**        Alice       | 43
**        Bob         | 28
**        Cindy       | 21
** </pre></blockquote>
**
** There are two columns (M==2) and three rows (N==3).  Thus the
** result table has 8 entries.  Suppose the result table is stored
** in an array named azResult.  Then azResult holds this content:
**
** <blockquote><pre>
**        azResult&#91;0] = "Name";
**        azResult&#91;1] = "Age";
**        azResult&#91;2] = "Alice";
**        azResult&#91;3] = "43";
**        azResult&#91;4] = "Bob";
**        azResult&#91;5] = "28";
**        azResult&#91;6] = "Cindy";
**        azResult&#91;7] = "21";
** </pre></blockquote>)^
**
** ^The sqlite3_get_table() function evaluates one or more
** semicolon-separated SQL statements in the zero-terminated UTF-8
** string of its 2nd parameter and returns a result table to the
** pointer given in its 3rd parameter.
**
** After the application has finished with the result from sqlite3_get_table(),
** it must pass the result table pointer to sqlite3_free_table() in order to
** release the memory that was malloced.  Because of the way the
** [sqlite3_malloc()] happens within sqlite3_get_table(), the calling
** function must not try to call [sqlite3_free()] directly.  Only
** [sqlite3_free_table()] is able to release the memory properly and safely.
**
** The sqlite3_get_table() interface is implemented as a wrapper around
** [sqlite3_exec()].  The sqlite3_get_table() routine does not have access
** to any internal data structures of SQLite.  It uses only the public
** interface defined here.  As a consequence, errors that occur in the
** wrapper layer outside of the internal [sqlite3_exec()] call are not
** reflected in subsequent calls to [sqlite3_errcode()] or
** [sqlite3_errmsg()].
 */
// llgo:link (*Sqlite3).GetTable C.sqlite3_get_table
func (recv_ *Sqlite3) GetTable(zSql *c.Char, pazResult ***c.Char, pnRow *c.Int, pnColumn *c.Int, pzErrmsg **c.Char) c.Int {
	return 0
}

//go:linkname FreeTable C.sqlite3_free_table
func FreeTable(result **c.Char)

/*
** CAPI3REF: Formatted String Printing Functions
**
** These routines are work-alikes of the "printf()" family of functions
** from the standard C library.
** These routines understand most of the common formatting options from
** the standard library printf()
** plus some additional non-standard formats ([%q], [%Q], [%w], and [%z]).
** See the [built-in printf()] documentation for details.
**
** ^The sqlite3_mprintf() and sqlite3_vmprintf() routines write their
** results into memory obtained from [sqlite3_malloc64()].
** The strings returned by these two routines should be
** released by [sqlite3_free()].  ^Both routines return a
** NULL pointer if [sqlite3_malloc64()] is unable to allocate enough
** memory to hold the resulting string.
**
** ^(The sqlite3_snprintf() routine is similar to "snprintf()" from
** the standard C library.  The result is written into the
** buffer supplied as the second parameter whose size is given by
** the first parameter. Note that the order of the
** first two parameters is reversed from snprintf().)^  This is an
** historical accident that cannot be fixed without breaking
** backwards compatibility.  ^(Note also that sqlite3_snprintf()
** returns a pointer to its buffer instead of the number of
** characters actually written into the buffer.)^  We admit that
** the number of characters written would be a more useful return
** value but we cannot change the implementation of sqlite3_snprintf()
** now without breaking compatibility.
**
** ^As long as the buffer size is greater than zero, sqlite3_snprintf()
** guarantees that the buffer is always zero-terminated.  ^The first
** parameter "n" is the total size of the buffer, including space for
** the zero terminator.  So the longest string that can be completely
** written will be n-1 characters.
**
** ^The sqlite3_vsnprintf() routine is a varargs version of sqlite3_snprintf().
**
** See also:  [built-in printf()], [printf() SQL function]
 */
//go:linkname Mprintf C.sqlite3_mprintf
func Mprintf(__llgo_arg_0 *c.Char, __llgo_va_list ...interface{}) *c.Char

//go:linkname Vmprintf C.sqlite3_vmprintf
func Vmprintf(*c.Char, c.VaList) *c.Char

//go:linkname Snprintf C.sqlite3_snprintf
func Snprintf(__llgo_arg_0 c.Int, __llgo_arg_1 *c.Char, __llgo_arg_2 *c.Char, __llgo_va_list ...interface{}) *c.Char

//go:linkname Vsnprintf C.sqlite3_vsnprintf
func Vsnprintf(c.Int, *c.Char, *c.Char, c.VaList) *c.Char

/*
** CAPI3REF: Memory Allocation Subsystem
**
** The SQLite core uses these three routines for all of its own
** internal memory allocation needs. "Core" in the previous sentence
** does not include operating-system specific [VFS] implementation.  The
** Windows VFS uses native malloc() and free() for some operations.
**
** ^The sqlite3_malloc() routine returns a pointer to a block
** of memory at least N bytes in length, where N is the parameter.
** ^If sqlite3_malloc() is unable to obtain sufficient free
** memory, it returns a NULL pointer.  ^If the parameter N to
** sqlite3_malloc() is zero or negative then sqlite3_malloc() returns
** a NULL pointer.
**
** ^The sqlite3_malloc64(N) routine works just like
** sqlite3_malloc(N) except that N is an unsigned 64-bit integer instead
** of a signed 32-bit integer.
**
** ^Calling sqlite3_free() with a pointer previously returned
** by sqlite3_malloc() or sqlite3_realloc() releases that memory so
** that it might be reused.  ^The sqlite3_free() routine is
** a no-op if is called with a NULL pointer.  Passing a NULL pointer
** to sqlite3_free() is harmless.  After being freed, memory
** should neither be read nor written.  Even reading previously freed
** memory might result in a segmentation fault or other severe error.
** Memory corruption, a segmentation fault, or other severe error
** might result if sqlite3_free() is called with a non-NULL pointer that
** was not obtained from sqlite3_malloc() or sqlite3_realloc().
**
** ^The sqlite3_realloc(X,N) interface attempts to resize a
** prior memory allocation X to be at least N bytes.
** ^If the X parameter to sqlite3_realloc(X,N)
** is a NULL pointer then its behavior is identical to calling
** sqlite3_malloc(N).
** ^If the N parameter to sqlite3_realloc(X,N) is zero or
** negative then the behavior is exactly the same as calling
** sqlite3_free(X).
** ^sqlite3_realloc(X,N) returns a pointer to a memory allocation
** of at least N bytes in size or NULL if insufficient memory is available.
** ^If M is the size of the prior allocation, then min(N,M) bytes
** of the prior allocation are copied into the beginning of buffer returned
** by sqlite3_realloc(X,N) and the prior allocation is freed.
** ^If sqlite3_realloc(X,N) returns NULL and N is positive, then the
** prior allocation is not freed.
**
** ^The sqlite3_realloc64(X,N) interfaces works the same as
** sqlite3_realloc(X,N) except that N is a 64-bit unsigned integer instead
** of a 32-bit signed integer.
**
** ^If X is a memory allocation previously obtained from sqlite3_malloc(),
** sqlite3_malloc64(), sqlite3_realloc(), or sqlite3_realloc64(), then
** sqlite3_msize(X) returns the size of that memory allocation in bytes.
** ^The value returned by sqlite3_msize(X) might be larger than the number
** of bytes requested when X was allocated.  ^If X is a NULL pointer then
** sqlite3_msize(X) returns zero.  If X points to something that is not
** the beginning of memory allocation, or if it points to a formerly
** valid memory allocation that has now been freed, then the behavior
** of sqlite3_msize(X) is undefined and possibly harmful.
**
** ^The memory returned by sqlite3_malloc(), sqlite3_realloc(),
** sqlite3_malloc64(), and sqlite3_realloc64()
** is always aligned to at least an 8 byte boundary, or to a
** 4 byte boundary if the [SQLITE_4_BYTE_ALIGNED_MALLOC] compile-time
** option is used.
**
** The pointer arguments to [sqlite3_free()] and [sqlite3_realloc()]
** must be either NULL or else pointers obtained from a prior
** invocation of [sqlite3_malloc()] or [sqlite3_realloc()] that have
** not yet been released.
**
** The application must not read or write any part of
** a block of memory after it has been released using
** [sqlite3_free()] or [sqlite3_realloc()].
 */
//go:linkname Malloc C.sqlite3_malloc
func Malloc(c.Int) c.Pointer

// llgo:link Uint64.Malloc64 C.sqlite3_malloc64
func (recv_ Uint64) Malloc64() c.Pointer {
	return nil
}

//go:linkname Realloc C.sqlite3_realloc
func Realloc(c.Pointer, c.Int) c.Pointer

//go:linkname Realloc64 C.sqlite3_realloc64
func Realloc64(c.Pointer, Uint64) c.Pointer

//go:linkname Free C.sqlite3_free
func Free(c.Pointer)

//go:linkname Msize C.sqlite3_msize
func Msize(c.Pointer) Uint64

/*
** CAPI3REF: Memory Allocator Statistics
**
** SQLite provides these two interfaces for reporting on the status
** of the [sqlite3_malloc()], [sqlite3_free()], and [sqlite3_realloc()]
** routines, which form the built-in memory allocation subsystem.
**
** ^The [sqlite3_memory_used()] routine returns the number of bytes
** of memory currently outstanding (malloced but not freed).
** ^The [sqlite3_memory_highwater()] routine returns the maximum
** value of [sqlite3_memory_used()] since the high-water mark
** was last reset.  ^The values returned by [sqlite3_memory_used()] and
** [sqlite3_memory_highwater()] include any overhead
** added by SQLite in its implementation of [sqlite3_malloc()],
** but not overhead added by the any underlying system library
** routines that [sqlite3_malloc()] may call.
**
** ^The memory high-water mark is reset to the current value of
** [sqlite3_memory_used()] if and only if the parameter to
** [sqlite3_memory_highwater()] is true.  ^The value returned
** by [sqlite3_memory_highwater(1)] is the high-water mark
** prior to the reset.
 */
//go:linkname MemoryUsed C.sqlite3_memory_used
func MemoryUsed() Int64

//go:linkname MemoryHighwater C.sqlite3_memory_highwater
func MemoryHighwater(resetFlag c.Int) Int64

/*
** CAPI3REF: Pseudo-Random Number Generator
**
** SQLite contains a high-quality pseudo-random number generator (PRNG) used to
** select random [ROWID | ROWIDs] when inserting new records into a table that
** already uses the largest possible [ROWID].  The PRNG is also used for
** the built-in random() and randomblob() SQL functions.  This interface allows
** applications to access the same PRNG for other purposes.
**
** ^A call to this routine stores N bytes of randomness into buffer P.
** ^The P parameter can be a NULL pointer.
**
** ^If this routine has not been previously called or if the previous
** call had N less than one or a NULL pointer for P, then the PRNG is
** seeded using randomness obtained from the xRandomness method of
** the default [sqlite3_vfs] object.
** ^If the previous call to this routine had an N of 1 or more and a
** non-NULL P then the pseudo-randomness is generated
** internally and without recourse to the [sqlite3_vfs] xRandomness
** method.
 */
//go:linkname Randomness C.sqlite3_randomness
func Randomness(N c.Int, P c.Pointer)

/*
** CAPI3REF: Compile-Time Authorization Callbacks
** METHOD: sqlite3
** KEYWORDS: {authorizer callback}
**
** ^This routine registers an authorizer callback with a particular
** [database connection], supplied in the first argument.
** ^The authorizer callback is invoked as SQL statements are being compiled
** by [sqlite3_prepare()] or its variants [sqlite3_prepare_v2()],
** [sqlite3_prepare_v3()], [sqlite3_prepare16()], [sqlite3_prepare16_v2()],
** and [sqlite3_prepare16_v3()].  ^At various
** points during the compilation process, as logic is being created
** to perform various actions, the authorizer callback is invoked to
** see if those actions are allowed.  ^The authorizer callback should
** return [SQLITE_OK] to allow the action, [SQLITE_IGNORE] to disallow the
** specific action but allow the SQL statement to continue to be
** compiled, or [SQLITE_DENY] to cause the entire SQL statement to be
** rejected with an error.  ^If the authorizer callback returns
** any value other than [SQLITE_IGNORE], [SQLITE_OK], or [SQLITE_DENY]
** then the [sqlite3_prepare_v2()] or equivalent call that triggered
** the authorizer will fail with an error message.
**
** When the callback returns [SQLITE_OK], that means the operation
** requested is ok.  ^When the callback returns [SQLITE_DENY], the
** [sqlite3_prepare_v2()] or equivalent call that triggered the
** authorizer will fail with an error message explaining that
** access is denied.
**
** ^The first parameter to the authorizer callback is a copy of the third
** parameter to the sqlite3_set_authorizer() interface. ^The second parameter
** to the callback is an integer [SQLITE_COPY | action code] that specifies
** the particular action to be authorized. ^The third through sixth parameters
** to the callback are either NULL pointers or zero-terminated strings
** that contain additional details about the action to be authorized.
** Applications must always be prepared to encounter a NULL pointer in any
** of the third through the sixth parameters of the authorization callback.
**
** ^If the action code is [SQLITE_READ]
** and the callback returns [SQLITE_IGNORE] then the
** [prepared statement] statement is constructed to substitute
** a NULL value in place of the table column that would have
** been read if [SQLITE_OK] had been returned.  The [SQLITE_IGNORE]
** return can be used to deny an untrusted user access to individual
** columns of a table.
** ^When a table is referenced by a [SELECT] but no column values are
** extracted from that table (for example in a query like
** "SELECT count(*) FROM tab") then the [SQLITE_READ] authorizer callback
** is invoked once for that table with a column name that is an empty string.
** ^If the action code is [SQLITE_DELETE] and the callback returns
** [SQLITE_IGNORE] then the [DELETE] operation proceeds but the
** [truncate optimization] is disabled and all rows are deleted individually.
**
** An authorizer is used when [sqlite3_prepare | preparing]
** SQL statements from an untrusted source, to ensure that the SQL statements
** do not try to access data they are not allowed to see, or that they do not
** try to execute malicious statements that damage the database.  For
** example, an application may allow a user to enter arbitrary
** SQL queries for evaluation by a database.  But the application does
** not want the user to be able to make arbitrary changes to the
** database.  An authorizer could then be put in place while the
** user-entered SQL is being [sqlite3_prepare | prepared] that
** disallows everything except [SELECT] statements.
**
** Applications that need to process SQL from untrusted sources
** might also consider lowering resource limits using [sqlite3_limit()]
** and limiting database size using the [max_page_count] [PRAGMA]
** in addition to using an authorizer.
**
** ^(Only a single authorizer can be in place on a database connection
** at a time.  Each call to sqlite3_set_authorizer overrides the
** previous call.)^  ^Disable the authorizer by installing a NULL callback.
** The authorizer is disabled by default.
**
** The authorizer callback must not do anything that will modify
** the database connection that invoked the authorizer callback.
** Note that [sqlite3_prepare_v2()] and [sqlite3_step()] both modify their
** database connections for the meaning of "modify" in this paragraph.
**
** ^When [sqlite3_prepare_v2()] is used to prepare a statement, the
** statement might be re-prepared during [sqlite3_step()] due to a
** schema change.  Hence, the application should ensure that the
** correct authorizer callback remains in place during the [sqlite3_step()].
**
** ^Note that the authorizer callback is invoked only during
** [sqlite3_prepare()] or its variants.  Authorization is not
** performed during statement evaluation in [sqlite3_step()], unless
** as stated in the previous paragraph, sqlite3_step() invokes
** sqlite3_prepare_v2() to reprepare a statement after a schema change.
 */
// llgo:link (*Sqlite3).SetAuthorizer C.sqlite3_set_authorizer
func (recv_ *Sqlite3) SetAuthorizer(xAuth func(c.Pointer, c.Int, *c.Char, *c.Char, *c.Char, *c.Char) c.Int, pUserData c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: Deprecated Tracing And Profiling Functions
** DEPRECATED
**
** These routines are deprecated. Use the [sqlite3_trace_v2()] interface
** instead of the routines described here.
**
** These routines register callback functions that can be used for
** tracing and profiling the execution of SQL statements.
**
** ^The callback function registered by sqlite3_trace() is invoked at
** various times when an SQL statement is being run by [sqlite3_step()].
** ^The sqlite3_trace() callback is invoked with a UTF-8 rendering of the
** SQL statement text as the statement first begins executing.
** ^(Additional sqlite3_trace() callbacks might occur
** as each triggered subprogram is entered.  The callbacks for triggers
** contain a UTF-8 SQL comment that identifies the trigger.)^
**
** The [SQLITE_TRACE_SIZE_LIMIT] compile-time option can be used to limit
** the length of [bound parameter] expansion in the output of sqlite3_trace().
**
** ^The callback function registered by sqlite3_profile() is invoked
** as each SQL statement finishes.  ^The profile callback contains
** the original statement text and an estimate of wall-clock time
** of how long that statement took to run.  ^The profile callback
** time is in units of nanoseconds, however the current implementation
** is only capable of millisecond resolution so the six least significant
** digits in the time are meaningless.  Future versions of SQLite
** might provide greater resolution on the profiler callback.  Invoking
** either [sqlite3_trace()] or [sqlite3_trace_v2()] will cancel the
** profile callback.
 */
// llgo:link (*Sqlite3).Trace C.sqlite3_trace
func (recv_ *Sqlite3) Trace(xTrace func(c.Pointer, *c.Char), __llgo_arg_1 c.Pointer) c.Pointer {
	return nil
}

// llgo:link (*Sqlite3).Profile C.sqlite3_profile
func (recv_ *Sqlite3) Profile(xProfile func(c.Pointer, *c.Char, Uint64), __llgo_arg_1 c.Pointer) c.Pointer {
	return nil
}

/*
** CAPI3REF: SQL Trace Hook
** METHOD: sqlite3
**
** ^The sqlite3_trace_v2(D,M,X,P) interface registers a trace callback
** function X against [database connection] D, using property mask M
** and context pointer P.  ^If the X callback is
** NULL or if the M mask is zero, then tracing is disabled.  The
** M argument should be the bitwise OR-ed combination of
** zero or more [SQLITE_TRACE] constants.
**
** ^Each call to either sqlite3_trace(D,X,P) or sqlite3_trace_v2(D,M,X,P)
** overrides (cancels) all prior calls to sqlite3_trace(D,X,P) or
** sqlite3_trace_v2(D,M,X,P) for the [database connection] D.  Each
** database connection may have at most one trace callback.
**
** ^The X callback is invoked whenever any of the events identified by
** mask M occur.  ^The integer return value from the callback is currently
** ignored, though this may change in future releases.  Callback
** implementations should return zero to ensure future compatibility.
**
** ^A trace callback is invoked with four arguments: callback(T,C,P,X).
** ^The T argument is one of the [SQLITE_TRACE]
** constants to indicate why the callback was invoked.
** ^The C argument is a copy of the context pointer.
** The P and X arguments are pointers whose meanings depend on T.
**
** The sqlite3_trace_v2() interface is intended to replace the legacy
** interfaces [sqlite3_trace()] and [sqlite3_profile()], both of which
** are deprecated.
 */
// llgo:link (*Sqlite3).TraceV2 C.sqlite3_trace_v2
func (recv_ *Sqlite3) TraceV2(uMask c.Uint, xCallback func(c.Uint, c.Pointer, c.Pointer, c.Pointer) c.Int, pCtx c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: Query Progress Callbacks
** METHOD: sqlite3
**
** ^The sqlite3_progress_handler(D,N,X,P) interface causes the callback
** function X to be invoked periodically during long running calls to
** [sqlite3_step()] and [sqlite3_prepare()] and similar for
** database connection D.  An example use for this
** interface is to keep a GUI updated during a large query.
**
** ^The parameter P is passed through as the only parameter to the
** callback function X.  ^The parameter N is the approximate number of
** [virtual machine instructions] that are evaluated between successive
** invocations of the callback X.  ^If N is less than one then the progress
** handler is disabled.
**
** ^Only a single progress handler may be defined at one time per
** [database connection]; setting a new progress handler cancels the
** old one.  ^Setting parameter X to NULL disables the progress handler.
** ^The progress handler is also disabled by setting N to a value less
** than 1.
**
** ^If the progress callback returns non-zero, the operation is
** interrupted.  This feature can be used to implement a
** "Cancel" button on a GUI progress dialog box.
**
** The progress handler callback must not do anything that will modify
** the database connection that invoked the progress handler.
** Note that [sqlite3_prepare_v2()] and [sqlite3_step()] both modify their
** database connections for the meaning of "modify" in this paragraph.
**
** The progress handler callback would originally only be invoked from the
** bytecode engine.  It still might be invoked during [sqlite3_prepare()]
** and similar because those routines might force a reparse of the schema
** which involves running the bytecode engine.  However, beginning with
** SQLite version 3.41.0, the progress handler callback might also be
** invoked directly from [sqlite3_prepare()] while analyzing and generating
** code for complex queries.
 */
// llgo:link (*Sqlite3).ProgressHandler C.sqlite3_progress_handler
func (recv_ *Sqlite3) ProgressHandler(c.Int, func(c.Pointer) c.Int, c.Pointer) {
}

/*
** CAPI3REF: Opening A New Database Connection
** CONSTRUCTOR: sqlite3
**
** ^These routines open an SQLite database file as specified by the
** filename argument. ^The filename argument is interpreted as UTF-8 for
** sqlite3_open() and sqlite3_open_v2() and as UTF-16 in the native byte
** order for sqlite3_open16(). ^(A [database connection] handle is usually
** returned in *ppDb, even if an error occurs.  The only exception is that
** if SQLite is unable to allocate memory to hold the [sqlite3] object,
** a NULL will be written into *ppDb instead of a pointer to the [sqlite3]
** object.)^ ^(If the database is opened (and/or created) successfully, then
** [SQLITE_OK] is returned.  Otherwise an [error code] is returned.)^ ^The
** [sqlite3_errmsg()] or [sqlite3_errmsg16()] routines can be used to obtain
** an English language description of the error following a failure of any
** of the sqlite3_open() routines.
**
** ^The default encoding will be UTF-8 for databases created using
** sqlite3_open() or sqlite3_open_v2().  ^The default encoding for databases
** created using sqlite3_open16() will be UTF-16 in the native byte order.
**
** Whether or not an error occurs when it is opened, resources
** associated with the [database connection] handle should be released by
** passing it to [sqlite3_close()] when it is no longer required.
**
** The sqlite3_open_v2() interface works like sqlite3_open()
** except that it accepts two additional parameters for additional control
** over the new database connection.  ^(The flags parameter to
** sqlite3_open_v2() must include, at a minimum, one of the following
** three flag combinations:)^
**
** <dl>
** ^(<dt>[SQLITE_OPEN_READONLY]</dt>
** <dd>The database is opened in read-only mode.  If the database does
** not already exist, an error is returned.</dd>)^
**
** ^(<dt>[SQLITE_OPEN_READWRITE]</dt>
** <dd>The database is opened for reading and writing if possible, or
** reading only if the file is write protected by the operating
** system.  In either case the database must already exist, otherwise
** an error is returned.  For historical reasons, if opening in
** read-write mode fails due to OS-level permissions, an attempt is
** made to open it in read-only mode. [sqlite3_db_readonly()] can be
** used to determine whether the database is actually
** read-write.</dd>)^
**
** ^(<dt>[SQLITE_OPEN_READWRITE] | [SQLITE_OPEN_CREATE]</dt>
** <dd>The database is opened for reading and writing, and is created if
** it does not already exist. This is the behavior that is always used for
** sqlite3_open() and sqlite3_open16().</dd>)^
** </dl>
**
** In addition to the required flags, the following optional flags are
** also supported:
**
** <dl>
** ^(<dt>[SQLITE_OPEN_URI]</dt>
** <dd>The filename can be interpreted as a URI if this flag is set.</dd>)^
**
** ^(<dt>[SQLITE_OPEN_MEMORY]</dt>
** <dd>The database will be opened as an in-memory database.  The database
** is named by the "filename" argument for the purposes of cache-sharing,
** if shared cache mode is enabled, but the "filename" is otherwise ignored.
** </dd>)^
**
** ^(<dt>[SQLITE_OPEN_NOMUTEX]</dt>
** <dd>The new database connection will use the "multi-thread"
** [threading mode].)^  This means that separate threads are allowed
** to use SQLite at the same time, as long as each thread is using
** a different [database connection].
**
** ^(<dt>[SQLITE_OPEN_FULLMUTEX]</dt>
** <dd>The new database connection will use the "serialized"
** [threading mode].)^  This means the multiple threads can safely
** attempt to use the same database connection at the same time.
** (Mutexes will block any actual concurrency, but in this mode
** there is no harm in trying.)
**
** ^(<dt>[SQLITE_OPEN_SHAREDCACHE]</dt>
** <dd>The database is opened [shared cache] enabled, overriding
** the default shared cache setting provided by
** [sqlite3_enable_shared_cache()].)^
** The [use of shared cache mode is discouraged] and hence shared cache
** capabilities may be omitted from many builds of SQLite.  In such cases,
** this option is a no-op.
**
** ^(<dt>[SQLITE_OPEN_PRIVATECACHE]</dt>
** <dd>The database is opened [shared cache] disabled, overriding
** the default shared cache setting provided by
** [sqlite3_enable_shared_cache()].)^
**
** [[OPEN_EXRESCODE]] ^(<dt>[SQLITE_OPEN_EXRESCODE]</dt>
** <dd>The database connection comes up in "extended result code mode".
** In other words, the database behaves as if
** [sqlite3_extended_result_codes(db,1)] were called on the database
** connection as soon as the connection is created. In addition to setting
** the extended result code mode, this flag also causes [sqlite3_open_v2()]
** to return an extended result code.</dd>
**
** [[OPEN_NOFOLLOW]] ^(<dt>[SQLITE_OPEN_NOFOLLOW]</dt>
** <dd>The database filename is not allowed to contain a symbolic link</dd>
** </dl>)^
**
** If the 3rd parameter to sqlite3_open_v2() is not one of the
** required combinations shown above optionally combined with other
** [SQLITE_OPEN_READONLY | SQLITE_OPEN_* bits]
** then the behavior is undefined.  Historic versions of SQLite
** have silently ignored surplus bits in the flags parameter to
** sqlite3_open_v2(), however that behavior might not be carried through
** into future versions of SQLite and so applications should not rely
** upon it.  Note in particular that the SQLITE_OPEN_EXCLUSIVE flag is a no-op
** for sqlite3_open_v2().  The SQLITE_OPEN_EXCLUSIVE does *not* cause
** the open to fail if the database already exists.  The SQLITE_OPEN_EXCLUSIVE
** flag is intended for use by the [sqlite3_vfs|VFS interface] only, and not
** by sqlite3_open_v2().
**
** ^The fourth parameter to sqlite3_open_v2() is the name of the
** [sqlite3_vfs] object that defines the operating system interface that
** the new database connection should use.  ^If the fourth parameter is
** a NULL pointer then the default [sqlite3_vfs] object is used.
**
** ^If the filename is ":memory:", then a private, temporary in-memory database
** is created for the connection.  ^This in-memory database will vanish when
** the database connection is closed.  Future versions of SQLite might
** make use of additional special filenames that begin with the ":" character.
** It is recommended that when a database filename actually does begin with
** a ":" character you should prefix the filename with a pathname such as
** "./" to avoid ambiguity.
**
** ^If the filename is an empty string, then a private, temporary
** on-disk database will be created.  ^This private database will be
** automatically deleted as soon as the database connection is closed.
**
** [[URI filenames in sqlite3_open()]] <h3>URI Filenames</h3>
**
** ^If [URI filename] interpretation is enabled, and the filename argument
** begins with "file:", then the filename is interpreted as a URI. ^URI
** filename interpretation is enabled if the [SQLITE_OPEN_URI] flag is
** set in the third argument to sqlite3_open_v2(), or if it has
** been enabled globally using the [SQLITE_CONFIG_URI] option with the
** [sqlite3_config()] method or by the [SQLITE_USE_URI] compile-time option.
** URI filename interpretation is turned off
** by default, but future releases of SQLite might enable URI filename
** interpretation by default.  See "[URI filenames]" for additional
** information.
**
** URI filenames are parsed according to RFC 3986. ^If the URI contains an
** authority, then it must be either an empty string or the string
** "localhost". ^If the authority is not an empty string or "localhost", an
** error is returned to the caller. ^The fragment component of a URI, if
** present, is ignored.
**
** ^SQLite uses the path component of the URI as the name of the disk file
** which contains the database. ^If the path begins with a '/' character,
** then it is interpreted as an absolute path. ^If the path does not begin
** with a '/' (meaning that the authority section is omitted from the URI)
** then the path is interpreted as a relative path.
** ^(On windows, the first component of an absolute path
** is a drive specification (e.g. "C:").)^
**
** [[core URI query parameters]]
** The query component of a URI may contain parameters that are interpreted
** either by SQLite itself, or by a [VFS | custom VFS implementation].
** SQLite and its built-in [VFSes] interpret the
** following query parameters:
**
** <ul>
**   <li> <b>vfs</b>: ^The "vfs" parameter may be used to specify the name of
**     a VFS object that provides the operating system interface that should
**     be used to access the database file on disk. ^If this option is set to
**     an empty string the default VFS object is used. ^Specifying an unknown
**     VFS is an error. ^If sqlite3_open_v2() is used and the vfs option is
**     present, then the VFS specified by the option takes precedence over
**     the value passed as the fourth parameter to sqlite3_open_v2().
**
**   <li> <b>mode</b>: ^(The mode parameter may be set to either "ro", "rw",
**     "rwc", or "memory". Attempting to set it to any other value is
**     an error)^.
**     ^If "ro" is specified, then the database is opened for read-only
**     access, just as if the [SQLITE_OPEN_READONLY] flag had been set in the
**     third argument to sqlite3_open_v2(). ^If the mode option is set to
**     "rw", then the database is opened for read-write (but not create)
**     access, as if SQLITE_OPEN_READWRITE (but not SQLITE_OPEN_CREATE) had
**     been set. ^Value "rwc" is equivalent to setting both
**     SQLITE_OPEN_READWRITE and SQLITE_OPEN_CREATE.  ^If the mode option is
**     set to "memory" then a pure [in-memory database] that never reads
**     or writes from disk is used. ^It is an error to specify a value for
**     the mode parameter that is less restrictive than that specified by
**     the flags passed in the third parameter to sqlite3_open_v2().
**
**   <li> <b>cache</b>: ^The cache parameter may be set to either "shared" or
**     "private". ^Setting it to "shared" is equivalent to setting the
**     SQLITE_OPEN_SHAREDCACHE bit in the flags argument passed to
**     sqlite3_open_v2(). ^Setting the cache parameter to "private" is
**     equivalent to setting the SQLITE_OPEN_PRIVATECACHE bit.
**     ^If sqlite3_open_v2() is used and the "cache" parameter is present in
**     a URI filename, its value overrides any behavior requested by setting
**     SQLITE_OPEN_PRIVATECACHE or SQLITE_OPEN_SHAREDCACHE flag.
**
**  <li> <b>psow</b>: ^The psow parameter indicates whether or not the
**     [powersafe overwrite] property does or does not apply to the
**     storage media on which the database file resides.
**
**  <li> <b>nolock</b>: ^The nolock parameter is a boolean query parameter
**     which if set disables file locking in rollback journal modes.  This
**     is useful for accessing a database on a filesystem that does not
**     support locking.  Caution:  Database corruption might result if two
**     or more processes write to the same database and any one of those
**     processes uses nolock=1.
**
**  <li> <b>immutable</b>: ^The immutable parameter is a boolean query
**     parameter that indicates that the database file is stored on
**     read-only media.  ^When immutable is set, SQLite assumes that the
**     database file cannot be changed, even by a process with higher
**     privilege, and so the database is opened read-only and all locking
**     and change detection is disabled.  Caution: Setting the immutable
**     property on a database file that does in fact change can result
**     in incorrect query results and/or [SQLITE_CORRUPT] errors.
**     See also: [SQLITE_IOCAP_IMMUTABLE].
**
** </ul>
**
** ^Specifying an unknown parameter in the query component of a URI is not an
** error.  Future versions of SQLite might understand additional query
** parameters.  See "[query parameters with special meaning to SQLite]" for
** additional information.
**
** [[URI filename examples]] <h3>URI filename examples</h3>
**
** <table border="1" align=center cellpadding=5>
** <tr><th> URI filenames <th> Results
** <tr><td> file:data.db <td>
**          Open the file "data.db" in the current directory.
** <tr><td> file:/home/fred/data.db<br>
**          file:///home/fred/data.db <br>
**          file://localhost/home/fred/data.db <br> <td>
**          Open the database file "/home/fred/data.db".
** <tr><td> file://darkstar/home/fred/data.db <td>
**          An error. "darkstar" is not a recognized authority.
** <tr><td style="white-space:nowrap">
**          file:///C:/Documents%20and%20Settings/fred/Desktop/data.db
**     <td> Windows only: Open the file "data.db" on fred's desktop on drive
**          C:. Note that the %20 escaping in this example is not strictly
**          necessary - space characters can be used literally
**          in URI filenames.
** <tr><td> file:data.db?mode=ro&cache=private <td>
**          Open file "data.db" in the current directory for read-only access.
**          Regardless of whether or not shared-cache mode is enabled by
**          default, use a private cache.
** <tr><td> file:/home/fred/data.db?vfs=unix-dotfile <td>
**          Open file "/home/fred/data.db". Use the special VFS "unix-dotfile"
**          that uses dot-files in place of posix advisory locking.
** <tr><td> file:data.db?mode=readonly <td>
**          An error. "readonly" is not a valid option for the "mode" parameter.
**          Use "ro" instead:  "file:data.db?mode=ro".
** </table>
**
** ^URI hexadecimal escape sequences (%HH) are supported within the path and
** query components of a URI. A hexadecimal escape sequence consists of a
** percent sign - "%" - followed by exactly two hexadecimal digits
** specifying an octet value. ^Before the path or query components of a
** URI filename are interpreted, they are encoded using UTF-8 and all
** hexadecimal escape sequences replaced by a single byte containing the
** corresponding octet. If this process generates an invalid UTF-8 encoding,
** the results are undefined.
**
** <b>Note to Windows users:</b>  The encoding used for the filename argument
** of sqlite3_open() and sqlite3_open_v2() must be UTF-8, not whatever
** codepage is currently defined.  Filenames containing international
** characters must be converted to UTF-8 prior to passing them into
** sqlite3_open() or sqlite3_open_v2().
**
** <b>Note to Windows Runtime users:</b>  The temporary directory must be set
** prior to calling sqlite3_open() or sqlite3_open_v2().  Otherwise, various
** features that require the use of temporary files may fail.
**
** See also: [sqlite3_temp_directory]
 */
//go:linkname DoOpen C.sqlite3_open
func DoOpen(filename *c.Char, ppDb **Sqlite3) c.Int

//go:linkname DoOpen16 C.sqlite3_open16
func DoOpen16(filename c.Pointer, ppDb **Sqlite3) c.Int

//go:linkname DoOpenV2 C.sqlite3_open_v2
func DoOpenV2(filename *c.Char, ppDb **Sqlite3, flags c.Int, zVfs *c.Char) c.Int

/*
** CAPI3REF: Obtain Values For URI Parameters
**
** These are utility routines, useful to [VFS|custom VFS implementations],
** that check if a database file was a URI that contained a specific query
** parameter, and if so obtains the value of that query parameter.
**
** The first parameter to these interfaces (hereafter referred to
** as F) must be one of:
** <ul>
** <li> A database filename pointer created by the SQLite core and
** passed into the xOpen() method of a VFS implementation, or
** <li> A filename obtained from [sqlite3_db_filename()], or
** <li> A new filename constructed using [sqlite3_create_filename()].
** </ul>
** If the F parameter is not one of the above, then the behavior is
** undefined and probably undesirable.  Older versions of SQLite were
** more tolerant of invalid F parameters than newer versions.
**
** If F is a suitable filename (as described in the previous paragraph)
** and if P is the name of the query parameter, then
** sqlite3_uri_parameter(F,P) returns the value of the P
** parameter if it exists or a NULL pointer if P does not appear as a
** query parameter on F.  If P is a query parameter of F and it
** has no explicit value, then sqlite3_uri_parameter(F,P) returns
** a pointer to an empty string.
**
** The sqlite3_uri_boolean(F,P,B) routine assumes that P is a boolean
** parameter and returns true (1) or false (0) according to the value
** of P.  The sqlite3_uri_boolean(F,P,B) routine returns true (1) if the
** value of query parameter P is one of "yes", "true", or "on" in any
** case or if the value begins with a non-zero number.  The
** sqlite3_uri_boolean(F,P,B) routines returns false (0) if the value of
** query parameter P is one of "no", "false", or "off" in any case or
** if the value begins with a numeric zero.  If P is not a query
** parameter on F or if the value of P does not match any of the
** above, then sqlite3_uri_boolean(F,P,B) returns (B!=0).
**
** The sqlite3_uri_int64(F,P,D) routine converts the value of P into a
** 64-bit signed integer and returns that integer, or D if P does not
** exist.  If the value of P is something other than an integer, then
** zero is returned.
**
** The sqlite3_uri_key(F,N) returns a pointer to the name (not
** the value) of the N-th query parameter for filename F, or a NULL
** pointer if N is less than zero or greater than the number of query
** parameters minus 1.  The N value is zero-based so N should be 0 to obtain
** the name of the first query parameter, 1 for the second parameter, and
** so forth.
**
** If F is a NULL pointer, then sqlite3_uri_parameter(F,P) returns NULL and
** sqlite3_uri_boolean(F,P,B) returns B.  If F is not a NULL pointer and
** is not a database file pathname pointer that the SQLite core passed
** into the xOpen VFS method, then the behavior of this routine is undefined
** and probably undesirable.
**
** Beginning with SQLite [version 3.31.0] ([dateof:3.31.0]) the input F
** parameter can also be the name of a rollback journal file or WAL file
** in addition to the main database file.  Prior to version 3.31.0, these
** routines would only work if F was the name of the main database file.
** When the F parameter is the name of the rollback journal or WAL file,
** it has access to all the same query parameters as were found on the
** main database file.
**
** See the [URI filename] documentation for additional information.
 */
//go:linkname UriParameter C.sqlite3_uri_parameter
func UriParameter(z Filename, zParam *c.Char) *c.Char

//go:linkname UriBoolean C.sqlite3_uri_boolean
func UriBoolean(z Filename, zParam *c.Char, bDefault c.Int) c.Int

//go:linkname UriInt64 C.sqlite3_uri_int64
func UriInt64(Filename, *c.Char, Int64) Int64

//go:linkname UriKey C.sqlite3_uri_key
func UriKey(z Filename, N c.Int) *c.Char

/*
** CAPI3REF:  Translate filenames
**
** These routines are available to [VFS|custom VFS implementations] for
** translating filenames between the main database file, the journal file,
** and the WAL file.
**
** If F is the name of an sqlite database file, journal file, or WAL file
** passed by the SQLite core into the VFS, then sqlite3_filename_database(F)
** returns the name of the corresponding database file.
**
** If F is the name of an sqlite database file, journal file, or WAL file
** passed by the SQLite core into the VFS, or if F is a database filename
** obtained from [sqlite3_db_filename()], then sqlite3_filename_journal(F)
** returns the name of the corresponding rollback journal file.
**
** If F is the name of an sqlite database file, journal file, or WAL file
** that was passed by the SQLite core into the VFS, or if F is a database
** filename obtained from [sqlite3_db_filename()], then
** sqlite3_filename_wal(F) returns the name of the corresponding
** WAL file.
**
** In all of the above, if F is not the name of a database, journal or WAL
** filename passed into the VFS from the SQLite core and F is not the
** return value from [sqlite3_db_filename()], then the result is
** undefined and is likely a memory access violation.
 */
//go:linkname FilenameDatabase C.sqlite3_filename_database
func FilenameDatabase(Filename) *c.Char

//go:linkname FilenameJournal C.sqlite3_filename_journal
func FilenameJournal(Filename) *c.Char

//go:linkname FilenameWal C.sqlite3_filename_wal
func FilenameWal(Filename) *c.Char

/*
** CAPI3REF:  Database File Corresponding To A Journal
**
** ^If X is the name of a rollback or WAL-mode journal file that is
** passed into the xOpen method of [sqlite3_vfs], then
** sqlite3_database_file_object(X) returns a pointer to the [sqlite3_file]
** object that represents the main database file.
**
** This routine is intended for use in custom [VFS] implementations
** only.  It is not a general-purpose interface.
** The argument sqlite3_file_object(X) must be a filename pointer that
** has been passed into [sqlite3_vfs].xOpen method where the
** flags parameter to xOpen contains one of the bits
** [SQLITE_OPEN_MAIN_JOURNAL] or [SQLITE_OPEN_WAL].  Any other use
** of this routine results in undefined and probably undesirable
** behavior.
 */
//go:linkname DatabaseFileObject C.sqlite3_database_file_object
func DatabaseFileObject(*c.Char) *File

/*
** CAPI3REF: Create and Destroy VFS Filenames
**
** These interfaces are provided for use by [VFS shim] implementations and
** are not useful outside of that context.
**
** The sqlite3_create_filename(D,J,W,N,P) allocates memory to hold a version of
** database filename D with corresponding journal file J and WAL file W and
** with N URI parameters key/values pairs in the array P.  The result from
** sqlite3_create_filename(D,J,W,N,P) is a pointer to a database filename that
** is safe to pass to routines like:
** <ul>
** <li> [sqlite3_uri_parameter()],
** <li> [sqlite3_uri_boolean()],
** <li> [sqlite3_uri_int64()],
** <li> [sqlite3_uri_key()],
** <li> [sqlite3_filename_database()],
** <li> [sqlite3_filename_journal()], or
** <li> [sqlite3_filename_wal()].
** </ul>
** If a memory allocation error occurs, sqlite3_create_filename() might
** return a NULL pointer.  The memory obtained from sqlite3_create_filename(X)
** must be released by a corresponding call to sqlite3_free_filename(Y).
**
** The P parameter in sqlite3_create_filename(D,J,W,N,P) should be an array
** of 2*N pointers to strings.  Each pair of pointers in this array corresponds
** to a key and value for a query parameter.  The P parameter may be a NULL
** pointer if N is zero.  None of the 2*N pointers in the P array may be
** NULL pointers and key pointers should not be empty strings.
** None of the D, J, or W parameters to sqlite3_create_filename(D,J,W,N,P) may
** be NULL pointers, though they can be empty strings.
**
** The sqlite3_free_filename(Y) routine releases a memory allocation
** previously obtained from sqlite3_create_filename().  Invoking
** sqlite3_free_filename(Y) where Y is a NULL pointer is a harmless no-op.
**
** If the Y parameter to sqlite3_free_filename(Y) is anything other
** than a NULL pointer or a pointer previously acquired from
** sqlite3_create_filename(), then bad things such as heap
** corruption or segfaults may occur. The value Y should not be
** used again after sqlite3_free_filename(Y) has been called.  This means
** that if the [sqlite3_vfs.xOpen()] method of a VFS has been called using Y,
** then the corresponding [sqlite3_module.xClose() method should also be
** invoked prior to calling sqlite3_free_filename(Y).
 */
//go:linkname CreateFilename C.sqlite3_create_filename
func CreateFilename(zDatabase *c.Char, zJournal *c.Char, zWal *c.Char, nParam c.Int, azParam **c.Char) Filename

//go:linkname FreeFilename C.sqlite3_free_filename
func FreeFilename(Filename)

/*
** CAPI3REF: Error Codes And Messages
** METHOD: sqlite3
**
** ^If the most recent sqlite3_* API call associated with
** [database connection] D failed, then the sqlite3_errcode(D) interface
** returns the numeric [result code] or [extended result code] for that
** API call.
** ^The sqlite3_extended_errcode()
** interface is the same except that it always returns the
** [extended result code] even when extended result codes are
** disabled.
**
** The values returned by sqlite3_errcode() and/or
** sqlite3_extended_errcode() might change with each API call.
** Except, there are some interfaces that are guaranteed to never
** change the value of the error code.  The error-code preserving
** interfaces include the following:
**
** <ul>
** <li> sqlite3_errcode()
** <li> sqlite3_extended_errcode()
** <li> sqlite3_errmsg()
** <li> sqlite3_errmsg16()
** <li> sqlite3_error_offset()
** </ul>
**
** ^The sqlite3_errmsg() and sqlite3_errmsg16() return English-language
** text that describes the error, as either UTF-8 or UTF-16 respectively,
** or NULL if no error message is available.
** (See how SQLite handles [invalid UTF] for exceptions to this rule.)
** ^(Memory to hold the error message string is managed internally.
** The application does not need to worry about freeing the result.
** However, the error string might be overwritten or deallocated by
** subsequent calls to other SQLite interface functions.)^
**
** ^The sqlite3_errstr(E) interface returns the English-language text
** that describes the [result code] E, as UTF-8, or NULL if E is not an
** result code for which a text error message is available.
** ^(Memory to hold the error message string is managed internally
** and must not be freed by the application)^.
**
** ^If the most recent error references a specific token in the input
** SQL, the sqlite3_error_offset() interface returns the byte offset
** of the start of that token.  ^The byte offset returned by
** sqlite3_error_offset() assumes that the input SQL is UTF8.
** ^If the most recent error does not reference a specific token in the input
** SQL, then the sqlite3_error_offset() function returns -1.
**
** When the serialized [threading mode] is in use, it might be the
** case that a second error occurs on a separate thread in between
** the time of the first error and the call to these interfaces.
** When that happens, the second error will be reported since these
** interfaces always report the most recent result.  To avoid
** this, each thread can obtain exclusive use of the [database connection] D
** by invoking [sqlite3_mutex_enter]([sqlite3_db_mutex](D)) before beginning
** to use D and invoking [sqlite3_mutex_leave]([sqlite3_db_mutex](D)) after
** all calls to the interfaces listed here are completed.
**
** If an interface fails with SQLITE_MISUSE, that means the interface
** was invoked incorrectly by the application.  In that case, the
** error code and message may or may not be set.
 */
// llgo:link (*Sqlite3).Errcode C.sqlite3_errcode
func (recv_ *Sqlite3) Errcode() c.Int {
	return 0
}

// llgo:link (*Sqlite3).ExtendedErrcode C.sqlite3_extended_errcode
func (recv_ *Sqlite3) ExtendedErrcode() c.Int {
	return 0
}

// llgo:link (*Sqlite3).Errmsg C.sqlite3_errmsg
func (recv_ *Sqlite3) Errmsg() *c.Char {
	return nil
}

// llgo:link (*Sqlite3).Errmsg16 C.sqlite3_errmsg16
func (recv_ *Sqlite3) Errmsg16() c.Pointer {
	return nil
}

//go:linkname Errstr C.sqlite3_errstr
func Errstr(c.Int) *c.Char

// llgo:link (*Sqlite3).ErrorOffset C.sqlite3_error_offset
func (recv_ *Sqlite3) ErrorOffset() c.Int {
	return 0
}

type Stmt struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Run-time Limits
** METHOD: sqlite3
**
** ^(This interface allows the size of various constructs to be limited
** on a connection by connection basis.  The first parameter is the
** [database connection] whose limit is to be set or queried.  The
** second parameter is one of the [limit categories] that define a
** class of constructs to be size limited.  The third parameter is the
** new limit for that construct.)^
**
** ^If the new limit is a negative number, the limit is unchanged.
** ^(For each limit category SQLITE_LIMIT_<i>NAME</i> there is a
** [limits | hard upper bound]
** set at compile-time by a C preprocessor macro called
** [limits | SQLITE_MAX_<i>NAME</i>].
** (The "_LIMIT_" in the name is changed to "_MAX_".))^
** ^Attempts to increase a limit above its hard upper bound are
** silently truncated to the hard upper bound.
**
** ^Regardless of whether or not the limit was changed, the
** [sqlite3_limit()] interface returns the prior value of the limit.
** ^Hence, to find the current value of a limit without changing it,
** simply invoke this interface with the third parameter set to -1.
**
** Run-time limits are intended for use in applications that manage
** both their own internal database and also databases that are controlled
** by untrusted external sources.  An example application might be a
** web browser that has its own databases for storing history and
** separate databases controlled by JavaScript applications downloaded
** off the Internet.  The internal databases can be given the
** large, default limits.  Databases managed by external sources can
** be given much smaller limits designed to prevent a denial of service
** attack.  Developers might also want to use the [sqlite3_set_authorizer()]
** interface to further control untrusted SQL.  The size of the database
** created by an untrusted script can be contained using the
** [max_page_count] [PRAGMA].
**
** New run-time limit categories may be added in future releases.
 */
// llgo:link (*Sqlite3).Limit C.sqlite3_limit
func (recv_ *Sqlite3) Limit(id c.Int, newVal c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Compiling An SQL Statement
** KEYWORDS: {SQL statement compiler}
** METHOD: sqlite3
** CONSTRUCTOR: sqlite3_stmt
**
** To execute an SQL statement, it must first be compiled into a byte-code
** program using one of these routines.  Or, in other words, these routines
** are constructors for the [prepared statement] object.
**
** The preferred routine to use is [sqlite3_prepare_v2()].  The
** [sqlite3_prepare()] interface is legacy and should be avoided.
** [sqlite3_prepare_v3()] has an extra "prepFlags" option that is used
** for special purposes.
**
** The use of the UTF-8 interfaces is preferred, as SQLite currently
** does all parsing using UTF-8.  The UTF-16 interfaces are provided
** as a convenience.  The UTF-16 interfaces work by converting the
** input text into UTF-8, then invoking the corresponding UTF-8 interface.
**
** The first argument, "db", is a [database connection] obtained from a
** prior successful call to [sqlite3_open()], [sqlite3_open_v2()] or
** [sqlite3_open16()].  The database connection must not have been closed.
**
** The second argument, "zSql", is the statement to be compiled, encoded
** as either UTF-8 or UTF-16.  The sqlite3_prepare(), sqlite3_prepare_v2(),
** and sqlite3_prepare_v3()
** interfaces use UTF-8, and sqlite3_prepare16(), sqlite3_prepare16_v2(),
** and sqlite3_prepare16_v3() use UTF-16.
**
** ^If the nByte argument is negative, then zSql is read up to the
** first zero terminator. ^If nByte is positive, then it is the maximum
** number of bytes read from zSql.  When nByte is positive, zSql is read
** up to the first zero terminator or until the nByte bytes have been read,
** whichever comes first.  ^If nByte is zero, then no prepared
** statement is generated.
** If the caller knows that the supplied string is nul-terminated, then
** there is a small performance advantage to passing an nByte parameter that
** is the number of bytes in the input string <i>including</i>
** the nul-terminator.
** Note that nByte measure the length of the input in bytes, not
** characters, even for the UTF-16 interfaces.
**
** ^If pzTail is not NULL then *pzTail is made to point to the first byte
** past the end of the first SQL statement in zSql.  These routines only
** compile the first statement in zSql, so *pzTail is left pointing to
** what remains uncompiled.
**
** ^*ppStmt is left pointing to a compiled [prepared statement] that can be
** executed using [sqlite3_step()].  ^If there is an error, *ppStmt is set
** to NULL.  ^If the input text contains no SQL (if the input is an empty
** string or a comment) then *ppStmt is set to NULL.
** The calling procedure is responsible for deleting the compiled
** SQL statement using [sqlite3_finalize()] after it has finished with it.
** ppStmt may not be NULL.
**
** ^On success, the sqlite3_prepare() family of routines return [SQLITE_OK];
** otherwise an [error code] is returned.
**
** The sqlite3_prepare_v2(), sqlite3_prepare_v3(), sqlite3_prepare16_v2(),
** and sqlite3_prepare16_v3() interfaces are recommended for all new programs.
** The older interfaces (sqlite3_prepare() and sqlite3_prepare16())
** are retained for backwards compatibility, but their use is discouraged.
** ^In the "vX" interfaces, the prepared statement
** that is returned (the [sqlite3_stmt] object) contains a copy of the
** original SQL text. This causes the [sqlite3_step()] interface to
** behave differently in three ways:
**
** <ol>
** <li>
** ^If the database schema changes, instead of returning [SQLITE_SCHEMA] as it
** always used to do, [sqlite3_step()] will automatically recompile the SQL
** statement and try to run it again. As many as [SQLITE_MAX_SCHEMA_RETRY]
** retries will occur before sqlite3_step() gives up and returns an error.
** </li>
**
** <li>
** ^When an error occurs, [sqlite3_step()] will return one of the detailed
** [error codes] or [extended error codes].  ^The legacy behavior was that
** [sqlite3_step()] would only return a generic [SQLITE_ERROR] result code
** and the application would have to make a second call to [sqlite3_reset()]
** in order to find the underlying cause of the problem. With the "v2" prepare
** interfaces, the underlying reason for the error is returned immediately.
** </li>
**
** <li>
** ^If the specific value bound to a [parameter | host parameter] in the
** WHERE clause might influence the choice of query plan for a statement,
** then the statement will be automatically recompiled, as if there had been
** a schema change, on the first [sqlite3_step()] call following any change
** to the [sqlite3_bind_text | bindings] of that [parameter].
** ^The specific value of a WHERE-clause [parameter] might influence the
** choice of query plan if the parameter is the left-hand side of a [LIKE]
** or [GLOB] operator or if the parameter is compared to an indexed column
** and the [SQLITE_ENABLE_STAT4] compile-time option is enabled.
** </li>
** </ol>
**
** <p>^sqlite3_prepare_v3() differs from sqlite3_prepare_v2() only in having
** the extra prepFlags parameter, which is a bit array consisting of zero or
** more of the [SQLITE_PREPARE_PERSISTENT|SQLITE_PREPARE_*] flags.  ^The
** sqlite3_prepare_v2() interface works exactly the same as
** sqlite3_prepare_v3() with a zero prepFlags parameter.
 */
// llgo:link (*Sqlite3).DoPrepare C.sqlite3_prepare
func (recv_ *Sqlite3) DoPrepare(zSql *c.Char, nByte c.Int, ppStmt **Stmt, pzTail **c.Char) c.Int {
	return 0
}

// llgo:link (*Sqlite3).DoPrepareV2 C.sqlite3_prepare_v2
func (recv_ *Sqlite3) DoPrepareV2(zSql *c.Char, nByte c.Int, ppStmt **Stmt, pzTail **c.Char) c.Int {
	return 0
}

// llgo:link (*Sqlite3).DoPrepareV3 C.sqlite3_prepare_v3
func (recv_ *Sqlite3) DoPrepareV3(zSql *c.Char, nByte c.Int, prepFlags c.Uint, ppStmt **Stmt, pzTail **c.Char) c.Int {
	return 0
}

// llgo:link (*Sqlite3).DoPrepare16 C.sqlite3_prepare16
func (recv_ *Sqlite3) DoPrepare16(zSql c.Pointer, nByte c.Int, ppStmt **Stmt, pzTail *c.Pointer) c.Int {
	return 0
}

// llgo:link (*Sqlite3).DoPrepare16V2 C.sqlite3_prepare16_v2
func (recv_ *Sqlite3) DoPrepare16V2(zSql c.Pointer, nByte c.Int, ppStmt **Stmt, pzTail *c.Pointer) c.Int {
	return 0
}

// llgo:link (*Sqlite3).DoPrepare16V3 C.sqlite3_prepare16_v3
func (recv_ *Sqlite3) DoPrepare16V3(zSql c.Pointer, nByte c.Int, prepFlags c.Uint, ppStmt **Stmt, pzTail *c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: Retrieving Statement SQL
** METHOD: sqlite3_stmt
**
** ^The sqlite3_sql(P) interface returns a pointer to a copy of the UTF-8
** SQL text used to create [prepared statement] P if P was
** created by [sqlite3_prepare_v2()], [sqlite3_prepare_v3()],
** [sqlite3_prepare16_v2()], or [sqlite3_prepare16_v3()].
** ^The sqlite3_expanded_sql(P) interface returns a pointer to a UTF-8
** string containing the SQL text of prepared statement P with
** [bound parameters] expanded.
** ^The sqlite3_normalized_sql(P) interface returns a pointer to a UTF-8
** string containing the normalized SQL text of prepared statement P.  The
** semantics used to normalize a SQL statement are unspecified and subject
** to change.  At a minimum, literal values will be replaced with suitable
** placeholders.
**
** ^(For example, if a prepared statement is created using the SQL
** text "SELECT $abc,:xyz" and if parameter $abc is bound to integer 2345
** and parameter :xyz is unbound, then sqlite3_sql() will return
** the original string, "SELECT $abc,:xyz" but sqlite3_expanded_sql()
** will return "SELECT 2345,NULL".)^
**
** ^The sqlite3_expanded_sql() interface returns NULL if insufficient memory
** is available to hold the result, or if the result would exceed the
** the maximum string length determined by the [SQLITE_LIMIT_LENGTH].
**
** ^The [SQLITE_TRACE_SIZE_LIMIT] compile-time option limits the size of
** bound parameter expansions.  ^The [SQLITE_OMIT_TRACE] compile-time
** option causes sqlite3_expanded_sql() to always return NULL.
**
** ^The strings returned by sqlite3_sql(P) and sqlite3_normalized_sql(P)
** are managed by SQLite and are automatically freed when the prepared
** statement is finalized.
** ^The string returned by sqlite3_expanded_sql(P), on the other hand,
** is obtained from [sqlite3_malloc()] and must be freed by the application
** by passing it to [sqlite3_free()].
**
** ^The sqlite3_normalized_sql() interface is only available if
** the [SQLITE_ENABLE_NORMALIZE] compile-time option is defined.
 */
// llgo:link (*Stmt).Sql C.sqlite3_sql
func (recv_ *Stmt) Sql() *c.Char {
	return nil
}

// llgo:link (*Stmt).ExpandedSql C.sqlite3_expanded_sql
func (recv_ *Stmt) ExpandedSql() *c.Char {
	return nil
}

/*
** CAPI3REF: Determine If An SQL Statement Writes The Database
** METHOD: sqlite3_stmt
**
** ^The sqlite3_stmt_readonly(X) interface returns true (non-zero) if
** and only if the [prepared statement] X makes no direct changes to
** the content of the database file.
**
** Note that [application-defined SQL functions] or
** [virtual tables] might change the database indirectly as a side effect.
** ^(For example, if an application defines a function "eval()" that
** calls [sqlite3_exec()], then the following SQL statement would
** change the database file through side-effects:
**
** <blockquote><pre>
**    SELECT eval('DELETE FROM t1') FROM t2;
** </pre></blockquote>
**
** But because the [SELECT] statement does not change the database file
** directly, sqlite3_stmt_readonly() would still return true.)^
**
** ^Transaction control statements such as [BEGIN], [COMMIT], [ROLLBACK],
** [SAVEPOINT], and [RELEASE] cause sqlite3_stmt_readonly() to return true,
** since the statements themselves do not actually modify the database but
** rather they control the timing of when other statements modify the
** database.  ^The [ATTACH] and [DETACH] statements also cause
** sqlite3_stmt_readonly() to return true since, while those statements
** change the configuration of a database connection, they do not make
** changes to the content of the database files on disk.
** ^The sqlite3_stmt_readonly() interface returns true for [BEGIN] since
** [BEGIN] merely sets internal flags, but the [BEGIN|BEGIN IMMEDIATE] and
** [BEGIN|BEGIN EXCLUSIVE] commands do touch the database and so
** sqlite3_stmt_readonly() returns false for those commands.
**
** ^This routine returns false if there is any possibility that the
** statement might change the database file.  ^A false return does
** not guarantee that the statement will change the database file.
** ^For example, an UPDATE statement might have a WHERE clause that
** makes it a no-op, but the sqlite3_stmt_readonly() result would still
** be false.  ^Similarly, a CREATE TABLE IF NOT EXISTS statement is a
** read-only no-op if the table already exists, but
** sqlite3_stmt_readonly() still returns false for such a statement.
**
** ^If prepared statement X is an [EXPLAIN] or [EXPLAIN QUERY PLAN]
** statement, then sqlite3_stmt_readonly(X) returns the same value as
** if the EXPLAIN or EXPLAIN QUERY PLAN prefix were omitted.
 */
// llgo:link (*Stmt).StmtReadonly C.sqlite3_stmt_readonly
func (recv_ *Stmt) StmtReadonly() c.Int {
	return 0
}

/*
** CAPI3REF: Query The EXPLAIN Setting For A Prepared Statement
** METHOD: sqlite3_stmt
**
** ^The sqlite3_stmt_isexplain(S) interface returns 1 if the
** prepared statement S is an EXPLAIN statement, or 2 if the
** statement S is an EXPLAIN QUERY PLAN.
** ^The sqlite3_stmt_isexplain(S) interface returns 0 if S is
** an ordinary statement or a NULL pointer.
 */
// llgo:link (*Stmt).StmtIsexplain C.sqlite3_stmt_isexplain
func (recv_ *Stmt) StmtIsexplain() c.Int {
	return 0
}

/*
** CAPI3REF: Change The EXPLAIN Setting For A Prepared Statement
** METHOD: sqlite3_stmt
**
** The sqlite3_stmt_explain(S,E) interface changes the EXPLAIN
** setting for [prepared statement] S.  If E is zero, then S becomes
** a normal prepared statement.  If E is 1, then S behaves as if
** its SQL text began with "[EXPLAIN]".  If E is 2, then S behaves as if
** its SQL text began with "[EXPLAIN QUERY PLAN]".
**
** Calling sqlite3_stmt_explain(S,E) might cause S to be reprepared.
** SQLite tries to avoid a reprepare, but a reprepare might be necessary
** on the first transition into EXPLAIN or EXPLAIN QUERY PLAN mode.
**
** Because of the potential need to reprepare, a call to
** sqlite3_stmt_explain(S,E) will fail with SQLITE_ERROR if S cannot be
** reprepared because it was created using [sqlite3_prepare()] instead of
** the newer [sqlite3_prepare_v2()] or [sqlite3_prepare_v3()] interfaces and
** hence has no saved SQL text with which to reprepare.
**
** Changing the explain setting for a prepared statement does not change
** the original SQL text for the statement.  Hence, if the SQL text originally
** began with EXPLAIN or EXPLAIN QUERY PLAN, but sqlite3_stmt_explain(S,0)
** is called to convert the statement into an ordinary statement, the EXPLAIN
** or EXPLAIN QUERY PLAN keywords will still appear in the sqlite3_sql(S)
** output, even though the statement now acts like a normal SQL statement.
**
** This routine returns SQLITE_OK if the explain mode is successfully
** changed, or an error code if the explain mode could not be changed.
** The explain mode cannot be changed while a statement is active.
** Hence, it is good practice to call [sqlite3_reset(S)]
** immediately prior to calling sqlite3_stmt_explain(S,E).
 */
// llgo:link (*Stmt).StmtExplain C.sqlite3_stmt_explain
func (recv_ *Stmt) StmtExplain(eMode c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Determine If A Prepared Statement Has Been Reset
** METHOD: sqlite3_stmt
**
** ^The sqlite3_stmt_busy(S) interface returns true (non-zero) if the
** [prepared statement] S has been stepped at least once using
** [sqlite3_step(S)] but has neither run to completion (returned
** [SQLITE_DONE] from [sqlite3_step(S)]) nor
** been reset using [sqlite3_reset(S)].  ^The sqlite3_stmt_busy(S)
** interface returns false if S is a NULL pointer.  If S is not a
** NULL pointer and is not a pointer to a valid [prepared statement]
** object, then the behavior is undefined and probably undesirable.
**
** This interface can be used in combination [sqlite3_next_stmt()]
** to locate all prepared statements associated with a database
** connection that are in need of being reset.  This can be used,
** for example, in diagnostic routines to search for prepared
** statements that are holding a transaction open.
 */
// llgo:link (*Stmt).StmtBusy C.sqlite3_stmt_busy
func (recv_ *Stmt) StmtBusy() c.Int {
	return 0
}

type Value struct {
	Unused [8]uint8
}

type Context struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Binding Values To Prepared Statements
** KEYWORDS: {host parameter} {host parameters} {host parameter name}
** KEYWORDS: {SQL parameter} {SQL parameters} {parameter binding}
** METHOD: sqlite3_stmt
**
** ^(In the SQL statement text input to [sqlite3_prepare_v2()] and its variants,
** literals may be replaced by a [parameter] that matches one of following
** templates:
**
** <ul>
** <li>  ?
** <li>  ?NNN
** <li>  :VVV
** <li>  @VVV
** <li>  $VVV
** </ul>
**
** In the templates above, NNN represents an integer literal,
** and VVV represents an alphanumeric identifier.)^  ^The values of these
** parameters (also called "host parameter names" or "SQL parameters")
** can be set using the sqlite3_bind_*() routines defined here.
**
** ^The first argument to the sqlite3_bind_*() routines is always
** a pointer to the [sqlite3_stmt] object returned from
** [sqlite3_prepare_v2()] or its variants.
**
** ^The second argument is the index of the SQL parameter to be set.
** ^The leftmost SQL parameter has an index of 1.  ^When the same named
** SQL parameter is used more than once, second and subsequent
** occurrences have the same index as the first occurrence.
** ^The index for named parameters can be looked up using the
** [sqlite3_bind_parameter_index()] API if desired.  ^The index
** for "?NNN" parameters is the value of NNN.
** ^The NNN value must be between 1 and the [sqlite3_limit()]
** parameter [SQLITE_LIMIT_VARIABLE_NUMBER] (default value: 32766).
**
** ^The third argument is the value to bind to the parameter.
** ^If the third parameter to sqlite3_bind_text() or sqlite3_bind_text16()
** or sqlite3_bind_blob() is a NULL pointer then the fourth parameter
** is ignored and the end result is the same as sqlite3_bind_null().
** ^If the third parameter to sqlite3_bind_text() is not NULL, then
** it should be a pointer to well-formed UTF8 text.
** ^If the third parameter to sqlite3_bind_text16() is not NULL, then
** it should be a pointer to well-formed UTF16 text.
** ^If the third parameter to sqlite3_bind_text64() is not NULL, then
** it should be a pointer to a well-formed unicode string that is
** either UTF8 if the sixth parameter is SQLITE_UTF8, or UTF16
** otherwise.
**
** [[byte-order determination rules]] ^The byte-order of
** UTF16 input text is determined by the byte-order mark (BOM, U+FEFF)
** found in first character, which is removed, or in the absence of a BOM
** the byte order is the native byte order of the host
** machine for sqlite3_bind_text16() or the byte order specified in
** the 6th parameter for sqlite3_bind_text64().)^
** ^If UTF16 input text contains invalid unicode
** characters, then SQLite might change those invalid characters
** into the unicode replacement character: U+FFFD.
**
** ^(In those routines that have a fourth argument, its value is the
** number of bytes in the parameter.  To be clear: the value is the
** number of <u>bytes</u> in the value, not the number of characters.)^
** ^If the fourth parameter to sqlite3_bind_text() or sqlite3_bind_text16()
** is negative, then the length of the string is
** the number of bytes up to the first zero terminator.
** If the fourth parameter to sqlite3_bind_blob() is negative, then
** the behavior is undefined.
** If a non-negative fourth parameter is provided to sqlite3_bind_text()
** or sqlite3_bind_text16() or sqlite3_bind_text64() then
** that parameter must be the byte offset
** where the NUL terminator would occur assuming the string were NUL
** terminated.  If any NUL characters occurs at byte offsets less than
** the value of the fourth parameter then the resulting string value will
** contain embedded NULs.  The result of expressions involving strings
** with embedded NULs is undefined.
**
** ^The fifth argument to the BLOB and string binding interfaces controls
** or indicates the lifetime of the object referenced by the third parameter.
** These three options exist:
** ^ (1) A destructor to dispose of the BLOB or string after SQLite has finished
** with it may be passed. ^It is called to dispose of the BLOB or string even
** if the call to the bind API fails, except the destructor is not called if
** the third parameter is a NULL pointer or the fourth parameter is negative.
** ^ (2) The special constant, [SQLITE_STATIC], may be passed to indicate that
** the application remains responsible for disposing of the object. ^In this
** case, the object and the provided pointer to it must remain valid until
** either the prepared statement is finalized or the same SQL parameter is
** bound to something else, whichever occurs sooner.
** ^ (3) The constant, [SQLITE_TRANSIENT], may be passed to indicate that the
** object is to be copied prior to the return from sqlite3_bind_*(). ^The
** object and pointer to it must remain valid until then. ^SQLite will then
** manage the lifetime of its private copy.
**
** ^The sixth argument to sqlite3_bind_text64() must be one of
** [SQLITE_UTF8], [SQLITE_UTF16], [SQLITE_UTF16BE], or [SQLITE_UTF16LE]
** to specify the encoding of the text in the third parameter.  If
** the sixth argument to sqlite3_bind_text64() is not one of the
** allowed values shown above, or if the text encoding is different
** from the encoding specified by the sixth parameter, then the behavior
** is undefined.
**
** ^The sqlite3_bind_zeroblob() routine binds a BLOB of length N that
** is filled with zeroes.  ^A zeroblob uses a fixed amount of memory
** (just an integer to hold its size) while it is being processed.
** Zeroblobs are intended to serve as placeholders for BLOBs whose
** content is later written using
** [sqlite3_blob_open | incremental BLOB I/O] routines.
** ^A negative value for the zeroblob results in a zero-length BLOB.
**
** ^The sqlite3_bind_pointer(S,I,P,T,D) routine causes the I-th parameter in
** [prepared statement] S to have an SQL value of NULL, but to also be
** associated with the pointer P of type T.  ^D is either a NULL pointer or
** a pointer to a destructor function for P. ^SQLite will invoke the
** destructor D with a single argument of P when it is finished using
** P.  The T parameter should be a static string, preferably a string
** literal. The sqlite3_bind_pointer() routine is part of the
** [pointer passing interface] added for SQLite 3.20.0.
**
** ^If any of the sqlite3_bind_*() routines are called with a NULL pointer
** for the [prepared statement] or with a prepared statement for which
** [sqlite3_step()] has been called more recently than [sqlite3_reset()],
** then the call will return [SQLITE_MISUSE].  If any sqlite3_bind_()
** routine is passed a [prepared statement] that has been finalized, the
** result is undefined and probably harmful.
**
** ^Bindings are not cleared by the [sqlite3_reset()] routine.
** ^Unbound parameters are interpreted as NULL.
**
** ^The sqlite3_bind_* routines return [SQLITE_OK] on success or an
** [error code] if anything goes wrong.
** ^[SQLITE_TOOBIG] might be returned if the size of a string or BLOB
** exceeds limits imposed by [sqlite3_limit]([SQLITE_LIMIT_LENGTH]) or
** [SQLITE_MAX_LENGTH].
** ^[SQLITE_RANGE] is returned if the parameter
** index is out of range.  ^[SQLITE_NOMEM] is returned if malloc() fails.
**
** See also: [sqlite3_bind_parameter_count()],
** [sqlite3_bind_parameter_name()], and [sqlite3_bind_parameter_index()].
 */
// llgo:link (*Stmt).BindBlob C.sqlite3_bind_blob
func (recv_ *Stmt) BindBlob(__llgo_arg_0 c.Int, __llgo_arg_1 c.Pointer, n c.Int, __llgo_arg_3 func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Stmt).BindBlob64 C.sqlite3_bind_blob64
func (recv_ *Stmt) BindBlob64(c.Int, c.Pointer, Uint64, func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Stmt).BindDouble C.sqlite3_bind_double
func (recv_ *Stmt) BindDouble(c.Int, c.Double) c.Int {
	return 0
}

// llgo:link (*Stmt).BindInt C.sqlite3_bind_int
func (recv_ *Stmt) BindInt(c.Int, c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).BindInt64 C.sqlite3_bind_int64
func (recv_ *Stmt) BindInt64(c.Int, Int64) c.Int {
	return 0
}

// llgo:link (*Stmt).BindNull C.sqlite3_bind_null
func (recv_ *Stmt) BindNull(c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).BindText C.sqlite3_bind_text
func (recv_ *Stmt) BindText(c.Int, *c.Char, c.Int, func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Stmt).BindText16 C.sqlite3_bind_text16
func (recv_ *Stmt) BindText16(c.Int, c.Pointer, c.Int, func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Stmt).BindText64 C.sqlite3_bind_text64
func (recv_ *Stmt) BindText64(__llgo_arg_0 c.Int, __llgo_arg_1 *c.Char, __llgo_arg_2 Uint64, __llgo_arg_3 func(c.Pointer), encoding c.Char) c.Int {
	return 0
}

// llgo:link (*Stmt).BindValue C.sqlite3_bind_value
func (recv_ *Stmt) BindValue(c.Int, *Value) c.Int {
	return 0
}

// llgo:link (*Stmt).BindPointer C.sqlite3_bind_pointer
func (recv_ *Stmt) BindPointer(c.Int, c.Pointer, *c.Char, func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Stmt).BindZeroblob C.sqlite3_bind_zeroblob
func (recv_ *Stmt) BindZeroblob(__llgo_arg_0 c.Int, n c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).BindZeroblob64 C.sqlite3_bind_zeroblob64
func (recv_ *Stmt) BindZeroblob64(c.Int, Uint64) c.Int {
	return 0
}

/*
** CAPI3REF: Number Of SQL Parameters
** METHOD: sqlite3_stmt
**
** ^This routine can be used to find the number of [SQL parameters]
** in a [prepared statement].  SQL parameters are tokens of the
** form "?", "?NNN", ":AAA", "$AAA", or "@AAA" that serve as
** placeholders for values that are [sqlite3_bind_blob | bound]
** to the parameters at a later time.
**
** ^(This routine actually returns the index of the largest (rightmost)
** parameter. For all forms except ?NNN, this will correspond to the
** number of unique parameters.  If parameters of the ?NNN form are used,
** there may be gaps in the list.)^
**
** See also: [sqlite3_bind_blob|sqlite3_bind()],
** [sqlite3_bind_parameter_name()], and
** [sqlite3_bind_parameter_index()].
 */
// llgo:link (*Stmt).BindParameterCount C.sqlite3_bind_parameter_count
func (recv_ *Stmt) BindParameterCount() c.Int {
	return 0
}

/*
** CAPI3REF: Name Of A Host Parameter
** METHOD: sqlite3_stmt
**
** ^The sqlite3_bind_parameter_name(P,N) interface returns
** the name of the N-th [SQL parameter] in the [prepared statement] P.
** ^(SQL parameters of the form "?NNN" or ":AAA" or "@AAA" or "$AAA"
** have a name which is the string "?NNN" or ":AAA" or "@AAA" or "$AAA"
** respectively.
** In other words, the initial ":" or "$" or "@" or "?"
** is included as part of the name.)^
** ^Parameters of the form "?" without a following integer have no name
** and are referred to as "nameless" or "anonymous parameters".
**
** ^The first host parameter has an index of 1, not 0.
**
** ^If the value N is out of range or if the N-th parameter is
** nameless, then NULL is returned.  ^The returned string is
** always in UTF-8 encoding even if the named parameter was
** originally specified as UTF-16 in [sqlite3_prepare16()],
** [sqlite3_prepare16_v2()], or [sqlite3_prepare16_v3()].
**
** See also: [sqlite3_bind_blob|sqlite3_bind()],
** [sqlite3_bind_parameter_count()], and
** [sqlite3_bind_parameter_index()].
 */
// llgo:link (*Stmt).BindParameterName C.sqlite3_bind_parameter_name
func (recv_ *Stmt) BindParameterName(c.Int) *c.Char {
	return nil
}

/*
** CAPI3REF: Index Of A Parameter With A Given Name
** METHOD: sqlite3_stmt
**
** ^Return the index of an SQL parameter given its name.  ^The
** index value returned is suitable for use as the second
** parameter to [sqlite3_bind_blob|sqlite3_bind()].  ^A zero
** is returned if no matching parameter is found.  ^The parameter
** name must be given in UTF-8 even if the original statement
** was prepared from UTF-16 text using [sqlite3_prepare16_v2()] or
** [sqlite3_prepare16_v3()].
**
** See also: [sqlite3_bind_blob|sqlite3_bind()],
** [sqlite3_bind_parameter_count()], and
** [sqlite3_bind_parameter_name()].
 */
// llgo:link (*Stmt).BindParameterIndex C.sqlite3_bind_parameter_index
func (recv_ *Stmt) BindParameterIndex(zName *c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Reset All Bindings On A Prepared Statement
** METHOD: sqlite3_stmt
**
** ^Contrary to the intuition of many, [sqlite3_reset()] does not reset
** the [sqlite3_bind_blob | bindings] on a [prepared statement].
** ^Use this routine to reset all host parameters to NULL.
 */
// llgo:link (*Stmt).ClearBindings C.sqlite3_clear_bindings
func (recv_ *Stmt) ClearBindings() c.Int {
	return 0
}

/*
** CAPI3REF: Number Of Columns In A Result Set
** METHOD: sqlite3_stmt
**
** ^Return the number of columns in the result set returned by the
** [prepared statement]. ^If this routine returns 0, that means the
** [prepared statement] returns no data (for example an [UPDATE]).
** ^However, just because this routine returns a positive number does not
** mean that one or more rows of data will be returned.  ^A SELECT statement
** will always have a positive sqlite3_column_count() but depending on the
** WHERE clause constraints and the table content, it might return no rows.
**
** See also: [sqlite3_data_count()]
 */
// llgo:link (*Stmt).ColumnCount C.sqlite3_column_count
func (recv_ *Stmt) ColumnCount() c.Int {
	return 0
}

/*
** CAPI3REF: Column Names In A Result Set
** METHOD: sqlite3_stmt
**
** ^These routines return the name assigned to a particular column
** in the result set of a [SELECT] statement.  ^The sqlite3_column_name()
** interface returns a pointer to a zero-terminated UTF-8 string
** and sqlite3_column_name16() returns a pointer to a zero-terminated
** UTF-16 string.  ^The first parameter is the [prepared statement]
** that implements the [SELECT] statement. ^The second parameter is the
** column number.  ^The leftmost column is number 0.
**
** ^The returned string pointer is valid until either the [prepared statement]
** is destroyed by [sqlite3_finalize()] or until the statement is automatically
** reprepared by the first call to [sqlite3_step()] for a particular run
** or until the next call to
** sqlite3_column_name() or sqlite3_column_name16() on the same column.
**
** ^If sqlite3_malloc() fails during the processing of either routine
** (for example during a conversion from UTF-8 to UTF-16) then a
** NULL pointer is returned.
**
** ^The name of a result column is the value of the "AS" clause for
** that column, if there is an AS clause.  If there is no AS clause
** then the name of the column is unspecified and may change from
** one release of SQLite to the next.
 */
// llgo:link (*Stmt).ColumnName C.sqlite3_column_name
func (recv_ *Stmt) ColumnName(N c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnName16 C.sqlite3_column_name16
func (recv_ *Stmt) ColumnName16(N c.Int) c.Pointer {
	return nil
}

/*
** CAPI3REF: Source Of Data In A Query Result
** METHOD: sqlite3_stmt
**
** ^These routines provide a means to determine the database, table, and
** table column that is the origin of a particular result column in
** [SELECT] statement.
** ^The name of the database or table or column can be returned as
** either a UTF-8 or UTF-16 string.  ^The _database_ routines return
** the database name, the _table_ routines return the table name, and
** the origin_ routines return the column name.
** ^The returned string is valid until the [prepared statement] is destroyed
** using [sqlite3_finalize()] or until the statement is automatically
** reprepared by the first call to [sqlite3_step()] for a particular run
** or until the same information is requested
** again in a different encoding.
**
** ^The names returned are the original un-aliased names of the
** database, table, and column.
**
** ^The first argument to these interfaces is a [prepared statement].
** ^These functions return information about the Nth result column returned by
** the statement, where N is the second function argument.
** ^The left-most column is column 0 for these routines.
**
** ^If the Nth column returned by the statement is an expression or
** subquery and is not a column value, then all of these functions return
** NULL.  ^These routines might also return NULL if a memory allocation error
** occurs.  ^Otherwise, they return the name of the attached database, table,
** or column that query result column was extracted from.
**
** ^As with all other SQLite APIs, those whose names end with "16" return
** UTF-16 encoded strings and the other functions return UTF-8.
**
** ^These APIs are only available if the library was compiled with the
** [SQLITE_ENABLE_COLUMN_METADATA] C-preprocessor symbol.
**
** If two or more threads call one or more
** [sqlite3_column_database_name | column metadata interfaces]
** for the same [prepared statement] and result column
** at the same time then the results are undefined.
 */
// llgo:link (*Stmt).ColumnDatabaseName C.sqlite3_column_database_name
func (recv_ *Stmt) ColumnDatabaseName(c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnDatabaseName16 C.sqlite3_column_database_name16
func (recv_ *Stmt) ColumnDatabaseName16(c.Int) c.Pointer {
	return nil
}

// llgo:link (*Stmt).ColumnTableName C.sqlite3_column_table_name
func (recv_ *Stmt) ColumnTableName(c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnTableName16 C.sqlite3_column_table_name16
func (recv_ *Stmt) ColumnTableName16(c.Int) c.Pointer {
	return nil
}

// llgo:link (*Stmt).ColumnOriginName C.sqlite3_column_origin_name
func (recv_ *Stmt) ColumnOriginName(c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnOriginName16 C.sqlite3_column_origin_name16
func (recv_ *Stmt) ColumnOriginName16(c.Int) c.Pointer {
	return nil
}

/*
** CAPI3REF: Declared Datatype Of A Query Result
** METHOD: sqlite3_stmt
**
** ^(The first parameter is a [prepared statement].
** If this statement is a [SELECT] statement and the Nth column of the
** returned result set of that [SELECT] is a table column (not an
** expression or subquery) then the declared type of the table
** column is returned.)^  ^If the Nth column of the result set is an
** expression or subquery, then a NULL pointer is returned.
** ^The returned string is always UTF-8 encoded.
**
** ^(For example, given the database schema:
**
** CREATE TABLE t1(c1 VARIANT);
**
** and the following statement to be compiled:
**
** SELECT c1 + 1, c1 FROM t1;
**
** this routine would return the string "VARIANT" for the second result
** column (i==1), and a NULL pointer for the first result column (i==0).)^
**
** ^SQLite uses dynamic run-time typing.  ^So just because a column
** is declared to contain a particular type does not mean that the
** data stored in that column is of the declared type.  SQLite is
** strongly typed, but the typing is dynamic not static.  ^Type
** is associated with individual values, not with the containers
** used to hold those values.
 */
// llgo:link (*Stmt).ColumnDecltype C.sqlite3_column_decltype
func (recv_ *Stmt) ColumnDecltype(c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnDecltype16 C.sqlite3_column_decltype16
func (recv_ *Stmt) ColumnDecltype16(c.Int) c.Pointer {
	return nil
}

/*
** CAPI3REF: Evaluate An SQL Statement
** METHOD: sqlite3_stmt
**
** After a [prepared statement] has been prepared using any of
** [sqlite3_prepare_v2()], [sqlite3_prepare_v3()], [sqlite3_prepare16_v2()],
** or [sqlite3_prepare16_v3()] or one of the legacy
** interfaces [sqlite3_prepare()] or [sqlite3_prepare16()], this function
** must be called one or more times to evaluate the statement.
**
** The details of the behavior of the sqlite3_step() interface depend
** on whether the statement was prepared using the newer "vX" interfaces
** [sqlite3_prepare_v3()], [sqlite3_prepare_v2()], [sqlite3_prepare16_v3()],
** [sqlite3_prepare16_v2()] or the older legacy
** interfaces [sqlite3_prepare()] and [sqlite3_prepare16()].  The use of the
** new "vX" interface is recommended for new applications but the legacy
** interface will continue to be supported.
**
** ^In the legacy interface, the return value will be either [SQLITE_BUSY],
** [SQLITE_DONE], [SQLITE_ROW], [SQLITE_ERROR], or [SQLITE_MISUSE].
** ^With the "v2" interface, any of the other [result codes] or
** [extended result codes] might be returned as well.
**
** ^[SQLITE_BUSY] means that the database engine was unable to acquire the
** database locks it needs to do its job.  ^If the statement is a [COMMIT]
** or occurs outside of an explicit transaction, then you can retry the
** statement.  If the statement is not a [COMMIT] and occurs within an
** explicit transaction then you should rollback the transaction before
** continuing.
**
** ^[SQLITE_DONE] means that the statement has finished executing
** successfully.  sqlite3_step() should not be called again on this virtual
** machine without first calling [sqlite3_reset()] to reset the virtual
** machine back to its initial state.
**
** ^If the SQL statement being executed returns any data, then [SQLITE_ROW]
** is returned each time a new row of data is ready for processing by the
** caller. The values may be accessed using the [column access functions].
** sqlite3_step() is called again to retrieve the next row of data.
**
** ^[SQLITE_ERROR] means that a run-time error (such as a constraint
** violation) has occurred.  sqlite3_step() should not be called again on
** the VM. More information may be found by calling [sqlite3_errmsg()].
** ^With the legacy interface, a more specific error code (for example,
** [SQLITE_INTERRUPT], [SQLITE_SCHEMA], [SQLITE_CORRUPT], and so forth)
** can be obtained by calling [sqlite3_reset()] on the
** [prepared statement].  ^In the "v2" interface,
** the more specific error code is returned directly by sqlite3_step().
**
** [SQLITE_MISUSE] means that the this routine was called inappropriately.
** Perhaps it was called on a [prepared statement] that has
** already been [sqlite3_finalize | finalized] or on one that had
** previously returned [SQLITE_ERROR] or [SQLITE_DONE].  Or it could
** be the case that the same database connection is being used by two or
** more threads at the same moment in time.
**
** For all versions of SQLite up to and including 3.6.23.1, a call to
** [sqlite3_reset()] was required after sqlite3_step() returned anything
** other than [SQLITE_ROW] before any subsequent invocation of
** sqlite3_step().  Failure to reset the prepared statement using
** [sqlite3_reset()] would result in an [SQLITE_MISUSE] return from
** sqlite3_step().  But after [version 3.6.23.1] ([dateof:3.6.23.1],
** sqlite3_step() began
** calling [sqlite3_reset()] automatically in this circumstance rather
** than returning [SQLITE_MISUSE].  This is not considered a compatibility
** break because any application that ever receives an SQLITE_MISUSE error
** is broken by definition.  The [SQLITE_OMIT_AUTORESET] compile-time option
** can be used to restore the legacy behavior.
**
** <b>Goofy Interface Alert:</b> In the legacy interface, the sqlite3_step()
** API always returns a generic error code, [SQLITE_ERROR], following any
** error other than [SQLITE_BUSY] and [SQLITE_MISUSE].  You must call
** [sqlite3_reset()] or [sqlite3_finalize()] in order to find one of the
** specific [error codes] that better describes the error.
** We admit that this is a goofy design.  The problem has been fixed
** with the "v2" interface.  If you prepare all of your SQL statements
** using [sqlite3_prepare_v3()] or [sqlite3_prepare_v2()]
** or [sqlite3_prepare16_v2()] or [sqlite3_prepare16_v3()] instead
** of the legacy [sqlite3_prepare()] and [sqlite3_prepare16()] interfaces,
** then the more specific [error codes] are returned directly
** by sqlite3_step().  The use of the "vX" interfaces is recommended.
 */
// llgo:link (*Stmt).Step C.sqlite3_step
func (recv_ *Stmt) Step() c.Int {
	return 0
}

/*
** CAPI3REF: Number of columns in a result set
** METHOD: sqlite3_stmt
**
** ^The sqlite3_data_count(P) interface returns the number of columns in the
** current row of the result set of [prepared statement] P.
** ^If prepared statement P does not have results ready to return
** (via calls to the [sqlite3_column_int | sqlite3_column()] family of
** interfaces) then sqlite3_data_count(P) returns 0.
** ^The sqlite3_data_count(P) routine also returns 0 if P is a NULL pointer.
** ^The sqlite3_data_count(P) routine returns 0 if the previous call to
** [sqlite3_step](P) returned [SQLITE_DONE].  ^The sqlite3_data_count(P)
** will return non-zero if previous call to [sqlite3_step](P) returned
** [SQLITE_ROW], except in the case of the [PRAGMA incremental_vacuum]
** where it always returns zero since each step of that multi-step
** pragma returns 0 columns of data.
**
** See also: [sqlite3_column_count()]
 */
// llgo:link (*Stmt).DataCount C.sqlite3_data_count
func (recv_ *Stmt) DataCount() c.Int {
	return 0
}

/*
** CAPI3REF: Result Values From A Query
** KEYWORDS: {column access functions}
** METHOD: sqlite3_stmt
**
** <b>Summary:</b>
** <blockquote><table border=0 cellpadding=0 cellspacing=0>
** <tr><td><b>sqlite3_column_blob</b><td>&rarr;<td>BLOB result
** <tr><td><b>sqlite3_column_double</b><td>&rarr;<td>REAL result
** <tr><td><b>sqlite3_column_int</b><td>&rarr;<td>32-bit INTEGER result
** <tr><td><b>sqlite3_column_int64</b><td>&rarr;<td>64-bit INTEGER result
** <tr><td><b>sqlite3_column_text</b><td>&rarr;<td>UTF-8 TEXT result
** <tr><td><b>sqlite3_column_text16</b><td>&rarr;<td>UTF-16 TEXT result
** <tr><td><b>sqlite3_column_value</b><td>&rarr;<td>The result as an
** [sqlite3_value|unprotected sqlite3_value] object.
** <tr><td>&nbsp;<td>&nbsp;<td>&nbsp;
** <tr><td><b>sqlite3_column_bytes</b><td>&rarr;<td>Size of a BLOB
** or a UTF-8 TEXT result in bytes
** <tr><td><b>sqlite3_column_bytes16&nbsp;&nbsp;</b>
** <td>&rarr;&nbsp;&nbsp;<td>Size of UTF-16
** TEXT in bytes
** <tr><td><b>sqlite3_column_type</b><td>&rarr;<td>Default
** datatype of the result
** </table></blockquote>
**
** <b>Details:</b>
**
** ^These routines return information about a single column of the current
** result row of a query.  ^In every case the first argument is a pointer
** to the [prepared statement] that is being evaluated (the [sqlite3_stmt*]
** that was returned from [sqlite3_prepare_v2()] or one of its variants)
** and the second argument is the index of the column for which information
** should be returned. ^The leftmost column of the result set has the index 0.
** ^The number of columns in the result can be determined using
** [sqlite3_column_count()].
**
** If the SQL statement does not currently point to a valid row, or if the
** column index is out of range, the result is undefined.
** These routines may only be called when the most recent call to
** [sqlite3_step()] has returned [SQLITE_ROW] and neither
** [sqlite3_reset()] nor [sqlite3_finalize()] have been called subsequently.
** If any of these routines are called after [sqlite3_reset()] or
** [sqlite3_finalize()] or after [sqlite3_step()] has returned
** something other than [SQLITE_ROW], the results are undefined.
** If [sqlite3_step()] or [sqlite3_reset()] or [sqlite3_finalize()]
** are called from a different thread while any of these routines
** are pending, then the results are undefined.
**
** The first six interfaces (_blob, _double, _int, _int64, _text, and _text16)
** each return the value of a result column in a specific data format.  If
** the result column is not initially in the requested format (for example,
** if the query returns an integer but the sqlite3_column_text() interface
** is used to extract the value) then an automatic type conversion is performed.
**
** ^The sqlite3_column_type() routine returns the
** [SQLITE_INTEGER | datatype code] for the initial data type
** of the result column.  ^The returned value is one of [SQLITE_INTEGER],
** [SQLITE_FLOAT], [SQLITE_TEXT], [SQLITE_BLOB], or [SQLITE_NULL].
** The return value of sqlite3_column_type() can be used to decide which
** of the first six interface should be used to extract the column value.
** The value returned by sqlite3_column_type() is only meaningful if no
** automatic type conversions have occurred for the value in question.
** After a type conversion, the result of calling sqlite3_column_type()
** is undefined, though harmless.  Future
** versions of SQLite may change the behavior of sqlite3_column_type()
** following a type conversion.
**
** If the result is a BLOB or a TEXT string, then the sqlite3_column_bytes()
** or sqlite3_column_bytes16() interfaces can be used to determine the size
** of that BLOB or string.
**
** ^If the result is a BLOB or UTF-8 string then the sqlite3_column_bytes()
** routine returns the number of bytes in that BLOB or string.
** ^If the result is a UTF-16 string, then sqlite3_column_bytes() converts
** the string to UTF-8 and then returns the number of bytes.
** ^If the result is a numeric value then sqlite3_column_bytes() uses
** [sqlite3_snprintf()] to convert that value to a UTF-8 string and returns
** the number of bytes in that string.
** ^If the result is NULL, then sqlite3_column_bytes() returns zero.
**
** ^If the result is a BLOB or UTF-16 string then the sqlite3_column_bytes16()
** routine returns the number of bytes in that BLOB or string.
** ^If the result is a UTF-8 string, then sqlite3_column_bytes16() converts
** the string to UTF-16 and then returns the number of bytes.
** ^If the result is a numeric value then sqlite3_column_bytes16() uses
** [sqlite3_snprintf()] to convert that value to a UTF-16 string and returns
** the number of bytes in that string.
** ^If the result is NULL, then sqlite3_column_bytes16() returns zero.
**
** ^The values returned by [sqlite3_column_bytes()] and
** [sqlite3_column_bytes16()] do not include the zero terminators at the end
** of the string.  ^For clarity: the values returned by
** [sqlite3_column_bytes()] and [sqlite3_column_bytes16()] are the number of
** bytes in the string, not the number of characters.
**
** ^Strings returned by sqlite3_column_text() and sqlite3_column_text16(),
** even empty strings, are always zero-terminated.  ^The return
** value from sqlite3_column_blob() for a zero-length BLOB is a NULL pointer.
**
** ^Strings returned by sqlite3_column_text16() always have the endianness
** which is native to the platform, regardless of the text encoding set
** for the database.
**
** <b>Warning:</b> ^The object returned by [sqlite3_column_value()] is an
** [unprotected sqlite3_value] object.  In a multithreaded environment,
** an unprotected sqlite3_value object may only be used safely with
** [sqlite3_bind_value()] and [sqlite3_result_value()].
** If the [unprotected sqlite3_value] object returned by
** [sqlite3_column_value()] is used in any other way, including calls
** to routines like [sqlite3_value_int()], [sqlite3_value_text()],
** or [sqlite3_value_bytes()], the behavior is not threadsafe.
** Hence, the sqlite3_column_value() interface
** is normally only useful within the implementation of
** [application-defined SQL functions] or [virtual tables], not within
** top-level application code.
**
** These routines may attempt to convert the datatype of the result.
** ^For example, if the internal representation is FLOAT and a text result
** is requested, [sqlite3_snprintf()] is used internally to perform the
** conversion automatically.  ^(The following table details the conversions
** that are applied:
**
** <blockquote>
** <table border="1">
** <tr><th> Internal<br>Type <th> Requested<br>Type <th>  Conversion
**
** <tr><td>  NULL    <td> INTEGER   <td> Result is 0
** <tr><td>  NULL    <td>  FLOAT    <td> Result is 0.0
** <tr><td>  NULL    <td>   TEXT    <td> Result is a NULL pointer
** <tr><td>  NULL    <td>   BLOB    <td> Result is a NULL pointer
** <tr><td> INTEGER  <td>  FLOAT    <td> Convert from integer to float
** <tr><td> INTEGER  <td>   TEXT    <td> ASCII rendering of the integer
** <tr><td> INTEGER  <td>   BLOB    <td> Same as INTEGER->TEXT
** <tr><td>  FLOAT   <td> INTEGER   <td> [CAST] to INTEGER
** <tr><td>  FLOAT   <td>   TEXT    <td> ASCII rendering of the float
** <tr><td>  FLOAT   <td>   BLOB    <td> [CAST] to BLOB
** <tr><td>  TEXT    <td> INTEGER   <td> [CAST] to INTEGER
** <tr><td>  TEXT    <td>  FLOAT    <td> [CAST] to REAL
** <tr><td>  TEXT    <td>   BLOB    <td> No change
** <tr><td>  BLOB    <td> INTEGER   <td> [CAST] to INTEGER
** <tr><td>  BLOB    <td>  FLOAT    <td> [CAST] to REAL
** <tr><td>  BLOB    <td>   TEXT    <td> [CAST] to TEXT, ensure zero terminator
** </table>
** </blockquote>)^
**
** Note that when type conversions occur, pointers returned by prior
** calls to sqlite3_column_blob(), sqlite3_column_text(), and/or
** sqlite3_column_text16() may be invalidated.
** Type conversions and pointer invalidations might occur
** in the following cases:
**
** <ul>
** <li> The initial content is a BLOB and sqlite3_column_text() or
**      sqlite3_column_text16() is called.  A zero-terminator might
**      need to be added to the string.</li>
** <li> The initial content is UTF-8 text and sqlite3_column_bytes16() or
**      sqlite3_column_text16() is called.  The content must be converted
**      to UTF-16.</li>
** <li> The initial content is UTF-16 text and sqlite3_column_bytes() or
**      sqlite3_column_text() is called.  The content must be converted
**      to UTF-8.</li>
** </ul>
**
** ^Conversions between UTF-16be and UTF-16le are always done in place and do
** not invalidate a prior pointer, though of course the content of the buffer
** that the prior pointer references will have been modified.  Other kinds
** of conversion are done in place when it is possible, but sometimes they
** are not possible and in those cases prior pointers are invalidated.
**
** The safest policy is to invoke these routines
** in one of the following ways:
**
** <ul>
**  <li>sqlite3_column_text() followed by sqlite3_column_bytes()</li>
**  <li>sqlite3_column_blob() followed by sqlite3_column_bytes()</li>
**  <li>sqlite3_column_text16() followed by sqlite3_column_bytes16()</li>
** </ul>
**
** In other words, you should call sqlite3_column_text(),
** sqlite3_column_blob(), or sqlite3_column_text16() first to force the result
** into the desired format, then invoke sqlite3_column_bytes() or
** sqlite3_column_bytes16() to find the size of the result.  Do not mix calls
** to sqlite3_column_text() or sqlite3_column_blob() with calls to
** sqlite3_column_bytes16(), and do not mix calls to sqlite3_column_text16()
** with calls to sqlite3_column_bytes().
**
** ^The pointers returned are valid until a type conversion occurs as
** described above, or until [sqlite3_step()] or [sqlite3_reset()] or
** [sqlite3_finalize()] is called.  ^The memory space used to hold strings
** and BLOBs is freed automatically.  Do not pass the pointers returned
** from [sqlite3_column_blob()], [sqlite3_column_text()], etc. into
** [sqlite3_free()].
**
** As long as the input parameters are correct, these routines will only
** fail if an out-of-memory error occurs during a format conversion.
** Only the following subset of interfaces are subject to out-of-memory
** errors:
**
** <ul>
** <li> sqlite3_column_blob()
** <li> sqlite3_column_text()
** <li> sqlite3_column_text16()
** <li> sqlite3_column_bytes()
** <li> sqlite3_column_bytes16()
** </ul>
**
** If an out-of-memory error occurs, then the return value from these
** routines is the same as if the column had contained an SQL NULL value.
** Valid SQL NULL returns can be distinguished from out-of-memory errors
** by invoking the [sqlite3_errcode()] immediately after the suspect
** return value is obtained and before any
** other SQLite interface is called on the same [database connection].
 */
// llgo:link (*Stmt).ColumnBlob C.sqlite3_column_blob
func (recv_ *Stmt) ColumnBlob(iCol c.Int) c.Pointer {
	return nil
}

// llgo:link (*Stmt).ColumnDouble C.sqlite3_column_double
func (recv_ *Stmt) ColumnDouble(iCol c.Int) c.Double {
	return 0
}

// llgo:link (*Stmt).ColumnInt C.sqlite3_column_int
func (recv_ *Stmt) ColumnInt(iCol c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).ColumnInt64 C.sqlite3_column_int64
func (recv_ *Stmt) ColumnInt64(iCol c.Int) Int64 {
	return 0
}

// llgo:link (*Stmt).ColumnText C.sqlite3_column_text
func (recv_ *Stmt) ColumnText(iCol c.Int) *c.Char {
	return nil
}

// llgo:link (*Stmt).ColumnText16 C.sqlite3_column_text16
func (recv_ *Stmt) ColumnText16(iCol c.Int) c.Pointer {
	return nil
}

// llgo:link (*Stmt).ColumnValue C.sqlite3_column_value
func (recv_ *Stmt) ColumnValue(iCol c.Int) *Value {
	return nil
}

// llgo:link (*Stmt).ColumnBytes C.sqlite3_column_bytes
func (recv_ *Stmt) ColumnBytes(iCol c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).ColumnBytes16 C.sqlite3_column_bytes16
func (recv_ *Stmt) ColumnBytes16(iCol c.Int) c.Int {
	return 0
}

// llgo:link (*Stmt).ColumnType C.sqlite3_column_type
func (recv_ *Stmt) ColumnType(iCol c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Destroy A Prepared Statement Object
** DESTRUCTOR: sqlite3_stmt
**
** ^The sqlite3_finalize() function is called to delete a [prepared statement].
** ^If the most recent evaluation of the statement encountered no errors
** or if the statement is never been evaluated, then sqlite3_finalize() returns
** SQLITE_OK.  ^If the most recent evaluation of statement S failed, then
** sqlite3_finalize(S) returns the appropriate [error code] or
** [extended error code].
**
** ^The sqlite3_finalize(S) routine can be called at any point during
** the life cycle of [prepared statement] S:
** before statement S is ever evaluated, after
** one or more calls to [sqlite3_reset()], or after any call
** to [sqlite3_step()] regardless of whether or not the statement has
** completed execution.
**
** ^Invoking sqlite3_finalize() on a NULL pointer is a harmless no-op.
**
** The application must finalize every [prepared statement] in order to avoid
** resource leaks.  It is a grievous error for the application to try to use
** a prepared statement after it has been finalized.  Any use of a prepared
** statement after it has been finalized can result in undefined and
** undesirable behavior such as segfaults and heap corruption.
 */
// llgo:link (*Stmt).Close C.sqlite3_finalize
func (recv_ *Stmt) Close() c.Int {
	return 0
}

/*
** CAPI3REF: Reset A Prepared Statement Object
** METHOD: sqlite3_stmt
**
** The sqlite3_reset() function is called to reset a [prepared statement]
** object back to its initial state, ready to be re-executed.
** ^Any SQL statement variables that had values bound to them using
** the [sqlite3_bind_blob | sqlite3_bind_*() API] retain their values.
** Use [sqlite3_clear_bindings()] to reset the bindings.
**
** ^The [sqlite3_reset(S)] interface resets the [prepared statement] S
** back to the beginning of its program.
**
** ^The return code from [sqlite3_reset(S)] indicates whether or not
** the previous evaluation of prepared statement S completed successfully.
** ^If [sqlite3_step(S)] has never before been called on S or if
** [sqlite3_step(S)] has not been called since the previous call
** to [sqlite3_reset(S)], then [sqlite3_reset(S)] will return
** [SQLITE_OK].
**
** ^If the most recent call to [sqlite3_step(S)] for the
** [prepared statement] S indicated an error, then
** [sqlite3_reset(S)] returns an appropriate [error code].
** ^The [sqlite3_reset(S)] interface might also return an [error code]
** if there were no prior errors but the process of resetting
** the prepared statement caused a new error. ^For example, if an
** [INSERT] statement with a [RETURNING] clause is only stepped one time,
** that one call to [sqlite3_step(S)] might return SQLITE_ROW but
** the overall statement might still fail and the [sqlite3_reset(S)] call
** might return SQLITE_BUSY if locking constraints prevent the
** database change from committing.  Therefore, it is important that
** applications check the return code from [sqlite3_reset(S)] even if
** no prior call to [sqlite3_step(S)] indicated a problem.
**
** ^The [sqlite3_reset(S)] interface does not change the values
** of any [sqlite3_bind_blob|bindings] on the [prepared statement] S.
 */
// llgo:link (*Stmt).Reset C.sqlite3_reset
func (recv_ *Stmt) Reset() c.Int {
	return 0
}

/*
** CAPI3REF: Create Or Redefine SQL Functions
** KEYWORDS: {function creation routines}
** METHOD: sqlite3
**
** ^These functions (collectively known as "function creation routines")
** are used to add SQL functions or aggregates or to redefine the behavior
** of existing SQL functions or aggregates. The only differences between
** the three "sqlite3_create_function*" routines are the text encoding
** expected for the second parameter (the name of the function being
** created) and the presence or absence of a destructor callback for
** the application data pointer. Function sqlite3_create_window_function()
** is similar, but allows the user to supply the extra callback functions
** needed by [aggregate window functions].
**
** ^The first parameter is the [database connection] to which the SQL
** function is to be added.  ^If an application uses more than one database
** connection then application-defined SQL functions must be added
** to each database connection separately.
**
** ^The second parameter is the name of the SQL function to be created or
** redefined.  ^The length of the name is limited to 255 bytes in a UTF-8
** representation, exclusive of the zero-terminator.  ^Note that the name
** length limit is in UTF-8 bytes, not characters nor UTF-16 bytes.
** ^Any attempt to create a function with a longer name
** will result in [SQLITE_MISUSE] being returned.
**
** ^The third parameter (nArg)
** is the number of arguments that the SQL function or
** aggregate takes. ^If this parameter is -1, then the SQL function or
** aggregate may take any number of arguments between 0 and the limit
** set by [sqlite3_limit]([SQLITE_LIMIT_FUNCTION_ARG]).  If the third
** parameter is less than -1 or greater than 127 then the behavior is
** undefined.
**
** ^The fourth parameter, eTextRep, specifies what
** [SQLITE_UTF8 | text encoding] this SQL function prefers for
** its parameters.  The application should set this parameter to
** [SQLITE_UTF16LE] if the function implementation invokes
** [sqlite3_value_text16le()] on an input, or [SQLITE_UTF16BE] if the
** implementation invokes [sqlite3_value_text16be()] on an input, or
** [SQLITE_UTF16] if [sqlite3_value_text16()] is used, or [SQLITE_UTF8]
** otherwise.  ^The same SQL function may be registered multiple times using
** different preferred text encodings, with different implementations for
** each encoding.
** ^When multiple implementations of the same function are available, SQLite
** will pick the one that involves the least amount of data conversion.
**
** ^The fourth parameter may optionally be ORed with [SQLITE_DETERMINISTIC]
** to signal that the function will always return the same result given
** the same inputs within a single SQL statement.  Most SQL functions are
** deterministic.  The built-in [random()] SQL function is an example of a
** function that is not deterministic.  The SQLite query planner is able to
** perform additional optimizations on deterministic functions, so use
** of the [SQLITE_DETERMINISTIC] flag is recommended where possible.
**
** ^The fourth parameter may also optionally include the [SQLITE_DIRECTONLY]
** flag, which if present prevents the function from being invoked from
** within VIEWs, TRIGGERs, CHECK constraints, generated column expressions,
** index expressions, or the WHERE clause of partial indexes.
**
** For best security, the [SQLITE_DIRECTONLY] flag is recommended for
** all application-defined SQL functions that do not need to be
** used inside of triggers, view, CHECK constraints, or other elements of
** the database schema.  This flags is especially recommended for SQL
** functions that have side effects or reveal internal application state.
** Without this flag, an attacker might be able to modify the schema of
** a database file to include invocations of the function with parameters
** chosen by the attacker, which the application will then execute when
** the database file is opened and read.
**
** ^(The fifth parameter is an arbitrary pointer.  The implementation of the
** function can gain access to this pointer using [sqlite3_user_data()].)^
**
** ^The sixth, seventh and eighth parameters passed to the three
** "sqlite3_create_function*" functions, xFunc, xStep and xFinal, are
** pointers to C-language functions that implement the SQL function or
** aggregate. ^A scalar SQL function requires an implementation of the xFunc
** callback only; NULL pointers must be passed as the xStep and xFinal
** parameters. ^An aggregate SQL function requires an implementation of xStep
** and xFinal and NULL pointer must be passed for xFunc. ^To delete an existing
** SQL function or aggregate, pass NULL pointers for all three function
** callbacks.
**
** ^The sixth, seventh, eighth and ninth parameters (xStep, xFinal, xValue
** and xInverse) passed to sqlite3_create_window_function are pointers to
** C-language callbacks that implement the new function. xStep and xFinal
** must both be non-NULL. xValue and xInverse may either both be NULL, in
** which case a regular aggregate function is created, or must both be
** non-NULL, in which case the new function may be used as either an aggregate
** or aggregate window function. More details regarding the implementation
** of aggregate window functions are
** [user-defined window functions|available here].
**
** ^(If the final parameter to sqlite3_create_function_v2() or
** sqlite3_create_window_function() is not NULL, then it is destructor for
** the application data pointer. The destructor is invoked when the function
** is deleted, either by being overloaded or when the database connection
** closes.)^ ^The destructor is also invoked if the call to
** sqlite3_create_function_v2() fails.  ^When the destructor callback is
** invoked, it is passed a single argument which is a copy of the application
** data pointer which was the fifth parameter to sqlite3_create_function_v2().
**
** ^It is permitted to register multiple implementations of the same
** functions with the same name but with either differing numbers of
** arguments or differing preferred text encodings.  ^SQLite will use
** the implementation that most closely matches the way in which the
** SQL function is used.  ^A function implementation with a non-negative
** nArg parameter is a better match than a function implementation with
** a negative nArg.  ^A function where the preferred text encoding
** matches the database encoding is a better
** match than a function where the encoding is different.
** ^A function where the encoding difference is between UTF16le and UTF16be
** is a closer match than a function where the encoding difference is
** between UTF8 and UTF16.
**
** ^Built-in functions may be overloaded by new application-defined functions.
**
** ^An application-defined function is permitted to call other
** SQLite interfaces.  However, such calls must not
** close the database connection nor finalize or reset the prepared
** statement in which the function is running.
 */
// llgo:link (*Sqlite3).CreateFunction C.sqlite3_create_function
func (recv_ *Sqlite3) CreateFunction(zFunctionName *c.Char, nArg c.Int, eTextRep c.Int, pApp c.Pointer, xFunc func(*Context, c.Int, **Value), xStep func(*Context, c.Int, **Value), xFinal func(*Context)) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateFunction16 C.sqlite3_create_function16
func (recv_ *Sqlite3) CreateFunction16(zFunctionName c.Pointer, nArg c.Int, eTextRep c.Int, pApp c.Pointer, xFunc func(*Context, c.Int, **Value), xStep func(*Context, c.Int, **Value), xFinal func(*Context)) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateFunctionV2 C.sqlite3_create_function_v2
func (recv_ *Sqlite3) CreateFunctionV2(zFunctionName *c.Char, nArg c.Int, eTextRep c.Int, pApp c.Pointer, xFunc func(*Context, c.Int, **Value), xStep func(*Context, c.Int, **Value), xFinal func(*Context), xDestroy func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateWindowFunction C.sqlite3_create_window_function
func (recv_ *Sqlite3) CreateWindowFunction(zFunctionName *c.Char, nArg c.Int, eTextRep c.Int, pApp c.Pointer, xStep func(*Context, c.Int, **Value), xFinal func(*Context), xValue func(*Context), xInverse func(*Context, c.Int, **Value), xDestroy func(c.Pointer)) c.Int {
	return 0
}

/*
** CAPI3REF: Deprecated Functions
** DEPRECATED
**
** These functions are [deprecated].  In order to maintain
** backwards compatibility with older code, these functions continue
** to be supported.  However, new applications should avoid
** the use of these functions.  To encourage programmers to avoid
** these functions, we will not explain what they do.
 */
// llgo:link (*Context).AggregateCount C.sqlite3_aggregate_count
func (recv_ *Context) AggregateCount() c.Int {
	return 0
}

// llgo:link (*Stmt).Expired C.sqlite3_expired
func (recv_ *Stmt) Expired() c.Int {
	return 0
}

// llgo:link (*Stmt).TransferBindings C.sqlite3_transfer_bindings
func (recv_ *Stmt) TransferBindings(*Stmt) c.Int {
	return 0
}

//go:linkname GlobalRecover C.sqlite3_global_recover
func GlobalRecover() c.Int

//go:linkname ThreadCleanup C.sqlite3_thread_cleanup
func ThreadCleanup()

//go:linkname MemoryAlarm C.sqlite3_memory_alarm
func MemoryAlarm(func(c.Pointer, Int64, c.Int), c.Pointer, Int64) c.Int

/*
** CAPI3REF: Obtaining SQL Values
** METHOD: sqlite3_value
**
** <b>Summary:</b>
** <blockquote><table border=0 cellpadding=0 cellspacing=0>
** <tr><td><b>sqlite3_value_blob</b><td>&rarr;<td>BLOB value
** <tr><td><b>sqlite3_value_double</b><td>&rarr;<td>REAL value
** <tr><td><b>sqlite3_value_int</b><td>&rarr;<td>32-bit INTEGER value
** <tr><td><b>sqlite3_value_int64</b><td>&rarr;<td>64-bit INTEGER value
** <tr><td><b>sqlite3_value_pointer</b><td>&rarr;<td>Pointer value
** <tr><td><b>sqlite3_value_text</b><td>&rarr;<td>UTF-8 TEXT value
** <tr><td><b>sqlite3_value_text16</b><td>&rarr;<td>UTF-16 TEXT value in
** the native byteorder
** <tr><td><b>sqlite3_value_text16be</b><td>&rarr;<td>UTF-16be TEXT value
** <tr><td><b>sqlite3_value_text16le</b><td>&rarr;<td>UTF-16le TEXT value
** <tr><td>&nbsp;<td>&nbsp;<td>&nbsp;
** <tr><td><b>sqlite3_value_bytes</b><td>&rarr;<td>Size of a BLOB
** or a UTF-8 TEXT in bytes
** <tr><td><b>sqlite3_value_bytes16&nbsp;&nbsp;</b>
** <td>&rarr;&nbsp;&nbsp;<td>Size of UTF-16
** TEXT in bytes
** <tr><td><b>sqlite3_value_type</b><td>&rarr;<td>Default
** datatype of the value
** <tr><td><b>sqlite3_value_numeric_type&nbsp;&nbsp;</b>
** <td>&rarr;&nbsp;&nbsp;<td>Best numeric datatype of the value
** <tr><td><b>sqlite3_value_nochange&nbsp;&nbsp;</b>
** <td>&rarr;&nbsp;&nbsp;<td>True if the column is unchanged in an UPDATE
** against a virtual table.
** <tr><td><b>sqlite3_value_frombind&nbsp;&nbsp;</b>
** <td>&rarr;&nbsp;&nbsp;<td>True if value originated from a [bound parameter]
** </table></blockquote>
**
** <b>Details:</b>
**
** These routines extract type, size, and content information from
** [protected sqlite3_value] objects.  Protected sqlite3_value objects
** are used to pass parameter information into the functions that
** implement [application-defined SQL functions] and [virtual tables].
**
** These routines work only with [protected sqlite3_value] objects.
** Any attempt to use these routines on an [unprotected sqlite3_value]
** is not threadsafe.
**
** ^These routines work just like the corresponding [column access functions]
** except that these routines take a single [protected sqlite3_value] object
** pointer instead of a [sqlite3_stmt*] pointer and an integer column number.
**
** ^The sqlite3_value_text16() interface extracts a UTF-16 string
** in the native byte-order of the host machine.  ^The
** sqlite3_value_text16be() and sqlite3_value_text16le() interfaces
** extract UTF-16 strings as big-endian and little-endian respectively.
**
** ^If [sqlite3_value] object V was initialized
** using [sqlite3_bind_pointer(S,I,P,X,D)] or [sqlite3_result_pointer(C,P,X,D)]
** and if X and Y are strings that compare equal according to strcmp(X,Y),
** then sqlite3_value_pointer(V,Y) will return the pointer P.  ^Otherwise,
** sqlite3_value_pointer(V,Y) returns a NULL. The sqlite3_bind_pointer()
** routine is part of the [pointer passing interface] added for SQLite 3.20.0.
**
** ^(The sqlite3_value_type(V) interface returns the
** [SQLITE_INTEGER | datatype code] for the initial datatype of the
** [sqlite3_value] object V. The returned value is one of [SQLITE_INTEGER],
** [SQLITE_FLOAT], [SQLITE_TEXT], [SQLITE_BLOB], or [SQLITE_NULL].)^
** Other interfaces might change the datatype for an sqlite3_value object.
** For example, if the datatype is initially SQLITE_INTEGER and
** sqlite3_value_text(V) is called to extract a text value for that
** integer, then subsequent calls to sqlite3_value_type(V) might return
** SQLITE_TEXT.  Whether or not a persistent internal datatype conversion
** occurs is undefined and may change from one release of SQLite to the next.
**
** ^(The sqlite3_value_numeric_type() interface attempts to apply
** numeric affinity to the value.  This means that an attempt is
** made to convert the value to an integer or floating point.  If
** such a conversion is possible without loss of information (in other
** words, if the value is a string that looks like a number)
** then the conversion is performed.  Otherwise no conversion occurs.
** The [SQLITE_INTEGER | datatype] after conversion is returned.)^
**
** ^Within the [xUpdate] method of a [virtual table], the
** sqlite3_value_nochange(X) interface returns true if and only if
** the column corresponding to X is unchanged by the UPDATE operation
** that the xUpdate method call was invoked to implement and if
** and the prior [xColumn] method call that was invoked to extracted
** the value for that column returned without setting a result (probably
** because it queried [sqlite3_vtab_nochange()] and found that the column
** was unchanging).  ^Within an [xUpdate] method, any value for which
** sqlite3_value_nochange(X) is true will in all other respects appear
** to be a NULL value.  If sqlite3_value_nochange(X) is invoked anywhere other
** than within an [xUpdate] method call for an UPDATE statement, then
** the return value is arbitrary and meaningless.
**
** ^The sqlite3_value_frombind(X) interface returns non-zero if the
** value X originated from one of the [sqlite3_bind_int|sqlite3_bind()]
** interfaces.  ^If X comes from an SQL literal value, or a table column,
** or an expression, then sqlite3_value_frombind(X) returns zero.
**
** Please pay particular attention to the fact that the pointer returned
** from [sqlite3_value_blob()], [sqlite3_value_text()], or
** [sqlite3_value_text16()] can be invalidated by a subsequent call to
** [sqlite3_value_bytes()], [sqlite3_value_bytes16()], [sqlite3_value_text()],
** or [sqlite3_value_text16()].
**
** These routines must be called from the same thread as
** the SQL function that supplied the [sqlite3_value*] parameters.
**
** As long as the input parameter is correct, these routines can only
** fail if an out-of-memory error occurs during a format conversion.
** Only the following subset of interfaces are subject to out-of-memory
** errors:
**
** <ul>
** <li> sqlite3_value_blob()
** <li> sqlite3_value_text()
** <li> sqlite3_value_text16()
** <li> sqlite3_value_text16le()
** <li> sqlite3_value_text16be()
** <li> sqlite3_value_bytes()
** <li> sqlite3_value_bytes16()
** </ul>
**
** If an out-of-memory error occurs, then the return value from these
** routines is the same as if the column had contained an SQL NULL value.
** Valid SQL NULL returns can be distinguished from out-of-memory errors
** by invoking the [sqlite3_errcode()] immediately after the suspect
** return value is obtained and before any
** other SQLite interface is called on the same [database connection].
 */
// llgo:link (*Value).ValueBlob C.sqlite3_value_blob
func (recv_ *Value) ValueBlob() c.Pointer {
	return nil
}

// llgo:link (*Value).ValueDouble C.sqlite3_value_double
func (recv_ *Value) ValueDouble() c.Double {
	return 0
}

// llgo:link (*Value).ValueInt C.sqlite3_value_int
func (recv_ *Value) ValueInt() c.Int {
	return 0
}

// llgo:link (*Value).ValueInt64 C.sqlite3_value_int64
func (recv_ *Value) ValueInt64() Int64 {
	return 0
}

// llgo:link (*Value).ValuePointer C.sqlite3_value_pointer
func (recv_ *Value) ValuePointer(*c.Char) c.Pointer {
	return nil
}

// llgo:link (*Value).ValueText C.sqlite3_value_text
func (recv_ *Value) ValueText() *c.Char {
	return nil
}

// llgo:link (*Value).ValueText16 C.sqlite3_value_text16
func (recv_ *Value) ValueText16() c.Pointer {
	return nil
}

// llgo:link (*Value).ValueText16le C.sqlite3_value_text16le
func (recv_ *Value) ValueText16le() c.Pointer {
	return nil
}

// llgo:link (*Value).ValueText16be C.sqlite3_value_text16be
func (recv_ *Value) ValueText16be() c.Pointer {
	return nil
}

// llgo:link (*Value).ValueBytes C.sqlite3_value_bytes
func (recv_ *Value) ValueBytes() c.Int {
	return 0
}

// llgo:link (*Value).ValueBytes16 C.sqlite3_value_bytes16
func (recv_ *Value) ValueBytes16() c.Int {
	return 0
}

// llgo:link (*Value).ValueType C.sqlite3_value_type
func (recv_ *Value) ValueType() c.Int {
	return 0
}

// llgo:link (*Value).ValueNumericType C.sqlite3_value_numeric_type
func (recv_ *Value) ValueNumericType() c.Int {
	return 0
}

// llgo:link (*Value).ValueNochange C.sqlite3_value_nochange
func (recv_ *Value) ValueNochange() c.Int {
	return 0
}

// llgo:link (*Value).ValueFrombind C.sqlite3_value_frombind
func (recv_ *Value) ValueFrombind() c.Int {
	return 0
}

/*
** CAPI3REF: Report the internal text encoding state of an sqlite3_value object
** METHOD: sqlite3_value
**
** ^(The sqlite3_value_encoding(X) interface returns one of [SQLITE_UTF8],
** [SQLITE_UTF16BE], or [SQLITE_UTF16LE] according to the current text encoding
** of the value X, assuming that X has type TEXT.)^  If sqlite3_value_type(X)
** returns something other than SQLITE_TEXT, then the return value from
** sqlite3_value_encoding(X) is meaningless.  ^Calls to
** [sqlite3_value_text(X)], [sqlite3_value_text16(X)], [sqlite3_value_text16be(X)],
** [sqlite3_value_text16le(X)], [sqlite3_value_bytes(X)], or
** [sqlite3_value_bytes16(X)] might change the encoding of the value X and
** thus change the return from subsequent calls to sqlite3_value_encoding(X).
**
** This routine is intended for used by applications that test and validate
** the SQLite implementation.  This routine is inquiring about the opaque
** internal state of an [sqlite3_value] object.  Ordinary applications should
** not need to know what the internal state of an sqlite3_value object is and
** hence should not need to use this interface.
 */
// llgo:link (*Value).ValueEncoding C.sqlite3_value_encoding
func (recv_ *Value) ValueEncoding() c.Int {
	return 0
}

/*
** CAPI3REF: Finding The Subtype Of SQL Values
** METHOD: sqlite3_value
**
** The sqlite3_value_subtype(V) function returns the subtype for
** an [application-defined SQL function] argument V.  The subtype
** information can be used to pass a limited amount of context from
** one SQL function to another.  Use the [sqlite3_result_subtype()]
** routine to set the subtype for the return value of an SQL function.
**
** Every [application-defined SQL function] that invokes this interface
** should include the [SQLITE_SUBTYPE] property in the text
** encoding argument when the function is [sqlite3_create_function|registered].
** If the [SQLITE_SUBTYPE] property is omitted, then sqlite3_value_subtype()
** might return zero instead of the upstream subtype in some corner cases.
 */
// llgo:link (*Value).ValueSubtype C.sqlite3_value_subtype
func (recv_ *Value) ValueSubtype() c.Uint {
	return 0
}

/*
** CAPI3REF: Copy And Free SQL Values
** METHOD: sqlite3_value
**
** ^The sqlite3_value_dup(V) interface makes a copy of the [sqlite3_value]
** object D and returns a pointer to that copy.  ^The [sqlite3_value] returned
** is a [protected sqlite3_value] object even if the input is not.
** ^The sqlite3_value_dup(V) interface returns NULL if V is NULL or if a
** memory allocation fails. ^If V is a [pointer value], then the result
** of sqlite3_value_dup(V) is a NULL value.
**
** ^The sqlite3_value_free(V) interface frees an [sqlite3_value] object
** previously obtained from [sqlite3_value_dup()].  ^If V is a NULL pointer
** then sqlite3_value_free(V) is a harmless no-op.
 */
// llgo:link (*Value).ValueDup C.sqlite3_value_dup
func (recv_ *Value) ValueDup() *Value {
	return nil
}

// llgo:link (*Value).ValueFree C.sqlite3_value_free
func (recv_ *Value) ValueFree() {
}

/*
** CAPI3REF: Obtain Aggregate Function Context
** METHOD: sqlite3_context
**
** Implementations of aggregate SQL functions use this
** routine to allocate memory for storing their state.
**
** ^The first time the sqlite3_aggregate_context(C,N) routine is called
** for a particular aggregate function, SQLite allocates
** N bytes of memory, zeroes out that memory, and returns a pointer
** to the new memory. ^On second and subsequent calls to
** sqlite3_aggregate_context() for the same aggregate function instance,
** the same buffer is returned.  Sqlite3_aggregate_context() is normally
** called once for each invocation of the xStep callback and then one
** last time when the xFinal callback is invoked.  ^(When no rows match
** an aggregate query, the xStep() callback of the aggregate function
** implementation is never called and xFinal() is called exactly once.
** In those cases, sqlite3_aggregate_context() might be called for the
** first time from within xFinal().)^
**
** ^The sqlite3_aggregate_context(C,N) routine returns a NULL pointer
** when first called if N is less than or equal to zero or if a memory
** allocation error occurs.
**
** ^(The amount of space allocated by sqlite3_aggregate_context(C,N) is
** determined by the N parameter on first successful call.  Changing the
** value of N in any subsequent call to sqlite3_aggregate_context() within
** the same aggregate function instance will not resize the memory
** allocation.)^  Within the xFinal callback, it is customary to set
** N=0 in calls to sqlite3_aggregate_context(C,N) so that no
** pointless memory allocations occur.
**
** ^SQLite automatically frees the memory allocated by
** sqlite3_aggregate_context() when the aggregate query concludes.
**
** The first parameter must be a copy of the
** [sqlite3_context | SQL function context] that is the first parameter
** to the xStep or xFinal callback routine that implements the aggregate
** function.
**
** This routine must be called from the same thread in which
** the aggregate SQL function is running.
 */
// llgo:link (*Context).AggregateContext C.sqlite3_aggregate_context
func (recv_ *Context) AggregateContext(nBytes c.Int) c.Pointer {
	return nil
}

/*
** CAPI3REF: User Data For Functions
** METHOD: sqlite3_context
**
** ^The sqlite3_user_data() interface returns a copy of
** the pointer that was the pUserData parameter (the 5th parameter)
** of the [sqlite3_create_function()]
** and [sqlite3_create_function16()] routines that originally
** registered the application defined function.
**
** This routine must be called from the same thread in which
** the application-defined function is running.
 */
// llgo:link (*Context).UserData C.sqlite3_user_data
func (recv_ *Context) UserData() c.Pointer {
	return nil
}

/*
** CAPI3REF: Database Connection For Functions
** METHOD: sqlite3_context
**
** ^The sqlite3_context_db_handle() interface returns a copy of
** the pointer to the [database connection] (the 1st parameter)
** of the [sqlite3_create_function()]
** and [sqlite3_create_function16()] routines that originally
** registered the application defined function.
 */
// llgo:link (*Context).ContextDbHandle C.sqlite3_context_db_handle
func (recv_ *Context) ContextDbHandle() *Sqlite3 {
	return nil
}

/*
** CAPI3REF: Function Auxiliary Data
** METHOD: sqlite3_context
**
** These functions may be used by (non-aggregate) SQL functions to
** associate auxiliary data with argument values. If the same argument
** value is passed to multiple invocations of the same SQL function during
** query execution, under some circumstances the associated auxiliary data
** might be preserved.  An example of where this might be useful is in a
** regular-expression matching function. The compiled version of the regular
** expression can be stored as auxiliary data associated with the pattern string.
** Then as long as the pattern string remains the same,
** the compiled regular expression can be reused on multiple
** invocations of the same function.
**
** ^The sqlite3_get_auxdata(C,N) interface returns a pointer to the auxiliary data
** associated by the sqlite3_set_auxdata(C,N,P,X) function with the Nth argument
** value to the application-defined function.  ^N is zero for the left-most
** function argument.  ^If there is no auxiliary data
** associated with the function argument, the sqlite3_get_auxdata(C,N) interface
** returns a NULL pointer.
**
** ^The sqlite3_set_auxdata(C,N,P,X) interface saves P as auxiliary data for the
** N-th argument of the application-defined function.  ^Subsequent
** calls to sqlite3_get_auxdata(C,N) return P from the most recent
** sqlite3_set_auxdata(C,N,P,X) call if the auxiliary data is still valid or
** NULL if the auxiliary data has been discarded.
** ^After each call to sqlite3_set_auxdata(C,N,P,X) where X is not NULL,
** SQLite will invoke the destructor function X with parameter P exactly
** once, when the auxiliary data is discarded.
** SQLite is free to discard the auxiliary data at any time, including: <ul>
** <li> ^(when the corresponding function parameter changes)^, or
** <li> ^(when [sqlite3_reset()] or [sqlite3_finalize()] is called for the
**      SQL statement)^, or
** <li> ^(when sqlite3_set_auxdata() is invoked again on the same
**       parameter)^, or
** <li> ^(during the original sqlite3_set_auxdata() call when a memory
**      allocation error occurs.)^
** <li> ^(during the original sqlite3_set_auxdata() call if the function
**      is evaluated during query planning instead of during query execution,
**      as sometimes happens with [SQLITE_ENABLE_STAT4].)^ </ul>
**
** Note the last two bullets in particular.  The destructor X in
** sqlite3_set_auxdata(C,N,P,X) might be called immediately, before the
** sqlite3_set_auxdata() interface even returns.  Hence sqlite3_set_auxdata()
** should be called near the end of the function implementation and the
** function implementation should not make any use of P after
** sqlite3_set_auxdata() has been called.  Furthermore, a call to
** sqlite3_get_auxdata() that occurs immediately after a corresponding call
** to sqlite3_set_auxdata() might still return NULL if an out-of-memory
** condition occurred during the sqlite3_set_auxdata() call or if the
** function is being evaluated during query planning rather than during
** query execution.
**
** ^(In practice, auxiliary data is preserved between function calls for
** function parameters that are compile-time constants, including literal
** values and [parameters] and expressions composed from the same.)^
**
** The value of the N parameter to these interfaces should be non-negative.
** Future enhancements may make use of negative N values to define new
** kinds of function caching behavior.
**
** These routines must be called from the same thread in which
** the SQL function is running.
**
** See also: [sqlite3_get_clientdata()] and [sqlite3_set_clientdata()].
 */
// llgo:link (*Context).GetAuxdata C.sqlite3_get_auxdata
func (recv_ *Context) GetAuxdata(N c.Int) c.Pointer {
	return nil
}

// llgo:link (*Context).SetAuxdata C.sqlite3_set_auxdata
func (recv_ *Context) SetAuxdata(N c.Int, __llgo_arg_1 c.Pointer, __llgo_arg_2 func(c.Pointer)) {
}

/*
** CAPI3REF: Database Connection Client Data
** METHOD: sqlite3
**
** These functions are used to associate one or more named pointers
** with a [database connection].
** A call to sqlite3_set_clientdata(D,N,P,X) causes the pointer P
** to be attached to [database connection] D using name N.  Subsequent
** calls to sqlite3_get_clientdata(D,N) will return a copy of pointer P
** or a NULL pointer if there were no prior calls to
** sqlite3_set_clientdata() with the same values of D and N.
** Names are compared using strcmp() and are thus case sensitive.
**
** If P and X are both non-NULL, then the destructor X is invoked with
** argument P on the first of the following occurrences:
** <ul>
** <li> An out-of-memory error occurs during the call to
**      sqlite3_set_clientdata() which attempts to register pointer P.
** <li> A subsequent call to sqlite3_set_clientdata(D,N,P,X) is made
**      with the same D and N parameters.
** <li> The database connection closes.  SQLite does not make any guarantees
**      about the order in which destructors are called, only that all
**      destructors will be called exactly once at some point during the
**      database connection closing process.
** </ul>
**
** SQLite does not do anything with client data other than invoke
** destructors on the client data at the appropriate time.  The intended
** use for client data is to provide a mechanism for wrapper libraries
** to store additional information about an SQLite database connection.
**
** There is no limit (other than available memory) on the number of different
** client data pointers (with different names) that can be attached to a
** single database connection.  However, the implementation is optimized
** for the case of having only one or two different client data names.
** Applications and wrapper libraries are discouraged from using more than
** one client data name each.
**
** There is no way to enumerate the client data pointers
** associated with a database connection.  The N parameter can be thought
** of as a secret key such that only code that knows the secret key is able
** to access the associated data.
**
** Security Warning:  These interfaces should not be exposed in scripting
** languages or in other circumstances where it might be possible for an
** an attacker to invoke them.  Any agent that can invoke these interfaces
** can probably also take control of the process.
**
** Database connection client data is only available for SQLite
** version 3.44.0 ([dateof:3.44.0]) and later.
**
** See also: [sqlite3_set_auxdata()] and [sqlite3_get_auxdata()].
 */
// llgo:link (*Sqlite3).GetClientdata C.sqlite3_get_clientdata
func (recv_ *Sqlite3) GetClientdata(*c.Char) c.Pointer {
	return nil
}

// llgo:link (*Sqlite3).SetClientdata C.sqlite3_set_clientdata
func (recv_ *Sqlite3) SetClientdata(*c.Char, c.Pointer, func(c.Pointer)) c.Int {
	return 0
}

// llgo:type C
type DestructorType func(c.Pointer)

/*
** CAPI3REF: Setting The Result Of An SQL Function
** METHOD: sqlite3_context
**
** These routines are used by the xFunc or xFinal callbacks that
** implement SQL functions and aggregates.  See
** [sqlite3_create_function()] and [sqlite3_create_function16()]
** for additional information.
**
** These functions work very much like the [parameter binding] family of
** functions used to bind values to host parameters in prepared statements.
** Refer to the [SQL parameter] documentation for additional information.
**
** ^The sqlite3_result_blob() interface sets the result from
** an application-defined function to be the BLOB whose content is pointed
** to by the second parameter and which is N bytes long where N is the
** third parameter.
**
** ^The sqlite3_result_zeroblob(C,N) and sqlite3_result_zeroblob64(C,N)
** interfaces set the result of the application-defined function to be
** a BLOB containing all zero bytes and N bytes in size.
**
** ^The sqlite3_result_double() interface sets the result from
** an application-defined function to be a floating point value specified
** by its 2nd argument.
**
** ^The sqlite3_result_error() and sqlite3_result_error16() functions
** cause the implemented SQL function to throw an exception.
** ^SQLite uses the string pointed to by the
** 2nd parameter of sqlite3_result_error() or sqlite3_result_error16()
** as the text of an error message.  ^SQLite interprets the error
** message string from sqlite3_result_error() as UTF-8. ^SQLite
** interprets the string from sqlite3_result_error16() as UTF-16 using
** the same [byte-order determination rules] as [sqlite3_bind_text16()].
** ^If the third parameter to sqlite3_result_error()
** or sqlite3_result_error16() is negative then SQLite takes as the error
** message all text up through the first zero character.
** ^If the third parameter to sqlite3_result_error() or
** sqlite3_result_error16() is non-negative then SQLite takes that many
** bytes (not characters) from the 2nd parameter as the error message.
** ^The sqlite3_result_error() and sqlite3_result_error16()
** routines make a private copy of the error message text before
** they return.  Hence, the calling function can deallocate or
** modify the text after they return without harm.
** ^The sqlite3_result_error_code() function changes the error code
** returned by SQLite as a result of an error in a function.  ^By default,
** the error code is SQLITE_ERROR.  ^A subsequent call to sqlite3_result_error()
** or sqlite3_result_error16() resets the error code to SQLITE_ERROR.
**
** ^The sqlite3_result_error_toobig() interface causes SQLite to throw an
** error indicating that a string or BLOB is too long to represent.
**
** ^The sqlite3_result_error_nomem() interface causes SQLite to throw an
** error indicating that a memory allocation failed.
**
** ^The sqlite3_result_int() interface sets the return value
** of the application-defined function to be the 32-bit signed integer
** value given in the 2nd argument.
** ^The sqlite3_result_int64() interface sets the return value
** of the application-defined function to be the 64-bit signed integer
** value given in the 2nd argument.
**
** ^The sqlite3_result_null() interface sets the return value
** of the application-defined function to be NULL.
**
** ^The sqlite3_result_text(), sqlite3_result_text16(),
** sqlite3_result_text16le(), and sqlite3_result_text16be() interfaces
** set the return value of the application-defined function to be
** a text string which is represented as UTF-8, UTF-16 native byte order,
** UTF-16 little endian, or UTF-16 big endian, respectively.
** ^The sqlite3_result_text64() interface sets the return value of an
** application-defined function to be a text string in an encoding
** specified by the fifth (and last) parameter, which must be one
** of [SQLITE_UTF8], [SQLITE_UTF16], [SQLITE_UTF16BE], or [SQLITE_UTF16LE].
** ^SQLite takes the text result from the application from
** the 2nd parameter of the sqlite3_result_text* interfaces.
** ^If the 3rd parameter to any of the sqlite3_result_text* interfaces
** other than sqlite3_result_text64() is negative, then SQLite computes
** the string length itself by searching the 2nd parameter for the first
** zero character.
** ^If the 3rd parameter to the sqlite3_result_text* interfaces
** is non-negative, then as many bytes (not characters) of the text
** pointed to by the 2nd parameter are taken as the application-defined
** function result.  If the 3rd parameter is non-negative, then it
** must be the byte offset into the string where the NUL terminator would
** appear if the string where NUL terminated.  If any NUL characters occur
** in the string at a byte offset that is less than the value of the 3rd
** parameter, then the resulting string will contain embedded NULs and the
** result of expressions operating on strings with embedded NULs is undefined.
** ^If the 4th parameter to the sqlite3_result_text* interfaces
** or sqlite3_result_blob is a non-NULL pointer, then SQLite calls that
** function as the destructor on the text or BLOB result when it has
** finished using that result.
** ^If the 4th parameter to the sqlite3_result_text* interfaces or to
** sqlite3_result_blob is the special constant SQLITE_STATIC, then SQLite
** assumes that the text or BLOB result is in constant space and does not
** copy the content of the parameter nor call a destructor on the content
** when it has finished using that result.
** ^If the 4th parameter to the sqlite3_result_text* interfaces
** or sqlite3_result_blob is the special constant SQLITE_TRANSIENT
** then SQLite makes a copy of the result into space obtained
** from [sqlite3_malloc()] before it returns.
**
** ^For the sqlite3_result_text16(), sqlite3_result_text16le(), and
** sqlite3_result_text16be() routines, and for sqlite3_result_text64()
** when the encoding is not UTF8, if the input UTF16 begins with a
** byte-order mark (BOM, U+FEFF) then the BOM is removed from the
** string and the rest of the string is interpreted according to the
** byte-order specified by the BOM.  ^The byte-order specified by
** the BOM at the beginning of the text overrides the byte-order
** specified by the interface procedure.  ^So, for example, if
** sqlite3_result_text16le() is invoked with text that begins
** with bytes 0xfe, 0xff (a big-endian byte-order mark) then the
** first two bytes of input are skipped and the remaining input
** is interpreted as UTF16BE text.
**
** ^For UTF16 input text to the sqlite3_result_text16(),
** sqlite3_result_text16be(), sqlite3_result_text16le(), and
** sqlite3_result_text64() routines, if the text contains invalid
** UTF16 characters, the invalid characters might be converted
** into the unicode replacement character, U+FFFD.
**
** ^The sqlite3_result_value() interface sets the result of
** the application-defined function to be a copy of the
** [unprotected sqlite3_value] object specified by the 2nd parameter.  ^The
** sqlite3_result_value() interface makes a copy of the [sqlite3_value]
** so that the [sqlite3_value] specified in the parameter may change or
** be deallocated after sqlite3_result_value() returns without harm.
** ^A [protected sqlite3_value] object may always be used where an
** [unprotected sqlite3_value] object is required, so either
** kind of [sqlite3_value] object can be used with this interface.
**
** ^The sqlite3_result_pointer(C,P,T,D) interface sets the result to an
** SQL NULL value, just like [sqlite3_result_null(C)], except that it
** also associates the host-language pointer P or type T with that
** NULL value such that the pointer can be retrieved within an
** [application-defined SQL function] using [sqlite3_value_pointer()].
** ^If the D parameter is not NULL, then it is a pointer to a destructor
** for the P parameter.  ^SQLite invokes D with P as its only argument
** when SQLite is finished with P.  The T parameter should be a static
** string and preferably a string literal. The sqlite3_result_pointer()
** routine is part of the [pointer passing interface] added for SQLite 3.20.0.
**
** If these routines are called from within the different thread
** than the one containing the application-defined function that received
** the [sqlite3_context] pointer, the results are undefined.
 */
// llgo:link (*Context).ResultBlob C.sqlite3_result_blob
func (recv_ *Context) ResultBlob(c.Pointer, c.Int, func(c.Pointer)) {
}

// llgo:link (*Context).ResultBlob64 C.sqlite3_result_blob64
func (recv_ *Context) ResultBlob64(c.Pointer, Uint64, func(c.Pointer)) {
}

// llgo:link (*Context).ResultDouble C.sqlite3_result_double
func (recv_ *Context) ResultDouble(c.Double) {
}

// llgo:link (*Context).ResultError C.sqlite3_result_error
func (recv_ *Context) ResultError(*c.Char, c.Int) {
}

// llgo:link (*Context).ResultError16 C.sqlite3_result_error16
func (recv_ *Context) ResultError16(c.Pointer, c.Int) {
}

// llgo:link (*Context).ResultErrorToobig C.sqlite3_result_error_toobig
func (recv_ *Context) ResultErrorToobig() {
}

// llgo:link (*Context).ResultErrorNomem C.sqlite3_result_error_nomem
func (recv_ *Context) ResultErrorNomem() {
}

// llgo:link (*Context).ResultErrorCode C.sqlite3_result_error_code
func (recv_ *Context) ResultErrorCode(c.Int) {
}

// llgo:link (*Context).ResultInt C.sqlite3_result_int
func (recv_ *Context) ResultInt(c.Int) {
}

// llgo:link (*Context).ResultInt64 C.sqlite3_result_int64
func (recv_ *Context) ResultInt64(Int64) {
}

// llgo:link (*Context).ResultNull C.sqlite3_result_null
func (recv_ *Context) ResultNull() {
}

// llgo:link (*Context).ResultText C.sqlite3_result_text
func (recv_ *Context) ResultText(*c.Char, c.Int, func(c.Pointer)) {
}

// llgo:link (*Context).ResultText64 C.sqlite3_result_text64
func (recv_ *Context) ResultText64(__llgo_arg_0 *c.Char, __llgo_arg_1 Uint64, __llgo_arg_2 func(c.Pointer), encoding c.Char) {
}

// llgo:link (*Context).ResultText16 C.sqlite3_result_text16
func (recv_ *Context) ResultText16(c.Pointer, c.Int, func(c.Pointer)) {
}

// llgo:link (*Context).ResultText16le C.sqlite3_result_text16le
func (recv_ *Context) ResultText16le(c.Pointer, c.Int, func(c.Pointer)) {
}

// llgo:link (*Context).ResultText16be C.sqlite3_result_text16be
func (recv_ *Context) ResultText16be(c.Pointer, c.Int, func(c.Pointer)) {
}

// llgo:link (*Context).ResultValue C.sqlite3_result_value
func (recv_ *Context) ResultValue(*Value) {
}

// llgo:link (*Context).ResultPointer C.sqlite3_result_pointer
func (recv_ *Context) ResultPointer(c.Pointer, *c.Char, func(c.Pointer)) {
}

// llgo:link (*Context).ResultZeroblob C.sqlite3_result_zeroblob
func (recv_ *Context) ResultZeroblob(n c.Int) {
}

// llgo:link (*Context).ResultZeroblob64 C.sqlite3_result_zeroblob64
func (recv_ *Context) ResultZeroblob64(n Uint64) c.Int {
	return 0
}

/*
** CAPI3REF: Setting The Subtype Of An SQL Function
** METHOD: sqlite3_context
**
** The sqlite3_result_subtype(C,T) function causes the subtype of
** the result from the [application-defined SQL function] with
** [sqlite3_context] C to be the value T.  Only the lower 8 bits
** of the subtype T are preserved in current versions of SQLite;
** higher order bits are discarded.
** The number of subtype bytes preserved by SQLite might increase
** in future releases of SQLite.
**
** Every [application-defined SQL function] that invokes this interface
** should include the [SQLITE_RESULT_SUBTYPE] property in its
** text encoding argument when the SQL function is
** [sqlite3_create_function|registered].  If the [SQLITE_RESULT_SUBTYPE]
** property is omitted from the function that invokes sqlite3_result_subtype(),
** then in some cases the sqlite3_result_subtype() might fail to set
** the result subtype.
**
** If SQLite is compiled with -DSQLITE_STRICT_SUBTYPE=1, then any
** SQL function that invokes the sqlite3_result_subtype() interface
** and that does not have the SQLITE_RESULT_SUBTYPE property will raise
** an error.  Future versions of SQLite might enable -DSQLITE_STRICT_SUBTYPE=1
** by default.
 */
// llgo:link (*Context).ResultSubtype C.sqlite3_result_subtype
func (recv_ *Context) ResultSubtype(c.Uint) {
}

/*
** CAPI3REF: Define New Collating Sequences
** METHOD: sqlite3
**
** ^These functions add, remove, or modify a [collation] associated
** with the [database connection] specified as the first argument.
**
** ^The name of the collation is a UTF-8 string
** for sqlite3_create_collation() and sqlite3_create_collation_v2()
** and a UTF-16 string in native byte order for sqlite3_create_collation16().
** ^Collation names that compare equal according to [sqlite3_strnicmp()] are
** considered to be the same name.
**
** ^(The third argument (eTextRep) must be one of the constants:
** <ul>
** <li> [SQLITE_UTF8],
** <li> [SQLITE_UTF16LE],
** <li> [SQLITE_UTF16BE],
** <li> [SQLITE_UTF16], or
** <li> [SQLITE_UTF16_ALIGNED].
** </ul>)^
** ^The eTextRep argument determines the encoding of strings passed
** to the collating function callback, xCompare.
** ^The [SQLITE_UTF16] and [SQLITE_UTF16_ALIGNED] values for eTextRep
** force strings to be UTF16 with native byte order.
** ^The [SQLITE_UTF16_ALIGNED] value for eTextRep forces strings to begin
** on an even byte address.
**
** ^The fourth argument, pArg, is an application data pointer that is passed
** through as the first argument to the collating function callback.
**
** ^The fifth argument, xCompare, is a pointer to the collating function.
** ^Multiple collating functions can be registered using the same name but
** with different eTextRep parameters and SQLite will use whichever
** function requires the least amount of data transformation.
** ^If the xCompare argument is NULL then the collating function is
** deleted.  ^When all collating functions having the same name are deleted,
** that collation is no longer usable.
**
** ^The collating function callback is invoked with a copy of the pArg
** application data pointer and with two strings in the encoding specified
** by the eTextRep argument.  The two integer parameters to the collating
** function callback are the length of the two strings, in bytes. The collating
** function must return an integer that is negative, zero, or positive
** if the first string is less than, equal to, or greater than the second,
** respectively.  A collating function must always return the same answer
** given the same inputs.  If two or more collating functions are registered
** to the same collation name (using different eTextRep values) then all
** must give an equivalent answer when invoked with equivalent strings.
** The collating function must obey the following properties for all
** strings A, B, and C:
**
** <ol>
** <li> If A==B then B==A.
** <li> If A==B and B==C then A==C.
** <li> If A&lt;B THEN B&gt;A.
** <li> If A&lt;B and B&lt;C then A&lt;C.
** </ol>
**
** If a collating function fails any of the above constraints and that
** collating function is registered and used, then the behavior of SQLite
** is undefined.
**
** ^The sqlite3_create_collation_v2() works like sqlite3_create_collation()
** with the addition that the xDestroy callback is invoked on pArg when
** the collating function is deleted.
** ^Collating functions are deleted when they are overridden by later
** calls to the collation creation functions or when the
** [database connection] is closed using [sqlite3_close()].
**
** ^The xDestroy callback is <u>not</u> called if the
** sqlite3_create_collation_v2() function fails.  Applications that invoke
** sqlite3_create_collation_v2() with a non-NULL xDestroy argument should
** check the return code and dispose of the application data pointer
** themselves rather than expecting SQLite to deal with it for them.
** This is different from every other SQLite interface.  The inconsistency
** is unfortunate but cannot be changed without breaking backwards
** compatibility.
**
** See also:  [sqlite3_collation_needed()] and [sqlite3_collation_needed16()].
 */
// llgo:link (*Sqlite3).CreateCollation C.sqlite3_create_collation
func (recv_ *Sqlite3) CreateCollation(zName *c.Char, eTextRep c.Int, pArg c.Pointer, xCompare func(c.Pointer, c.Int, c.Pointer, c.Int, c.Pointer) c.Int) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateCollationV2 C.sqlite3_create_collation_v2
func (recv_ *Sqlite3) CreateCollationV2(zName *c.Char, eTextRep c.Int, pArg c.Pointer, xCompare func(c.Pointer, c.Int, c.Pointer, c.Int, c.Pointer) c.Int, xDestroy func(c.Pointer)) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateCollation16 C.sqlite3_create_collation16
func (recv_ *Sqlite3) CreateCollation16(zName c.Pointer, eTextRep c.Int, pArg c.Pointer, xCompare func(c.Pointer, c.Int, c.Pointer, c.Int, c.Pointer) c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Collation Needed Callbacks
** METHOD: sqlite3
**
** ^To avoid having to register all collation sequences before a database
** can be used, a single callback function may be registered with the
** [database connection] to be invoked whenever an undefined collation
** sequence is required.
**
** ^If the function is registered using the sqlite3_collation_needed() API,
** then it is passed the names of undefined collation sequences as strings
** encoded in UTF-8. ^If sqlite3_collation_needed16() is used,
** the names are passed as UTF-16 in machine native byte order.
** ^A call to either function replaces the existing collation-needed callback.
**
** ^(When the callback is invoked, the first argument passed is a copy
** of the second argument to sqlite3_collation_needed() or
** sqlite3_collation_needed16().  The second argument is the database
** connection.  The third argument is one of [SQLITE_UTF8], [SQLITE_UTF16BE],
** or [SQLITE_UTF16LE], indicating the most desirable form of the collation
** sequence function required.  The fourth parameter is the name of the
** required collation sequence.)^
**
** The callback function should register the desired collation using
** [sqlite3_create_collation()], [sqlite3_create_collation16()], or
** [sqlite3_create_collation_v2()].
 */
// llgo:link (*Sqlite3).CollationNeeded C.sqlite3_collation_needed
func (recv_ *Sqlite3) CollationNeeded(c.Pointer, func(c.Pointer, *Sqlite3, c.Int, *c.Char)) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CollationNeeded16 C.sqlite3_collation_needed16
func (recv_ *Sqlite3) CollationNeeded16(c.Pointer, func(c.Pointer, *Sqlite3, c.Int, c.Pointer)) c.Int {
	return 0
}

/*
** CAPI3REF: Suspend Execution For A Short Time
**
** The sqlite3_sleep() function causes the current thread to suspend execution
** for at least a number of milliseconds specified in its parameter.
**
** If the operating system does not support sleep requests with
** millisecond time resolution, then the time will be rounded up to
** the nearest second. The number of milliseconds of sleep actually
** requested from the operating system is returned.
**
** ^SQLite implements this interface by calling the xSleep()
** method of the default [sqlite3_vfs] object.  If the xSleep() method
** of the default VFS is not implemented correctly, or not implemented at
** all, then the behavior of sqlite3_sleep() may deviate from the description
** in the previous paragraphs.
**
** If a negative argument is passed to sqlite3_sleep() the results vary by
** VFS and operating system.  Some system treat a negative argument as an
** instruction to sleep forever.  Others understand it to mean do not sleep
** at all. ^In SQLite version 3.42.0 and later, a negative
** argument passed into sqlite3_sleep() is changed to zero before it is relayed
** down into the xSleep method of the VFS.
 */
//go:linkname Sleep C.sqlite3_sleep
func Sleep(c.Int) c.Int

/*
** CAPI3REF: Test For Auto-Commit Mode
** KEYWORDS: {autocommit mode}
** METHOD: sqlite3
**
** ^The sqlite3_get_autocommit() interface returns non-zero or
** zero if the given database connection is or is not in autocommit mode,
** respectively.  ^Autocommit mode is on by default.
** ^Autocommit mode is disabled by a [BEGIN] statement.
** ^Autocommit mode is re-enabled by a [COMMIT] or [ROLLBACK].
**
** If certain kinds of errors occur on a statement within a multi-statement
** transaction (errors including [SQLITE_FULL], [SQLITE_IOERR],
** [SQLITE_NOMEM], [SQLITE_BUSY], and [SQLITE_INTERRUPT]) then the
** transaction might be rolled back automatically.  The only way to
** find out whether SQLite automatically rolled back the transaction after
** an error is to use this function.
**
** If another thread changes the autocommit status of the database
** connection while this routine is running, then the return value
** is undefined.
 */
// llgo:link (*Sqlite3).GetAutocommit C.sqlite3_get_autocommit
func (recv_ *Sqlite3) GetAutocommit() c.Int {
	return 0
}

/*
** CAPI3REF: Find The Database Handle Of A Prepared Statement
** METHOD: sqlite3_stmt
**
** ^The sqlite3_db_handle interface returns the [database connection] handle
** to which a [prepared statement] belongs.  ^The [database connection]
** returned by sqlite3_db_handle is the same [database connection]
** that was the first argument
** to the [sqlite3_prepare_v2()] call (or its variants) that was used to
** create the statement in the first place.
 */
// llgo:link (*Stmt).DbHandle C.sqlite3_db_handle
func (recv_ *Stmt) DbHandle() *Sqlite3 {
	return nil
}

/*
** CAPI3REF: Return The Schema Name For A Database Connection
** METHOD: sqlite3
**
** ^The sqlite3_db_name(D,N) interface returns a pointer to the schema name
** for the N-th database on database connection D, or a NULL pointer of N is
** out of range.  An N value of 0 means the main database file.  An N of 1 is
** the "temp" schema.  Larger values of N correspond to various ATTACH-ed
** databases.
**
** Space to hold the string that is returned by sqlite3_db_name() is managed
** by SQLite itself.  The string might be deallocated by any operation that
** changes the schema, including [ATTACH] or [DETACH] or calls to
** [sqlite3_serialize()] or [sqlite3_deserialize()], even operations that
** occur on a different thread.  Applications that need to
** remember the string long-term should make their own copy.  Applications that
** are accessing the same database connection simultaneously on multiple
** threads should mutex-protect calls to this API and should make their own
** private copy of the result prior to releasing the mutex.
 */
// llgo:link (*Sqlite3).DbName C.sqlite3_db_name
func (recv_ *Sqlite3) DbName(N c.Int) *c.Char {
	return nil
}

/*
** CAPI3REF: Return The Filename For A Database Connection
** METHOD: sqlite3
**
** ^The sqlite3_db_filename(D,N) interface returns a pointer to the filename
** associated with database N of connection D.
** ^If there is no attached database N on the database
** connection D, or if database N is a temporary or in-memory database, then
** this function will return either a NULL pointer or an empty string.
**
** ^The string value returned by this routine is owned and managed by
** the database connection.  ^The value will be valid until the database N
** is [DETACH]-ed or until the database connection closes.
**
** ^The filename returned by this function is the output of the
** xFullPathname method of the [VFS].  ^In other words, the filename
** will be an absolute pathname, even if the filename used
** to open the database originally was a URI or relative pathname.
**
** If the filename pointer returned by this routine is not NULL, then it
** can be used as the filename input parameter to these routines:
** <ul>
** <li> [sqlite3_uri_parameter()]
** <li> [sqlite3_uri_boolean()]
** <li> [sqlite3_uri_int64()]
** <li> [sqlite3_filename_database()]
** <li> [sqlite3_filename_journal()]
** <li> [sqlite3_filename_wal()]
** </ul>
 */
// llgo:link (*Sqlite3).DbFilename C.sqlite3_db_filename
func (recv_ *Sqlite3) DbFilename(zDbName *c.Char) Filename {
	return nil
}

/*
** CAPI3REF: Determine if a database is read-only
** METHOD: sqlite3
**
** ^The sqlite3_db_readonly(D,N) interface returns 1 if the database N
** of connection D is read-only, 0 if it is read/write, or -1 if N is not
** the name of a database on connection D.
 */
// llgo:link (*Sqlite3).DbReadonly C.sqlite3_db_readonly
func (recv_ *Sqlite3) DbReadonly(zDbName *c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Determine the transaction state of a database
** METHOD: sqlite3
**
** ^The sqlite3_txn_state(D,S) interface returns the current
** [transaction state] of schema S in database connection D.  ^If S is NULL,
** then the highest transaction state of any schema on database connection D
** is returned.  Transaction states are (in order of lowest to highest):
** <ol>
** <li value="0"> SQLITE_TXN_NONE
** <li value="1"> SQLITE_TXN_READ
** <li value="2"> SQLITE_TXN_WRITE
** </ol>
** ^If the S argument to sqlite3_txn_state(D,S) is not the name of
** a valid schema, then -1 is returned.
 */
// llgo:link (*Sqlite3).TxnState C.sqlite3_txn_state
func (recv_ *Sqlite3) TxnState(zSchema *c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Find the next prepared statement
** METHOD: sqlite3
**
** ^This interface returns a pointer to the next [prepared statement] after
** pStmt associated with the [database connection] pDb.  ^If pStmt is NULL
** then this interface returns a pointer to the first prepared statement
** associated with the database connection pDb.  ^If no prepared statement
** satisfies the conditions of this routine, it returns NULL.
**
** The [database connection] pointer D in a call to
** [sqlite3_next_stmt(D,S)] must refer to an open database
** connection and in particular must not be a NULL pointer.
 */
// llgo:link (*Sqlite3).NextStmt C.sqlite3_next_stmt
func (recv_ *Sqlite3) NextStmt(pStmt *Stmt) *Stmt {
	return nil
}

/*
** CAPI3REF: Commit And Rollback Notification Callbacks
** METHOD: sqlite3
**
** ^The sqlite3_commit_hook() interface registers a callback
** function to be invoked whenever a transaction is [COMMIT | committed].
** ^Any callback set by a previous call to sqlite3_commit_hook()
** for the same database connection is overridden.
** ^The sqlite3_rollback_hook() interface registers a callback
** function to be invoked whenever a transaction is [ROLLBACK | rolled back].
** ^Any callback set by a previous call to sqlite3_rollback_hook()
** for the same database connection is overridden.
** ^The pArg argument is passed through to the callback.
** ^If the callback on a commit hook function returns non-zero,
** then the commit is converted into a rollback.
**
** ^The sqlite3_commit_hook(D,C,P) and sqlite3_rollback_hook(D,C,P) functions
** return the P argument from the previous call of the same function
** on the same [database connection] D, or NULL for
** the first call for each function on D.
**
** The commit and rollback hook callbacks are not reentrant.
** The callback implementation must not do anything that will modify
** the database connection that invoked the callback.  Any actions
** to modify the database connection must be deferred until after the
** completion of the [sqlite3_step()] call that triggered the commit
** or rollback hook in the first place.
** Note that running any other SQL statements, including SELECT statements,
** or merely calling [sqlite3_prepare_v2()] and [sqlite3_step()] will modify
** the database connections for the meaning of "modify" in this paragraph.
**
** ^Registering a NULL function disables the callback.
**
** ^When the commit hook callback routine returns zero, the [COMMIT]
** operation is allowed to continue normally.  ^If the commit hook
** returns non-zero, then the [COMMIT] is converted into a [ROLLBACK].
** ^The rollback hook is invoked on a rollback that results from a commit
** hook returning non-zero, just as it would be with any other rollback.
**
** ^For the purposes of this API, a transaction is said to have been
** rolled back if an explicit "ROLLBACK" statement is executed, or
** an error or constraint causes an implicit rollback to occur.
** ^The rollback callback is not invoked if a transaction is
** automatically rolled back because the database connection is closed.
**
** See also the [sqlite3_update_hook()] interface.
 */
// llgo:link (*Sqlite3).CommitHook C.sqlite3_commit_hook
func (recv_ *Sqlite3) CommitHook(func(c.Pointer) c.Int, c.Pointer) c.Pointer {
	return nil
}

// llgo:link (*Sqlite3).RollbackHook C.sqlite3_rollback_hook
func (recv_ *Sqlite3) RollbackHook(func(c.Pointer), c.Pointer) c.Pointer {
	return nil
}

/*
** CAPI3REF: Autovacuum Compaction Amount Callback
** METHOD: sqlite3
**
** ^The sqlite3_autovacuum_pages(D,C,P,X) interface registers a callback
** function C that is invoked prior to each autovacuum of the database
** file.  ^The callback is passed a copy of the generic data pointer (P),
** the schema-name of the attached database that is being autovacuumed,
** the size of the database file in pages, the number of free pages,
** and the number of bytes per page, respectively.  The callback should
** return the number of free pages that should be removed by the
** autovacuum.  ^If the callback returns zero, then no autovacuum happens.
** ^If the value returned is greater than or equal to the number of
** free pages, then a complete autovacuum happens.
**
** <p>^If there are multiple ATTACH-ed database files that are being
** modified as part of a transaction commit, then the autovacuum pages
** callback is invoked separately for each file.
**
** <p><b>The callback is not reentrant.</b> The callback function should
** not attempt to invoke any other SQLite interface.  If it does, bad
** things may happen, including segmentation faults and corrupt database
** files.  The callback function should be a simple function that
** does some arithmetic on its input parameters and returns a result.
**
** ^The X parameter to sqlite3_autovacuum_pages(D,C,P,X) is an optional
** destructor for the P parameter.  ^If X is not NULL, then X(P) is
** invoked whenever the database connection closes or when the callback
** is overwritten by another invocation of sqlite3_autovacuum_pages().
**
** <p>^There is only one autovacuum pages callback per database connection.
** ^Each call to the sqlite3_autovacuum_pages() interface overrides all
** previous invocations for that database connection.  ^If the callback
** argument (C) to sqlite3_autovacuum_pages(D,C,P,X) is a NULL pointer,
** then the autovacuum steps callback is canceled.  The return value
** from sqlite3_autovacuum_pages() is normally SQLITE_OK, but might
** be some other error code if something goes wrong.  The current
** implementation will only return SQLITE_OK or SQLITE_MISUSE, but other
** return codes might be added in future releases.
**
** <p>If no autovacuum pages callback is specified (the usual case) or
** a NULL pointer is provided for the callback,
** then the default behavior is to vacuum all free pages.  So, in other
** words, the default behavior is the same as if the callback function
** were something like this:
**
** <blockquote><pre>
** &nbsp;   unsigned int demonstration_autovac_pages_callback(
** &nbsp;     void *pClientData,
** &nbsp;     const char *zSchema,
** &nbsp;     unsigned int nDbPage,
** &nbsp;     unsigned int nFreePage,
** &nbsp;     unsigned int nBytePerPage
** &nbsp;   ){
** &nbsp;     return nFreePage;
** &nbsp;   }
** </pre></blockquote>
 */
// llgo:link (*Sqlite3).AutovacuumPages C.sqlite3_autovacuum_pages
func (recv_ *Sqlite3) AutovacuumPages(func(c.Pointer, *c.Char, c.Uint, c.Uint, c.Uint) c.Uint, c.Pointer, func(c.Pointer)) c.Int {
	return 0
}

/*
** CAPI3REF: Data Change Notification Callbacks
** METHOD: sqlite3
**
** ^The sqlite3_update_hook() interface registers a callback function
** with the [database connection] identified by the first argument
** to be invoked whenever a row is updated, inserted or deleted in
** a [rowid table].
** ^Any callback set by a previous call to this function
** for the same database connection is overridden.
**
** ^The second argument is a pointer to the function to invoke when a
** row is updated, inserted or deleted in a rowid table.
** ^The first argument to the callback is a copy of the third argument
** to sqlite3_update_hook().
** ^The second callback argument is one of [SQLITE_INSERT], [SQLITE_DELETE],
** or [SQLITE_UPDATE], depending on the operation that caused the callback
** to be invoked.
** ^The third and fourth arguments to the callback contain pointers to the
** database and table name containing the affected row.
** ^The final callback parameter is the [rowid] of the row.
** ^In the case of an update, this is the [rowid] after the update takes place.
**
** ^(The update hook is not invoked when internal system tables are
** modified (i.e. sqlite_sequence).)^
** ^The update hook is not invoked when [WITHOUT ROWID] tables are modified.
**
** ^In the current implementation, the update hook
** is not invoked when conflicting rows are deleted because of an
** [ON CONFLICT | ON CONFLICT REPLACE] clause.  ^Nor is the update hook
** invoked when rows are deleted using the [truncate optimization].
** The exceptions defined in this paragraph might change in a future
** release of SQLite.
**
** Whether the update hook is invoked before or after the
** corresponding change is currently unspecified and may differ
** depending on the type of change. Do not rely on the order of the
** hook call with regards to the final result of the operation which
** triggers the hook.
**
** The update hook implementation must not do anything that will modify
** the database connection that invoked the update hook.  Any actions
** to modify the database connection must be deferred until after the
** completion of the [sqlite3_step()] call that triggered the update hook.
** Note that [sqlite3_prepare_v2()] and [sqlite3_step()] both modify their
** database connections for the meaning of "modify" in this paragraph.
**
** ^The sqlite3_update_hook(D,C,P) function
** returns the P argument from the previous call
** on the same [database connection] D, or NULL for
** the first call on D.
**
** See also the [sqlite3_commit_hook()], [sqlite3_rollback_hook()],
** and [sqlite3_preupdate_hook()] interfaces.
 */
// llgo:link (*Sqlite3).UpdateHook C.sqlite3_update_hook
func (recv_ *Sqlite3) UpdateHook(func(c.Pointer, c.Int, *c.Char, *c.Char, Int64), c.Pointer) c.Pointer {
	return nil
}

/*
** CAPI3REF: Enable Or Disable Shared Pager Cache
**
** ^(This routine enables or disables the sharing of the database cache
** and schema data structures between [database connection | connections]
** to the same database. Sharing is enabled if the argument is true
** and disabled if the argument is false.)^
**
** This interface is omitted if SQLite is compiled with
** [-DSQLITE_OMIT_SHARED_CACHE].  The [-DSQLITE_OMIT_SHARED_CACHE]
** compile-time option is recommended because the
** [use of shared cache mode is discouraged].
**
** ^Cache sharing is enabled and disabled for an entire process.
** This is a change as of SQLite [version 3.5.0] ([dateof:3.5.0]).
** In prior versions of SQLite,
** sharing was enabled or disabled for each thread separately.
**
** ^(The cache sharing mode set by this interface effects all subsequent
** calls to [sqlite3_open()], [sqlite3_open_v2()], and [sqlite3_open16()].
** Existing database connections continue to use the sharing mode
** that was in effect at the time they were opened.)^
**
** ^(This routine returns [SQLITE_OK] if shared cache was enabled or disabled
** successfully.  An [error code] is returned otherwise.)^
**
** ^Shared cache is disabled by default. It is recommended that it stay
** that way.  In other words, do not use this routine.  This interface
** continues to be provided for historical compatibility, but its use is
** discouraged.  Any use of shared cache is discouraged.  If shared cache
** must be used, it is recommended that shared cache only be enabled for
** individual database connections using the [sqlite3_open_v2()] interface
** with the [SQLITE_OPEN_SHAREDCACHE] flag.
**
** Note: This method is disabled on MacOS X 10.7 and iOS version 5.0
** and will always return SQLITE_MISUSE. On those systems,
** shared cache mode should be enabled per-database connection via
** [sqlite3_open_v2()] with [SQLITE_OPEN_SHAREDCACHE].
**
** This interface is threadsafe on processors where writing a
** 32-bit integer is atomic.
**
** See Also:  [SQLite Shared-Cache Mode]
 */
//go:linkname EnableSharedCache C.sqlite3_enable_shared_cache
func EnableSharedCache(c.Int) c.Int

/*
** CAPI3REF: Attempt To Free Heap Memory
**
** ^The sqlite3_release_memory() interface attempts to free N bytes
** of heap memory by deallocating non-essential memory allocations
** held by the database library.   Memory used to cache database
** pages to improve performance is an example of non-essential memory.
** ^sqlite3_release_memory() returns the number of bytes actually freed,
** which might be more or less than the amount requested.
** ^The sqlite3_release_memory() routine is a no-op returning zero
** if SQLite is not compiled with [SQLITE_ENABLE_MEMORY_MANAGEMENT].
**
** See also: [sqlite3_db_release_memory()]
 */
//go:linkname ReleaseMemory C.sqlite3_release_memory
func ReleaseMemory(c.Int) c.Int

/*
** CAPI3REF: Free Memory Used By A Database Connection
** METHOD: sqlite3
**
** ^The sqlite3_db_release_memory(D) interface attempts to free as much heap
** memory as possible from database connection D. Unlike the
** [sqlite3_release_memory()] interface, this interface is in effect even
** when the [SQLITE_ENABLE_MEMORY_MANAGEMENT] compile-time option is
** omitted.
**
** See also: [sqlite3_release_memory()]
 */
// llgo:link (*Sqlite3).DbReleaseMemory C.sqlite3_db_release_memory
func (recv_ *Sqlite3) DbReleaseMemory() c.Int {
	return 0
}

/*
** CAPI3REF: Impose A Limit On Heap Size
**
** These interfaces impose limits on the amount of heap memory that will be
** by all database connections within a single process.
**
** ^The sqlite3_soft_heap_limit64() interface sets and/or queries the
** soft limit on the amount of heap memory that may be allocated by SQLite.
** ^SQLite strives to keep heap memory utilization below the soft heap
** limit by reducing the number of pages held in the page cache
** as heap memory usages approaches the limit.
** ^The soft heap limit is "soft" because even though SQLite strives to stay
** below the limit, it will exceed the limit rather than generate
** an [SQLITE_NOMEM] error.  In other words, the soft heap limit
** is advisory only.
**
** ^The sqlite3_hard_heap_limit64(N) interface sets a hard upper bound of
** N bytes on the amount of memory that will be allocated.  ^The
** sqlite3_hard_heap_limit64(N) interface is similar to
** sqlite3_soft_heap_limit64(N) except that memory allocations will fail
** when the hard heap limit is reached.
**
** ^The return value from both sqlite3_soft_heap_limit64() and
** sqlite3_hard_heap_limit64() is the size of
** the heap limit prior to the call, or negative in the case of an
** error.  ^If the argument N is negative
** then no change is made to the heap limit.  Hence, the current
** size of heap limits can be determined by invoking
** sqlite3_soft_heap_limit64(-1) or sqlite3_hard_heap_limit(-1).
**
** ^Setting the heap limits to zero disables the heap limiter mechanism.
**
** ^The soft heap limit may not be greater than the hard heap limit.
** ^If the hard heap limit is enabled and if sqlite3_soft_heap_limit(N)
** is invoked with a value of N that is greater than the hard heap limit,
** the soft heap limit is set to the value of the hard heap limit.
** ^The soft heap limit is automatically enabled whenever the hard heap
** limit is enabled. ^When sqlite3_hard_heap_limit64(N) is invoked and
** the soft heap limit is outside the range of 1..N, then the soft heap
** limit is set to N.  ^Invoking sqlite3_soft_heap_limit64(0) when the
** hard heap limit is enabled makes the soft heap limit equal to the
** hard heap limit.
**
** The memory allocation limits can also be adjusted using
** [PRAGMA soft_heap_limit] and [PRAGMA hard_heap_limit].
**
** ^(The heap limits are not enforced in the current implementation
** if one or more of following conditions are true:
**
** <ul>
** <li> The limit value is set to zero.
** <li> Memory accounting is disabled using a combination of the
**      [sqlite3_config]([SQLITE_CONFIG_MEMSTATUS],...) start-time option and
**      the [SQLITE_DEFAULT_MEMSTATUS] compile-time option.
** <li> An alternative page cache implementation is specified using
**      [sqlite3_config]([SQLITE_CONFIG_PCACHE2],...).
** <li> The page cache allocates from its own memory pool supplied
**      by [sqlite3_config]([SQLITE_CONFIG_PAGECACHE],...) rather than
**      from the heap.
** </ul>)^
**
** The circumstances under which SQLite will enforce the heap limits may
** changes in future releases of SQLite.
 */
// llgo:link Int64.SoftHeapLimit64 C.sqlite3_soft_heap_limit64
func (recv_ Int64) SoftHeapLimit64() Int64 {
	return 0
}

// llgo:link Int64.HardHeapLimit64 C.sqlite3_hard_heap_limit64
func (recv_ Int64) HardHeapLimit64() Int64 {
	return 0
}

/*
** CAPI3REF: Deprecated Soft Heap Limit Interface
** DEPRECATED
**
** This is a deprecated version of the [sqlite3_soft_heap_limit64()]
** interface.  This routine is provided for historical compatibility
** only.  All new applications should use the
** [sqlite3_soft_heap_limit64()] interface rather than this one.
 */
//go:linkname SoftHeapLimit C.sqlite3_soft_heap_limit
func SoftHeapLimit(N c.Int)

/*
** CAPI3REF: Extract Metadata About A Column Of A Table
** METHOD: sqlite3
**
** ^(The sqlite3_table_column_metadata(X,D,T,C,....) routine returns
** information about column C of table T in database D
** on [database connection] X.)^  ^The sqlite3_table_column_metadata()
** interface returns SQLITE_OK and fills in the non-NULL pointers in
** the final five arguments with appropriate values if the specified
** column exists.  ^The sqlite3_table_column_metadata() interface returns
** SQLITE_ERROR if the specified column does not exist.
** ^If the column-name parameter to sqlite3_table_column_metadata() is a
** NULL pointer, then this routine simply checks for the existence of the
** table and returns SQLITE_OK if the table exists and SQLITE_ERROR if it
** does not.  If the table name parameter T in a call to
** sqlite3_table_column_metadata(X,D,T,C,...) is NULL then the result is
** undefined behavior.
**
** ^The column is identified by the second, third and fourth parameters to
** this function. ^(The second parameter is either the name of the database
** (i.e. "main", "temp", or an attached database) containing the specified
** table or NULL.)^ ^If it is NULL, then all attached databases are searched
** for the table using the same algorithm used by the database engine to
** resolve unqualified table references.
**
** ^The third and fourth parameters to this function are the table and column
** name of the desired column, respectively.
**
** ^Metadata is returned by writing to the memory locations passed as the 5th
** and subsequent parameters to this function. ^Any of these arguments may be
** NULL, in which case the corresponding element of metadata is omitted.
**
** ^(<blockquote>
** <table border="1">
** <tr><th> Parameter <th> Output<br>Type <th>  Description
**
** <tr><td> 5th <td> const char* <td> Data type
** <tr><td> 6th <td> const char* <td> Name of default collation sequence
** <tr><td> 7th <td> int         <td> True if column has a NOT NULL constraint
** <tr><td> 8th <td> int         <td> True if column is part of the PRIMARY KEY
** <tr><td> 9th <td> int         <td> True if column is [AUTOINCREMENT]
** </table>
** </blockquote>)^
**
** ^The memory pointed to by the character pointers returned for the
** declaration type and collation sequence is valid until the next
** call to any SQLite API function.
**
** ^If the specified table is actually a view, an [error code] is returned.
**
** ^If the specified column is "rowid", "oid" or "_rowid_" and the table
** is not a [WITHOUT ROWID] table and an
** [INTEGER PRIMARY KEY] column has been explicitly declared, then the output
** parameters are set for the explicitly declared column. ^(If there is no
** [INTEGER PRIMARY KEY] column, then the outputs
** for the [rowid] are set as follows:
**
** <pre>
**     data type: "INTEGER"
**     collation sequence: "BINARY"
**     not null: 0
**     primary key: 1
**     auto increment: 0
** </pre>)^
**
** ^This function causes all database schemas to be read from disk and
** parsed, if that has not already been done, and returns an error if
** any errors are encountered while loading the schema.
 */
// llgo:link (*Sqlite3).TableColumnMetadata C.sqlite3_table_column_metadata
func (recv_ *Sqlite3) TableColumnMetadata(zDbName *c.Char, zTableName *c.Char, zColumnName *c.Char, pzDataType **c.Char, pzCollSeq **c.Char, pNotNull *c.Int, pPrimaryKey *c.Int, pAutoinc *c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Load An Extension
** METHOD: sqlite3
**
** ^This interface loads an SQLite extension library from the named file.
**
** ^The sqlite3_load_extension() interface attempts to load an
** [SQLite extension] library contained in the file zFile.  If
** the file cannot be loaded directly, attempts are made to load
** with various operating-system specific extensions added.
** So for example, if "samplelib" cannot be loaded, then names like
** "samplelib.so" or "samplelib.dylib" or "samplelib.dll" might
** be tried also.
**
** ^The entry point is zProc.
** ^(zProc may be 0, in which case SQLite will try to come up with an
** entry point name on its own.  It first tries "sqlite3_extension_init".
** If that does not work, it constructs a name "sqlite3_X_init" where the
** X is consists of the lower-case equivalent of all ASCII alphabetic
** characters in the filename from the last "/" to the first following
** "." and omitting any initial "lib".)^
** ^The sqlite3_load_extension() interface returns
** [SQLITE_OK] on success and [SQLITE_ERROR] if something goes wrong.
** ^If an error occurs and pzErrMsg is not 0, then the
** [sqlite3_load_extension()] interface shall attempt to
** fill *pzErrMsg with error message text stored in memory
** obtained from [sqlite3_malloc()]. The calling function
** should free this memory by calling [sqlite3_free()].
**
** ^Extension loading must be enabled using
** [sqlite3_enable_load_extension()] or
** [sqlite3_db_config](db,[SQLITE_DBCONFIG_ENABLE_LOAD_EXTENSION],1,NULL)
** prior to calling this API,
** otherwise an error will be returned.
**
** <b>Security warning:</b> It is recommended that the
** [SQLITE_DBCONFIG_ENABLE_LOAD_EXTENSION] method be used to enable only this
** interface.  The use of the [sqlite3_enable_load_extension()] interface
** should be avoided.  This will keep the SQL function [load_extension()]
** disabled and prevent SQL injections from giving attackers
** access to extension loading capabilities.
**
** See also the [load_extension() SQL function].
 */
// llgo:link (*Sqlite3).LoadExtension C.sqlite3_load_extension
func (recv_ *Sqlite3) LoadExtension(zFile *c.Char, zProc *c.Char, pzErrMsg **c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Enable Or Disable Extension Loading
** METHOD: sqlite3
**
** ^So as not to open security holes in older applications that are
** unprepared to deal with [extension loading], and as a means of disabling
** [extension loading] while evaluating user-entered SQL, the following API
** is provided to turn the [sqlite3_load_extension()] mechanism on and off.
**
** ^Extension loading is off by default.
** ^Call the sqlite3_enable_load_extension() routine with onoff==1
** to turn extension loading on and call it with onoff==0 to turn
** it back off again.
**
** ^This interface enables or disables both the C-API
** [sqlite3_load_extension()] and the SQL function [load_extension()].
** ^(Use [sqlite3_db_config](db,[SQLITE_DBCONFIG_ENABLE_LOAD_EXTENSION],..)
** to enable or disable only the C-API.)^
**
** <b>Security warning:</b> It is recommended that extension loading
** be enabled using the [SQLITE_DBCONFIG_ENABLE_LOAD_EXTENSION] method
** rather than this interface, so the [load_extension()] SQL function
** remains disabled. This will prevent SQL injections from giving attackers
** access to extension loading capabilities.
 */
// llgo:link (*Sqlite3).EnableLoadExtension C.sqlite3_enable_load_extension
func (recv_ *Sqlite3) EnableLoadExtension(onoff c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Automatically Load Statically Linked Extensions
**
** ^This interface causes the xEntryPoint() function to be invoked for
** each new [database connection] that is created.  The idea here is that
** xEntryPoint() is the entry point for a statically linked [SQLite extension]
** that is to be automatically loaded into all new database connections.
**
** ^(Even though the function prototype shows that xEntryPoint() takes
** no arguments and returns void, SQLite invokes xEntryPoint() with three
** arguments and expects an integer result as if the signature of the
** entry point where as follows:
**
** <blockquote><pre>
** &nbsp;  int xEntryPoint(
** &nbsp;    sqlite3 *db,
** &nbsp;    const char **pzErrMsg,
** &nbsp;    const struct sqlite3_api_routines *pThunk
** &nbsp;  );
** </pre></blockquote>)^
**
** If the xEntryPoint routine encounters an error, it should make *pzErrMsg
** point to an appropriate error message (obtained from [sqlite3_mprintf()])
** and return an appropriate [error code].  ^SQLite ensures that *pzErrMsg
** is NULL before calling the xEntryPoint().  ^SQLite will invoke
** [sqlite3_free()] on *pzErrMsg after xEntryPoint() returns.  ^If any
** xEntryPoint() returns an error, the [sqlite3_open()], [sqlite3_open16()],
** or [sqlite3_open_v2()] call that provoked the xEntryPoint() will fail.
**
** ^Calling sqlite3_auto_extension(X) with an entry point X that is already
** on the list of automatic extensions is a harmless no-op. ^No entry point
** will be called more than once for each database connection that is opened.
**
** See also: [sqlite3_reset_auto_extension()]
** and [sqlite3_cancel_auto_extension()]
 */
//go:linkname AutoExtension C.sqlite3_auto_extension
func AutoExtension(xEntryPoint func()) c.Int

/*
** CAPI3REF: Cancel Automatic Extension Loading
**
** ^The [sqlite3_cancel_auto_extension(X)] interface unregisters the
** initialization routine X that was registered using a prior call to
** [sqlite3_auto_extension(X)].  ^The [sqlite3_cancel_auto_extension(X)]
** routine returns 1 if initialization routine X was successfully
** unregistered and it returns 0 if X was not on the list of initialization
** routines.
 */
//go:linkname CancelAutoExtension C.sqlite3_cancel_auto_extension
func CancelAutoExtension(xEntryPoint func()) c.Int

/*
** CAPI3REF: Reset Automatic Extension Loading
**
** ^This interface disables all automatic extensions previously
** registered using [sqlite3_auto_extension()].
 */
//go:linkname ResetAutoExtension C.sqlite3_reset_auto_extension
func ResetAutoExtension()

type Vtab struct {
	PModule *Module
	NRef    c.Int
	ZErrMsg *c.Char
}

type IndexInfo struct {
	NConstraint      c.Int
	AConstraint      *IndexConstraint
	NOrderBy         c.Int
	AOrderBy         *IndexOrderby
	AConstraintUsage *IndexConstraintUsage
	IdxNum           c.Int
	IdxStr           *c.Char
	NeedToFreeIdxStr c.Int
	OrderByConsumed  c.Int
	EstimatedCost    c.Double
	EstimatedRows    Int64
	IdxFlags         c.Int
	ColUsed          Uint64
}

type VtabCursor struct {
	PVtab *Vtab
}

type Module struct {
	IVersion      c.Int
	XCreate       c.Pointer
	XConnect      c.Pointer
	XBestIndex    c.Pointer
	XDisconnect   c.Pointer
	XDestroy      c.Pointer
	XOpen         c.Pointer
	XClose        c.Pointer
	XFilter       c.Pointer
	XNext         c.Pointer
	XEof          c.Pointer
	XColumn       c.Pointer
	XRowid        c.Pointer
	XUpdate       c.Pointer
	XBegin        c.Pointer
	XSync         c.Pointer
	XCommit       c.Pointer
	XRollback     c.Pointer
	XFindFunction c.Pointer
	XRename       c.Pointer
	XSavepoint    c.Pointer
	XRelease      c.Pointer
	XRollbackTo   c.Pointer
	XShadowName   c.Pointer
	XIntegrity    c.Pointer
}

type IndexConstraint struct {
	Unused [8]uint8
}

type IndexOrderby struct {
	Unused [8]uint8
}

type IndexConstraintUsage struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Register A Virtual Table Implementation
** METHOD: sqlite3
**
** ^These routines are used to register a new [virtual table module] name.
** ^Module names must be registered before
** creating a new [virtual table] using the module and before using a
** preexisting [virtual table] for the module.
**
** ^The module name is registered on the [database connection] specified
** by the first parameter.  ^The name of the module is given by the
** second parameter.  ^The third parameter is a pointer to
** the implementation of the [virtual table module].   ^The fourth
** parameter is an arbitrary client data pointer that is passed through
** into the [xCreate] and [xConnect] methods of the virtual table module
** when a new virtual table is be being created or reinitialized.
**
** ^The sqlite3_create_module_v2() interface has a fifth parameter which
** is a pointer to a destructor for the pClientData.  ^SQLite will
** invoke the destructor function (if it is not NULL) when SQLite
** no longer needs the pClientData pointer.  ^The destructor will also
** be invoked if the call to sqlite3_create_module_v2() fails.
** ^The sqlite3_create_module()
** interface is equivalent to sqlite3_create_module_v2() with a NULL
** destructor.
**
** ^If the third parameter (the pointer to the sqlite3_module object) is
** NULL then no new module is created and any existing modules with the
** same name are dropped.
**
** See also: [sqlite3_drop_modules()]
 */
// llgo:link (*Sqlite3).CreateModule C.sqlite3_create_module
func (recv_ *Sqlite3) CreateModule(zName *c.Char, p *Module, pClientData c.Pointer) c.Int {
	return 0
}

// llgo:link (*Sqlite3).CreateModuleV2 C.sqlite3_create_module_v2
func (recv_ *Sqlite3) CreateModuleV2(zName *c.Char, p *Module, pClientData c.Pointer, xDestroy func(c.Pointer)) c.Int {
	return 0
}

/*
** CAPI3REF: Remove Unnecessary Virtual Table Implementations
** METHOD: sqlite3
**
** ^The sqlite3_drop_modules(D,L) interface removes all virtual
** table modules from database connection D except those named on list L.
** The L parameter must be either NULL or a pointer to an array of pointers
** to strings where the array is terminated by a single NULL pointer.
** ^If the L parameter is NULL, then all virtual table modules are removed.
**
** See also: [sqlite3_create_module()]
 */
// llgo:link (*Sqlite3).DropModules C.sqlite3_drop_modules
func (recv_ *Sqlite3) DropModules(azKeep **c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Declare The Schema Of A Virtual Table
**
** ^The [xCreate] and [xConnect] methods of a
** [virtual table module] call this interface
** to declare the format (the names and datatypes of the columns) of
** the virtual tables they implement.
 */
// llgo:link (*Sqlite3).DeclareVtab C.sqlite3_declare_vtab
func (recv_ *Sqlite3) DeclareVtab(zSQL *c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Overload A Function For A Virtual Table
** METHOD: sqlite3
**
** ^(Virtual tables can provide alternative implementations of functions
** using the [xFindFunction] method of the [virtual table module].
** But global versions of those functions
** must exist in order to be overloaded.)^
**
** ^(This API makes sure a global version of a function with a particular
** name and number of parameters exists.  If no such function exists
** before this API is called, a new function is created.)^  ^The implementation
** of the new function always causes an exception to be thrown.  So
** the new function is not good for anything by itself.  Its only
** purpose is to be a placeholder function that can be overloaded
** by a [virtual table].
 */
// llgo:link (*Sqlite3).OverloadFunction C.sqlite3_overload_function
func (recv_ *Sqlite3) OverloadFunction(zFuncName *c.Char, nArg c.Int) c.Int {
	return 0
}

type Blob struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Open A BLOB For Incremental I/O
** METHOD: sqlite3
** CONSTRUCTOR: sqlite3_blob
**
** ^(This interfaces opens a [BLOB handle | handle] to the BLOB located
** in row iRow, column zColumn, table zTable in database zDb;
** in other words, the same BLOB that would be selected by:
**
** <pre>
**     SELECT zColumn FROM zDb.zTable WHERE [rowid] = iRow;
** </pre>)^
**
** ^(Parameter zDb is not the filename that contains the database, but
** rather the symbolic name of the database. For attached databases, this is
** the name that appears after the AS keyword in the [ATTACH] statement.
** For the main database file, the database name is "main". For TEMP
** tables, the database name is "temp".)^
**
** ^If the flags parameter is non-zero, then the BLOB is opened for read
** and write access. ^If the flags parameter is zero, the BLOB is opened for
** read-only access.
**
** ^(On success, [SQLITE_OK] is returned and the new [BLOB handle] is stored
** in *ppBlob. Otherwise an [error code] is returned and, unless the error
** code is SQLITE_MISUSE, *ppBlob is set to NULL.)^ ^This means that, provided
** the API is not misused, it is always safe to call [sqlite3_blob_close()]
** on *ppBlob after this function it returns.
**
** This function fails with SQLITE_ERROR if any of the following are true:
** <ul>
**   <li> ^(Database zDb does not exist)^,
**   <li> ^(Table zTable does not exist within database zDb)^,
**   <li> ^(Table zTable is a WITHOUT ROWID table)^,
**   <li> ^(Column zColumn does not exist)^,
**   <li> ^(Row iRow is not present in the table)^,
**   <li> ^(The specified column of row iRow contains a value that is not
**         a TEXT or BLOB value)^,
**   <li> ^(Column zColumn is part of an index, PRIMARY KEY or UNIQUE
**         constraint and the blob is being opened for read/write access)^,
**   <li> ^([foreign key constraints | Foreign key constraints] are enabled,
**         column zColumn is part of a [child key] definition and the blob is
**         being opened for read/write access)^.
** </ul>
**
** ^Unless it returns SQLITE_MISUSE, this function sets the
** [database connection] error code and message accessible via
** [sqlite3_errcode()] and [sqlite3_errmsg()] and related functions.
**
** A BLOB referenced by sqlite3_blob_open() may be read using the
** [sqlite3_blob_read()] interface and modified by using
** [sqlite3_blob_write()].  The [BLOB handle] can be moved to a
** different row of the same table using the [sqlite3_blob_reopen()]
** interface.  However, the column, table, or database of a [BLOB handle]
** cannot be changed after the [BLOB handle] is opened.
**
** ^(If the row that a BLOB handle points to is modified by an
** [UPDATE], [DELETE], or by [ON CONFLICT] side-effects
** then the BLOB handle is marked as "expired".
** This is true if any column of the row is changed, even a column
** other than the one the BLOB handle is open on.)^
** ^Calls to [sqlite3_blob_read()] and [sqlite3_blob_write()] for
** an expired BLOB handle fail with a return code of [SQLITE_ABORT].
** ^(Changes written into a BLOB prior to the BLOB expiring are not
** rolled back by the expiration of the BLOB.  Such changes will eventually
** commit if the transaction continues to completion.)^
**
** ^Use the [sqlite3_blob_bytes()] interface to determine the size of
** the opened blob.  ^The size of a blob may not be changed by this
** interface.  Use the [UPDATE] SQL command to change the size of a
** blob.
**
** ^The [sqlite3_bind_zeroblob()] and [sqlite3_result_zeroblob()] interfaces
** and the built-in [zeroblob] SQL function may be used to create a
** zero-filled blob to read or write using the incremental-blob interface.
**
** To avoid a resource leak, every open [BLOB handle] should eventually
** be released by a call to [sqlite3_blob_close()].
**
** See also: [sqlite3_blob_close()],
** [sqlite3_blob_reopen()], [sqlite3_blob_read()],
** [sqlite3_blob_bytes()], [sqlite3_blob_write()].
 */
// llgo:link (*Sqlite3).BlobOpen C.sqlite3_blob_open
func (recv_ *Sqlite3) BlobOpen(zDb *c.Char, zTable *c.Char, zColumn *c.Char, iRow Int64, flags c.Int, ppBlob **Blob) c.Int {
	return 0
}

/*
** CAPI3REF: Move a BLOB Handle to a New Row
** METHOD: sqlite3_blob
**
** ^This function is used to move an existing [BLOB handle] so that it points
** to a different row of the same database table. ^The new row is identified
** by the rowid value passed as the second argument. Only the row can be
** changed. ^The database, table and column on which the blob handle is open
** remain the same. Moving an existing [BLOB handle] to a new row is
** faster than closing the existing handle and opening a new one.
**
** ^(The new row must meet the same criteria as for [sqlite3_blob_open()] -
** it must exist and there must be either a blob or text value stored in
** the nominated column.)^ ^If the new row is not present in the table, or if
** it does not contain a blob or text value, or if another error occurs, an
** SQLite error code is returned and the blob handle is considered aborted.
** ^All subsequent calls to [sqlite3_blob_read()], [sqlite3_blob_write()] or
** [sqlite3_blob_reopen()] on an aborted blob handle immediately return
** SQLITE_ABORT. ^Calling [sqlite3_blob_bytes()] on an aborted blob handle
** always returns zero.
**
** ^This function sets the database handle error code and message.
 */
// llgo:link (*Blob).BlobReopen C.sqlite3_blob_reopen
func (recv_ *Blob) BlobReopen(Int64) c.Int {
	return 0
}

/*
** CAPI3REF: Close A BLOB Handle
** DESTRUCTOR: sqlite3_blob
**
** ^This function closes an open [BLOB handle]. ^(The BLOB handle is closed
** unconditionally.  Even if this routine returns an error code, the
** handle is still closed.)^
**
** ^If the blob handle being closed was opened for read-write access, and if
** the database is in auto-commit mode and there are no other open read-write
** blob handles or active write statements, the current transaction is
** committed. ^If an error occurs while committing the transaction, an error
** code is returned and the transaction rolled back.
**
** Calling this function with an argument that is not a NULL pointer or an
** open blob handle results in undefined behavior. ^Calling this routine
** with a null pointer (such as would be returned by a failed call to
** [sqlite3_blob_open()]) is a harmless no-op. ^Otherwise, if this function
** is passed a valid open blob handle, the values returned by the
** sqlite3_errcode() and sqlite3_errmsg() functions are set before returning.
 */
// llgo:link (*Blob).BlobClose C.sqlite3_blob_close
func (recv_ *Blob) BlobClose() c.Int {
	return 0
}

/*
** CAPI3REF: Return The Size Of An Open BLOB
** METHOD: sqlite3_blob
**
** ^Returns the size in bytes of the BLOB accessible via the
** successfully opened [BLOB handle] in its only argument.  ^The
** incremental blob I/O routines can only read or overwriting existing
** blob content; they cannot change the size of a blob.
**
** This routine only works on a [BLOB handle] which has been created
** by a prior successful call to [sqlite3_blob_open()] and which has not
** been closed by [sqlite3_blob_close()].  Passing any other pointer in
** to this routine results in undefined and probably undesirable behavior.
 */
// llgo:link (*Blob).BlobBytes C.sqlite3_blob_bytes
func (recv_ *Blob) BlobBytes() c.Int {
	return 0
}

/*
** CAPI3REF: Read Data From A BLOB Incrementally
** METHOD: sqlite3_blob
**
** ^(This function is used to read data from an open [BLOB handle] into a
** caller-supplied buffer. N bytes of data are copied into buffer Z
** from the open BLOB, starting at offset iOffset.)^
**
** ^If offset iOffset is less than N bytes from the end of the BLOB,
** [SQLITE_ERROR] is returned and no data is read.  ^If N or iOffset is
** less than zero, [SQLITE_ERROR] is returned and no data is read.
** ^The size of the blob (and hence the maximum value of N+iOffset)
** can be determined using the [sqlite3_blob_bytes()] interface.
**
** ^An attempt to read from an expired [BLOB handle] fails with an
** error code of [SQLITE_ABORT].
**
** ^(On success, sqlite3_blob_read() returns SQLITE_OK.
** Otherwise, an [error code] or an [extended error code] is returned.)^
**
** This routine only works on a [BLOB handle] which has been created
** by a prior successful call to [sqlite3_blob_open()] and which has not
** been closed by [sqlite3_blob_close()].  Passing any other pointer in
** to this routine results in undefined and probably undesirable behavior.
**
** See also: [sqlite3_blob_write()].
 */
// llgo:link (*Blob).BlobRead C.sqlite3_blob_read
func (recv_ *Blob) BlobRead(Z c.Pointer, N c.Int, iOffset c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Write Data Into A BLOB Incrementally
** METHOD: sqlite3_blob
**
** ^(This function is used to write data into an open [BLOB handle] from a
** caller-supplied buffer. N bytes of data are copied from the buffer Z
** into the open BLOB, starting at offset iOffset.)^
**
** ^(On success, sqlite3_blob_write() returns SQLITE_OK.
** Otherwise, an  [error code] or an [extended error code] is returned.)^
** ^Unless SQLITE_MISUSE is returned, this function sets the
** [database connection] error code and message accessible via
** [sqlite3_errcode()] and [sqlite3_errmsg()] and related functions.
**
** ^If the [BLOB handle] passed as the first argument was not opened for
** writing (the flags parameter to [sqlite3_blob_open()] was zero),
** this function returns [SQLITE_READONLY].
**
** This function may only modify the contents of the BLOB; it is
** not possible to increase the size of a BLOB using this API.
** ^If offset iOffset is less than N bytes from the end of the BLOB,
** [SQLITE_ERROR] is returned and no data is written. The size of the
** BLOB (and hence the maximum value of N+iOffset) can be determined
** using the [sqlite3_blob_bytes()] interface. ^If N or iOffset are less
** than zero [SQLITE_ERROR] is returned and no data is written.
**
** ^An attempt to write to an expired [BLOB handle] fails with an
** error code of [SQLITE_ABORT].  ^Writes to the BLOB that occurred
** before the [BLOB handle] expired are not rolled back by the
** expiration of the handle, though of course those changes might
** have been overwritten by the statement that expired the BLOB handle
** or by other independent statements.
**
** This routine only works on a [BLOB handle] which has been created
** by a prior successful call to [sqlite3_blob_open()] and which has not
** been closed by [sqlite3_blob_close()].  Passing any other pointer in
** to this routine results in undefined and probably undesirable behavior.
**
** See also: [sqlite3_blob_read()].
 */
// llgo:link (*Blob).BlobWrite C.sqlite3_blob_write
func (recv_ *Blob) BlobWrite(z c.Pointer, n c.Int, iOffset c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Virtual File System Objects
**
** A virtual filesystem (VFS) is an [sqlite3_vfs] object
** that SQLite uses to interact
** with the underlying operating system.  Most SQLite builds come with a
** single default VFS that is appropriate for the host computer.
** New VFSes can be registered and existing VFSes can be unregistered.
** The following interfaces are provided.
**
** ^The sqlite3_vfs_find() interface returns a pointer to a VFS given its name.
** ^Names are case sensitive.
** ^Names are zero-terminated UTF-8 strings.
** ^If there is no match, a NULL pointer is returned.
** ^If zVfsName is NULL then the default VFS is returned.
**
** ^New VFSes are registered with sqlite3_vfs_register().
** ^Each new VFS becomes the default VFS if the makeDflt flag is set.
** ^The same VFS can be registered multiple times without injury.
** ^To make an existing VFS into the default VFS, register it again
** with the makeDflt flag set.  If two different VFSes with the
** same name are registered, the behavior is undefined.  If a
** VFS is registered with a name that is NULL or an empty string,
** then the behavior is undefined.
**
** ^Unregister a VFS with the sqlite3_vfs_unregister() interface.
** ^(If the default VFS is unregistered, another VFS is chosen as
** the default.  The choice for the new VFS is arbitrary.)^
 */
//go:linkname VfsFind C.sqlite3_vfs_find
func VfsFind(zVfsName *c.Char) *Vfs

// llgo:link (*Vfs).VfsRegister C.sqlite3_vfs_register
func (recv_ *Vfs) VfsRegister(makeDflt c.Int) c.Int {
	return 0
}

// llgo:link (*Vfs).VfsUnregister C.sqlite3_vfs_unregister
func (recv_ *Vfs) VfsUnregister() c.Int {
	return 0
}

/*
** CAPI3REF: Mutexes
**
** The SQLite core uses these routines for thread
** synchronization. Though they are intended for internal
** use by SQLite, code that links against SQLite is
** permitted to use any of these routines.
**
** The SQLite source code contains multiple implementations
** of these mutex routines.  An appropriate implementation
** is selected automatically at compile-time.  The following
** implementations are available in the SQLite core:
**
** <ul>
** <li>   SQLITE_MUTEX_PTHREADS
** <li>   SQLITE_MUTEX_W32
** <li>   SQLITE_MUTEX_NOOP
** </ul>
**
** The SQLITE_MUTEX_NOOP implementation is a set of routines
** that does no real locking and is appropriate for use in
** a single-threaded application.  The SQLITE_MUTEX_PTHREADS and
** SQLITE_MUTEX_W32 implementations are appropriate for use on Unix
** and Windows.
**
** If SQLite is compiled with the SQLITE_MUTEX_APPDEF preprocessor
** macro defined (with "-DSQLITE_MUTEX_APPDEF=1"), then no mutex
** implementation is included with the library. In this case the
** application must supply a custom mutex implementation using the
** [SQLITE_CONFIG_MUTEX] option of the sqlite3_config() function
** before calling sqlite3_initialize() or any other public sqlite3_
** function that calls sqlite3_initialize().
**
** ^The sqlite3_mutex_alloc() routine allocates a new
** mutex and returns a pointer to it. ^The sqlite3_mutex_alloc()
** routine returns NULL if it is unable to allocate the requested
** mutex.  The argument to sqlite3_mutex_alloc() must one of these
** integer constants:
**
** <ul>
** <li>  SQLITE_MUTEX_FAST
** <li>  SQLITE_MUTEX_RECURSIVE
** <li>  SQLITE_MUTEX_STATIC_MAIN
** <li>  SQLITE_MUTEX_STATIC_MEM
** <li>  SQLITE_MUTEX_STATIC_OPEN
** <li>  SQLITE_MUTEX_STATIC_PRNG
** <li>  SQLITE_MUTEX_STATIC_LRU
** <li>  SQLITE_MUTEX_STATIC_PMEM
** <li>  SQLITE_MUTEX_STATIC_APP1
** <li>  SQLITE_MUTEX_STATIC_APP2
** <li>  SQLITE_MUTEX_STATIC_APP3
** <li>  SQLITE_MUTEX_STATIC_VFS1
** <li>  SQLITE_MUTEX_STATIC_VFS2
** <li>  SQLITE_MUTEX_STATIC_VFS3
** </ul>
**
** ^The first two constants (SQLITE_MUTEX_FAST and SQLITE_MUTEX_RECURSIVE)
** cause sqlite3_mutex_alloc() to create
** a new mutex.  ^The new mutex is recursive when SQLITE_MUTEX_RECURSIVE
** is used but not necessarily so when SQLITE_MUTEX_FAST is used.
** The mutex implementation does not need to make a distinction
** between SQLITE_MUTEX_RECURSIVE and SQLITE_MUTEX_FAST if it does
** not want to.  SQLite will only request a recursive mutex in
** cases where it really needs one.  If a faster non-recursive mutex
** implementation is available on the host platform, the mutex subsystem
** might return such a mutex in response to SQLITE_MUTEX_FAST.
**
** ^The other allowed parameters to sqlite3_mutex_alloc() (anything other
** than SQLITE_MUTEX_FAST and SQLITE_MUTEX_RECURSIVE) each return
** a pointer to a static preexisting mutex.  ^Nine static mutexes are
** used by the current version of SQLite.  Future versions of SQLite
** may add additional static mutexes.  Static mutexes are for internal
** use by SQLite only.  Applications that use SQLite mutexes should
** use only the dynamic mutexes returned by SQLITE_MUTEX_FAST or
** SQLITE_MUTEX_RECURSIVE.
**
** ^Note that if one of the dynamic mutex parameters (SQLITE_MUTEX_FAST
** or SQLITE_MUTEX_RECURSIVE) is used then sqlite3_mutex_alloc()
** returns a different mutex on every call.  ^For the static
** mutex types, the same mutex is returned on every call that has
** the same type number.
**
** ^The sqlite3_mutex_free() routine deallocates a previously
** allocated dynamic mutex.  Attempting to deallocate a static
** mutex results in undefined behavior.
**
** ^The sqlite3_mutex_enter() and sqlite3_mutex_try() routines attempt
** to enter a mutex.  ^If another thread is already within the mutex,
** sqlite3_mutex_enter() will block and sqlite3_mutex_try() will return
** SQLITE_BUSY.  ^The sqlite3_mutex_try() interface returns [SQLITE_OK]
** upon successful entry.  ^(Mutexes created using
** SQLITE_MUTEX_RECURSIVE can be entered multiple times by the same thread.
** In such cases, the
** mutex must be exited an equal number of times before another thread
** can enter.)^  If the same thread tries to enter any mutex other
** than an SQLITE_MUTEX_RECURSIVE more than once, the behavior is undefined.
**
** ^(Some systems (for example, Windows 95) do not support the operation
** implemented by sqlite3_mutex_try().  On those systems, sqlite3_mutex_try()
** will always return SQLITE_BUSY. In most cases the SQLite core only uses
** sqlite3_mutex_try() as an optimization, so this is acceptable
** behavior. The exceptions are unix builds that set the
** SQLITE_ENABLE_SETLK_TIMEOUT build option. In that case a working
** sqlite3_mutex_try() is required.)^
**
** ^The sqlite3_mutex_leave() routine exits a mutex that was
** previously entered by the same thread.   The behavior
** is undefined if the mutex is not currently entered by the
** calling thread or is not currently allocated.
**
** ^If the argument to sqlite3_mutex_enter(), sqlite3_mutex_try(),
** sqlite3_mutex_leave(), or sqlite3_mutex_free() is a NULL pointer,
** then any of the four routines behaves as a no-op.
**
** See also: [sqlite3_mutex_held()] and [sqlite3_mutex_notheld()].
 */
//go:linkname MutexAlloc C.sqlite3_mutex_alloc
func MutexAlloc(c.Int) *Mutex

// llgo:link (*Mutex).MutexFree C.sqlite3_mutex_free
func (recv_ *Mutex) MutexFree() {
}

// llgo:link (*Mutex).MutexEnter C.sqlite3_mutex_enter
func (recv_ *Mutex) MutexEnter() {
}

// llgo:link (*Mutex).MutexTry C.sqlite3_mutex_try
func (recv_ *Mutex) MutexTry() c.Int {
	return 0
}

// llgo:link (*Mutex).MutexLeave C.sqlite3_mutex_leave
func (recv_ *Mutex) MutexLeave() {
}

type MutexMethods struct {
	XMutexInit    c.Pointer
	XMutexEnd     c.Pointer
	XMutexAlloc   c.Pointer
	XMutexFree    c.Pointer
	XMutexEnter   c.Pointer
	XMutexTry     c.Pointer
	XMutexLeave   c.Pointer
	XMutexHeld    c.Pointer
	XMutexNotheld c.Pointer
}

/*
** CAPI3REF: Retrieve the mutex for a database connection
** METHOD: sqlite3
**
** ^This interface returns a pointer the [sqlite3_mutex] object that
** serializes access to the [database connection] given in the argument
** when the [threading mode] is Serialized.
** ^If the [threading mode] is Single-thread or Multi-thread then this
** routine returns a NULL pointer.
 */
// llgo:link (*Sqlite3).DbMutex C.sqlite3_db_mutex
func (recv_ *Sqlite3) DbMutex() *Mutex {
	return nil
}

/*
** CAPI3REF: Low-Level Control Of Database Files
** METHOD: sqlite3
** KEYWORDS: {file control}
**
** ^The [sqlite3_file_control()] interface makes a direct call to the
** xFileControl method for the [sqlite3_io_methods] object associated
** with a particular database identified by the second argument. ^The
** name of the database is "main" for the main database or "temp" for the
** TEMP database, or the name that appears after the AS keyword for
** databases that are added using the [ATTACH] SQL command.
** ^A NULL pointer can be used in place of "main" to refer to the
** main database file.
** ^The third and fourth parameters to this routine
** are passed directly through to the second and third parameters of
** the xFileControl method.  ^The return value of the xFileControl
** method becomes the return value of this routine.
**
** A few opcodes for [sqlite3_file_control()] are handled directly
** by the SQLite core and never invoke the
** sqlite3_io_methods.xFileControl method.
** ^The [SQLITE_FCNTL_FILE_POINTER] value for the op parameter causes
** a pointer to the underlying [sqlite3_file] object to be written into
** the space pointed to by the 4th parameter.  The
** [SQLITE_FCNTL_JOURNAL_POINTER] works similarly except that it returns
** the [sqlite3_file] object associated with the journal file instead of
** the main database.  The [SQLITE_FCNTL_VFS_POINTER] opcode returns
** a pointer to the underlying [sqlite3_vfs] object for the file.
** The [SQLITE_FCNTL_DATA_VERSION] returns the data version counter
** from the pager.
**
** ^If the second parameter (zDbName) does not match the name of any
** open database file, then SQLITE_ERROR is returned.  ^This error
** code is not remembered and will not be recalled by [sqlite3_errcode()]
** or [sqlite3_errmsg()].  The underlying xFileControl method might
** also return SQLITE_ERROR.  There is no way to distinguish between
** an incorrect zDbName and an SQLITE_ERROR return from the underlying
** xFileControl method.
**
** See also: [file control opcodes]
 */
// llgo:link (*Sqlite3).FileControl C.sqlite3_file_control
func (recv_ *Sqlite3) FileControl(zDbName *c.Char, op c.Int, __llgo_arg_2 c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: Testing Interface
**
** ^The sqlite3_test_control() interface is used to read out internal
** state of SQLite and to inject faults into SQLite for testing
** purposes.  ^The first parameter is an operation code that determines
** the number, meaning, and operation of all subsequent parameters.
**
** This interface is not for use by applications.  It exists solely
** for verifying the correct operation of the SQLite library.  Depending
** on how the SQLite library is compiled, this interface might not exist.
**
** The details of the operation codes, their meanings, the parameters
** they take, and what they do are all subject to change without notice.
** Unlike most of the SQLite API, this function is not guaranteed to
** operate consistently from one release to the next.
 */
//go:linkname TestControl C.sqlite3_test_control
func TestControl(op c.Int, __llgo_va_list ...interface{}) c.Int

/*
** CAPI3REF: SQL Keyword Checking
**
** These routines provide access to the set of SQL language keywords
** recognized by SQLite.  Applications can uses these routines to determine
** whether or not a specific identifier needs to be escaped (for example,
** by enclosing in double-quotes) so as not to confuse the parser.
**
** The sqlite3_keyword_count() interface returns the number of distinct
** keywords understood by SQLite.
**
** The sqlite3_keyword_name(N,Z,L) interface finds the 0-based N-th keyword and
** makes *Z point to that keyword expressed as UTF8 and writes the number
** of bytes in the keyword into *L.  The string that *Z points to is not
** zero-terminated.  The sqlite3_keyword_name(N,Z,L) routine returns
** SQLITE_OK if N is within bounds and SQLITE_ERROR if not. If either Z
** or L are NULL or invalid pointers then calls to
** sqlite3_keyword_name(N,Z,L) result in undefined behavior.
**
** The sqlite3_keyword_check(Z,L) interface checks to see whether or not
** the L-byte UTF8 identifier that Z points to is a keyword, returning non-zero
** if it is and zero if not.
**
** The parser used by SQLite is forgiving.  It is often possible to use
** a keyword as an identifier as long as such use does not result in a
** parsing ambiguity.  For example, the statement
** "CREATE TABLE BEGIN(REPLACE,PRAGMA,END);" is accepted by SQLite, and
** creates a new table named "BEGIN" with three columns named
** "REPLACE", "PRAGMA", and "END".  Nevertheless, best practice is to avoid
** using keywords as identifiers.  Common techniques used to avoid keyword
** name collisions include:
** <ul>
** <li> Put all identifier names inside double-quotes.  This is the official
**      SQL way to escape identifier names.
** <li> Put identifier names inside &#91;...&#93;.  This is not standard SQL,
**      but it is what SQL Server does and so lots of programmers use this
**      technique.
** <li> Begin every identifier with the letter "Z" as no SQL keywords start
**      with "Z".
** <li> Include a digit somewhere in every identifier name.
** </ul>
**
** Note that the number of keywords understood by SQLite can depend on
** compile-time options.  For example, "VACUUM" is not a keyword if
** SQLite is compiled with the [-DSQLITE_OMIT_VACUUM] option.  Also,
** new keywords may be added to future releases of SQLite.
 */
//go:linkname KeywordCount C.sqlite3_keyword_count
func KeywordCount() c.Int

//go:linkname KeywordName C.sqlite3_keyword_name
func KeywordName(c.Int, **c.Char, *c.Int) c.Int

//go:linkname KeywordCheck C.sqlite3_keyword_check
func KeywordCheck(*c.Char, c.Int) c.Int

type Str struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Create A New Dynamic String Object
** CONSTRUCTOR: sqlite3_str
**
** ^The [sqlite3_str_new(D)] interface allocates and initializes
** a new [sqlite3_str] object.  To avoid memory leaks, the object returned by
** [sqlite3_str_new()] must be freed by a subsequent call to
** [sqlite3_str_finish(X)].
**
** ^The [sqlite3_str_new(D)] interface always returns a pointer to a
** valid [sqlite3_str] object, though in the event of an out-of-memory
** error the returned object might be a special singleton that will
** silently reject new text, always return SQLITE_NOMEM from
** [sqlite3_str_errcode()], always return 0 for
** [sqlite3_str_length()], and always return NULL from
** [sqlite3_str_finish(X)].  It is always safe to use the value
** returned by [sqlite3_str_new(D)] as the sqlite3_str parameter
** to any of the other [sqlite3_str] methods.
**
** The D parameter to [sqlite3_str_new(D)] may be NULL.  If the
** D parameter in [sqlite3_str_new(D)] is not NULL, then the maximum
** length of the string contained in the [sqlite3_str] object will be
** the value set for [sqlite3_limit](D,[SQLITE_LIMIT_LENGTH]) instead
** of [SQLITE_MAX_LENGTH].
 */
// llgo:link (*Sqlite3).StrNew C.sqlite3_str_new
func (recv_ *Sqlite3) StrNew() *Str {
	return nil
}

/*
** CAPI3REF: Finalize A Dynamic String
** DESTRUCTOR: sqlite3_str
**
** ^The [sqlite3_str_finish(X)] interface destroys the sqlite3_str object X
** and returns a pointer to a memory buffer obtained from [sqlite3_malloc64()]
** that contains the constructed string.  The calling application should
** pass the returned value to [sqlite3_free()] to avoid a memory leak.
** ^The [sqlite3_str_finish(X)] interface may return a NULL pointer if any
** errors were encountered during construction of the string.  ^The
** [sqlite3_str_finish(X)] interface will also return a NULL pointer if the
** string in [sqlite3_str] object X is zero bytes long.
 */
// llgo:link (*Str).StrFinish C.sqlite3_str_finish
func (recv_ *Str) StrFinish() *c.Char {
	return nil
}

/*
** CAPI3REF: Add Content To A Dynamic String
** METHOD: sqlite3_str
**
** These interfaces add content to an sqlite3_str object previously obtained
** from [sqlite3_str_new()].
**
** ^The [sqlite3_str_appendf(X,F,...)] and
** [sqlite3_str_vappendf(X,F,V)] interfaces uses the [built-in printf]
** functionality of SQLite to append formatted text onto the end of
** [sqlite3_str] object X.
**
** ^The [sqlite3_str_append(X,S,N)] method appends exactly N bytes from string S
** onto the end of the [sqlite3_str] object X.  N must be non-negative.
** S must contain at least N non-zero bytes of content.  To append a
** zero-terminated string in its entirety, use the [sqlite3_str_appendall()]
** method instead.
**
** ^The [sqlite3_str_appendall(X,S)] method appends the complete content of
** zero-terminated string S onto the end of [sqlite3_str] object X.
**
** ^The [sqlite3_str_appendchar(X,N,C)] method appends N copies of the
** single-byte character C onto the end of [sqlite3_str] object X.
** ^This method can be used, for example, to add whitespace indentation.
**
** ^The [sqlite3_str_reset(X)] method resets the string under construction
** inside [sqlite3_str] object X back to zero bytes in length.
**
** These methods do not return a result code.  ^If an error occurs, that fact
** is recorded in the [sqlite3_str] object and can be recovered by a
** subsequent call to [sqlite3_str_errcode(X)].
 */
//go:linkname StrAppendf C.sqlite3_str_appendf
func StrAppendf(__llgo_arg_0 *Str, zFormat *c.Char, __llgo_va_list ...interface{})

// llgo:link (*Str).StrVappendf C.sqlite3_str_vappendf
func (recv_ *Str) StrVappendf(zFormat *c.Char, __llgo_arg_1 c.VaList) {
}

// llgo:link (*Str).StrAppend C.sqlite3_str_append
func (recv_ *Str) StrAppend(zIn *c.Char, N c.Int) {
}

// llgo:link (*Str).StrAppendall C.sqlite3_str_appendall
func (recv_ *Str) StrAppendall(zIn *c.Char) {
}

// llgo:link (*Str).StrAppendchar C.sqlite3_str_appendchar
func (recv_ *Str) StrAppendchar(N c.Int, C c.Char) {
}

// llgo:link (*Str).StrReset C.sqlite3_str_reset
func (recv_ *Str) StrReset() {
}

/*
** CAPI3REF: Status Of A Dynamic String
** METHOD: sqlite3_str
**
** These interfaces return the current status of an [sqlite3_str] object.
**
** ^If any prior errors have occurred while constructing the dynamic string
** in sqlite3_str X, then the [sqlite3_str_errcode(X)] method will return
** an appropriate error code.  ^The [sqlite3_str_errcode(X)] method returns
** [SQLITE_NOMEM] following any out-of-memory error, or
** [SQLITE_TOOBIG] if the size of the dynamic string exceeds
** [SQLITE_MAX_LENGTH], or [SQLITE_OK] if there have been no errors.
**
** ^The [sqlite3_str_length(X)] method returns the current length, in bytes,
** of the dynamic string under construction in [sqlite3_str] object X.
** ^The length returned by [sqlite3_str_length(X)] does not include the
** zero-termination byte.
**
** ^The [sqlite3_str_value(X)] method returns a pointer to the current
** content of the dynamic string under construction in X.  The value
** returned by [sqlite3_str_value(X)] is managed by the sqlite3_str object X
** and might be freed or altered by any subsequent method on the same
** [sqlite3_str] object.  Applications must not used the pointer returned
** [sqlite3_str_value(X)] after any subsequent method call on the same
** object.  ^Applications may change the content of the string returned
** by [sqlite3_str_value(X)] as long as they do not write into any bytes
** outside the range of 0 to [sqlite3_str_length(X)] and do not read or
** write any byte after any subsequent sqlite3_str method call.
 */
// llgo:link (*Str).StrErrcode C.sqlite3_str_errcode
func (recv_ *Str) StrErrcode() c.Int {
	return 0
}

// llgo:link (*Str).StrLength C.sqlite3_str_length
func (recv_ *Str) StrLength() c.Int {
	return 0
}

// llgo:link (*Str).StrValue C.sqlite3_str_value
func (recv_ *Str) StrValue() *c.Char {
	return nil
}

/*
** CAPI3REF: SQLite Runtime Status
**
** ^These interfaces are used to retrieve runtime status information
** about the performance of SQLite, and optionally to reset various
** highwater marks.  ^The first argument is an integer code for
** the specific parameter to measure.  ^(Recognized integer codes
** are of the form [status parameters | SQLITE_STATUS_...].)^
** ^The current value of the parameter is returned into *pCurrent.
** ^The highest recorded value is returned in *pHighwater.  ^If the
** resetFlag is true, then the highest record value is reset after
** *pHighwater is written.  ^(Some parameters do not record the highest
** value.  For those parameters
** nothing is written into *pHighwater and the resetFlag is ignored.)^
** ^(Other parameters record only the highwater mark and not the current
** value.  For these latter parameters nothing is written into *pCurrent.)^
**
** ^The sqlite3_status() and sqlite3_status64() routines return
** SQLITE_OK on success and a non-zero [error code] on failure.
**
** If either the current value or the highwater mark is too large to
** be represented by a 32-bit integer, then the values returned by
** sqlite3_status() are undefined.
**
** See also: [sqlite3_db_status()]
 */
//go:linkname Status C.sqlite3_status
func Status(op c.Int, pCurrent *c.Int, pHighwater *c.Int, resetFlag c.Int) c.Int

//go:linkname Status64 C.sqlite3_status64
func Status64(op c.Int, pCurrent *Int64, pHighwater *Int64, resetFlag c.Int) c.Int

/*
** CAPI3REF: Database Connection Status
** METHOD: sqlite3
**
** ^This interface is used to retrieve runtime status information
** about a single [database connection].  ^The first argument is the
** database connection object to be interrogated.  ^The second argument
** is an integer constant, taken from the set of
** [SQLITE_DBSTATUS options], that
** determines the parameter to interrogate.  The set of
** [SQLITE_DBSTATUS options] is likely
** to grow in future releases of SQLite.
**
** ^The current value of the requested parameter is written into *pCur
** and the highest instantaneous value is written into *pHiwtr.  ^If
** the resetFlg is true, then the highest instantaneous value is
** reset back down to the current value.
**
** ^The sqlite3_db_status() routine returns SQLITE_OK on success and a
** non-zero [error code] on failure.
**
** See also: [sqlite3_status()] and [sqlite3_stmt_status()].
 */
// llgo:link (*Sqlite3).DbStatus C.sqlite3_db_status
func (recv_ *Sqlite3) DbStatus(op c.Int, pCur *c.Int, pHiwtr *c.Int, resetFlg c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Prepared Statement Status
** METHOD: sqlite3_stmt
**
** ^(Each prepared statement maintains various
** [SQLITE_STMTSTATUS counters] that measure the number
** of times it has performed specific operations.)^  These counters can
** be used to monitor the performance characteristics of the prepared
** statements.  For example, if the number of table steps greatly exceeds
** the number of table searches or result rows, that would tend to indicate
** that the prepared statement is using a full table scan rather than
** an index.
**
** ^(This interface is used to retrieve and reset counter values from
** a [prepared statement].  The first argument is the prepared statement
** object to be interrogated.  The second argument
** is an integer code for a specific [SQLITE_STMTSTATUS counter]
** to be interrogated.)^
** ^The current value of the requested counter is returned.
** ^If the resetFlg is true, then the counter is reset to zero after this
** interface call returns.
**
** See also: [sqlite3_status()] and [sqlite3_db_status()].
 */
// llgo:link (*Stmt).StmtStatus C.sqlite3_stmt_status
func (recv_ *Stmt) StmtStatus(op c.Int, resetFlg c.Int) c.Int {
	return 0
}

type Pcache struct {
	Unused [8]uint8
}

type PcachePage struct {
	PBuf   c.Pointer
	PExtra c.Pointer
}

type PcacheMethods2 struct {
	IVersion   c.Int
	PArg       c.Pointer
	XInit      c.Pointer
	XShutdown  c.Pointer
	XCreate    c.Pointer
	XCachesize c.Pointer
	XPagecount c.Pointer
	XFetch     c.Pointer
	XUnpin     c.Pointer
	XRekey     c.Pointer
	XTruncate  c.Pointer
	XDestroy   c.Pointer
	XShrink    c.Pointer
}

type PcacheMethods struct {
	PArg       c.Pointer
	XInit      c.Pointer
	XShutdown  c.Pointer
	XCreate    c.Pointer
	XCachesize c.Pointer
	XPagecount c.Pointer
	XFetch     c.Pointer
	XUnpin     c.Pointer
	XRekey     c.Pointer
	XTruncate  c.Pointer
	XDestroy   c.Pointer
}

type Backup struct {
	Unused [8]uint8
}

/*
** CAPI3REF: Online Backup API.
**
** The backup API copies the content of one database into another.
** It is useful either for creating backups of databases or
** for copying in-memory databases to or from persistent files.
**
** See Also: [Using the SQLite Online Backup API]
**
** ^SQLite holds a write transaction open on the destination database file
** for the duration of the backup operation.
** ^The source database is read-locked only while it is being read;
** it is not locked continuously for the entire backup operation.
** ^Thus, the backup may be performed on a live source database without
** preventing other database connections from
** reading or writing to the source database while the backup is underway.
**
** ^(To perform a backup operation:
**   <ol>
**     <li><b>sqlite3_backup_init()</b> is called once to initialize the
**         backup,
**     <li><b>sqlite3_backup_step()</b> is called one or more times to transfer
**         the data between the two databases, and finally
**     <li><b>sqlite3_backup_finish()</b> is called to release all resources
**         associated with the backup operation.
**   </ol>)^
** There should be exactly one call to sqlite3_backup_finish() for each
** successful call to sqlite3_backup_init().
**
** [[sqlite3_backup_init()]] <b>sqlite3_backup_init()</b>
**
** ^The D and N arguments to sqlite3_backup_init(D,N,S,M) are the
** [database connection] associated with the destination database
** and the database name, respectively.
** ^The database name is "main" for the main database, "temp" for the
** temporary database, or the name specified after the AS keyword in
** an [ATTACH] statement for an attached database.
** ^The S and M arguments passed to
** sqlite3_backup_init(D,N,S,M) identify the [database connection]
** and database name of the source database, respectively.
** ^The source and destination [database connections] (parameters S and D)
** must be different or else sqlite3_backup_init(D,N,S,M) will fail with
** an error.
**
** ^A call to sqlite3_backup_init() will fail, returning NULL, if
** there is already a read or read-write transaction open on the
** destination database.
**
** ^If an error occurs within sqlite3_backup_init(D,N,S,M), then NULL is
** returned and an error code and error message are stored in the
** destination [database connection] D.
** ^The error code and message for the failed call to sqlite3_backup_init()
** can be retrieved using the [sqlite3_errcode()], [sqlite3_errmsg()], and/or
** [sqlite3_errmsg16()] functions.
** ^A successful call to sqlite3_backup_init() returns a pointer to an
** [sqlite3_backup] object.
** ^The [sqlite3_backup] object may be used with the sqlite3_backup_step() and
** sqlite3_backup_finish() functions to perform the specified backup
** operation.
**
** [[sqlite3_backup_step()]] <b>sqlite3_backup_step()</b>
**
** ^Function sqlite3_backup_step(B,N) will copy up to N pages between
** the source and destination databases specified by [sqlite3_backup] object B.
** ^If N is negative, all remaining source pages are copied.
** ^If sqlite3_backup_step(B,N) successfully copies N pages and there
** are still more pages to be copied, then the function returns [SQLITE_OK].
** ^If sqlite3_backup_step(B,N) successfully finishes copying all pages
** from source to destination, then it returns [SQLITE_DONE].
** ^If an error occurs while running sqlite3_backup_step(B,N),
** then an [error code] is returned. ^As well as [SQLITE_OK] and
** [SQLITE_DONE], a call to sqlite3_backup_step() may return [SQLITE_READONLY],
** [SQLITE_NOMEM], [SQLITE_BUSY], [SQLITE_LOCKED], or an
** [SQLITE_IOERR_ACCESS | SQLITE_IOERR_XXX] extended error code.
**
** ^(The sqlite3_backup_step() might return [SQLITE_READONLY] if
** <ol>
** <li> the destination database was opened read-only, or
** <li> the destination database is using write-ahead-log journaling
** and the destination and source page sizes differ, or
** <li> the destination database is an in-memory database and the
** destination and source page sizes differ.
** </ol>)^
**
** ^If sqlite3_backup_step() cannot obtain a required file-system lock, then
** the [sqlite3_busy_handler | busy-handler function]
** is invoked (if one is specified). ^If the
** busy-handler returns non-zero before the lock is available, then
** [SQLITE_BUSY] is returned to the caller. ^In this case the call to
** sqlite3_backup_step() can be retried later. ^If the source
** [database connection]
** is being used to write to the source database when sqlite3_backup_step()
** is called, then [SQLITE_LOCKED] is returned immediately. ^Again, in this
** case the call to sqlite3_backup_step() can be retried later on. ^(If
** [SQLITE_IOERR_ACCESS | SQLITE_IOERR_XXX], [SQLITE_NOMEM], or
** [SQLITE_READONLY] is returned, then
** there is no point in retrying the call to sqlite3_backup_step(). These
** errors are considered fatal.)^  The application must accept
** that the backup operation has failed and pass the backup operation handle
** to the sqlite3_backup_finish() to release associated resources.
**
** ^The first call to sqlite3_backup_step() obtains an exclusive lock
** on the destination file. ^The exclusive lock is not released until either
** sqlite3_backup_finish() is called or the backup operation is complete
** and sqlite3_backup_step() returns [SQLITE_DONE].  ^Every call to
** sqlite3_backup_step() obtains a [shared lock] on the source database that
** lasts for the duration of the sqlite3_backup_step() call.
** ^Because the source database is not locked between calls to
** sqlite3_backup_step(), the source database may be modified mid-way
** through the backup process.  ^If the source database is modified by an
** external process or via a database connection other than the one being
** used by the backup operation, then the backup will be automatically
** restarted by the next call to sqlite3_backup_step(). ^If the source
** database is modified by the using the same database connection as is used
** by the backup operation, then the backup database is automatically
** updated at the same time.
**
** [[sqlite3_backup_finish()]] <b>sqlite3_backup_finish()</b>
**
** When sqlite3_backup_step() has returned [SQLITE_DONE], or when the
** application wishes to abandon the backup operation, the application
** should destroy the [sqlite3_backup] by passing it to sqlite3_backup_finish().
** ^The sqlite3_backup_finish() interfaces releases all
** resources associated with the [sqlite3_backup] object.
** ^If sqlite3_backup_step() has not yet returned [SQLITE_DONE], then any
** active write-transaction on the destination database is rolled back.
** The [sqlite3_backup] object is invalid
** and may not be used following a call to sqlite3_backup_finish().
**
** ^The value returned by sqlite3_backup_finish is [SQLITE_OK] if no
** sqlite3_backup_step() errors occurred, regardless or whether or not
** sqlite3_backup_step() completed.
** ^If an out-of-memory condition or IO error occurred during any prior
** sqlite3_backup_step() call on the same [sqlite3_backup] object, then
** sqlite3_backup_finish() returns the corresponding [error code].
**
** ^A return of [SQLITE_BUSY] or [SQLITE_LOCKED] from sqlite3_backup_step()
** is not a permanent error and does not affect the return value of
** sqlite3_backup_finish().
**
** [[sqlite3_backup_remaining()]] [[sqlite3_backup_pagecount()]]
** <b>sqlite3_backup_remaining() and sqlite3_backup_pagecount()</b>
**
** ^The sqlite3_backup_remaining() routine returns the number of pages still
** to be backed up at the conclusion of the most recent sqlite3_backup_step().
** ^The sqlite3_backup_pagecount() routine returns the total number of pages
** in the source database at the conclusion of the most recent
** sqlite3_backup_step().
** ^(The values returned by these functions are only updated by
** sqlite3_backup_step(). If the source database is modified in a way that
** changes the size of the source database or the number of pages remaining,
** those changes are not reflected in the output of sqlite3_backup_pagecount()
** and sqlite3_backup_remaining() until after the next
** sqlite3_backup_step().)^
**
** <b>Concurrent Usage of Database Handles</b>
**
** ^The source [database connection] may be used by the application for other
** purposes while a backup operation is underway or being initialized.
** ^If SQLite is compiled and configured to support threadsafe database
** connections, then the source database connection may be used concurrently
** from within other threads.
**
** However, the application must guarantee that the destination
** [database connection] is not passed to any other API (by any thread) after
** sqlite3_backup_init() is called and before the corresponding call to
** sqlite3_backup_finish().  SQLite does not currently check to see
** if the application incorrectly accesses the destination [database connection]
** and so no error code is reported, but the operations may malfunction
** nevertheless.  Use of the destination database connection while a
** backup is in progress might also cause a mutex deadlock.
**
** If running in [shared cache mode], the application must
** guarantee that the shared cache used by the destination database
** is not accessed while the backup is running. In practice this means
** that the application must guarantee that the disk file being
** backed up to is not accessed by any connection within the process,
** not just the specific connection that was passed to sqlite3_backup_init().
**
** The [sqlite3_backup] object itself is partially threadsafe. Multiple
** threads may safely make multiple concurrent calls to sqlite3_backup_step().
** However, the sqlite3_backup_remaining() and sqlite3_backup_pagecount()
** APIs are not strictly speaking threadsafe. If they are invoked at the
** same time as another thread is invoking sqlite3_backup_step() it is
** possible that they return invalid values.
**
** <b>Alternatives To Using The Backup API</b>
**
** Other techniques for safely creating a consistent backup of an SQLite
** database include:
**
** <ul>
** <li> The [VACUUM INTO] command.
** <li> The [sqlite3_rsync] utility program.
** </ul>
 */
// llgo:link (*Sqlite3).BackupInit C.sqlite3_backup_init
func (recv_ *Sqlite3) BackupInit(zDestName *c.Char, pSource *Sqlite3, zSourceName *c.Char) *Backup {
	return nil
}

// llgo:link (*Backup).BackupStep C.sqlite3_backup_step
func (recv_ *Backup) BackupStep(nPage c.Int) c.Int {
	return 0
}

// llgo:link (*Backup).BackupFinish C.sqlite3_backup_finish
func (recv_ *Backup) BackupFinish() c.Int {
	return 0
}

// llgo:link (*Backup).BackupRemaining C.sqlite3_backup_remaining
func (recv_ *Backup) BackupRemaining() c.Int {
	return 0
}

// llgo:link (*Backup).BackupPagecount C.sqlite3_backup_pagecount
func (recv_ *Backup) BackupPagecount() c.Int {
	return 0
}

/*
** CAPI3REF: Unlock Notification
** METHOD: sqlite3
**
** ^When running in shared-cache mode, a database operation may fail with
** an [SQLITE_LOCKED] error if the required locks on the shared-cache or
** individual tables within the shared-cache cannot be obtained. See
** [SQLite Shared-Cache Mode] for a description of shared-cache locking.
** ^This API may be used to register a callback that SQLite will invoke
** when the connection currently holding the required lock relinquishes it.
** ^This API is only available if the library was compiled with the
** [SQLITE_ENABLE_UNLOCK_NOTIFY] C-preprocessor symbol defined.
**
** See Also: [Using the SQLite Unlock Notification Feature].
**
** ^Shared-cache locks are released when a database connection concludes
** its current transaction, either by committing it or rolling it back.
**
** ^When a connection (known as the blocked connection) fails to obtain a
** shared-cache lock and SQLITE_LOCKED is returned to the caller, the
** identity of the database connection (the blocking connection) that
** has locked the required resource is stored internally. ^After an
** application receives an SQLITE_LOCKED error, it may call the
** sqlite3_unlock_notify() method with the blocked connection handle as
** the first argument to register for a callback that will be invoked
** when the blocking connections current transaction is concluded. ^The
** callback is invoked from within the [sqlite3_step] or [sqlite3_close]
** call that concludes the blocking connection's transaction.
**
** ^(If sqlite3_unlock_notify() is called in a multi-threaded application,
** there is a chance that the blocking connection will have already
** concluded its transaction by the time sqlite3_unlock_notify() is invoked.
** If this happens, then the specified callback is invoked immediately,
** from within the call to sqlite3_unlock_notify().)^
**
** ^If the blocked connection is attempting to obtain a write-lock on a
** shared-cache table, and more than one other connection currently holds
** a read-lock on the same table, then SQLite arbitrarily selects one of
** the other connections to use as the blocking connection.
**
** ^(There may be at most one unlock-notify callback registered by a
** blocked connection. If sqlite3_unlock_notify() is called when the
** blocked connection already has a registered unlock-notify callback,
** then the new callback replaces the old.)^ ^If sqlite3_unlock_notify() is
** called with a NULL pointer as its second argument, then any existing
** unlock-notify callback is canceled. ^The blocked connections
** unlock-notify callback may also be canceled by closing the blocked
** connection using [sqlite3_close()].
**
** The unlock-notify callback is not reentrant. If an application invokes
** any sqlite3_xxx API functions from within an unlock-notify callback, a
** crash or deadlock may be the result.
**
** ^Unless deadlock is detected (see below), sqlite3_unlock_notify() always
** returns SQLITE_OK.
**
** <b>Callback Invocation Details</b>
**
** When an unlock-notify callback is registered, the application provides a
** single void* pointer that is passed to the callback when it is invoked.
** However, the signature of the callback function allows SQLite to pass
** it an array of void* context pointers. The first argument passed to
** an unlock-notify callback is a pointer to an array of void* pointers,
** and the second is the number of entries in the array.
**
** When a blocking connection's transaction is concluded, there may be
** more than one blocked connection that has registered for an unlock-notify
** callback. ^If two or more such blocked connections have specified the
** same callback function, then instead of invoking the callback function
** multiple times, it is invoked once with the set of void* context pointers
** specified by the blocked connections bundled together into an array.
** This gives the application an opportunity to prioritize any actions
** related to the set of unblocked database connections.
**
** <b>Deadlock Detection</b>
**
** Assuming that after registering for an unlock-notify callback a
** database waits for the callback to be issued before taking any further
** action (a reasonable assumption), then using this API may cause the
** application to deadlock. For example, if connection X is waiting for
** connection Y's transaction to be concluded, and similarly connection
** Y is waiting on connection X's transaction, then neither connection
** will proceed and the system may remain deadlocked indefinitely.
**
** To avoid this scenario, the sqlite3_unlock_notify() performs deadlock
** detection. ^If a given call to sqlite3_unlock_notify() would put the
** system in a deadlocked state, then SQLITE_LOCKED is returned and no
** unlock-notify callback is registered. The system is said to be in
** a deadlocked state if connection A has registered for an unlock-notify
** callback on the conclusion of connection B's transaction, and connection
** B has itself registered for an unlock-notify callback when connection
** A's transaction is concluded. ^Indirect deadlock is also detected, so
** the system is also considered to be deadlocked if connection B has
** registered for an unlock-notify callback on the conclusion of connection
** C's transaction, where connection C is waiting on connection A. ^Any
** number of levels of indirection are allowed.
**
** <b>The "DROP TABLE" Exception</b>
**
** When a call to [sqlite3_step()] returns SQLITE_LOCKED, it is almost
** always appropriate to call sqlite3_unlock_notify(). There is however,
** one exception. When executing a "DROP TABLE" or "DROP INDEX" statement,
** SQLite checks if there are any currently executing SELECT statements
** that belong to the same connection. If there are, SQLITE_LOCKED is
** returned. In this case there is no "blocking connection", so invoking
** sqlite3_unlock_notify() results in the unlock-notify callback being
** invoked immediately. If the application then re-attempts the "DROP TABLE"
** or "DROP INDEX" query, an infinite loop might be the result.
**
** One way around this problem is to check the extended error code returned
** by an sqlite3_step() call. ^(If there is a blocking connection, then the
** extended error code is set to SQLITE_LOCKED_SHAREDCACHE. Otherwise, in
** the special "DROP TABLE/INDEX" case, the extended error code is just
** SQLITE_LOCKED.)^
 */
// llgo:link (*Sqlite3).UnlockNotify C.sqlite3_unlock_notify
func (recv_ *Sqlite3) UnlockNotify(xNotify func(*c.Pointer, c.Int), pNotifyArg c.Pointer) c.Int {
	return 0
}

/*
** CAPI3REF: String Comparison
**
** ^The [sqlite3_stricmp()] and [sqlite3_strnicmp()] APIs allow applications
** and extensions to compare the contents of two buffers containing UTF-8
** strings in a case-independent fashion, using the same definition of "case
** independence" that SQLite uses internally when comparing identifiers.
 */
//go:linkname Stricmp C.sqlite3_stricmp
func Stricmp(*c.Char, *c.Char) c.Int

//go:linkname Strnicmp C.sqlite3_strnicmp
func Strnicmp(*c.Char, *c.Char, c.Int) c.Int

/*
** CAPI3REF: String Globbing
*
** ^The [sqlite3_strglob(P,X)] interface returns zero if and only if
** string X matches the [GLOB] pattern P.
** ^The definition of [GLOB] pattern matching used in
** [sqlite3_strglob(P,X)] is the same as for the "X GLOB P" operator in the
** SQL dialect understood by SQLite.  ^The [sqlite3_strglob(P,X)] function
** is case sensitive.
**
** Note that this routine returns zero on a match and non-zero if the strings
** do not match, the same as [sqlite3_stricmp()] and [sqlite3_strnicmp()].
**
** See also: [sqlite3_strlike()].
 */
//go:linkname Strglob C.sqlite3_strglob
func Strglob(zGlob *c.Char, zStr *c.Char) c.Int

/*
** CAPI3REF: String LIKE Matching
*
** ^The [sqlite3_strlike(P,X,E)] interface returns zero if and only if
** string X matches the [LIKE] pattern P with escape character E.
** ^The definition of [LIKE] pattern matching used in
** [sqlite3_strlike(P,X,E)] is the same as for the "X LIKE P ESCAPE E"
** operator in the SQL dialect understood by SQLite.  ^For "X LIKE P" without
** the ESCAPE clause, set the E parameter of [sqlite3_strlike(P,X,E)] to 0.
** ^As with the LIKE operator, the [sqlite3_strlike(P,X,E)] function is case
** insensitive - equivalent upper and lower case ASCII characters match
** one another.
**
** ^The [sqlite3_strlike(P,X,E)] function matches Unicode characters, though
** only ASCII characters are case folded.
**
** Note that this routine returns zero on a match and non-zero if the strings
** do not match, the same as [sqlite3_stricmp()] and [sqlite3_strnicmp()].
**
** See also: [sqlite3_strglob()].
 */
//go:linkname Strlike C.sqlite3_strlike
func Strlike(zGlob *c.Char, zStr *c.Char, cEsc c.Uint) c.Int

/*
** CAPI3REF: Error Logging Interface
**
** ^The [sqlite3_log()] interface writes a message into the [error log]
** established by the [SQLITE_CONFIG_LOG] option to [sqlite3_config()].
** ^If logging is enabled, the zFormat string and subsequent arguments are
** used with [sqlite3_snprintf()] to generate the final output string.
**
** The sqlite3_log() interface is intended for use by extensions such as
** virtual tables, collating functions, and SQL functions.  While there is
** nothing to prevent an application from calling sqlite3_log(), doing so
** is considered bad form.
**
** The zFormat string must not be NULL.
**
** To avoid deadlocks and other threading problems, the sqlite3_log() routine
** will not use dynamically allocated memory.  The log message is stored in
** a fixed-length buffer on the stack.  If the log message is longer than
** a few hundred characters, it will be truncated to the length of the
** buffer.
 */
//go:linkname Log C.sqlite3_log
func Log(iErrCode c.Int, zFormat *c.Char, __llgo_va_list ...interface{})

/*
** CAPI3REF: Write-Ahead Log Commit Hook
** METHOD: sqlite3
**
** ^The [sqlite3_wal_hook()] function is used to register a callback that
** is invoked each time data is committed to a database in wal mode.
**
** ^(The callback is invoked by SQLite after the commit has taken place and
** the associated write-lock on the database released)^, so the implementation
** may read, write or [checkpoint] the database as required.
**
** ^The first parameter passed to the callback function when it is invoked
** is a copy of the third parameter passed to sqlite3_wal_hook() when
** registering the callback. ^The second is a copy of the database handle.
** ^The third parameter is the name of the database that was written to -
** either "main" or the name of an [ATTACH]-ed database. ^The fourth parameter
** is the number of pages currently in the write-ahead log file,
** including those that were just committed.
**
** The callback function should normally return [SQLITE_OK].  ^If an error
** code is returned, that error will propagate back up through the
** SQLite code base to cause the statement that provoked the callback
** to report an error, though the commit will have still occurred. If the
** callback returns [SQLITE_ROW] or [SQLITE_DONE], or if it returns a value
** that does not correspond to any valid SQLite error code, the results
** are undefined.
**
** A single database handle may have at most a single write-ahead log callback
** registered at one time. ^Calling [sqlite3_wal_hook()] replaces any
** previously registered write-ahead log callback. ^The return value is
** a copy of the third parameter from the previous call, if any, or 0.
** ^Note that the [sqlite3_wal_autocheckpoint()] interface and the
** [wal_autocheckpoint pragma] both invoke [sqlite3_wal_hook()] and will
** overwrite any prior [sqlite3_wal_hook()] settings.
 */
// llgo:link (*Sqlite3).WalHook C.sqlite3_wal_hook
func (recv_ *Sqlite3) WalHook(func(c.Pointer, *Sqlite3, *c.Char, c.Int) c.Int, c.Pointer) c.Pointer {
	return nil
}

/*
** CAPI3REF: Configure an auto-checkpoint
** METHOD: sqlite3
**
** ^The [sqlite3_wal_autocheckpoint(D,N)] is a wrapper around
** [sqlite3_wal_hook()] that causes any database on [database connection] D
** to automatically [checkpoint]
** after committing a transaction if there are N or
** more frames in the [write-ahead log] file.  ^Passing zero or
** a negative value as the nFrame parameter disables automatic
** checkpoints entirely.
**
** ^The callback registered by this function replaces any existing callback
** registered using [sqlite3_wal_hook()].  ^Likewise, registering a callback
** using [sqlite3_wal_hook()] disables the automatic checkpoint mechanism
** configured by this function.
**
** ^The [wal_autocheckpoint pragma] can be used to invoke this interface
** from SQL.
**
** ^Checkpoints initiated by this mechanism are
** [sqlite3_wal_checkpoint_v2|PASSIVE].
**
** ^Every new [database connection] defaults to having the auto-checkpoint
** enabled with a threshold of 1000 or [SQLITE_DEFAULT_WAL_AUTOCHECKPOINT]
** pages.  The use of this interface
** is only necessary if the default setting is found to be suboptimal
** for a particular application.
 */
// llgo:link (*Sqlite3).WalAutocheckpoint C.sqlite3_wal_autocheckpoint
func (recv_ *Sqlite3) WalAutocheckpoint(N c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Checkpoint a database
** METHOD: sqlite3
**
** ^(The sqlite3_wal_checkpoint(D,X) is equivalent to
** [sqlite3_wal_checkpoint_v2](D,X,[SQLITE_CHECKPOINT_PASSIVE],0,0).)^
**
** In brief, sqlite3_wal_checkpoint(D,X) causes the content in the
** [write-ahead log] for database X on [database connection] D to be
** transferred into the database file and for the write-ahead log to
** be reset.  See the [checkpointing] documentation for addition
** information.
**
** This interface used to be the only way to cause a checkpoint to
** occur.  But then the newer and more powerful [sqlite3_wal_checkpoint_v2()]
** interface was added.  This interface is retained for backwards
** compatibility and as a convenience for applications that need to manually
** start a callback but which do not need the full power (and corresponding
** complication) of [sqlite3_wal_checkpoint_v2()].
 */
// llgo:link (*Sqlite3).WalCheckpoint C.sqlite3_wal_checkpoint
func (recv_ *Sqlite3) WalCheckpoint(zDb *c.Char) c.Int {
	return 0
}

/*
** CAPI3REF: Checkpoint a database
** METHOD: sqlite3
**
** ^(The sqlite3_wal_checkpoint_v2(D,X,M,L,C) interface runs a checkpoint
** operation on database X of [database connection] D in mode M.  Status
** information is written back into integers pointed to by L and C.)^
** ^(The M parameter must be a valid [checkpoint mode]:)^
**
** <dl>
** <dt>SQLITE_CHECKPOINT_PASSIVE<dd>
**   ^Checkpoint as many frames as possible without waiting for any database
**   readers or writers to finish, then sync the database file if all frames
**   in the log were checkpointed. ^The [busy-handler callback]
**   is never invoked in the SQLITE_CHECKPOINT_PASSIVE mode.
**   ^On the other hand, passive mode might leave the checkpoint unfinished
**   if there are concurrent readers or writers.
**
** <dt>SQLITE_CHECKPOINT_FULL<dd>
**   ^This mode blocks (it invokes the
**   [sqlite3_busy_handler|busy-handler callback]) until there is no
**   database writer and all readers are reading from the most recent database
**   snapshot. ^It then checkpoints all frames in the log file and syncs the
**   database file. ^This mode blocks new database writers while it is pending,
**   but new database readers are allowed to continue unimpeded.
**
** <dt>SQLITE_CHECKPOINT_RESTART<dd>
**   ^This mode works the same way as SQLITE_CHECKPOINT_FULL with the addition
**   that after checkpointing the log file it blocks (calls the
**   [busy-handler callback])
**   until all readers are reading from the database file only. ^This ensures
**   that the next writer will restart the log file from the beginning.
**   ^Like SQLITE_CHECKPOINT_FULL, this mode blocks new
**   database writer attempts while it is pending, but does not impede readers.
**
** <dt>SQLITE_CHECKPOINT_TRUNCATE<dd>
**   ^This mode works the same way as SQLITE_CHECKPOINT_RESTART with the
**   addition that it also truncates the log file to zero bytes just prior
**   to a successful return.
** </dl>
**
** ^If pnLog is not NULL, then *pnLog is set to the total number of frames in
** the log file or to -1 if the checkpoint could not run because
** of an error or because the database is not in [WAL mode]. ^If pnCkpt is not
** NULL,then *pnCkpt is set to the total number of checkpointed frames in the
** log file (including any that were already checkpointed before the function
** was called) or to -1 if the checkpoint could not run due to an error or
** because the database is not in WAL mode. ^Note that upon successful
** completion of an SQLITE_CHECKPOINT_TRUNCATE, the log file will have been
** truncated to zero bytes and so both *pnLog and *pnCkpt will be set to zero.
**
** ^All calls obtain an exclusive "checkpoint" lock on the database file. ^If
** any other process is running a checkpoint operation at the same time, the
** lock cannot be obtained and SQLITE_BUSY is returned. ^Even if there is a
** busy-handler configured, it will not be invoked in this case.
**
** ^The SQLITE_CHECKPOINT_FULL, RESTART and TRUNCATE modes also obtain the
** exclusive "writer" lock on the database file. ^If the writer lock cannot be
** obtained immediately, and a busy-handler is configured, it is invoked and
** the writer lock retried until either the busy-handler returns 0 or the lock
** is successfully obtained. ^The busy-handler is also invoked while waiting for
** database readers as described above. ^If the busy-handler returns 0 before
** the writer lock is obtained or while waiting for database readers, the
** checkpoint operation proceeds from that point in the same way as
** SQLITE_CHECKPOINT_PASSIVE - checkpointing as many frames as possible
** without blocking any further. ^SQLITE_BUSY is returned in this case.
**
** ^If parameter zDb is NULL or points to a zero length string, then the
** specified operation is attempted on all WAL databases [attached] to
** [database connection] db.  In this case the
** values written to output parameters *pnLog and *pnCkpt are undefined. ^If
** an SQLITE_BUSY error is encountered when processing one or more of the
** attached WAL databases, the operation is still attempted on any remaining
** attached databases and SQLITE_BUSY is returned at the end. ^If any other
** error occurs while processing an attached database, processing is abandoned
** and the error code is returned to the caller immediately. ^If no error
** (SQLITE_BUSY or otherwise) is encountered while processing the attached
** databases, SQLITE_OK is returned.
**
** ^If database zDb is the name of an attached database that is not in WAL
** mode, SQLITE_OK is returned and both *pnLog and *pnCkpt set to -1. ^If
** zDb is not NULL (or a zero length string) and is not the name of any
** attached database, SQLITE_ERROR is returned to the caller.
**
** ^Unless it returns SQLITE_MISUSE,
** the sqlite3_wal_checkpoint_v2() interface
** sets the error information that is queried by
** [sqlite3_errcode()] and [sqlite3_errmsg()].
**
** ^The [PRAGMA wal_checkpoint] command can be used to invoke this interface
** from SQL.
 */
// llgo:link (*Sqlite3).WalCheckpointV2 C.sqlite3_wal_checkpoint_v2
func (recv_ *Sqlite3) WalCheckpointV2(zDb *c.Char, eMode c.Int, pnLog *c.Int, pnCkpt *c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Virtual Table Interface Configuration
**
** This function may be called by either the [xConnect] or [xCreate] method
** of a [virtual table] implementation to configure
** various facets of the virtual table interface.
**
** If this interface is invoked outside the context of an xConnect or
** xCreate virtual table method then the behavior is undefined.
**
** In the call sqlite3_vtab_config(D,C,...) the D parameter is the
** [database connection] in which the virtual table is being created and
** which is passed in as the first argument to the [xConnect] or [xCreate]
** method that is invoking sqlite3_vtab_config().  The C parameter is one
** of the [virtual table configuration options].  The presence and meaning
** of parameters after C depend on which [virtual table configuration option]
** is used.
 */
//go:linkname VtabConfig C.sqlite3_vtab_config
func VtabConfig(__llgo_arg_0 *Sqlite3, op c.Int, __llgo_va_list ...interface{}) c.Int

/*
** CAPI3REF: Determine The Virtual Table Conflict Policy
**
** This function may only be called from within a call to the [xUpdate] method
** of a [virtual table] implementation for an INSERT or UPDATE operation. ^The
** value returned is one of [SQLITE_ROLLBACK], [SQLITE_IGNORE], [SQLITE_FAIL],
** [SQLITE_ABORT], or [SQLITE_REPLACE], according to the [ON CONFLICT] mode
** of the SQL statement that triggered the call to the [xUpdate] method of the
** [virtual table].
 */
// llgo:link (*Sqlite3).VtabOnConflict C.sqlite3_vtab_on_conflict
func (recv_ *Sqlite3) VtabOnConflict() c.Int {
	return 0
}

/*
** CAPI3REF: Determine If Virtual Table Column Access Is For UPDATE
**
** If the sqlite3_vtab_nochange(X) routine is called within the [xColumn]
** method of a [virtual table], then it might return true if the
** column is being fetched as part of an UPDATE operation during which the
** column value will not change.  The virtual table implementation can use
** this hint as permission to substitute a return value that is less
** expensive to compute and that the corresponding
** [xUpdate] method understands as a "no-change" value.
**
** If the [xColumn] method calls sqlite3_vtab_nochange() and finds that
** the column is not changed by the UPDATE statement, then the xColumn
** method can optionally return without setting a result, without calling
** any of the [sqlite3_result_int|sqlite3_result_xxxxx() interfaces].
** In that case, [sqlite3_value_nochange(X)] will return true for the
** same column in the [xUpdate] method.
**
** The sqlite3_vtab_nochange() routine is an optimization.  Virtual table
** implementations should continue to give a correct answer even if the
** sqlite3_vtab_nochange() interface were to always return false.  In the
** current implementation, the sqlite3_vtab_nochange() interface does always
** returns false for the enhanced [UPDATE FROM] statement.
 */
// llgo:link (*Context).VtabNochange C.sqlite3_vtab_nochange
func (recv_ *Context) VtabNochange() c.Int {
	return 0
}

/*
** CAPI3REF: Determine The Collation For a Virtual Table Constraint
** METHOD: sqlite3_index_info
**
** This function may only be called from within a call to the [xBestIndex]
** method of a [virtual table].  This function returns a pointer to a string
** that is the name of the appropriate collation sequence to use for text
** comparisons on the constraint identified by its arguments.
**
** The first argument must be the pointer to the [sqlite3_index_info] object
** that is the first parameter to the xBestIndex() method. The second argument
** must be an index into the aConstraint[] array belonging to the
** sqlite3_index_info structure passed to xBestIndex.
**
** Important:
** The first parameter must be the same pointer that is passed into the
** xBestMethod() method.  The first parameter may not be a pointer to a
** different [sqlite3_index_info] object, even an exact copy.
**
** The return value is computed as follows:
**
** <ol>
** <li><p> If the constraint comes from a WHERE clause expression that contains
**         a [COLLATE operator], then the name of the collation specified by
**         that COLLATE operator is returned.
** <li><p> If there is no COLLATE operator, but the column that is the subject
**         of the constraint specifies an alternative collating sequence via
**         a [COLLATE clause] on the column definition within the CREATE TABLE
**         statement that was passed into [sqlite3_declare_vtab()], then the
**         name of that alternative collating sequence is returned.
** <li><p> Otherwise, "BINARY" is returned.
** </ol>
 */
// llgo:link (*IndexInfo).VtabCollation C.sqlite3_vtab_collation
func (recv_ *IndexInfo) VtabCollation(c.Int) *c.Char {
	return nil
}

/*
** CAPI3REF: Determine if a virtual table query is DISTINCT
** METHOD: sqlite3_index_info
**
** This API may only be used from within an [xBestIndex|xBestIndex method]
** of a [virtual table] implementation. The result of calling this
** interface from outside of xBestIndex() is undefined and probably harmful.
**
** ^The sqlite3_vtab_distinct() interface returns an integer between 0 and
** 3.  The integer returned by sqlite3_vtab_distinct()
** gives the virtual table additional information about how the query
** planner wants the output to be ordered. As long as the virtual table
** can meet the ordering requirements of the query planner, it may set
** the "orderByConsumed" flag.
**
** <ol><li value="0"><p>
** ^If the sqlite3_vtab_distinct() interface returns 0, that means
** that the query planner needs the virtual table to return all rows in the
** sort order defined by the "nOrderBy" and "aOrderBy" fields of the
** [sqlite3_index_info] object.  This is the default expectation.  If the
** virtual table outputs all rows in sorted order, then it is always safe for
** the xBestIndex method to set the "orderByConsumed" flag, regardless of
** the return value from sqlite3_vtab_distinct().
** <li value="1"><p>
** ^(If the sqlite3_vtab_distinct() interface returns 1, that means
** that the query planner does not need the rows to be returned in sorted order
** as long as all rows with the same values in all columns identified by the
** "aOrderBy" field are adjacent.)^  This mode is used when the query planner
** is doing a GROUP BY.
** <li value="2"><p>
** ^(If the sqlite3_vtab_distinct() interface returns 2, that means
** that the query planner does not need the rows returned in any particular
** order, as long as rows with the same values in all columns identified
** by "aOrderBy" are adjacent.)^  ^(Furthermore, when two or more rows
** contain the same values for all columns identified by "colUsed", all but
** one such row may optionally be omitted from the result.)^
** The virtual table is not required to omit rows that are duplicates
** over the "colUsed" columns, but if the virtual table can do that without
** too much extra effort, it could potentially help the query to run faster.
** This mode is used for a DISTINCT query.
** <li value="3"><p>
** ^(If the sqlite3_vtab_distinct() interface returns 3, that means the
** virtual table must return rows in the order defined by "aOrderBy" as
** if the sqlite3_vtab_distinct() interface had returned 0.  However if
** two or more rows in the result have the same values for all columns
** identified by "colUsed", then all but one such row may optionally be
** omitted.)^  Like when the return value is 2, the virtual table
** is not required to omit rows that are duplicates over the "colUsed"
** columns, but if the virtual table can do that without
** too much extra effort, it could potentially help the query to run faster.
** This mode is used for queries
** that have both DISTINCT and ORDER BY clauses.
** </ol>
**
** <p>The following table summarizes the conditions under which the
** virtual table is allowed to set the "orderByConsumed" flag based on
** the value returned by sqlite3_vtab_distinct().  This table is a
** restatement of the previous four paragraphs:
**
** <table border=1 cellspacing=0 cellpadding=10 width="90%">
** <tr>
** <td valign="top">sqlite3_vtab_distinct() return value
** <td valign="top">Rows are returned in aOrderBy order
** <td valign="top">Rows with the same value in all aOrderBy columns are adjacent
** <td valign="top">Duplicates over all colUsed columns may be omitted
** <tr><td>0<td>yes<td>yes<td>no
** <tr><td>1<td>no<td>yes<td>no
** <tr><td>2<td>no<td>yes<td>yes
** <tr><td>3<td>yes<td>yes<td>yes
** </table>
**
** ^For the purposes of comparing virtual table output values to see if the
** values are same value for sorting purposes, two NULL values are considered
** to be the same.  In other words, the comparison operator is "IS"
** (or "IS NOT DISTINCT FROM") and not "==".
**
** If a virtual table implementation is unable to meet the requirements
** specified above, then it must not set the "orderByConsumed" flag in the
** [sqlite3_index_info] object or an incorrect answer may result.
**
** ^A virtual table implementation is always free to return rows in any order
** it wants, as long as the "orderByConsumed" flag is not set.  ^When the
** the "orderByConsumed" flag is unset, the query planner will add extra
** [bytecode] to ensure that the final results returned by the SQL query are
** ordered correctly.  The use of the "orderByConsumed" flag and the
** sqlite3_vtab_distinct() interface is merely an optimization.  ^Careful
** use of the sqlite3_vtab_distinct() interface and the "orderByConsumed"
** flag might help queries against a virtual table to run faster.  Being
** overly aggressive and setting the "orderByConsumed" flag when it is not
** valid to do so, on the other hand, might cause SQLite to return incorrect
** results.
 */
// llgo:link (*IndexInfo).VtabDistinct C.sqlite3_vtab_distinct
func (recv_ *IndexInfo) VtabDistinct() c.Int {
	return 0
}

/*
** CAPI3REF: Identify and handle IN constraints in xBestIndex
**
** This interface may only be used from within an
** [xBestIndex|xBestIndex() method] of a [virtual table] implementation.
** The result of invoking this interface from any other context is
** undefined and probably harmful.
**
** ^(A constraint on a virtual table of the form
** "[IN operator|column IN (...)]" is
** communicated to the xBestIndex method as a
** [SQLITE_INDEX_CONSTRAINT_EQ] constraint.)^  If xBestIndex wants to use
** this constraint, it must set the corresponding
** aConstraintUsage[].argvIndex to a positive integer.  ^(Then, under
** the usual mode of handling IN operators, SQLite generates [bytecode]
** that invokes the [xFilter|xFilter() method] once for each value
** on the right-hand side of the IN operator.)^  Thus the virtual table
** only sees a single value from the right-hand side of the IN operator
** at a time.
**
** In some cases, however, it would be advantageous for the virtual
** table to see all values on the right-hand of the IN operator all at
** once.  The sqlite3_vtab_in() interfaces facilitates this in two ways:
**
** <ol>
** <li><p>
**   ^A call to sqlite3_vtab_in(P,N,-1) will return true (non-zero)
**   if and only if the [sqlite3_index_info|P->aConstraint][N] constraint
**   is an [IN operator] that can be processed all at once.  ^In other words,
**   sqlite3_vtab_in() with -1 in the third argument is a mechanism
**   by which the virtual table can ask SQLite if all-at-once processing
**   of the IN operator is even possible.
**
** <li><p>
**   ^A call to sqlite3_vtab_in(P,N,F) with F==1 or F==0 indicates
**   to SQLite that the virtual table does or does not want to process
**   the IN operator all-at-once, respectively.  ^Thus when the third
**   parameter (F) is non-negative, this interface is the mechanism by
**   which the virtual table tells SQLite how it wants to process the
**   IN operator.
** </ol>
**
** ^The sqlite3_vtab_in(P,N,F) interface can be invoked multiple times
** within the same xBestIndex method call.  ^For any given P,N pair,
** the return value from sqlite3_vtab_in(P,N,F) will always be the same
** within the same xBestIndex call.  ^If the interface returns true
** (non-zero), that means that the constraint is an IN operator
** that can be processed all-at-once.  ^If the constraint is not an IN
** operator or cannot be processed all-at-once, then the interface returns
** false.
**
** ^(All-at-once processing of the IN operator is selected if both of the
** following conditions are met:
**
** <ol>
** <li><p> The P->aConstraintUsage[N].argvIndex value is set to a positive
** integer.  This is how the virtual table tells SQLite that it wants to
** use the N-th constraint.
**
** <li><p> The last call to sqlite3_vtab_in(P,N,F) for which F was
** non-negative had F>=1.
** </ol>)^
**
** ^If either or both of the conditions above are false, then SQLite uses
** the traditional one-at-a-time processing strategy for the IN constraint.
** ^If both conditions are true, then the argvIndex-th parameter to the
** xFilter method will be an [sqlite3_value] that appears to be NULL,
** but which can be passed to [sqlite3_vtab_in_first()] and
** [sqlite3_vtab_in_next()] to find all values on the right-hand side
** of the IN constraint.
 */
// llgo:link (*IndexInfo).VtabIn C.sqlite3_vtab_in
func (recv_ *IndexInfo) VtabIn(iCons c.Int, bHandle c.Int) c.Int {
	return 0
}

/*
** CAPI3REF: Find all elements on the right-hand side of an IN constraint.
**
** These interfaces are only useful from within the
** [xFilter|xFilter() method] of a [virtual table] implementation.
** The result of invoking these interfaces from any other context
** is undefined and probably harmful.
**
** The X parameter in a call to sqlite3_vtab_in_first(X,P) or
** sqlite3_vtab_in_next(X,P) should be one of the parameters to the
** xFilter method which invokes these routines, and specifically
** a parameter that was previously selected for all-at-once IN constraint
** processing use the [sqlite3_vtab_in()] interface in the
** [xBestIndex|xBestIndex method].  ^(If the X parameter is not
** an xFilter argument that was selected for all-at-once IN constraint
** processing, then these routines return [SQLITE_ERROR].)^
**
** ^(Use these routines to access all values on the right-hand side
** of the IN constraint using code like the following:
**
** <blockquote><pre>
** &nbsp;  for(rc=sqlite3_vtab_in_first(pList, &pVal);
** &nbsp;      rc==SQLITE_OK && pVal;
** &nbsp;      rc=sqlite3_vtab_in_next(pList, &pVal)
** &nbsp;  ){
** &nbsp;    // do something with pVal
** &nbsp;  }
** &nbsp;  if( rc!=SQLITE_OK ){
** &nbsp;    // an error has occurred
** &nbsp;  }
** </pre></blockquote>)^
**
** ^On success, the sqlite3_vtab_in_first(X,P) and sqlite3_vtab_in_next(X,P)
** routines return SQLITE_OK and set *P to point to the first or next value
** on the RHS of the IN constraint.  ^If there are no more values on the
** right hand side of the IN constraint, then *P is set to NULL and these
** routines return [SQLITE_DONE].  ^The return value might be
** some other value, such as SQLITE_NOMEM, in the event of a malfunction.
**
** The *ppOut values returned by these routines are only valid until the
** next call to either of these routines or until the end of the xFilter
** method from which these routines were called.  If the virtual table
** implementation needs to retain the *ppOut values for longer, it must make
** copies.  The *ppOut values are [protected sqlite3_value|protected].
 */
// llgo:link (*Value).VtabInFirst C.sqlite3_vtab_in_first
func (recv_ *Value) VtabInFirst(ppOut **Value) c.Int {
	return 0
}

// llgo:link (*Value).VtabInNext C.sqlite3_vtab_in_next
func (recv_ *Value) VtabInNext(ppOut **Value) c.Int {
	return 0
}

/*
** CAPI3REF: Constraint values in xBestIndex()
** METHOD: sqlite3_index_info
**
** This API may only be used from within the [xBestIndex|xBestIndex method]
** of a [virtual table] implementation. The result of calling this interface
** from outside of an xBestIndex method are undefined and probably harmful.
**
** ^When the sqlite3_vtab_rhs_value(P,J,V) interface is invoked from within
** the [xBestIndex] method of a [virtual table] implementation, with P being
** a copy of the [sqlite3_index_info] object pointer passed into xBestIndex and
** J being a 0-based index into P->aConstraint[], then this routine
** attempts to set *V to the value of the right-hand operand of
** that constraint if the right-hand operand is known.  ^If the
** right-hand operand is not known, then *V is set to a NULL pointer.
** ^The sqlite3_vtab_rhs_value(P,J,V) interface returns SQLITE_OK if
** and only if *V is set to a value.  ^The sqlite3_vtab_rhs_value(P,J,V)
** inteface returns SQLITE_NOTFOUND if the right-hand side of the J-th
** constraint is not available.  ^The sqlite3_vtab_rhs_value() interface
** can return an result code other than SQLITE_OK or SQLITE_NOTFOUND if
** something goes wrong.
**
** The sqlite3_vtab_rhs_value() interface is usually only successful if
** the right-hand operand of a constraint is a literal value in the original
** SQL statement.  If the right-hand operand is an expression or a reference
** to some other column or a [host parameter], then sqlite3_vtab_rhs_value()
** will probably return [SQLITE_NOTFOUND].
**
** ^(Some constraints, such as [SQLITE_INDEX_CONSTRAINT_ISNULL] and
** [SQLITE_INDEX_CONSTRAINT_ISNOTNULL], have no right-hand operand.  For such
** constraints, sqlite3_vtab_rhs_value() always returns SQLITE_NOTFOUND.)^
**
** ^The [sqlite3_value] object returned in *V is a protected sqlite3_value
** and remains valid for the duration of the xBestIndex method call.
** ^When xBestIndex returns, the sqlite3_value object returned by
** sqlite3_vtab_rhs_value() is automatically deallocated.
**
** The "_rhs_" in the name of this routine is an abbreviation for
** "Right-Hand Side".
 */
// llgo:link (*IndexInfo).VtabRhsValue C.sqlite3_vtab_rhs_value
func (recv_ *IndexInfo) VtabRhsValue(__llgo_arg_0 c.Int, ppVal **Value) c.Int {
	return 0
}

/*
** CAPI3REF: Flush caches to disk mid-transaction
** METHOD: sqlite3
**
** ^If a write-transaction is open on [database connection] D when the
** [sqlite3_db_cacheflush(D)] interface invoked, any dirty
** pages in the pager-cache that are not currently in use are written out
** to disk. A dirty page may be in use if a database cursor created by an
** active SQL statement is reading from it, or if it is page 1 of a database
** file (page 1 is always "in use").  ^The [sqlite3_db_cacheflush(D)]
** interface flushes caches for all schemas - "main", "temp", and
** any [attached] databases.
**
** ^If this function needs to obtain extra database locks before dirty pages
** can be flushed to disk, it does so. ^If those locks cannot be obtained
** immediately and there is a busy-handler callback configured, it is invoked
** in the usual manner. ^If the required lock still cannot be obtained, then
** the database is skipped and an attempt made to flush any dirty pages
** belonging to the next (if any) database. ^If any databases are skipped
** because locks cannot be obtained, but no other error occurs, this
** function returns SQLITE_BUSY.
**
** ^If any other error occurs while flushing dirty pages to disk (for
** example an IO error or out-of-memory condition), then processing is
** abandoned and an SQLite [error code] is returned to the caller immediately.
**
** ^Otherwise, if no error occurs, [sqlite3_db_cacheflush()] returns SQLITE_OK.
**
** ^This function does not set the database handle error code or message
** returned by the [sqlite3_errcode()] and [sqlite3_errmsg()] functions.
 */
// llgo:link (*Sqlite3).DbCacheflush C.sqlite3_db_cacheflush
func (recv_ *Sqlite3) DbCacheflush() c.Int {
	return 0
}

/*
** CAPI3REF: Low-level system error code
** METHOD: sqlite3
**
** ^Attempt to return the underlying operating system error code or error
** number that caused the most recent I/O error or failure to open a file.
** The return value is OS-dependent.  For example, on unix systems, after
** [sqlite3_open_v2()] returns [SQLITE_CANTOPEN], this interface could be
** called to get back the underlying "errno" that caused the problem, such
** as ENOSPC, EAUTH, EISDIR, and so forth.
 */
// llgo:link (*Sqlite3).SystemErrno C.sqlite3_system_errno
func (recv_ *Sqlite3) SystemErrno() c.Int {
	return 0
}

/*
** CAPI3REF: Database Snapshot
** KEYWORDS: {snapshot} {sqlite3_snapshot}
**
** An instance of the snapshot object records the state of a [WAL mode]
** database for some specific point in history.
**
** In [WAL mode], multiple [database connections] that are open on the
** same database file can each be reading a different historical version
** of the database file.  When a [database connection] begins a read
** transaction, that connection sees an unchanging copy of the database
** as it existed for the point in time when the transaction first started.
** Subsequent changes to the database from other connections are not seen
** by the reader until a new read transaction is started.
**
** The sqlite3_snapshot object records state information about an historical
** version of the database file so that it is possible to later open a new read
** transaction that sees that historical version of the database rather than
** the most recent version.
 */
type Snapshot struct {
	Hidden [48]c.Char
}

/*
** CAPI3REF: Serialize a database
**
** The sqlite3_serialize(D,S,P,F) interface returns a pointer to
** memory that is a serialization of the S database on
** [database connection] D.  If S is a NULL pointer, the main database is used.
** If P is not a NULL pointer, then the size of the database in bytes
** is written into *P.
**
** For an ordinary on-disk database file, the serialization is just a
** copy of the disk file.  For an in-memory database or a "TEMP" database,
** the serialization is the same sequence of bytes which would be written
** to disk if that database where backed up to disk.
**
** The usual case is that sqlite3_serialize() copies the serialization of
** the database into memory obtained from [sqlite3_malloc64()] and returns
** a pointer to that memory.  The caller is responsible for freeing the
** returned value to avoid a memory leak.  However, if the F argument
** contains the SQLITE_SERIALIZE_NOCOPY bit, then no memory allocations
** are made, and the sqlite3_serialize() function will return a pointer
** to the contiguous memory representation of the database that SQLite
** is currently using for that database, or NULL if the no such contiguous
** memory representation of the database exists.  A contiguous memory
** representation of the database will usually only exist if there has
** been a prior call to [sqlite3_deserialize(D,S,...)] with the same
** values of D and S.
** The size of the database is written into *P even if the
** SQLITE_SERIALIZE_NOCOPY bit is set but no contiguous copy
** of the database exists.
**
** After the call, if the SQLITE_SERIALIZE_NOCOPY bit had been set,
** the returned buffer content will remain accessible and unchanged
** until either the next write operation on the connection or when
** the connection is closed, and applications must not modify the
** buffer. If the bit had been clear, the returned buffer will not
** be accessed by SQLite after the call.
**
** A call to sqlite3_serialize(D,S,P,F) might return NULL even if the
** SQLITE_SERIALIZE_NOCOPY bit is omitted from argument F if a memory
** allocation error occurs.
**
** This interface is omitted if SQLite is compiled with the
** [SQLITE_OMIT_DESERIALIZE] option.
 */
// llgo:link (*Sqlite3).Serialize C.sqlite3_serialize
func (recv_ *Sqlite3) Serialize(zSchema *c.Char, piSize *Int64, mFlags c.Uint) *c.Char {
	return nil
}

/*
** CAPI3REF: Deserialize a database
**
** The sqlite3_deserialize(D,S,P,N,M,F) interface causes the
** [database connection] D to disconnect from database S and then
** reopen S as an in-memory database based on the serialization contained
** in P.  The serialized database P is N bytes in size.  M is the size of
** the buffer P, which might be larger than N.  If M is larger than N, and
** the SQLITE_DESERIALIZE_READONLY bit is not set in F, then SQLite is
** permitted to add content to the in-memory database as long as the total
** size does not exceed M bytes.
**
** If the SQLITE_DESERIALIZE_FREEONCLOSE bit is set in F, then SQLite will
** invoke sqlite3_free() on the serialization buffer when the database
** connection closes.  If the SQLITE_DESERIALIZE_RESIZEABLE bit is set, then
** SQLite will try to increase the buffer size using sqlite3_realloc64()
** if writes on the database cause it to grow larger than M bytes.
**
** Applications must not modify the buffer P or invalidate it before
** the database connection D is closed.
**
** The sqlite3_deserialize() interface will fail with SQLITE_BUSY if the
** database is currently in a read transaction or is involved in a backup
** operation.
**
** It is not possible to deserialized into the TEMP database.  If the
** S argument to sqlite3_deserialize(D,S,P,N,M,F) is "temp" then the
** function returns SQLITE_ERROR.
**
** The deserialized database should not be in [WAL mode].  If the database
** is in WAL mode, then any attempt to use the database file will result
** in an [SQLITE_CANTOPEN] error.  The application can set the
** [file format version numbers] (bytes 18 and 19) of the input database P
** to 0x01 prior to invoking sqlite3_deserialize(D,S,P,N,M,F) to force the
** database file into rollback mode and work around this limitation.
**
** If sqlite3_deserialize(D,S,P,N,M,F) fails for any reason and if the
** SQLITE_DESERIALIZE_FREEONCLOSE bit is set in argument F, then
** [sqlite3_free()] is invoked on argument P prior to returning.
**
** This interface is omitted if SQLite is compiled with the
** [SQLITE_OMIT_DESERIALIZE] option.
 */
// llgo:link (*Sqlite3).Deserialize C.sqlite3_deserialize
func (recv_ *Sqlite3) Deserialize(zSchema *c.Char, pData *c.Char, szDb Int64, szBuf Int64, mFlags c.Uint) c.Int {
	return 0
}

type RtreeGeometry struct {
	PContext c.Pointer
	NParam   c.Int
	AParam   *RtreeDbl
	PUser    c.Pointer
	XDelUser c.Pointer
}

type RtreeQueryInfo struct {
	PContext      c.Pointer
	NParam        c.Int
	AParam        *RtreeDbl
	PUser         c.Pointer
	XDelUser      c.Pointer
	ACoord        *RtreeDbl
	AnQueue       *c.Uint
	NCoord        c.Int
	ILevel        c.Int
	MxLevel       c.Int
	IRowid        Int64
	RParentScore  RtreeDbl
	EParentWithin c.Int
	EWithin       c.Int
	RScore        RtreeDbl
	ApSqlParam    **Value
}
type RtreeDbl c.Double

/*
** Register a geometry callback named zGeom that can be used as part of an
** R-Tree geometry query as follows:
**
**   SELECT ... FROM <rtree> WHERE <rtree col> MATCH $zGeom(... params ...)
 */
// llgo:link (*Sqlite3).RtreeGeometryCallback C.sqlite3_rtree_geometry_callback
func (recv_ *Sqlite3) RtreeGeometryCallback(zGeom *c.Char, xGeom func(*RtreeGeometry, c.Int, *RtreeDbl, *c.Int) c.Int, pContext c.Pointer) c.Int {
	return 0
}

/*
** Register a 2nd-generation geometry callback named zScore that can be
** used as part of an R-Tree geometry query as follows:
**
**   SELECT ... FROM <rtree> WHERE <rtree col> MATCH $zQueryFunc(... params ...)
 */
// llgo:link (*Sqlite3).RtreeQueryCallback C.sqlite3_rtree_query_callback
func (recv_ *Sqlite3) RtreeQueryCallback(zQueryFunc *c.Char, xQueryFunc func(*RtreeQueryInfo) c.Int, pContext c.Pointer, xDestructor func(c.Pointer)) c.Int {
	return 0
}

type Fts5ExtensionApi struct {
	IVersion           c.Int
	XUserData          c.Pointer
	XColumnCount       c.Pointer
	XRowCount          c.Pointer
	XColumnTotalSize   c.Pointer
	XTokenize          c.Pointer
	XPhraseCount       c.Pointer
	XPhraseSize        c.Pointer
	XInstCount         c.Pointer
	XInst              c.Pointer
	XRowid             c.Pointer
	XColumnText        c.Pointer
	XColumnSize        c.Pointer
	XQueryPhrase       c.Pointer
	XSetAuxdata        c.Pointer
	XGetAuxdata        c.Pointer
	XPhraseFirst       c.Pointer
	XPhraseNext        c.Pointer
	XPhraseFirstColumn c.Pointer
	XPhraseNextColumn  c.Pointer
	XQueryToken        c.Pointer
	XInstToken         c.Pointer
	XColumnLocale      c.Pointer
	XTokenizeV2        c.Pointer
}

type Fts5Context struct {
	Unused [8]uint8
}

type Fts5PhraseIter struct {
	A *c.Char
	B *c.Char
}

// llgo:type C
type Fts5ExtensionFunction func(*Fts5ExtensionApi, *Fts5Context, *Context, c.Int, **Value)

type Fts5Tokenizer struct {
	Unused [8]uint8
}

type Fts5TokenizerV2 struct {
	IVersion  c.Int
	XCreate   c.Pointer
	XDelete   c.Pointer
	XTokenize c.Pointer
}

type Fts5Tokenizer__1 struct {
	XCreate   c.Pointer
	XDelete   c.Pointer
	XTokenize c.Pointer
}

type Fts5Api struct {
	IVersion           c.Int
	XCreateTokenizer   c.Pointer
	XFindTokenizer     c.Pointer
	XCreateFunction    c.Pointer
	XCreateTokenizerV2 c.Pointer
	XFindTokenizerV2   c.Pointer
}
