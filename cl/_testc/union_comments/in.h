/** Named union documentation. */
union Value {
    int number;
    float real;
};

/** Typedef union documentation. */
typedef union {
    int code;
    float value;
} Result;

union InlineUnion {
    int value;
}; // Inline union type documentation.

/** Empty union documentation. */
union Empty {};
