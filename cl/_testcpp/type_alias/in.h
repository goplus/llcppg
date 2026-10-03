namespace std {
    template <typename CharT, typename Traits, typename Allocator>
    class basic_string {
    };

    using ptrdiff_t = decltype(static_cast<int*>(nullptr) - static_cast<int*>(nullptr));
    using string = basic_string<char, int, int>;
}
