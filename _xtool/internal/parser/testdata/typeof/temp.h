// https://github.com/xgo-dev/llcppg/issues/497
typedef struct {
    int x;
    union {
        long val;
    } clock;
} spi_mem_dev_t;

extern spi_mem_dev_t GPSPI2_t;

typedef typeof(GPSPI2_t.clock.val) gpspi_flash_ll_clock_reg_t;

typedef union {
    gpspi_flash_ll_clock_reg_t clock;
} gpspi_flash_ll_dev_t;
