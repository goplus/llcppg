template <class _Category, class _Tp, class _Distance, class _Pointer = _Tp*, class _Reference = _Tp&>
struct iterator {
  typedef _Tp value_type;
  typedef _Distance difference_type;
  typedef _Pointer pointer;
  typedef _Reference reference;
  typedef _Category iterator_category;
};

template <class _Derived, class _Category, class _Tp, class _Distance, class _Pointer, class _Reference>
using __iterator_base = iterator<_Category, _Tp, _Distance, _Pointer, _Reference>;

typedef int output_iterator_tag;

template <class _Container>
class back_insert_iterator
    : public __iterator_base<back_insert_iterator<_Container>, output_iterator_tag, void, void, void, void> {
protected:
  _Container* container;
};
