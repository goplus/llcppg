void foo();
namespace a {
void foo();
}
namespace a {
namespace b {
void foo();
}
} // namespace a
class Foo {
  public:
    void foo();
};
namespace a {
class Foo {
  public:
    void foo();
};
} // namespace a