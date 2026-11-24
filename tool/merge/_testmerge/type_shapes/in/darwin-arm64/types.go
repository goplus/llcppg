package example

type Slot[T any] struct {
	Value T
}

func (s Slot[T]) Get() T { return s.Value }

type Pair[K, V any] struct {
	Key   K
	Value V
}

func (p Pair[K, V]) Get() V { return p.Value }

func (Slot[T]) Reserve(size int) {}

type (
	Token    []byte
	Counters map[string]int
	Handler  func(...int)
	Frame    [1 << 3]byte
	Label    [len("label")]byte
	Row      Pair[int, string]
)

type Reader interface {
	Read(...string) ([]byte, [4]byte)
	Current() *int
	Entries() (Slot[int], Pair[int, string])
	Stats() (struct{ Hits int }, map[string]int)
}
