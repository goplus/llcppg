class A {
  public:
    int a;
    int b;
};

class B {
  public:
    static int a;
    int b;
    float foo(int a, double b);
    void vafoo(int a, ...);

  private:
    static void bar();

  protected:
    void bar2();
};

class C {
  public:
    C();
    explicit C();
    ~C();
    static inline void foo();
};

class Base {
  public:
    Base();
    virtual ~Base();
    virtual void foo();
};
class Derived : public Base {
  public:
    Derived();
    ~Derived() override;
    void foo() override;
};

namespace NSA {
class Foo {};
} // namespace NSA
