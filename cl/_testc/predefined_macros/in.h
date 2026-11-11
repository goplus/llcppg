#define COMPILER_CLANG __clang__
#define CHAR_BITS __CHAR_BIT__
#define INT_BITS (__SIZEOF_INT__ * __CHAR_BIT__)
#define INT_MAX_VALUE __INT_MAX__
#define BUILTIN_EXPR (__clang__ + 41)
#define LITERAL_ENCODING __clang_literal_encoding__

#if defined(__APPLE__)
#define PLATFORM_FLAG __APPLE__
#elif defined(__linux__)
#define PLATFORM_FLAG __linux__
#elif defined(_WIN32)
#define PLATFORM_FLAG _WIN32
#endif
