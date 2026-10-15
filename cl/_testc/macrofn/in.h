typedef struct gzFile_s *gzFile;

struct gzFile_s {
    unsigned have;
    unsigned char *next;
    long pos;
};

typedef struct z_stream_s z_stream;
typedef z_stream *z_streamp;

/// Read one byte, falling back to the real function.
int gzgetc_(gzFile file);

int deflateInit_(z_streamp strm, int level, const char *version, int stream_size);

void reset_(z_streamp strm, int flags, const char *version);

#define z_gzgetc(g) \
        ((g)->have ? ((g)->have--, (g)->pos++, *((g)->next)++) : (gzgetc_)(g))
// Read one byte from the gzip stream (macro fast path).

/// Initialize for deflation, filling in the version and stream size.
#define z_deflateInit(strm, level) \
        deflateInit_((strm), (level), "1.3", (int)sizeof(z_stream))

#define z_reset(strm, flags) \
        reset_((strm), (flags), "1.3")
// Reset the stream, restoring the version string.
