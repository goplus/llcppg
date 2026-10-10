// Named members retain their access level even when their union has no tag.
struct Packet {
    int kind;
    union {
        int i;
        union {
            int x;
            short s;
        } inner;
    } data;
};

// Fieldless union levels promote members without adding storage.
union Anonymous {
    int i;
    union {
        short s;
        float f;
    };
};

union Recursive {
    int i;
    union {
        union {
            short s;
            float f;
        };
    };
};

// Promotion stops at an actual named member.
union Mixed {
    int i;
    union {
        short s;
        union {
            int n;
            float f;
        } inner;
    };
};

// A struct embeds its fieldless union; the inner union still shares storage.
struct Envelope {
    int kind;
    union {
        int i;
        union {
            short s;
            float f;
        };
    };
};

// Existing named unions and typedefs keep their types.
union Tagged {
    int i;
    float f;
};
typedef union Tagged Alias;
union References {
    union Tagged tagged;
    Alias alias;
};
struct Holder {
    Alias value;
};
