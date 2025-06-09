package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type ErrorLevel c.Int

const (
	ERR_NONE    ErrorLevel = 0
	ERR_WARNING ErrorLevel = 1
	ERR_ERROR   ErrorLevel = 2
	ERR_FATAL   ErrorLevel = 3
)

type ErrorDomain c.Int

const (
	FROM_NONE        ErrorDomain = 0
	FROM_PARSER      ErrorDomain = 1
	FROM_TREE        ErrorDomain = 2
	FROM_NAMESPACE   ErrorDomain = 3
	FROM_DTD         ErrorDomain = 4
	FROM_HTML        ErrorDomain = 5
	FROM_MEMORY      ErrorDomain = 6
	FROM_OUTPUT      ErrorDomain = 7
	FROM_IO          ErrorDomain = 8
	FROM_FTP         ErrorDomain = 9
	FROM_HTTP        ErrorDomain = 10
	FROM_XINCLUDE    ErrorDomain = 11
	FROM_XPATH       ErrorDomain = 12
	FROM_XPOINTER    ErrorDomain = 13
	FROM_REGEXP      ErrorDomain = 14
	FROM_DATATYPE    ErrorDomain = 15
	FROM_SCHEMASP    ErrorDomain = 16
	FROM_SCHEMASV    ErrorDomain = 17
	FROM_RELAXNGP    ErrorDomain = 18
	FROM_RELAXNGV    ErrorDomain = 19
	FROM_CATALOG     ErrorDomain = 20
	FROM_C14N        ErrorDomain = 21
	FROM_XSLT        ErrorDomain = 22
	FROM_VALID       ErrorDomain = 23
	FROM_CHECK       ErrorDomain = 24
	FROM_WRITER      ErrorDomain = 25
	FROM_MODULE      ErrorDomain = 26
	FROM_I18N        ErrorDomain = 27
	FROM_SCHEMATRONV ErrorDomain = 28
	FROM_BUFFER      ErrorDomain = 29
	FROM_URI         ErrorDomain = 30
)

type X_xmlError struct {
	Domain  c.Int
	Code    c.Int
	Message *c.Char
	Level   ErrorLevel
	File    *c.Char
	Line    c.Int
	Str1    *c.Char
	Str2    *c.Char
	Str3    *c.Char
	Int1    c.Int
	Int2    c.Int
	Ctxt    c.Pointer
	Node    c.Pointer
}
type Error X_xmlError
type ErrorPtr *Error
type ParserErrors c.Int

