enum { enum1, enum2 };

enum { COLOR_DEFAULT = -1 };

enum spectrum { red, orange, yello, green, blue, violet };

enum kids { nippy, slats, skippy, nina, liz };

enum levels { low = 100, medium = 500, high = 2000 };

enum feline { cat, lynx = 10, puma, tiger };

typedef enum algorithm {
    UNKNOWN = 0,
    NULL = 1,
} algorithm_t;

typedef enum {
    UNKNOWN2 = 0,
    NULL2 = 1,
} algorithm_t2;

typedef algorithm_t algorithm;

typedef enum {
    GPG_ERR_NO_ERROR = 0,
    GPG_ERR_GENERAL = 1,
    GPG_ERR_UNKNOWN_PACKET = 2,
    GPG_ERR_UNKNOWN_VERSION = 3,
    GPG_ERR_PUBKEY_ALGO = 4,
    GPG_ERR_DIGEST_ALGO = 5,
    GPG_ERR_BAD_PUBKEY = 6,
    GPG_ERR_BAD_SECKEY = 7,
    GPG_ERR_BAD_SIGNATURE = 8,
    GPG_ERR_NO_PUBKEY = 9,
    GPG_ERR_CHECKSUM = 10,
    GPG_ERR_BAD_PASSPHRASE = 11,
    GPG_ERR_CIPHER_ALGO = 12,
    GPG_ERR_KEYRING_OPEN = 13,
    GPG_ERR_INV_PACKET = 14,
    GPG_ERR_INV_ARMOR = 15,
    GPG_ERR_NO_USER_ID = 16,
    GPG_ERR_NO_SECKEY = 17,
    GPG_ERR_WRONG_SECKEY = 18,
    GPG_ERR_BAD_KEY = 19,
} gpg_err_code_t;
