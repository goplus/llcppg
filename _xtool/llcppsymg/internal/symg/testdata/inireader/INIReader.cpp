#include "INIReader.h"

__attribute__((visibility("default"))) INIReader::INIReader(const char *filename) {}
INI_API INIReader::INIReader(const char *buffer, long buffer_size) {}
INIReader::~INIReader() {}
INI_API int INIReader::ParseError() const {}
INI_API const char *INIReader::Get(const char *section, const char *name,
                                   const char *default_value) const {}

const char *INIReader::MakeKey(const char *section, const char *name) {}
