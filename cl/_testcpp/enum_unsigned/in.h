enum class ByteFlags : unsigned char {
    High = 0x80,
    All = 0xff
};

enum class ShortFlags : unsigned short {
    High = 0x8000,
    All = 0xffff
};

enum class WideFlags : unsigned long long {
    High = 1ULL << 63,
    All = 0xffffffffffffffffULL
};

enum class SignedValues : long long {
    Negative = -1,
    Minimum = -9223372036854775807LL - 1
};