const (
	ERR_OK                                       ParserErrors = 0
	ERR_INTERNAL_ERROR                           ParserErrors = 1
	ERR_NO_MEMORY                                ParserErrors = 2
	ERR_DOCUMENT_START                           ParserErrors = 3
	ERR_DOCUMENT_EMPTY                           ParserErrors = 4
	ERR_DOCUMENT_END                             ParserErrors = 5
	ERR_INVALID_HEX_CHARREF                      ParserErrors = 6
	ERR_INVALID_DEC_CHARREF                      ParserErrors = 7
	ERR_INVALID_CHARREF                          ParserErrors = 8
	ERR_INVALID_CHAR                             ParserErrors = 9
	ERR_CHARREF_AT_EOF                           ParserErrors = 10
	ERR_CHARREF_IN_PROLOG                        ParserErrors = 11
	ERR_CHARREF_IN_EPILOG                        ParserErrors = 12
	ERR_CHARREF_IN_DTD                           ParserErrors = 13
	ERR_ENTITYREF_AT_EOF                         ParserErrors = 14
	ERR_ENTITYREF_IN_PROLOG                      ParserErrors = 15
	ERR_ENTITYREF_IN_EPILOG                      ParserErrors = 16
	ERR_ENTITYREF_IN_DTD                         ParserErrors = 17
	ERR_PEREF_AT_EOF                             ParserErrors = 18
	ERR_PEREF_IN_PROLOG                          ParserErrors = 19
	ERR_PEREF_IN_EPILOG                          ParserErrors = 20
	ERR_PEREF_IN_INT_SUBSET                      ParserErrors = 21
	ERR_ENTITYREF_NO_NAME                        ParserErrors = 22
	ERR_ENTITYREF_SEMICOL_MISSING                ParserErrors = 23
	ERR_PEREF_NO_NAME                            ParserErrors = 24
	ERR_PEREF_SEMICOL_MISSING                    ParserErrors = 25
	ERR_UNDECLARED_ENTITY                        ParserErrors = 26
	WAR_UNDECLARED_ENTITY                        ParserErrors = 27
	ERR_UNPARSED_ENTITY                          ParserErrors = 28
	ERR_ENTITY_IS_EXTERNAL                       ParserErrors = 29
	ERR_ENTITY_IS_PARAMETER                      ParserErrors = 30
	ERR_UNKNOWN_ENCODING                         ParserErrors = 31
	ERR_UNSUPPORTED_ENCODING                     ParserErrors = 32
	ERR_STRING_NOT_STARTED                       ParserErrors = 33
	ERR_STRING_NOT_CLOSED                        ParserErrors = 34
	ERR_NS_DECL_ERROR                            ParserErrors = 35
	ERR_ENTITY_NOT_STARTED                       ParserErrors = 36
	ERR_ENTITY_NOT_FINISHED                      ParserErrors = 37
	ERR_LT_IN_ATTRIBUTE                          ParserErrors = 38
	ERR_ATTRIBUTE_NOT_STARTED                    ParserErrors = 39
	ERR_ATTRIBUTE_NOT_FINISHED                   ParserErrors = 40
	ERR_ATTRIBUTE_WITHOUT_VALUE                  ParserErrors = 41
	ERR_ATTRIBUTE_REDEFINED                      ParserErrors = 42
	ERR_LITERAL_NOT_STARTED                      ParserErrors = 43
	ERR_LITERAL_NOT_FINISHED                     ParserErrors = 44
	ERR_COMMENT_NOT_FINISHED                     ParserErrors = 45
	ERR_PI_NOT_STARTED                           ParserErrors = 46
	ERR_PI_NOT_FINISHED                          ParserErrors = 47
	ERR_NOTATION_NOT_STARTED                     ParserErrors = 48
	ERR_NOTATION_NOT_FINISHED                    ParserErrors = 49
	ERR_ATTLIST_NOT_STARTED                      ParserErrors = 50
	ERR_ATTLIST_NOT_FINISHED                     ParserErrors = 51
	ERR_MIXED_NOT_STARTED                        ParserErrors = 52
	ERR_MIXED_NOT_FINISHED                       ParserErrors = 53
	ERR_ELEMCONTENT_NOT_STARTED                  ParserErrors = 54
	ERR_ELEMCONTENT_NOT_FINISHED                 ParserErrors = 55
	ERR_XMLDECL_NOT_STARTED                      ParserErrors = 56
	ERR_XMLDECL_NOT_FINISHED                     ParserErrors = 57
	ERR_CONDSEC_NOT_STARTED                      ParserErrors = 58
	ERR_CONDSEC_NOT_FINISHED                     ParserErrors = 59
	ERR_EXT_SUBSET_NOT_FINISHED                  ParserErrors = 60
	ERR_DOCTYPE_NOT_FINISHED                     ParserErrors = 61
	ERR_MISPLACED_CDATA_END                      ParserErrors = 62
	ERR_CDATA_NOT_FINISHED                       ParserErrors = 63
	ERR_RESERVED_XML_NAME                        ParserErrors = 64
	ERR_SPACE_REQUIRED                           ParserErrors = 65
	ERR_SEPARATOR_REQUIRED                       ParserErrors = 66
	ERR_NMTOKEN_REQUIRED                         ParserErrors = 67
	ERR_NAME_REQUIRED                            ParserErrors = 68
	ERR_PCDATA_REQUIRED                          ParserErrors = 69
	ERR_URI_REQUIRED                             ParserErrors = 70
	ERR_PUBID_REQUIRED                           ParserErrors = 71
	ERR_LT_REQUIRED                              ParserErrors = 72
	ERR_GT_REQUIRED                              ParserErrors = 73
	ERR_LTSLASH_REQUIRED                         ParserErrors = 74
	ERR_EQUAL_REQUIRED                           ParserErrors = 75
	ERR_TAG_NAME_MISMATCH                        ParserErrors = 76
	ERR_TAG_NOT_FINISHED                         ParserErrors = 77
	ERR_STANDALONE_VALUE                         ParserErrors = 78
	ERR_ENCODING_NAME                            ParserErrors = 79
	ERR_HYPHEN_IN_COMMENT                        ParserErrors = 80
	ERR_INVALID_ENCODING                         ParserErrors = 81
	ERR_EXT_ENTITY_STANDALONE                    ParserErrors = 82
	ERR_CONDSEC_INVALID                          ParserErrors = 83
	ERR_VALUE_REQUIRED                           ParserErrors = 84
	ERR_NOT_WELL_BALANCED                        ParserErrors = 85
	ERR_EXTRA_CONTENT                            ParserErrors = 86
	ERR_ENTITY_CHAR_ERROR                        ParserErrors = 87
	ERR_ENTITY_PE_INTERNAL                       ParserErrors = 88
	ERR_ENTITY_LOOP                              ParserErrors = 89
	ERR_ENTITY_BOUNDARY                          ParserErrors = 90
	ERR_INVALID_URI                              ParserErrors = 91
	ERR_URI_FRAGMENT                             ParserErrors = 92
	WAR_CATALOG_PI                               ParserErrors = 93
	ERR_NO_DTD                                   ParserErrors = 94
	ERR_CONDSEC_INVALID_KEYWORD                  ParserErrors = 95
	ERR_VERSION_MISSING                          ParserErrors = 96
	WAR_UNKNOWN_VERSION                          ParserErrors = 97
	WAR_LANG_VALUE                               ParserErrors = 98
	WAR_NS_URI                                   ParserErrors = 99
	WAR_NS_URI_RELATIVE                          ParserErrors = 100
	ERR_MISSING_ENCODING                         ParserErrors = 101
	WAR_SPACE_VALUE                              ParserErrors = 102
	ERR_NOT_STANDALONE                           ParserErrors = 103
	ERR_ENTITY_PROCESSING                        ParserErrors = 104
	ERR_NOTATION_PROCESSING                      ParserErrors = 105
	WAR_NS_COLUMN                                ParserErrors = 106
	WAR_ENTITY_REDEFINED                         ParserErrors = 107
	ERR_UNKNOWN_VERSION                          ParserErrors = 108
	ERR_VERSION_MISMATCH                         ParserErrors = 109
	ERR_NAME_TOO_LONG                            ParserErrors = 110
	ERR_USER_STOP                                ParserErrors = 111
	ERR_COMMENT_ABRUPTLY_ENDED                   ParserErrors = 112
	WAR_ENCODING_MISMATCH                        ParserErrors = 113
	ERR_RESOURCE_LIMIT                           ParserErrors = 114
	ERR_ARGUMENT                                 ParserErrors = 115
	ERR_SYSTEM                                   ParserErrors = 116
	ERR_REDECL_PREDEF_ENTITY                     ParserErrors = 117
	ERR_INT_SUBSET_NOT_FINISHED                  ParserErrors = 118
	NS_ERR_XML_NAMESPACE                         ParserErrors = 200
	NS_ERR_UNDEFINED_NAMESPACE                   ParserErrors = 201
	NS_ERR_QNAME                                 ParserErrors = 202
	NS_ERR_ATTRIBUTE_REDEFINED                   ParserErrors = 203
	NS_ERR_EMPTY                                 ParserErrors = 204
	NS_ERR_COLON                                 ParserErrors = 205
	DTD_ATTRIBUTE_DEFAULT                        ParserErrors = 500
	DTD_ATTRIBUTE_REDEFINED                      ParserErrors = 501
	DTD_ATTRIBUTE_VALUE                          ParserErrors = 502
	DTD_CONTENT_ERROR                            ParserErrors = 503
	DTD_CONTENT_MODEL                            ParserErrors = 504
	DTD_CONTENT_NOT_DETERMINIST                  ParserErrors = 505
	DTD_DIFFERENT_PREFIX                         ParserErrors = 506
	DTD_ELEM_DEFAULT_NAMESPACE                   ParserErrors = 507
	DTD_ELEM_NAMESPACE                           ParserErrors = 508
	DTD_ELEM_REDEFINED                           ParserErrors = 509
	DTD_EMPTY_NOTATION                           ParserErrors = 510
	DTD_ENTITY_TYPE                              ParserErrors = 511
	DTD_ID_FIXED                                 ParserErrors = 512
	DTD_ID_REDEFINED                             ParserErrors = 513
	DTD_ID_SUBSET                                ParserErrors = 514
	DTD_INVALID_CHILD                            ParserErrors = 515
	DTD_INVALID_DEFAULT                          ParserErrors = 516
	DTD_LOAD_ERROR                               ParserErrors = 517
	DTD_MISSING_ATTRIBUTE                        ParserErrors = 518
	DTD_MIXED_CORRUPT                            ParserErrors = 519
	DTD_MULTIPLE_ID                              ParserErrors = 520
	DTD_NO_DOC                                   ParserErrors = 521
	DTD_NO_DTD                                   ParserErrors = 522
	DTD_NO_ELEM_NAME                             ParserErrors = 523
	DTD_NO_PREFIX                                ParserErrors = 524
	DTD_NO_ROOT                                  ParserErrors = 525
	DTD_NOTATION_REDEFINED                       ParserErrors = 526
	DTD_NOTATION_VALUE                           ParserErrors = 527
	DTD_NOT_EMPTY                                ParserErrors = 528
	DTD_NOT_PCDATA                               ParserErrors = 529
	DTD_NOT_STANDALONE                           ParserErrors = 530
	DTD_ROOT_NAME                                ParserErrors = 531
	DTD_STANDALONE_WHITE_SPACE                   ParserErrors = 532
	DTD_UNKNOWN_ATTRIBUTE                        ParserErrors = 533
	DTD_UNKNOWN_ELEM                             ParserErrors = 534
	DTD_UNKNOWN_ENTITY                           ParserErrors = 535
	DTD_UNKNOWN_ID                               ParserErrors = 536
	DTD_UNKNOWN_NOTATION                         ParserErrors = 537
	DTD_STANDALONE_DEFAULTED                     ParserErrors = 538
	DTD_XMLID_VALUE                              ParserErrors = 539
	DTD_XMLID_TYPE                               ParserErrors = 540
	DTD_DUP_TOKEN                                ParserErrors = 541
	HTML_STRUCURE_ERROR                          ParserErrors = 800
	HTML_UNKNOWN_TAG                             ParserErrors = 801
	HTML_INCORRECTLY_OPENED_COMMENT              ParserErrors = 802
	RNGP_ANYNAME_ATTR_ANCESTOR                   ParserErrors = 1000
	RNGP_ATTR_CONFLICT                           ParserErrors = 1001
	RNGP_ATTRIBUTE_CHILDREN                      ParserErrors = 1002
	RNGP_ATTRIBUTE_CONTENT                       ParserErrors = 1003
	RNGP_ATTRIBUTE_EMPTY                         ParserErrors = 1004
	RNGP_ATTRIBUTE_NOOP                          ParserErrors = 1005
	RNGP_CHOICE_CONTENT                          ParserErrors = 1006
	RNGP_CHOICE_EMPTY                            ParserErrors = 1007
	RNGP_CREATE_FAILURE                          ParserErrors = 1008
	RNGP_DATA_CONTENT                            ParserErrors = 1009
	RNGP_DEF_CHOICE_AND_INTERLEAVE               ParserErrors = 1010
	RNGP_DEFINE_CREATE_FAILED                    ParserErrors = 1011
	RNGP_DEFINE_EMPTY                            ParserErrors = 1012
	RNGP_DEFINE_MISSING                          ParserErrors = 1013
	RNGP_DEFINE_NAME_MISSING                     ParserErrors = 1014
	RNGP_ELEM_CONTENT_EMPTY                      ParserErrors = 1015
	RNGP_ELEM_CONTENT_ERROR                      ParserErrors = 1016
	RNGP_ELEMENT_EMPTY                           ParserErrors = 1017
	RNGP_ELEMENT_CONTENT                         ParserErrors = 1018
	RNGP_ELEMENT_NAME                            ParserErrors = 1019
	RNGP_ELEMENT_NO_CONTENT                      ParserErrors = 1020
	RNGP_ELEM_TEXT_CONFLICT                      ParserErrors = 1021
	RNGP_EMPTY                                   ParserErrors = 1022
	RNGP_EMPTY_CONSTRUCT                         ParserErrors = 1023
	RNGP_EMPTY_CONTENT                           ParserErrors = 1024
	RNGP_EMPTY_NOT_EMPTY                         ParserErrors = 1025
	RNGP_ERROR_TYPE_LIB                          ParserErrors = 1026
	RNGP_EXCEPT_EMPTY                            ParserErrors = 1027
	RNGP_EXCEPT_MISSING                          ParserErrors = 1028
	RNGP_EXCEPT_MULTIPLE                         ParserErrors = 1029
	RNGP_EXCEPT_NO_CONTENT                       ParserErrors = 1030
	RNGP_EXTERNALREF_EMTPY                       ParserErrors = 1031
	RNGP_EXTERNAL_REF_FAILURE                    ParserErrors = 1032
	RNGP_EXTERNALREF_RECURSE                     ParserErrors = 1033
	RNGP_FORBIDDEN_ATTRIBUTE                     ParserErrors = 1034
	RNGP_FOREIGN_ELEMENT                         ParserErrors = 1035
	RNGP_GRAMMAR_CONTENT                         ParserErrors = 1036
	RNGP_GRAMMAR_EMPTY                           ParserErrors = 1037
	RNGP_GRAMMAR_MISSING                         ParserErrors = 1038
	RNGP_GRAMMAR_NO_START                        ParserErrors = 1039
	RNGP_GROUP_ATTR_CONFLICT                     ParserErrors = 1040
	RNGP_HREF_ERROR                              ParserErrors = 1041
	RNGP_INCLUDE_EMPTY                           ParserErrors = 1042
	RNGP_INCLUDE_FAILURE                         ParserErrors = 1043
	RNGP_INCLUDE_RECURSE                         ParserErrors = 1044
	RNGP_INTERLEAVE_ADD                          ParserErrors = 1045
	RNGP_INTERLEAVE_CREATE_FAILED                ParserErrors = 1046
	RNGP_INTERLEAVE_EMPTY                        ParserErrors = 1047
	RNGP_INTERLEAVE_NO_CONTENT                   ParserErrors = 1048
	RNGP_INVALID_DEFINE_NAME                     ParserErrors = 1049
	RNGP_INVALID_URI                             ParserErrors = 1050
	RNGP_INVALID_VALUE                           ParserErrors = 1051
	RNGP_MISSING_HREF                            ParserErrors = 1052
	RNGP_NAME_MISSING                            ParserErrors = 1053
	RNGP_NEED_COMBINE                            ParserErrors = 1054
	RNGP_NOTALLOWED_NOT_EMPTY                    ParserErrors = 1055
	RNGP_NSNAME_ATTR_ANCESTOR                    ParserErrors = 1056
	RNGP_NSNAME_NO_NS                            ParserErrors = 1057
	RNGP_PARAM_FORBIDDEN                         ParserErrors = 1058
	RNGP_PARAM_NAME_MISSING                      ParserErrors = 1059
	RNGP_PARENTREF_CREATE_FAILED                 ParserErrors = 1060
	RNGP_PARENTREF_NAME_INVALID                  ParserErrors = 1061
	RNGP_PARENTREF_NO_NAME                       ParserErrors = 1062
	RNGP_PARENTREF_NO_PARENT                     ParserErrors = 1063
	RNGP_PARENTREF_NOT_EMPTY                     ParserErrors = 1064
	RNGP_PARSE_ERROR                             ParserErrors = 1065
	RNGP_PAT_ANYNAME_EXCEPT_ANYNAME              ParserErrors = 1066
	RNGP_PAT_ATTR_ATTR                           ParserErrors = 1067
	RNGP_PAT_ATTR_ELEM                           ParserErrors = 1068
	RNGP_PAT_DATA_EXCEPT_ATTR                    ParserErrors = 1069
	RNGP_PAT_DATA_EXCEPT_ELEM                    ParserErrors = 1070
	RNGP_PAT_DATA_EXCEPT_EMPTY                   ParserErrors = 1071
	RNGP_PAT_DATA_EXCEPT_GROUP                   ParserErrors = 1072
	RNGP_PAT_DATA_EXCEPT_INTERLEAVE              ParserErrors = 1073
	RNGP_PAT_DATA_EXCEPT_LIST                    ParserErrors = 1074
	RNGP_PAT_DATA_EXCEPT_ONEMORE                 ParserErrors = 1075
	RNGP_PAT_DATA_EXCEPT_REF                     ParserErrors = 1076
	RNGP_PAT_DATA_EXCEPT_TEXT                    ParserErrors = 1077
	RNGP_PAT_LIST_ATTR                           ParserErrors = 1078
	RNGP_PAT_LIST_ELEM                           ParserErrors = 1079
	RNGP_PAT_LIST_INTERLEAVE                     ParserErrors = 1080
	RNGP_PAT_LIST_LIST                           ParserErrors = 1081
	RNGP_PAT_LIST_REF                            ParserErrors = 1082
	RNGP_PAT_LIST_TEXT                           ParserErrors = 1083
	RNGP_PAT_NSNAME_EXCEPT_ANYNAME               ParserErrors = 1084
	RNGP_PAT_NSNAME_EXCEPT_NSNAME                ParserErrors = 1085
	RNGP_PAT_ONEMORE_GROUP_ATTR                  ParserErrors = 1086
	RNGP_PAT_ONEMORE_INTERLEAVE_ATTR             ParserErrors = 1087
	RNGP_PAT_START_ATTR                          ParserErrors = 1088
	RNGP_PAT_START_DATA                          ParserErrors = 1089
	RNGP_PAT_START_EMPTY                         ParserErrors = 1090
	RNGP_PAT_START_GROUP                         ParserErrors = 1091
	RNGP_PAT_START_INTERLEAVE                    ParserErrors = 1092
	RNGP_PAT_START_LIST                          ParserErrors = 1093
	RNGP_PAT_START_ONEMORE                       ParserErrors = 1094
	RNGP_PAT_START_TEXT                          ParserErrors = 1095
	RNGP_PAT_START_VALUE                         ParserErrors = 1096
	RNGP_PREFIX_UNDEFINED                        ParserErrors = 1097
	RNGP_REF_CREATE_FAILED                       ParserErrors = 1098
	RNGP_REF_CYCLE                               ParserErrors = 1099
	RNGP_REF_NAME_INVALID                        ParserErrors = 1100
	RNGP_REF_NO_DEF                              ParserErrors = 1101
	RNGP_REF_NO_NAME                             ParserErrors = 1102
	RNGP_REF_NOT_EMPTY                           ParserErrors = 1103
	RNGP_START_CHOICE_AND_INTERLEAVE             ParserErrors = 1104
	RNGP_START_CONTENT                           ParserErrors = 1105
	RNGP_START_EMPTY                             ParserErrors = 1106
	RNGP_START_MISSING                           ParserErrors = 1107
	RNGP_TEXT_EXPECTED                           ParserErrors = 1108
	RNGP_TEXT_HAS_CHILD                          ParserErrors = 1109
	RNGP_TYPE_MISSING                            ParserErrors = 1110
	RNGP_TYPE_NOT_FOUND                          ParserErrors = 1111
	RNGP_TYPE_VALUE                              ParserErrors = 1112
	RNGP_UNKNOWN_ATTRIBUTE                       ParserErrors = 1113
	RNGP_UNKNOWN_COMBINE                         ParserErrors = 1114
	RNGP_UNKNOWN_CONSTRUCT                       ParserErrors = 1115
	RNGP_UNKNOWN_TYPE_LIB                        ParserErrors = 1116
	RNGP_URI_FRAGMENT                            ParserErrors = 1117
	RNGP_URI_NOT_ABSOLUTE                        ParserErrors = 1118
	RNGP_VALUE_EMPTY                             ParserErrors = 1119
	RNGP_VALUE_NO_CONTENT                        ParserErrors = 1120
	RNGP_XMLNS_NAME                              ParserErrors = 1121
	RNGP_XML_NS                                  ParserErrors = 1122
	XPATH_EXPRESSION_OK                          ParserErrors = 1200
	XPATH_NUMBER_ERROR                           ParserErrors = 1201
	XPATH_UNFINISHED_LITERAL_ERROR               ParserErrors = 1202
	XPATH_START_LITERAL_ERROR                    ParserErrors = 1203
	XPATH_VARIABLE_REF_ERROR                     ParserErrors = 1204
	XPATH_UNDEF_VARIABLE_ERROR                   ParserErrors = 1205
	XPATH_INVALID_PREDICATE_ERROR                ParserErrors = 1206
	XPATH_EXPR_ERROR                             ParserErrors = 1207
	XPATH_UNCLOSED_ERROR                         ParserErrors = 1208
	XPATH_UNKNOWN_FUNC_ERROR                     ParserErrors = 1209
	XPATH_INVALID_OPERAND                        ParserErrors = 1210
	XPATH_INVALID_TYPE                           ParserErrors = 1211
	XPATH_INVALID_ARITY                          ParserErrors = 1212
	XPATH_INVALID_CTXT_SIZE                      ParserErrors = 1213
	XPATH_INVALID_CTXT_POSITION                  ParserErrors = 1214
	XPATH_MEMORY_ERROR                           ParserErrors = 1215
	XPTR_SYNTAX_ERROR                            ParserErrors = 1216
	XPTR_RESOURCE_ERROR                          ParserErrors = 1217
	XPTR_SUB_RESOURCE_ERROR                      ParserErrors = 1218
	XPATH_UNDEF_PREFIX_ERROR                     ParserErrors = 1219
	XPATH_ENCODING_ERROR                         ParserErrors = 1220
	XPATH_INVALID_CHAR_ERROR                     ParserErrors = 1221
	TREE_INVALID_HEX                             ParserErrors = 1300
	TREE_INVALID_DEC                             ParserErrors = 1301
	TREE_UNTERMINATED_ENTITY                     ParserErrors = 1302
	TREE_NOT_UTF8                                ParserErrors = 1303
	SAVE_NOT_UTF8                                ParserErrors = 1400
	SAVE_CHAR_INVALID                            ParserErrors = 1401
	SAVE_NO_DOCTYPE                              ParserErrors = 1402
	SAVE_UNKNOWN_ENCODING                        ParserErrors = 1403
	REGEXP_COMPILE_ERROR                         ParserErrors = 1450
	IO_UNKNOWN                                   ParserErrors = 1500
	IO_EACCES                                    ParserErrors = 1501
	IO_EAGAIN                                    ParserErrors = 1502
	IO_EBADF                                     ParserErrors = 1503
	IO_EBADMSG                                   ParserErrors = 1504
	IO_EBUSY                                     ParserErrors = 1505
	IO_ECANCELED                                 ParserErrors = 1506
	IO_ECHILD                                    ParserErrors = 1507
	IO_EDEADLK                                   ParserErrors = 1508
	IO_EDOM                                      ParserErrors = 1509
	IO_EEXIST                                    ParserErrors = 1510
	IO_EFAULT                                    ParserErrors = 1511
	IO_EFBIG                                     ParserErrors = 1512
	IO_EINPROGRESS                               ParserErrors = 1513
	IO_EINTR                                     ParserErrors = 1514
	IO_EINVAL                                    ParserErrors = 1515
	IO_EIO                                       ParserErrors = 1516
	IO_EISDIR                                    ParserErrors = 1517
	IO_EMFILE                                    ParserErrors = 1518
	IO_EMLINK                                    ParserErrors = 1519
	IO_EMSGSIZE                                  ParserErrors = 1520
	IO_ENAMETOOLONG                              ParserErrors = 1521
	IO_ENFILE                                    ParserErrors = 1522
	IO_ENODEV                                    ParserErrors = 1523
	IO_ENOENT                                    ParserErrors = 1524
	IO_ENOEXEC                                   ParserErrors = 1525
	IO_ENOLCK                                    ParserErrors = 1526
	IO_ENOMEM                                    ParserErrors = 1527
	IO_ENOSPC                                    ParserErrors = 1528
	IO_ENOSYS                                    ParserErrors = 1529
	IO_ENOTDIR                                   ParserErrors = 1530
	IO_ENOTEMPTY                                 ParserErrors = 1531
	IO_ENOTSUP                                   ParserErrors = 1532
	IO_ENOTTY                                    ParserErrors = 1533
	IO_ENXIO                                     ParserErrors = 1534
	IO_EPERM                                     ParserErrors = 1535
	IO_EPIPE                                     ParserErrors = 1536
	IO_ERANGE                                    ParserErrors = 1537
	IO_EROFS                                     ParserErrors = 1538
	IO_ESPIPE                                    ParserErrors = 1539
	IO_ESRCH                                     ParserErrors = 1540
	IO_ETIMEDOUT                                 ParserErrors = 1541
	IO_EXDEV                                     ParserErrors = 1542
	IO_NETWORK_ATTEMPT                           ParserErrors = 1543
	IO_ENCODER                                   ParserErrors = 1544
	IO_FLUSH                                     ParserErrors = 1545
	IO_WRITE                                     ParserErrors = 1546
	IO_NO_INPUT                                  ParserErrors = 1547
	IO_BUFFER_FULL                               ParserErrors = 1548
	IO_LOAD_ERROR                                ParserErrors = 1549
	IO_ENOTSOCK                                  ParserErrors = 1550
	IO_EISCONN                                   ParserErrors = 1551
	IO_ECONNREFUSED                              ParserErrors = 1552
	IO_ENETUNREACH                               ParserErrors = 1553
	IO_EADDRINUSE                                ParserErrors = 1554
	IO_EALREADY                                  ParserErrors = 1555
	IO_EAFNOSUPPORT                              ParserErrors = 1556
	IO_UNSUPPORTED_PROTOCOL                      ParserErrors = 1557
	XINCLUDE_RECURSION                           ParserErrors = 1600
	XINCLUDE_PARSE_VALUE                         ParserErrors = 1601
	XINCLUDE_ENTITY_DEF_MISMATCH                 ParserErrors = 1602
	XINCLUDE_NO_HREF                             ParserErrors = 1603
	XINCLUDE_NO_FALLBACK                         ParserErrors = 1604
	XINCLUDE_HREF_URI                            ParserErrors = 1605
	XINCLUDE_TEXT_FRAGMENT                       ParserErrors = 1606
	XINCLUDE_TEXT_DOCUMENT                       ParserErrors = 1607
	XINCLUDE_INVALID_CHAR                        ParserErrors = 1608
	XINCLUDE_BUILD_FAILED                        ParserErrors = 1609
	XINCLUDE_UNKNOWN_ENCODING                    ParserErrors = 1610
	XINCLUDE_MULTIPLE_ROOT                       ParserErrors = 1611
	XINCLUDE_XPTR_FAILED                         ParserErrors = 1612
	XINCLUDE_XPTR_RESULT                         ParserErrors = 1613
	XINCLUDE_INCLUDE_IN_INCLUDE                  ParserErrors = 1614
	XINCLUDE_FALLBACKS_IN_INCLUDE                ParserErrors = 1615
	XINCLUDE_FALLBACK_NOT_IN_INCLUDE             ParserErrors = 1616
	XINCLUDE_DEPRECATED_NS                       ParserErrors = 1617
	XINCLUDE_FRAGMENT_ID                         ParserErrors = 1618
	CATALOG_MISSING_ATTR                         ParserErrors = 1650
	CATALOG_ENTRY_BROKEN                         ParserErrors = 1651
	CATALOG_PREFER_VALUE                         ParserErrors = 1652
	CATALOG_NOT_CATALOG                          ParserErrors = 1653
	CATALOG_RECURSION                            ParserErrors = 1654
	SCHEMAP_PREFIX_UNDEFINED                     ParserErrors = 1700
	SCHEMAP_ATTRFORMDEFAULT_VALUE                ParserErrors = 1701
	SCHEMAP_ATTRGRP_NONAME_NOREF                 ParserErrors = 1702
	SCHEMAP_ATTR_NONAME_NOREF                    ParserErrors = 1703
	SCHEMAP_COMPLEXTYPE_NONAME_NOREF             ParserErrors = 1704
	SCHEMAP_ELEMFORMDEFAULT_VALUE                ParserErrors = 1705
	SCHEMAP_ELEM_NONAME_NOREF                    ParserErrors = 1706
	SCHEMAP_EXTENSION_NO_BASE                    ParserErrors = 1707
	SCHEMAP_FACET_NO_VALUE                       ParserErrors = 1708
	SCHEMAP_FAILED_BUILD_IMPORT                  ParserErrors = 1709
	SCHEMAP_GROUP_NONAME_NOREF                   ParserErrors = 1710
	SCHEMAP_IMPORT_NAMESPACE_NOT_URI             ParserErrors = 1711
	SCHEMAP_IMPORT_REDEFINE_NSNAME               ParserErrors = 1712
	SCHEMAP_IMPORT_SCHEMA_NOT_URI                ParserErrors = 1713
	SCHEMAP_INVALID_BOOLEAN                      ParserErrors = 1714
	SCHEMAP_INVALID_ENUM                         ParserErrors = 1715
	SCHEMAP_INVALID_FACET                        ParserErrors = 1716
	SCHEMAP_INVALID_FACET_VALUE                  ParserErrors = 1717
	SCHEMAP_INVALID_MAXOCCURS                    ParserErrors = 1718
	SCHEMAP_INVALID_MINOCCURS                    ParserErrors = 1719
	SCHEMAP_INVALID_REF_AND_SUBTYPE              ParserErrors = 1720
	SCHEMAP_INVALID_WHITE_SPACE                  ParserErrors = 1721
	SCHEMAP_NOATTR_NOREF                         ParserErrors = 1722
	SCHEMAP_NOTATION_NO_NAME                     ParserErrors = 1723
	SCHEMAP_NOTYPE_NOREF                         ParserErrors = 1724
	SCHEMAP_REF_AND_SUBTYPE                      ParserErrors = 1725
	SCHEMAP_RESTRICTION_NONAME_NOREF             ParserErrors = 1726
	SCHEMAP_SIMPLETYPE_NONAME                    ParserErrors = 1727
	SCHEMAP_TYPE_AND_SUBTYPE                     ParserErrors = 1728
	SCHEMAP_UNKNOWN_ALL_CHILD                    ParserErrors = 1729
	SCHEMAP_UNKNOWN_ANYATTRIBUTE_CHILD           ParserErrors = 1730
	SCHEMAP_UNKNOWN_ATTR_CHILD                   ParserErrors = 1731
	SCHEMAP_UNKNOWN_ATTRGRP_CHILD                ParserErrors = 1732
	SCHEMAP_UNKNOWN_ATTRIBUTE_GROUP              ParserErrors = 1733
	SCHEMAP_UNKNOWN_BASE_TYPE                    ParserErrors = 1734
	SCHEMAP_UNKNOWN_CHOICE_CHILD                 ParserErrors = 1735
	SCHEMAP_UNKNOWN_COMPLEXCONTENT_CHILD         ParserErrors = 1736
	SCHEMAP_UNKNOWN_COMPLEXTYPE_CHILD            ParserErrors = 1737
	SCHEMAP_UNKNOWN_ELEM_CHILD                   ParserErrors = 1738
	SCHEMAP_UNKNOWN_EXTENSION_CHILD              ParserErrors = 1739
	SCHEMAP_UNKNOWN_FACET_CHILD                  ParserErrors = 1740
	SCHEMAP_UNKNOWN_FACET_TYPE                   ParserErrors = 1741
	SCHEMAP_UNKNOWN_GROUP_CHILD                  ParserErrors = 1742
	SCHEMAP_UNKNOWN_IMPORT_CHILD                 ParserErrors = 1743
	SCHEMAP_UNKNOWN_LIST_CHILD                   ParserErrors = 1744
	SCHEMAP_UNKNOWN_NOTATION_CHILD               ParserErrors = 1745
	SCHEMAP_UNKNOWN_PROCESSCONTENT_CHILD         ParserErrors = 1746
	SCHEMAP_UNKNOWN_REF                          ParserErrors = 1747
	SCHEMAP_UNKNOWN_RESTRICTION_CHILD            ParserErrors = 1748
	SCHEMAP_UNKNOWN_SCHEMAS_CHILD                ParserErrors = 1749
	SCHEMAP_UNKNOWN_SEQUENCE_CHILD               ParserErrors = 1750
	SCHEMAP_UNKNOWN_SIMPLECONTENT_CHILD          ParserErrors = 1751
	SCHEMAP_UNKNOWN_SIMPLETYPE_CHILD             ParserErrors = 1752
	SCHEMAP_UNKNOWN_TYPE                         ParserErrors = 1753
	SCHEMAP_UNKNOWN_UNION_CHILD                  ParserErrors = 1754
	SCHEMAP_ELEM_DEFAULT_FIXED                   ParserErrors = 1755
	SCHEMAP_REGEXP_INVALID                       ParserErrors = 1756
	SCHEMAP_FAILED_LOAD                          ParserErrors = 1757
	SCHEMAP_NOTHING_TO_PARSE                     ParserErrors = 1758
	SCHEMAP_NOROOT                               ParserErrors = 1759
	SCHEMAP_REDEFINED_GROUP                      ParserErrors = 1760
	SCHEMAP_REDEFINED_TYPE                       ParserErrors = 1761
	SCHEMAP_REDEFINED_ELEMENT                    ParserErrors = 1762
	SCHEMAP_REDEFINED_ATTRGROUP                  ParserErrors = 1763
	SCHEMAP_REDEFINED_ATTR                       ParserErrors = 1764
	SCHEMAP_REDEFINED_NOTATION                   ParserErrors = 1765
	SCHEMAP_FAILED_PARSE                         ParserErrors = 1766
	SCHEMAP_UNKNOWN_PREFIX                       ParserErrors = 1767
	SCHEMAP_DEF_AND_PREFIX                       ParserErrors = 1768
	SCHEMAP_UNKNOWN_INCLUDE_CHILD                ParserErrors = 1769
	SCHEMAP_INCLUDE_SCHEMA_NOT_URI               ParserErrors = 1770
	SCHEMAP_INCLUDE_SCHEMA_NO_URI                ParserErrors = 1771
	SCHEMAP_NOT_SCHEMA                           ParserErrors = 1772
	SCHEMAP_UNKNOWN_MEMBER_TYPE                  ParserErrors = 1773
	SCHEMAP_INVALID_ATTR_USE                     ParserErrors = 1774
	SCHEMAP_RECURSIVE                            ParserErrors = 1775
	SCHEMAP_SUPERNUMEROUS_LIST_ITEM_TYPE         ParserErrors = 1776
	SCHEMAP_INVALID_ATTR_COMBINATION             ParserErrors = 1777
	SCHEMAP_INVALID_ATTR_INLINE_COMBINATION      ParserErrors = 1778
	SCHEMAP_MISSING_SIMPLETYPE_CHILD             ParserErrors = 1779
	SCHEMAP_INVALID_ATTR_NAME                    ParserErrors = 1780
	SCHEMAP_REF_AND_CONTENT                      ParserErrors = 1781
	SCHEMAP_CT_PROPS_CORRECT_1                   ParserErrors = 1782
	SCHEMAP_CT_PROPS_CORRECT_2                   ParserErrors = 1783
	SCHEMAP_CT_PROPS_CORRECT_3                   ParserErrors = 1784
	SCHEMAP_CT_PROPS_CORRECT_4                   ParserErrors = 1785
	SCHEMAP_CT_PROPS_CORRECT_5                   ParserErrors = 1786
	SCHEMAP_DERIVATION_OK_RESTRICTION_1          ParserErrors = 1787
	SCHEMAP_DERIVATION_OK_RESTRICTION_2_1_1      ParserErrors = 1788
	SCHEMAP_DERIVATION_OK_RESTRICTION_2_1_2      ParserErrors = 1789
	SCHEMAP_DERIVATION_OK_RESTRICTION_2_2        ParserErrors = 1790
	SCHEMAP_DERIVATION_OK_RESTRICTION_3          ParserErrors = 1791
	SCHEMAP_WILDCARD_INVALID_NS_MEMBER           ParserErrors = 1792
	SCHEMAP_INTERSECTION_NOT_EXPRESSIBLE         ParserErrors = 1793
	SCHEMAP_UNION_NOT_EXPRESSIBLE                ParserErrors = 1794
	SCHEMAP_SRC_IMPORT_3_1                       ParserErrors = 1795
	SCHEMAP_SRC_IMPORT_3_2                       ParserErrors = 1796
	SCHEMAP_DERIVATION_OK_RESTRICTION_4_1        ParserErrors = 1797
	SCHEMAP_DERIVATION_OK_RESTRICTION_4_2        ParserErrors = 1798
	SCHEMAP_DERIVATION_OK_RESTRICTION_4_3        ParserErrors = 1799
	SCHEMAP_COS_CT_EXTENDS_1_3                   ParserErrors = 1800
	SCHEMAV_NOROOT                               ParserErrors = 1801
	SCHEMAV_UNDECLAREDELEM                       ParserErrors = 1802
	SCHEMAV_NOTTOPLEVEL                          ParserErrors = 1803
	SCHEMAV_MISSING                              ParserErrors = 1804
	SCHEMAV_WRONGELEM                            ParserErrors = 1805
	SCHEMAV_NOTYPE                               ParserErrors = 1806
	SCHEMAV_NOROLLBACK                           ParserErrors = 1807
	SCHEMAV_ISABSTRACT                           ParserErrors = 1808
	SCHEMAV_NOTEMPTY                             ParserErrors = 1809
	SCHEMAV_ELEMCONT                             ParserErrors = 1810
	SCHEMAV_HAVEDEFAULT                          ParserErrors = 1811
	SCHEMAV_NOTNILLABLE                          ParserErrors = 1812
	SCHEMAV_EXTRACONTENT                         ParserErrors = 1813
	SCHEMAV_INVALIDATTR                          ParserErrors = 1814
	SCHEMAV_INVALIDELEM                          ParserErrors = 1815
	SCHEMAV_NOTDETERMINIST                       ParserErrors = 1816
	SCHEMAV_CONSTRUCT                            ParserErrors = 1817
	SCHEMAV_INTERNAL                             ParserErrors = 1818
	SCHEMAV_NOTSIMPLE                            ParserErrors = 1819
	SCHEMAV_ATTRUNKNOWN                          ParserErrors = 1820
	SCHEMAV_ATTRINVALID                          ParserErrors = 1821
	SCHEMAV_VALUE                                ParserErrors = 1822
	SCHEMAV_FACET                                ParserErrors = 1823
	SCHEMAV_CVC_DATATYPE_VALID_1_2_1             ParserErrors = 1824
	SCHEMAV_CVC_DATATYPE_VALID_1_2_2             ParserErrors = 1825
	SCHEMAV_CVC_DATATYPE_VALID_1_2_3             ParserErrors = 1826
	SCHEMAV_CVC_TYPE_3_1_1                       ParserErrors = 1827
	SCHEMAV_CVC_TYPE_3_1_2                       ParserErrors = 1828
	SCHEMAV_CVC_FACET_VALID                      ParserErrors = 1829
	SCHEMAV_CVC_LENGTH_VALID                     ParserErrors = 1830
	SCHEMAV_CVC_MINLENGTH_VALID                  ParserErrors = 1831
	SCHEMAV_CVC_MAXLENGTH_VALID                  ParserErrors = 1832
	SCHEMAV_CVC_MININCLUSIVE_VALID               ParserErrors = 1833
	SCHEMAV_CVC_MAXINCLUSIVE_VALID               ParserErrors = 1834
	SCHEMAV_CVC_MINEXCLUSIVE_VALID               ParserErrors = 1835
	SCHEMAV_CVC_MAXEXCLUSIVE_VALID               ParserErrors = 1836
	SCHEMAV_CVC_TOTALDIGITS_VALID                ParserErrors = 1837
	SCHEMAV_CVC_FRACTIONDIGITS_VALID             ParserErrors = 1838
	SCHEMAV_CVC_PATTERN_VALID                    ParserErrors = 1839
	SCHEMAV_CVC_ENUMERATION_VALID                ParserErrors = 1840
	SCHEMAV_CVC_COMPLEX_TYPE_2_1                 ParserErrors = 1841
	SCHEMAV_CVC_COMPLEX_TYPE_2_2                 ParserErrors = 1842
	SCHEMAV_CVC_COMPLEX_TYPE_2_3                 ParserErrors = 1843
	SCHEMAV_CVC_COMPLEX_TYPE_2_4                 ParserErrors = 1844
	SCHEMAV_CVC_ELT_1                            ParserErrors = 1845
	SCHEMAV_CVC_ELT_2                            ParserErrors = 1846
	SCHEMAV_CVC_ELT_3_1                          ParserErrors = 1847
	SCHEMAV_CVC_ELT_3_2_1                        ParserErrors = 1848
	SCHEMAV_CVC_ELT_3_2_2                        ParserErrors = 1849
	SCHEMAV_CVC_ELT_4_1                          ParserErrors = 1850
	SCHEMAV_CVC_ELT_4_2                          ParserErrors = 1851
	SCHEMAV_CVC_ELT_4_3                          ParserErrors = 1852
	SCHEMAV_CVC_ELT_5_1_1                        ParserErrors = 1853
	SCHEMAV_CVC_ELT_5_1_2                        ParserErrors = 1854
	SCHEMAV_CVC_ELT_5_2_1                        ParserErrors = 1855
	SCHEMAV_CVC_ELT_5_2_2_1                      ParserErrors = 1856
	SCHEMAV_CVC_ELT_5_2_2_2_1                    ParserErrors = 1857
	SCHEMAV_CVC_ELT_5_2_2_2_2                    ParserErrors = 1858
	SCHEMAV_CVC_ELT_6                            ParserErrors = 1859
	SCHEMAV_CVC_ELT_7                            ParserErrors = 1860
	SCHEMAV_CVC_ATTRIBUTE_1                      ParserErrors = 1861
	SCHEMAV_CVC_ATTRIBUTE_2                      ParserErrors = 1862
	SCHEMAV_CVC_ATTRIBUTE_3                      ParserErrors = 1863
	SCHEMAV_CVC_ATTRIBUTE_4                      ParserErrors = 1864
	SCHEMAV_CVC_COMPLEX_TYPE_3_1                 ParserErrors = 1865
	SCHEMAV_CVC_COMPLEX_TYPE_3_2_1               ParserErrors = 1866
	SCHEMAV_CVC_COMPLEX_TYPE_3_2_2               ParserErrors = 1867
	SCHEMAV_CVC_COMPLEX_TYPE_4                   ParserErrors = 1868
	SCHEMAV_CVC_COMPLEX_TYPE_5_1                 ParserErrors = 1869
	SCHEMAV_CVC_COMPLEX_TYPE_5_2                 ParserErrors = 1870
	SCHEMAV_ELEMENT_CONTENT                      ParserErrors = 1871
	SCHEMAV_DOCUMENT_ELEMENT_MISSING             ParserErrors = 1872
	SCHEMAV_CVC_COMPLEX_TYPE_1                   ParserErrors = 1873
	SCHEMAV_CVC_AU                               ParserErrors = 1874
	SCHEMAV_CVC_TYPE_1                           ParserErrors = 1875
	SCHEMAV_CVC_TYPE_2                           ParserErrors = 1876
	SCHEMAV_CVC_IDC                              ParserErrors = 1877
	SCHEMAV_CVC_WILDCARD                         ParserErrors = 1878
	SCHEMAV_MISC                                 ParserErrors = 1879
	XPTR_UNKNOWN_SCHEME                          ParserErrors = 1900
	XPTR_CHILDSEQ_START                          ParserErrors = 1901
	XPTR_EVAL_FAILED                             ParserErrors = 1902
	XPTR_EXTRA_OBJECTS                           ParserErrors = 1903
	C14N_CREATE_CTXT                             ParserErrors = 1950
	C14N_REQUIRES_UTF8                           ParserErrors = 1951
	C14N_CREATE_STACK                            ParserErrors = 1952
	C14N_INVALID_NODE                            ParserErrors = 1953
	C14N_UNKNOW_NODE                             ParserErrors = 1954
	C14N_RELATIVE_NAMESPACE                      ParserErrors = 1955
	FTP_PASV_ANSWER                              ParserErrors = 2000
	FTP_EPSV_ANSWER                              ParserErrors = 2001
	FTP_ACCNT                                    ParserErrors = 2002
	FTP_URL_SYNTAX                               ParserErrors = 2003
	HTTP_URL_SYNTAX                              ParserErrors = 2020
	HTTP_USE_IP                                  ParserErrors = 2021
	HTTP_UNKNOWN_HOST                            ParserErrors = 2022
	SCHEMAP_SRC_SIMPLE_TYPE_1                    ParserErrors = 3000
	SCHEMAP_SRC_SIMPLE_TYPE_2                    ParserErrors = 3001
	SCHEMAP_SRC_SIMPLE_TYPE_3                    ParserErrors = 3002
	SCHEMAP_SRC_SIMPLE_TYPE_4                    ParserErrors = 3003
	SCHEMAP_SRC_RESOLVE                          ParserErrors = 3004
	SCHEMAP_SRC_RESTRICTION_BASE_OR_SIMPLETYPE   ParserErrors = 3005
	SCHEMAP_SRC_LIST_ITEMTYPE_OR_SIMPLETYPE      ParserErrors = 3006
	SCHEMAP_SRC_UNION_MEMBERTYPES_OR_SIMPLETYPES ParserErrors = 3007
	SCHEMAP_ST_PROPS_CORRECT_1                   ParserErrors = 3008
	SCHEMAP_ST_PROPS_CORRECT_2                   ParserErrors = 3009
	SCHEMAP_ST_PROPS_CORRECT_3                   ParserErrors = 3010
	SCHEMAP_COS_ST_RESTRICTS_1_1                 ParserErrors = 3011
	SCHEMAP_COS_ST_RESTRICTS_1_2                 ParserErrors = 3012
	SCHEMAP_COS_ST_RESTRICTS_1_3_1               ParserErrors = 3013
	SCHEMAP_COS_ST_RESTRICTS_1_3_2               ParserErrors = 3014
	SCHEMAP_COS_ST_RESTRICTS_2_1                 ParserErrors = 3015
	SCHEMAP_COS_ST_RESTRICTS_2_3_1_1             ParserErrors = 3016
	SCHEMAP_COS_ST_RESTRICTS_2_3_1_2             ParserErrors = 3017
	SCHEMAP_COS_ST_RESTRICTS_2_3_2_1             ParserErrors = 3018
	SCHEMAP_COS_ST_RESTRICTS_2_3_2_2             ParserErrors = 3019
	SCHEMAP_COS_ST_RESTRICTS_2_3_2_3             ParserErrors = 3020
	SCHEMAP_COS_ST_RESTRICTS_2_3_2_4             ParserErrors = 3021
	SCHEMAP_COS_ST_RESTRICTS_2_3_2_5             ParserErrors = 3022
	SCHEMAP_COS_ST_RESTRICTS_3_1                 ParserErrors = 3023
	SCHEMAP_COS_ST_RESTRICTS_3_3_1               ParserErrors = 3024
	SCHEMAP_COS_ST_RESTRICTS_3_3_1_2             ParserErrors = 3025
	SCHEMAP_COS_ST_RESTRICTS_3_3_2_2             ParserErrors = 3026
	SCHEMAP_COS_ST_RESTRICTS_3_3_2_1             ParserErrors = 3027
	SCHEMAP_COS_ST_RESTRICTS_3_3_2_3             ParserErrors = 3028
	SCHEMAP_COS_ST_RESTRICTS_3_3_2_4             ParserErrors = 3029
	SCHEMAP_COS_ST_RESTRICTS_3_3_2_5             ParserErrors = 3030
	SCHEMAP_COS_ST_DERIVED_OK_2_1                ParserErrors = 3031
	SCHEMAP_COS_ST_DERIVED_OK_2_2                ParserErrors = 3032
	SCHEMAP_S4S_ELEM_NOT_ALLOWED                 ParserErrors = 3033
	SCHEMAP_S4S_ELEM_MISSING                     ParserErrors = 3034
	SCHEMAP_S4S_ATTR_NOT_ALLOWED                 ParserErrors = 3035
	SCHEMAP_S4S_ATTR_MISSING                     ParserErrors = 3036
	SCHEMAP_S4S_ATTR_INVALID_VALUE               ParserErrors = 3037
	SCHEMAP_SRC_ELEMENT_1                        ParserErrors = 3038
	SCHEMAP_SRC_ELEMENT_2_1                      ParserErrors = 3039
	SCHEMAP_SRC_ELEMENT_2_2                      ParserErrors = 3040
	SCHEMAP_SRC_ELEMENT_3                        ParserErrors = 3041
	SCHEMAP_P_PROPS_CORRECT_1                    ParserErrors = 3042
	SCHEMAP_P_PROPS_CORRECT_2_1                  ParserErrors = 3043
	SCHEMAP_P_PROPS_CORRECT_2_2                  ParserErrors = 3044
	SCHEMAP_E_PROPS_CORRECT_2                    ParserErrors = 3045
	SCHEMAP_E_PROPS_CORRECT_3                    ParserErrors = 3046
	SCHEMAP_E_PROPS_CORRECT_4                    ParserErrors = 3047
	SCHEMAP_E_PROPS_CORRECT_5                    ParserErrors = 3048
	SCHEMAP_E_PROPS_CORRECT_6                    ParserErrors = 3049
	SCHEMAP_SRC_INCLUDE                          ParserErrors = 3050
	SCHEMAP_SRC_ATTRIBUTE_1                      ParserErrors = 3051
	SCHEMAP_SRC_ATTRIBUTE_2                      ParserErrors = 3052
	SCHEMAP_SRC_ATTRIBUTE_3_1                    ParserErrors = 3053
	SCHEMAP_SRC_ATTRIBUTE_3_2                    ParserErrors = 3054
	SCHEMAP_SRC_ATTRIBUTE_4                      ParserErrors = 3055
	SCHEMAP_NO_XMLNS                             ParserErrors = 3056
	SCHEMAP_NO_XSI                               ParserErrors = 3057
	SCHEMAP_COS_VALID_DEFAULT_1                  ParserErrors = 3058
	SCHEMAP_COS_VALID_DEFAULT_2_1                ParserErrors = 3059
	SCHEMAP_COS_VALID_DEFAULT_2_2_1              ParserErrors = 3060
	SCHEMAP_COS_VALID_DEFAULT_2_2_2              ParserErrors = 3061
	SCHEMAP_CVC_SIMPLE_TYPE                      ParserErrors = 3062
	SCHEMAP_COS_CT_EXTENDS_1_1                   ParserErrors = 3063
	SCHEMAP_SRC_IMPORT_1_1                       ParserErrors = 3064
	SCHEMAP_SRC_IMPORT_1_2                       ParserErrors = 3065
	SCHEMAP_SRC_IMPORT_2                         ParserErrors = 3066
	SCHEMAP_SRC_IMPORT_2_1                       ParserErrors = 3067
	SCHEMAP_SRC_IMPORT_2_2                       ParserErrors = 3068
	SCHEMAP_INTERNAL                             ParserErrors = 3069
	SCHEMAP_NOT_DETERMINISTIC                    ParserErrors = 3070
	SCHEMAP_SRC_ATTRIBUTE_GROUP_1                ParserErrors = 3071
	SCHEMAP_SRC_ATTRIBUTE_GROUP_2                ParserErrors = 3072
	SCHEMAP_SRC_ATTRIBUTE_GROUP_3                ParserErrors = 3073
	SCHEMAP_MG_PROPS_CORRECT_1                   ParserErrors = 3074
	SCHEMAP_MG_PROPS_CORRECT_2                   ParserErrors = 3075
	SCHEMAP_SRC_CT_1                             ParserErrors = 3076
	SCHEMAP_DERIVATION_OK_RESTRICTION_2_1_3      ParserErrors = 3077
	SCHEMAP_AU_PROPS_CORRECT_2                   ParserErrors = 3078
	SCHEMAP_A_PROPS_CORRECT_2                    ParserErrors = 3079
	SCHEMAP_C_PROPS_CORRECT                      ParserErrors = 3080
	SCHEMAP_SRC_REDEFINE                         ParserErrors = 3081
	SCHEMAP_SRC_IMPORT                           ParserErrors = 3082
	SCHEMAP_WARN_SKIP_SCHEMA                     ParserErrors = 3083
	SCHEMAP_WARN_UNLOCATED_SCHEMA                ParserErrors = 3084
	SCHEMAP_WARN_ATTR_REDECL_PROH                ParserErrors = 3085
	SCHEMAP_WARN_ATTR_POINTLESS_PROH             ParserErrors = 3086
	SCHEMAP_AG_PROPS_CORRECT                     ParserErrors = 3087
	SCHEMAP_COS_CT_EXTENDS_1_2                   ParserErrors = 3088
	SCHEMAP_AU_PROPS_CORRECT                     ParserErrors = 3089
	SCHEMAP_A_PROPS_CORRECT_3                    ParserErrors = 3090
	SCHEMAP_COS_ALL_LIMITED                      ParserErrors = 3091
	SCHEMATRONV_ASSERT                           ParserErrors = 4000
	SCHEMATRONV_REPORT                           ParserErrors = 4001
	MODULE_OPEN                                  ParserErrors = 4900
	MODULE_CLOSE                                 ParserErrors = 4901
	CHECK_FOUND_ELEMENT                          ParserErrors = 5000
	CHECK_FOUND_ATTRIBUTE                        ParserErrors = 5001
	CHECK_FOUND_TEXT                             ParserErrors = 5002
	CHECK_FOUND_CDATA                            ParserErrors = 5003
	CHECK_FOUND_ENTITYREF                        ParserErrors = 5004
	CHECK_FOUND_ENTITY                           ParserErrors = 5005
	CHECK_FOUND_PI                               ParserErrors = 5006
	CHECK_FOUND_COMMENT                          ParserErrors = 5007
	CHECK_FOUND_DOCTYPE                          ParserErrors = 5008
	CHECK_FOUND_FRAGMENT                         ParserErrors = 5009
	CHECK_FOUND_NOTATION                         ParserErrors = 5010
	CHECK_UNKNOWN_NODE                           ParserErrors = 5011
	CHECK_ENTITY_TYPE                            ParserErrors = 5012
	CHECK_NO_PARENT                              ParserErrors = 5013
	CHECK_NO_DOC                                 ParserErrors = 5014
	CHECK_NO_NAME                                ParserErrors = 5015
	CHECK_NO_ELEM                                ParserErrors = 5016
	CHECK_WRONG_DOC                              ParserErrors = 5017
	CHECK_NO_PREV                                ParserErrors = 5018
	CHECK_WRONG_PREV                             ParserErrors = 5019
	CHECK_NO_NEXT                                ParserErrors = 5020
	CHECK_WRONG_NEXT                             ParserErrors = 5021
	CHECK_NOT_DTD                                ParserErrors = 5022
	CHECK_NOT_ATTR                               ParserErrors = 5023
	CHECK_NOT_ATTR_DECL                          ParserErrors = 5024
	CHECK_NOT_ELEM_DECL                          ParserErrors = 5025
	CHECK_NOT_ENTITY_DECL                        ParserErrors = 5026
	CHECK_NOT_NS_DECL                            ParserErrors = 5027
	CHECK_NO_HREF                                ParserErrors = 5028
	CHECK_WRONG_PARENT                           ParserErrors = 5029
	CHECK_NS_SCOPE                               ParserErrors = 5030
	CHECK_NS_ANCESTOR                            ParserErrors = 5031
	CHECK_NOT_UTF8                               ParserErrors = 5032
	CHECK_NO_DICT                                ParserErrors = 5033
	CHECK_NOT_NCNAME                             ParserErrors = 5034
	CHECK_OUTSIDE_DICT                           ParserErrors = 5035
	CHECK_WRONG_NAME                             ParserErrors = 5036
	CHECK_NAME_NOT_NULL                          ParserErrors = 5037
	I18N_NO_NAME                                 ParserErrors = 6000
	I18N_NO_HANDLER                              ParserErrors = 6001
	I18N_EXCESS_HANDLER                          ParserErrors = 6002
	I18N_CONV_FAILED                             ParserErrors = 6003
	I18N_NO_OUTPUT                               ParserErrors = 6004
	BUF_OVERFLOW                                 ParserErrors = 7000
)

