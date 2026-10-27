typedef struct gzFile_s *gzFile;

struct gzFile_s {
    unsigned have;
    unsigned char *next;
    long pos;
};

typedef struct z_stream_s {
    unsigned char *next_in;
    unsigned avail_in;
    unsigned long total_in;
} z_stream;
typedef z_stream *z_streamp;

// An incomplete (opaque) struct: its layout is unknown, so sizeof over it must
// not be detected (rule 3).
typedef struct opaque_s opaque_t;

int openOpaque_(opaque_t *o, int size);

#define ZLIB_VERSION "1.3.1"

#define Z_DEFAULT_LEVEL (-1)

/// Read one byte, falling back to the real function.
int gzgetc_(gzFile file);

int deflateInit_(z_streamp strm, int level, const char *version, int stream_size);

int inflateInit_(z_streamp strm, const char *version, int stream_size);

/// Initialize for deflation with extra options.
int inflateBackInit_(z_streamp strm, int windowBits, unsigned char *window, const char *version, int stream_size);

void reset_(z_streamp strm, int flags);

/// Return the library version.
int version_(void);

// A package-level callee: its first parameter is a scalar, so it is not
// converted to a method. Used to cover the package-function call path and
// numeric/constant-macro arguments.
int configure_(int mode, int level);

// gzgetc has real logic, so it stays an explicit MacroFunc entry (C wrapper).
#define z_gzgetc(g) \
        ((g)->have ? ((g)->have--, (g)->pos++, *((g)->next)++) : (gzgetc_)(g))

/// Initialize for deflation, filling in the version and stream size.
#define z_deflateInit(strm, level) \
        deflateInit_((strm), (level), ZLIB_VERSION, sizeof(z_stream))

#define z_inflateInit(strm) \
        inflateInit_((strm), ZLIB_VERSION, (int)sizeof(z_stream))

#define z_inflateBackInit(strm, windowBits, window) \
        inflateBackInit_((strm), (windowBits), (window), \
                         ZLIB_VERSION, (int)sizeof(z_stream))

// void-returning forwarder: the body is a bare call with no result.
#define z_reset(strm) \
        reset_((strm), 0)

// No-argument forwarder
#define z_version() version_()

// Package-level callee with a numeric literal and an object-like constant macro.
#define z_configure(mode) \
        (configure_((mode), Z_DEFAULT_LEVEL))

// Skipped: sizeof over an incomplete type whose layout Go cannot know (rule 3).
#define z_openOpaque(o) \
        openOpaque_((o), (int)sizeof(opaque_t))

// Skipped: parameter "strm" is used twice (rule 4).
#define z_resetTwice(strm) \
        reset_((strm), (strm))

// Skipped: body is not a single forwarding call (conditional expression).
#define z_pick(strm, flags) \
        ((flags) ? reset_((strm), (flags)) : (void)0)

