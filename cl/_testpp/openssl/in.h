#define SSL_set_tlsext_status_type(ssl, type) \
    SSL_ctrl(ssl, 0, type, NULL)

typedef struct SSL SSL;

long SSL_ctrl(SSL *ssl, int cmd, long larg, void *parg);
