package hash

type Key interface {
	Hash() Val
	Equals(Key) bool
	String() string
}
