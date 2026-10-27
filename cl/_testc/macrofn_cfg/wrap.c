#include <foo.h>

int _llcppg_z_gzgetc(gzFile g) {
	return z_gzgetc(g);
}

int _llcppg_z_deflateInit(z_streamp strm, int level) {
	return z_deflateInit(strm, level);
}

void _llcppg_z_reset(z_streamp strm, int flags) {
	z_reset(strm, flags);
}

void _llcppg_z_flush(z_streamp strm) {
	z_flush(strm);
}
