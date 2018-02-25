package hash

// Key is an interface for values used a keys in fmap and set data structures.
type Key interface {
	Hash() Val
	Equals(Key) bool
	String() string
}
