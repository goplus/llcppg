/* This is the base type that holds just about everything :-) */
struct asn1_string_st;

typedef struct asn1_string_st ASN1_STRING;

ASN1_STRING *ASN1_STRING_new(void);
void ASN1_STRING_free(ASN1_STRING *a);
void ASN1_STRING_clear_free(ASN1_STRING *a);
int ASN1_STRING_copy(ASN1_STRING *dst, const ASN1_STRING *str);

typedef struct asn1_string_st ASN1_BIT_STRING;

ASN1_BIT_STRING *ASN1_BIT_STRING_new(void);
void ASN1_BIT_STRING_free(ASN1_BIT_STRING *a);

int ASN1_BIT_STRING_set(ASN1_BIT_STRING *a, unsigned char *d, int length);
int ASN1_BIT_STRING_set_bit(ASN1_BIT_STRING *a, int n, int value);
int ASN1_BIT_STRING_get_bit(const ASN1_BIT_STRING *a, int n);
int ASN1_BIT_STRING_check(const ASN1_BIT_STRING *a,
    const unsigned char *flags, int flags_len);
