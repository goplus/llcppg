enum UnsignedFlags {
    FlagHigh = 1u << 31,
    FlagAll = 0xffffffffu
};

enum UnsignedWide {
    WideHigh = 1ULL << 63,
    WideAll = 0xffffffffffffffffULL
};

enum SignedValues {
    Negative = -1,
    Positive = 1
};

enum MixedValues {
    MixedNegative = -1,
    MixedHigh = 1u << 31
};

enum {
    AnonymousHigh = 1u << 31,
    AnonymousAll = 0xffffffffu
};

enum {
    AnonymousWide = 0xffffffffffffffffULL
};
