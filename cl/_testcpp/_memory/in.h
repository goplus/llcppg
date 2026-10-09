typedef unsigned long size_t;

template <class _Tp>
class weak_ptr;

template <class _Tp>
class shared_ptr {
public:
  typedef weak_ptr<_Tp> weak_type;
  typedef _Tp element_type;

private:
  element_type* __ptr_;
  size_t* __cntrl_;
};

template <class _Tp>
class weak_ptr {
public:
  typedef _Tp element_type;

private:
  element_type* __ptr_;
  size_t* __cntrl_;
};
