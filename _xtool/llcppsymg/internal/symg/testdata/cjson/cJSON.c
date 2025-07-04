#include "cJSON.h"

CJSON_PUBLIC(char *)
cJSON_Print(const cJSON *item) {}
CJSON_PUBLIC(cJSON *)
cJSON_ParseWithLength(const char *value, size_t buffer_length) {}
/* Delete a cJSON entity and all subentities. */
CJSON_PUBLIC(void)
cJSON_Delete(cJSON *item) {}
