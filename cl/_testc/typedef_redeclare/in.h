// A forward typedef of a struct tag, mirroring id3v2lib's frame.h, followed by
// the defining typedef of the same tag (as in frame_header.h). Both emit the C
// name FrameHeader; the generator must reuse the single definition instead of
// redeclaring it and panicking. See issue goplus/llcppg#1001.
typedef struct _FrameHeader FrameHeader;

typedef struct _FrameHeader {
  int size;
} FrameHeader;

// A plain basic-type typedef repeated with the same underlying type must also
// be deduplicated rather than redeclared.
typedef int MyInt;

typedef int MyInt;
