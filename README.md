Functional Collections Library for the Go Language
==================================================

This library implements several data structures (also called collections)
that behave functionally. In this case "functional" means persistent and
immutable. Imutable means that every method that modifies a collection, for
instance inserting or deleting entries, returns a new instance of the
collection with the modification. Persistent means that each new instance of
a collection share all unmodified potions of the collection with the previous
version.

The following are currently implemented:

* A functional Map, called _fmap_, which uses a [HAMT][1] internally.
* A functional sorted Map, called _sorted_map_, which uses a [standard Red-Black Tree][2] internally.

I am planning on implementing:

* A functional Set, called _set_, which uses a [HAMT][1] internally.
* A functional Sorted Set, called _sorted_set_, which uses a [standard Red-Black Tree][2] internally.
* A functional Vector, called _vector_, which uses a clojure-like implementation of the same name.

[1]:https://en.wikipedia.org/wiki/Hash_array_mapped_trie
[2]:https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
