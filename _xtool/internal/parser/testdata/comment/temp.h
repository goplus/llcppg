// not read doc 1
void foo1();
/* not read doc 2 */
void foo2();
/// doc
void foo3();
/** doc */
void foo4();
/*! doc */
void foo5();
/// doc 1
/// doc 2
void foo6();
/*! doc 1 */
/*! doc 2 */
void foo7();
/** doc 1 */
/** doc 1 */
void foo8();
/**
 * doc 1
 * doc 2
 */
void foo9();
struct Foo {
    /// doc
    int x;
    int y; ///< comment
    /**
     * field doc (parse ignore with comment in same cursor)
     */
    int z; /*!< comment */
};
class Doc {
  public:
    /**
     * static field doc
     */
    static int x;
    static int y; /*!< static field comment */
    /**
     * field doc
     */
    int a;
    int b; ///< field comment
    /**
     * method doc
     */
    void Foo();

  protected:
    int value; /*!< protected field comment */
};