/* This is the base type that holds just about everything :-) */
struct asn1_string_st;

typedef struct asn1_string_st ASN1_STRING;

ASN1_STRING *ASN1_STRING_new(void);
void ASN1_STRING_free(ASN1_STRING *a);
void ASN1_STRING_clear_free(ASN1_STRING *a);
int ASN1_STRING_copy(ASN1_STRING *dst, const ASN1_STRING *str);

int i2d_DISPLAYTEXT(ASN1_STRING* a, char** pp);

typedef struct asn1_string_st ASN1_BIT_STRING;

ASN1_BIT_STRING *ASN1_BIT_STRING_new(void);
void ASN1_BIT_STRING_free(ASN1_BIT_STRING *a);

int ASN1_BIT_STRING_set(ASN1_BIT_STRING *a, unsigned char *d, int length);
int ASN1_BIT_STRING_set_bit(ASN1_BIT_STRING *a, int n, int value);
int ASN1_BIT_STRING_get_bit(const ASN1_BIT_STRING *a, int n);
int ASN1_BIT_STRING_check(const ASN1_BIT_STRING *a,
    const unsigned char *flags, int flags_len);

typedef struct asn1_type_st ASN1_TYPE;

ASN1_TYPE *ASN1_generate_nconf(const char *str, void *nconf);

typedef struct asn1_object_st ASN1_OBJECT;

int i2d_ASN1_OBJECT(ASN1_OBJECT* obj, char** pp);

int OBJ_add_object(ASN1_OBJECT *obj);
