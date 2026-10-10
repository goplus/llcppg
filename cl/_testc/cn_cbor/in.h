// Reduced from jimsch/cn-cbor 1.0.0, include/cn-cbor/cn-cbor.h.
// Preserve all members of the anonymous union and the recursive record fields.
typedef unsigned char uint8_t;

typedef enum cn_cbor_type {
    CN_CBOR_FALSE,
    CN_CBOR_TRUE,
    CN_CBOR_NULL,
    CN_CBOR_UNDEF,
    CN_CBOR_UINT,
    CN_CBOR_INT,
    CN_CBOR_BYTES,
    CN_CBOR_TEXT,
    CN_CBOR_BYTES_CHUNKED,
    CN_CBOR_TEXT_CHUNKED,
    CN_CBOR_ARRAY,
    CN_CBOR_MAP,
    CN_CBOR_TAG,
    CN_CBOR_SIMPLE,
    CN_CBOR_DOUBLE,
    CN_CBOR_FLOAT,
    CN_CBOR_INVALID
} cn_cbor_type;

typedef enum cn_cbor_flags {
    CN_CBOR_FL_COUNT = 1,
    CN_CBOR_FL_INDEF = 2,
    CN_CBOR_FL_OWNER = 0x80
} cn_cbor_flags;

typedef struct cn_cbor {
    cn_cbor_type type;
    cn_cbor_flags flags;
    union {
        const uint8_t *bytes;
        const char *str;
        long sint;
        unsigned long uint;
        double dbl;
        float f;
        unsigned long count;
    } v;
    int length;
    struct cn_cbor *first_child;
    struct cn_cbor *last_child;
    struct cn_cbor *next;
    struct cn_cbor *parent;
} cn_cbor;
