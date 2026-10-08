template <class _Category, class _Tp, class _Distance, class _Pointer = _Tp*, class _Reference = _Tp&>
struct iterator {
  typedef _Tp value_type;
  typedef _Distance difference_type;
  typedef _Pointer pointer;
  typedef _Reference reference;
  typedef _Category iterator_category;
};

template <class _Category, class _Tp, class _Distance, class _Pointer, class _Reference>
using __iterator_alias = iterator<_Category, _Tp, _Distance, _Pointer, _Reference>;

template <class _Derived, class _Category, class _Tp, class _Distance, class _Pointer, class _Reference>
using __iterator_ign = iterator<_Category, _Tp, _Distance, _Pointer, _Reference>;

template <class _Derived, class _Category, class _Tp, class _Distance, class _Pointer, class _Reference>
using __iterator_base = iterator<_Category, _Tp, _Distance, _Pointer, _Reference>;

typedef int output_iterator_tag;

template <class _Container>
class back_insert_iterator
    : public __iterator_base<back_insert_iterator<_Container>, output_iterator_tag, void, void, void, void> {
protected:
  _Container* container;
};

template <typename _Ip>
struct iterator_traits {
  using iterator_category = typename _Ip::iterator_category;
  using value_type        = typename _Ip::value_type;
  using difference_type   = typename _Ip::difference_type;
  using reference         = typename _Ip::reference;
  typedef value_type *pointer;
};

template <class _Iter>
class reverse_iterator
    : public __iterator_ign<reverse_iterator<_Iter>,
                             typename iterator_traits<_Iter>::iterator_category,
                             typename iterator_traits<_Iter>::value_type,
                             typename iterator_traits<_Iter>::difference_type,
                             typename iterator_traits<_Iter>::pointer,
                             typename iterator_traits<_Iter>::reference> {
protected:
  _Iter current;
};