// llgo:type C
type GenericErrorFunc func(__llgo_arg_0 c.Pointer, __llgo_arg_1 *c.Char, __llgo_va_list ...interface{})

// llgo:type C
type StructuredErrorFunc func(c.Pointer, *Error)

//go:linkname X__xmlLastError C.__xmlLastError
func X__xmlLastError() *Error

//go:linkname X__xmlGenericError C.__xmlGenericError
func X__xmlGenericError() GenericErrorFunc

//go:linkname X__xmlGenericErrorContext C.__xmlGenericErrorContext
func X__xmlGenericErrorContext() *c.Pointer

//go:linkname X__xmlStructuredError C.__xmlStructuredError
func X__xmlStructuredError() StructuredErrorFunc

//go:linkname X__xmlStructuredErrorContext C.__xmlStructuredErrorContext
func X__xmlStructuredErrorContext() *c.Pointer

/*
 * Use the following function to reset the two global variables
 * xmlGenericError and xmlGenericErrorContext.
 */
//go:linkname SetGenericErrorFunc C.xmlSetGenericErrorFunc
func SetGenericErrorFunc(ctx c.Pointer, handler GenericErrorFunc)

//go:linkname ThrDefSetGenericErrorFunc C.xmlThrDefSetGenericErrorFunc
func ThrDefSetGenericErrorFunc(ctx c.Pointer, handler GenericErrorFunc)

