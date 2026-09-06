const void *ares_dns_pton(const char *ipaddr, struct ares_addr *addr);
char *ares_dns_addr_to_ptr(const struct ares_addr *addr);


// case for https://github.com/xgo-dev/llcppg/issues/510
typedef enum {
    ESP_LOG_NONE = 0,
    ESP_LOG_ERROR = 1,
    ESP_LOG_WARN = 2,
    ESP_LOG_INFO = 3,
    ESP_LOG_DEBUG = 4,
    ESP_LOG_VERBOSE = 5,
    ESP_LOG_MAX = 6,
} esp_log_level_t;
void esp_log_write(esp_log_level_t level, const char *tag, const char *format, ...);
