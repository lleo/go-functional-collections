(temporary) Questions about the API design
==========================================

Note: several of the API choices were meant to follow the sync.Map design. For
example, LoadOrStore & Range. Other choices were to allow chained calls.

Example of chained calls:

    var m = fmap.New().
      Put(hash.StringKey("a"), 1).
      Put(hash.StringKey("b"), 2).
      Put(hash.StringKey("c"), 3)

0. Anything you would like to comment on.
1. Should the ``NumEntries()`` method on the collections be renamed? Possibly
   ``Count()``? My thought is that ``Count()`` seems to idicate an action which
   would not be O(1). Then again who cares if it **might** indicate an action if
   it is documented O(1); plus it is a shorter name and more common usage.
2. Should we keep the ``Iter()`` and/or ``Range()`` methods? ``Range()`` is
   implemented in terms of ``Iter()`` anyhow. Or just keep both for choice and
   because they are already implemented.
3. Should ``NumEntries()`` or ``Count()`` return uint or int? I made it uint
   because it should never return a negative number. Is that a good enough
   reason?

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
* A functional Set, called _set_, which uses a [HAMT][1] internally.
* A functional sorted Map, called _sorted_map_, which uses a
  [standard Red-Black Tree][2] internally.
* A functional sorted Set, called _sorted_set_, which uses a
  [standard Red-Black Tree][2] internally.

I am planning on implementing:

* A functional Vector, called _vector_, which uses a clojure-like
  implementation of the same name.

[1]:https://en.wikipedia.org/wiki/Hash_array_mapped_trie
[2]:https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
