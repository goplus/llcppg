/**
 * Categorizes how memory is being used by a translation unit.
 */
enum CXTUResourceUsageKind {
  CXTUResourceUsage_AST = 1,
  CXTUResourceUsage_Identifiers = 2
};

typedef struct CXTUResourceUsageEntry {
  /* The memory usage category. */
  enum CXTUResourceUsageKind kind;
  /* Amount of resources used.
      The units will depend on the resource kind. */
  unsigned long amount;
} CXTUResourceUsageEntry;

/**
 * The memory usage of a CXTranslationUnit, broken into categories.
 */
typedef struct CXTUResourceUsage {
  /* Private data member, used for queries. */
  void *data;

  /* The number of entries in the 'entries' array. */
  unsigned numEntries;

  /* An array of key-value pairs, representing the breakdown of memory
            usage. */
  CXTUResourceUsageEntry *entries;

} CXTUResourceUsage;
