struct NestedEnum
{
    enum
    {
        APR_BUCKET_DATA1 = 0,
        APR_BUCKET_METADATA2 = 1
    } is_metadata1_t;

    struct a
    {
        enum
        {
            APR_BUCKET_DATA_A1 = 0,
            APR_BUCKET_METADATA_A2 = 1
        } is_metadata1_t;
    } a_t;
};

struct NestedEnum2
{
    enum
    {
        APR_BUCKET_DATA3 = 0,
        APR_BUCKET_METADATA4 = 1
    };
};

struct NestedEnum3
{
    enum is_metadata3
    {
        APR_BUCKET_DATA5 = 0,
        APR_BUCKET_METADATA6 = 1
    };
};

struct NestedEnum4
{
    enum is_metadata4
    {
        APR_BUCKET_DATA7 = 0,
        APR_BUCKET_METADATA8 = 1
    } key;
};

enum OuterEnum
{
    APR_BUCKET_DATA9 = 0,
    APR_BUCKET_METADATA10 = 1
};
struct Enum
{
    enum OuterEnum k;
};