//go:linkname InitGenericErrorDefaultFunc C.initGenericErrorDefaultFunc
func InitGenericErrorDefaultFunc(handler GenericErrorFunc)

//go:linkname SetStructuredErrorFunc C.xmlSetStructuredErrorFunc
func SetStructuredErrorFunc(ctx c.Pointer, handler StructuredErrorFunc)

//go:linkname ThrDefSetStructuredErrorFunc C.xmlThrDefSetStructuredErrorFunc
func ThrDefSetStructuredErrorFunc(ctx c.Pointer, handler StructuredErrorFunc)

/*
 * Default message routines used by SAX and Valid context for error
 * and warning reporting.
 */
//go:linkname ParserError C.xmlParserError
func ParserError(ctx c.Pointer, msg *c.Char, __llgo_va_list ...interface{})

//go:linkname ParserWarning C.xmlParserWarning
func ParserWarning(ctx c.Pointer, msg *c.Char, __llgo_va_list ...interface{})

//go:linkname ParserValidityError C.xmlParserValidityError
func ParserValidityError(ctx c.Pointer, msg *c.Char, __llgo_va_list ...interface{})

//go:linkname ParserValidityWarning C.xmlParserValidityWarning
func ParserValidityWarning(ctx c.Pointer, msg *c.Char, __llgo_va_list ...interface{})

