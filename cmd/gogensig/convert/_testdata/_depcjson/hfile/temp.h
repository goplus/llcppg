#include <cJSON.h>
#include <stddef.h>
#include <thirddep.h>
#include <thirddep2.h>
#include "type.h"
// This file is supposed to depend on cjson in its cflags, but for testing,
// we will simulate its API using libcjson instead.
//   "cflags" :"$(pkg-config --cflags libcjson)"
cJSON* create_response(int status_code, const char* message);

cJSON_bool parse_client_request(const char* json_string, char* error_buffer, size_t buffer_size);

cJSON_bool serialize_response(cJSON *response, char *buffer, const int length, const cJSON_bool pretty_print);

third_dep third_depfn(third_dep *a, third_dep2 *b,_depcjson_type c,basic_dep d);