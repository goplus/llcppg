class Foo
{
public:
    // constructor -> init
    Foo(const char *filename);
    // redefined name -> init__1
    Foo(const char *buffer, long buffer_size);
    // destructor -> Dispose
    ~Foo();
    // with custom name in llcppg.cfg/symMap
    int ParseBar() const;
    // not in llcppg.cfg/symMap,generate automatically
    const char *Get(const char *section, const char *name,
                            const char *default_value) const;
private:
    // not in output symbol table
    static const char *MakeBar(const char *section, const char *name);
};

// method out of class decl
bool Foo::HasBar();