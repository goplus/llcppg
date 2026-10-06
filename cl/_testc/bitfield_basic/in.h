// A simple struct with bit-fields, based on the proposal example.
// ready: bit 0, width 1, unsigned
// mode: bits 1-3, width 3, unsigned
// delta: bits 4-7, width 4, signed
// (zero-width field at bit 8)
// level: bits 32-37, width 6, unsigned
// id: regular field at byte 8
typedef struct {
    unsigned int ready : 1;
    unsigned int mode  : 3;
    int          delta : 4;
    unsigned int       : 0;
    unsigned int level : 6;
    int          id;
} Packet;
