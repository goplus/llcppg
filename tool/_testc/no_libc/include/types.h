#include <probe.h>
#include <dep.h>

enum Status { StatusReady = 0 };

typedef void Void;
typedef char Char;
typedef int Int;
typedef unsigned int Uint;
typedef long Long;
typedef unsigned long Ulong;
typedef long long LongLong;
typedef unsigned long long UlongLong;
typedef float Float;
typedef double Double;
typedef int ConfiguredInt;

#include "stdint.h"

typedef PROBE_INT_TYPE probe_int_t;
typedef dep_t External;

struct Sample {
    int native;
    int32_t fixed;
    char code;
    double measure;
    void *context;
};

void accept(int value, char code);
