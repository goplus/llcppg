// Struct with multiple bit-field runs separated by regular fields
typedef struct {
    unsigned int a : 2;      // run 0, bits 0-1
    unsigned int b : 3;      // run 0, bits 2-4
    int          c;          // regular field at byte 4
    unsigned int d : 4;      // run 1, bits 32-35
    unsigned int e : 5;      // run 1, bits 36-40
    int          f;          // regular field at byte 12
} MultiRun;
