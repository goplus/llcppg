#include "builtin/types.h"
#include "external/types.h"
#include "unconfigured/types.h"

#if __has_feature(modules)
#error Module compilation must remain disabled
#endif

struct BuiltinTypedefFields {
    __gnuc_va_list *args;
    size_t size;
    external_mapped mapped;
};

void take_va_list(__gnuc_va_list *args);
void take_external(external_va_list *args);
void take_unconfigured(unconfigured_va_list *args);
void take_ignored(external_ignored value);
