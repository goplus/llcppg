#define ZLIB_VERSION "1.3.2"

#define Z_NO_FLUSH      0
#define Z_PARTIAL_FLUSH 1
#define Z_SYNC_FLUSH    2
#define Z_FULL_FLUSH    3
#define Z_FINISH        4
#define Z_BLOCK         5
#define Z_TREES         6
/* Allowed flush values; see deflate() and inflate() below for details */

#define Z_OK            0
#define Z_STREAM_END    1
#define Z_NEED_DICT     2
#define Z_ERRNO        (-1)
#define Z_STREAM_ERROR (-2)
#define Z_DATA_ERROR   (-3)
#define Z_MEM_ERROR    (-4)
#define Z_BUF_ERROR    (-5)
#define Z_VERSION_ERROR (-6)
/* Return codes for the compression/decompression functions. Negative values
 * are errors, positive values are used for special but normal events.
 */

#define Z_NO_COMPRESSION         0
#define Z_BEST_SPEED             1
#define Z_BEST_COMPRESSION       9
#define Z_DEFAULT_COMPRESSION  (-1)
/* compression levels */

#define Z_FILTERED            1
#define Z_HUFFMAN_ONLY        2
#define Z_RLE                 3
#define Z_FIXED               4
#define Z_DEFAULT_STRATEGY    0
/* compression strategy; see deflateInit2() below for details */

#define Z_BINARY   0
#define Z_TEXT     1
#define Z_ASCII    Z_TEXT   /* for compatibility with 1.2.2 and earlier */
#define Z_UNKNOWN  2
/* Possible values of the data_type field for deflate() */

#define Z_DEFLATED   8
/* The deflate compression method (the only one supported in this version) */

#define Z_NULL  0  /* for initializing zalloc, zfree, opaque */

#define zlib_version zlibVersion()
/* for compatibility with versions < 1.0.2 */

#define ZEXTERN
#define ZEXPORT
#define FAR

typedef unsigned long uLong;
typedef unsigned long z_size_t;
typedef long long z_off64_t;

                        /* basic functions */

ZEXTERN const char * ZEXPORT zlibVersion(void);
/* The application can compare zlibVersion and ZLIB_VERSION for consistency.
   If the first character differs, the library code actually used is not
   compatible with the zlib.h header file used by the application.  This check
   is automatically made by deflateInit and inflateInit.
 */

typedef struct z_stream_s {
    int     data_type;  /* best guess about the data type: binary or text
                           for deflate, or the decoding state for inflate */
} z_stream;

typedef z_stream FAR *z_streamp;

ZEXTERN int ZEXPORT deflateCopy(z_streamp dest,
                                z_streamp source);
/*
     Sets the destination stream as a complete copy of the source stream.

     This function can be useful when several compression strategies will be
   tried, for example when there are several ways of pre-processing the input
   data with a filter.  The streams that will be discarded should then be freed
   by calling deflateEnd.  Note that deflateCopy duplicates the internal
   compression state which can be quite large, so this strategy is slow and can
   consume lots of memory.

     deflateCopy returns Z_OK if success, Z_MEM_ERROR if there was not
   enough memory, Z_STREAM_ERROR if the source stream state was inconsistent
   (such as zalloc being Z_NULL).  msg is left unchanged in both source and
   destination.
*/

ZEXTERN int ZEXPORT deflateReset(z_streamp strm);
/*
     This function is equivalent to deflateEnd followed by deflateInit, but
   does not free and reallocate the internal compression state.  The stream
   will leave the compression level and any other attributes that may have been
   set unchanged.  total_in, total_out, adler, and msg are initialized.

     deflateReset returns Z_OK if success, or Z_STREAM_ERROR if the source
   stream state was inconsistent (such as zalloc or state being Z_NULL).
*/

ZEXTERN uLong ZEXPORT deflateBound(z_streamp strm, uLong sourceLen);
ZEXTERN z_size_t ZEXPORT deflateBound_z(z_streamp strm, z_size_t sourceLen);
/*
     deflateBound() returns an upper bound on the compressed size after
   deflation of sourceLen bytes.  It must be called after deflateInit() or
   deflateInit2(), and after deflateSetHeader(), if used.  This would be used
   to allocate an output buffer for deflation in a single pass, and so would be
   called before deflate().  If that first deflate() call is provided the
   sourceLen input bytes, an output buffer allocated to the size returned by
   deflateBound(), and the flush value Z_FINISH, then deflate() is guaranteed
   to return Z_STREAM_END.  Note that it is possible for the compressed size to
   be larger than the value returned by deflateBound() if flush options other
   than Z_FINISH or Z_NO_FLUSH are used.

     delfateBound_z() is the same, but takes and returns a size_t length.  Note
   that a long is 32 bits on Windows.
*/

                        /* various hacks, don't look :) */

/* deflateInit and inflateInit are macros to allow checking the zlib version
 * and the compiler's view of z_stream:
 */
ZEXTERN int ZEXPORT deflateInit_(z_streamp strm, int level,
                                 const char *version, int stream_size);
ZEXTERN int ZEXPORT inflateInit_(z_streamp strm,
                                 const char *version, int stream_size);
ZEXTERN int ZEXPORT deflateInit2_(z_streamp strm, int  level, int  method,
                                  int windowBits, int memLevel,
                                  int strategy, const char *version,
                                  int stream_size);
