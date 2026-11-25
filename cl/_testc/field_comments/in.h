struct Record {
    /// The horizontal coordinate.
    int x; // In pixels.
    int y; // The vertical coordinate.
    int plain;

    /**
     * A documented value.
     * Kept on two lines.
     */
    int multiline;
    int values[2]; /* Two stored values. */
    int renamed; // A field renamed by configuration.
    int inline_block; /* A block starting on the field's line.
                        * Its continuation stays with this field. */
    int after_block;

    /// The field's leading documentation.
    int both_multiline; /* A second description.
                          * Kept with the leading documentation. */
    int after_both;
};

struct Trailing {
    int value;
    // The last field's following documentation.
};

struct Positioned {
    int first;
    // Only the second field's documentation.
    // Its second documentation line.
    int second; /* The second field's inline comment. */
    int third;
};

struct Nested {
    struct {
        int value; // A named record's member.
    } child; // The named record.
    struct {
        int value; // An array element's member.
    } entries[2]; // The record array.
    struct {
        int promoted; // A promoted struct member.
    }; // The embedded struct.
    union {
        int code; // A promoted union member.
        float value;
    }; // The embedded union.
    int tail; // After the embedded records.
};

union Value {
    /// An integer value.
    int number; // In units of one.
    float real; /* A floating-point value. */
    int plain;

    struct {
        int tag; // The named view's tag.
    } named; // The named struct view.
    struct {
        /// The first member's leading documentation.
        int first; // The first promoted member.
        int second; // The second promoted member.
        union {
            int nested; // A member promoted through both records.
        };
    };
    union {
        int other; // A member promoted through an anonymous union.
    };
};
