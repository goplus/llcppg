/* A po_file_t represents the contents of a PO file.  */
typedef struct po_file *po_file_t;
/* A po_message_iterator_t represents an iterator through a domain of a
   PO file.  */
typedef struct po_message_iterator *po_message_iterator_t;
/* Create an iterator for traversing a domain of a PO file in memory.
   The domain NULL denotes the default domain.  */
extern po_message_iterator_t po_message_iterator (po_file_t file, const char *domain);