ZEXTERN int ZEXPORT inflateInit2_(z_streamp strm, int  windowBits,
                                  const char *version, int stream_size);
ZEXTERN int ZEXPORT inflateBackInit_(z_streamp strm, int windowBits,
                                     unsigned char FAR *window,
                                     const char *version,
                                     int stream_size);
#ifdef Z_PREFIX_SET
#  define z_deflateInit(strm, level) \
          deflateInit_((strm), (level), ZLIB_VERSION, (int)sizeof(z_stream))
#  define z_inflateInit(strm) \
          inflateInit_((strm), ZLIB_VERSION, (int)sizeof(z_stream))
#  define z_deflateInit2(strm, level, method, windowBits, memLevel, strategy) \
          deflateInit2_((strm),(level),(method),(windowBits),(memLevel),\
                        (strategy), ZLIB_VERSION, (int)sizeof(z_stream))
#  define z_inflateInit2(strm, windowBits) \
          inflateInit2_((strm), (windowBits), ZLIB_VERSION, \
                        (int)sizeof(z_stream))
#  define z_inflateBackInit(strm, windowBits, window) \
          inflateBackInit_((strm), (windowBits), (window), \
                           ZLIB_VERSION, (int)sizeof(z_stream))
#else
#  define deflateInit(strm, level) \
          deflateInit_((strm), (level), ZLIB_VERSION, (int)sizeof(z_stream))
#  define inflateInit(strm) \
          inflateInit_((strm), ZLIB_VERSION, (int)sizeof(z_stream))
#  define deflateInit2(strm, level, method, windowBits, memLevel, strategy) \
          deflateInit2_((strm),(level),(method),(windowBits),(memLevel),\
                        (strategy), ZLIB_VERSION, (int)sizeof(z_stream))
#  define inflateInit2(strm, windowBits) \
          inflateInit2_((strm), (windowBits), ZLIB_VERSION, \
                        (int)sizeof(z_stream))
#  define inflateBackInit(strm, windowBits, window) \
          inflateBackInit_((strm), (windowBits), (window), \
                           ZLIB_VERSION, (int)sizeof(z_stream))
#endif

typedef struct gzFile_s *gzFile;    /* semi-opaque gzip file descriptor */

ZEXTERN gzFile ZEXPORT gzdopen(int fd, const char *mode);
/*
     Associate a gzFile with the file descriptor fd.  File descriptors are
   obtained from calls like open, dup, creat, pipe or fileno (if the file has
   been previously opened with fopen).  The mode parameter is as in gzopen. An
   'e' in mode will set fd's flag to close the file on an execve() call. An 'N'
   in mode will set fd's non-blocking flag.

     The next call of gzclose on the returned gzFile will also close the file
   descriptor fd, just like fclose(fdopen(fd, mode)) closes the file descriptor
   fd.  If you want to keep fd open, use fd = dup(fd_keep); gz = gzdopen(fd,
   mode);.  The duplicated descriptor should be saved to avoid a leak, since
   gzdopen does not close fd if it fails.  If you are using fileno() to get the
   file descriptor from a FILE *, then you will have to use dup() to avoid
   double-close()ing the file descriptor.  Both gzclose() and fclose() will
   close the associated file descriptor, so they need to have different file
   descriptors.

     gzdopen returns NULL if there was insufficient memory to allocate the
   gzFile state, if an invalid mode was specified (an 'r', 'w', or 'a' was not
   provided, or '+' was provided), or if fd is -1.  The file descriptor is not
   used until the next gz* read, write, seek, or close operation, so gzdopen
   will not detect if fd is invalid (unless fd is -1).
*/

ZEXTERN int ZEXPORT gzbuffer(gzFile file, unsigned size);
/*
     Set the internal buffer size used by this library's functions for file to
   size.  The default buffer size is 8192 bytes.  This function must be called
   after gzopen() or gzdopen(), and before any other calls that read or write
   the file.  The buffer memory allocation is always deferred to the first read
   or write.  Three times that size in buffer space is allocated.  A larger
   buffer size of, for example, 64K or 128K bytes will noticeably increase the
   speed of decompression (reading).

     The new buffer size also affects the maximum length for gzprintf().

     gzbuffer() returns 0 on success, or -1 on failure, such as being called
   too late.
*/

/* gzgetc() macro and its supporting function and exposed data structure.  Note
 * that the real internal state is much larger than the exposed structure.
 * This abbreviated structure exposes just enough for the gzgetc() macro.  The
 * user should not mess with these exposed elements, since their names or
 * behavior could change in the future, perhaps even capriciously.  They can
 * only be used by the gzgetc() macro.  You have been warned.
 */
struct gzFile_s {
    unsigned have;
    unsigned char *next;
    z_off64_t pos;
};
ZEXTERN int ZEXPORT gzgetc_(gzFile file);       /* backward compatibility */
#ifdef Z_PREFIX_SET
#  undef z_gzgetc
#  define z_gzgetc(g) \
          ((g)->have ? ((g)->have--, (g)->pos++, *((g)->next)++) : (gzgetc_)(g))
#else
#  define gzgetc(g) \
          ((g)->have ? ((g)->have--, (g)->pos++, *((g)->next)++) : (gzgetc_)(g))
#endif
