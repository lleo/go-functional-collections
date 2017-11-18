Red-Black Tree Properties
=========================

0) Each node is either red or black.

1) The root is black. This rule is sometimes omitted. Since the root can always
   be changed from red to black, but not necessarily vice versa, this rule has
   little effect on analysis.

2) All leaves (NIL) are black.

3) If a node is red, then both its children are black.

4) Every path from a given node to any of its descendant NIL nodes contains the
   same number of black nodes. Some definitions: the number of black nodes from
   the root to a node is the node's black depth; the uniform number of black
   nodes in all paths from root to the leaves is called the black-height of the
   red–black tree.