/** DOC_ENABLE */
// llgo:link (*X_xmlParserInput).ParserPrintFileInfo C.xmlParserPrintFileInfo
func (recv_ *X_xmlParserInput) ParserPrintFileInfo() {
}

// llgo:link (*X_xmlParserInput).ParserPrintFileContext C.xmlParserPrintFileContext
func (recv_ *X_xmlParserInput) ParserPrintFileContext() {
}

// llgo:link (*Error).FormatError C.xmlFormatError
func (recv_ *Error) FormatError(channel GenericErrorFunc, data c.Pointer) {
}

/*
 * Extended error information routines
 */
//go:linkname GetLastError C.xmlGetLastError
func GetLastError() *Error

//go:linkname ResetLastError C.xmlResetLastError
func ResetLastError()

//go:linkname CtxtGetLastError C.xmlCtxtGetLastError
func CtxtGetLastError(ctx c.Pointer) *Error

//go:linkname CtxtResetLastError C.xmlCtxtResetLastError
func CtxtResetLastError(ctx c.Pointer)

//go:linkname ResetError C.xmlResetError
func ResetError(err ErrorPtr)

// llgo:link (*Error).CopyError C.xmlCopyError
func (recv_ *Error) CopyError(to ErrorPtr) c.Int {
	return 0
}
