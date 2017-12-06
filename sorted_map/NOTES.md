If the inserted node is an inside grandchild, rotate out. In other words,
if the new node is inside on the left side of the grandparent rotate left, and
if the new node is inside on the right side of the grandparent rotate right.


Non-inside gradchild scenarios:

                   10                          10
                 /    \                      /    \
                5                                  15
              /  \                                /  \
             3*                                       20*

Inside gradchild scenarios:

                   10                          10
                 /    \                      /    \
                5                                  15
              /  \                                /  \
                  7*                             12*


Prep-rotate stage:

nn=12
rotate_right(15):


                   10                          10
                 /    \                      /    \
                       15        --->              12*
                      /  \                        /  \
                    12*                          x    15
                   /  \                              /  \
                  x    y                            y

nn=7
rotate_left(5):

                   10                          10
                 /    \                      /    \
                5                --->       7*
              /  \                         /  \
                  7*                      5    y
                 / \                    /  \
                x   y                       x

rotateRight with pre-rotateLeft:

          10                      10                      7*
        /    \                  /    \                  /    \
       5              --->     7*             --->     5      10
     /  \                    /   \                   /  \    /  \
          7*                 5    y                      x  y
         / \                / \
        x   y                  x

rotateRight; no pre-rotate occured:

                   10                          5
                 /    \                      /   \
                5      y         --->       3*    10
              /  \                         / \   /  \
             3*   x                             x    y
            / \

rotateLeft with pre-rotateRight:

          10                      10                      12*
        /    \                  /    \                  /    \
              15      --->            12*     --->    10      15
             /  \                    /  \            /  \    /  \
           12*                      x    15              x  y
          /  \                          /  \
         x    y                        y

rotateLeft; no pre-rotate occured:

                   10                          15
                 /    \                      /    \
                x      15        --->      10      20*
                      /  \                /  \    /  \
                     y    20*            x    y
                         /  \

              50
            /    \
         25        75
        /  \      /  \
      12    47  63    87
      /\    /\  /\    /\

                15b
              /     \
            5r       25r
          /  \      /   \
         3b  10b  20b    30b
                           \
                           35r

deleteCase1: Remove(k=25), swap(25r 21r), remove
=============

                  15b                          15b
              /         \                  /         \
            5r            25r             +            21r
          /  \          /     \                       /    \
         3b  10b      20b      30b                  20b     +
                        \         \                   \
                         21r       35r                 21r


delete_case2: where gp of on is nil; target 5b
==============================================

               10b
             /     \
           5b        15b
          /  \      /   \
         5b   7b  12b    20b

insertCase4.2
               100r                 50b
              /    \               /   \
            50b                  30b   100r
           /
         30r

insertCase4.1

             7940b                                        7940b
           /       \                                    /       \
      4930b         8090b      inserCase4.1        4930b         8090b
     /     \       /     \         --->           /     \       /     \
    a      7100r        10050r                   a     5310r         10050r
           /  \                                        /  \
        5310r  d               rotateRight(7100)      b   7100r
        /  \                                              /  \
       b    c                                            c    d

                                  7940b
                                /       \
    insertCase4.2           5310b        8090b
       --->                 /  \        /     \
                        4930r   7100r        10050r
    rotateLeft(4930)    /  \    /  \
                       a    b  c    d

insertCase4.1
ogp=4930b     = p
oparent=7100r = n
nn=5310r      = l

insertCase4.2
oggp=7940b    == p
ogp=4930b     == n
oparent=5310r == r == ngp
nn=7100r


insert(on=nil, nn=60r, path=[50r, 40b, 20b])
============================================

               20b                                20b
             /     \                            /     \
          10b       40b        insertCase3   10b       40r
         /  \       /  \          --->      /  \       /  \
                  30r   50r                          30b   50b
                        / \                                / \
                           60r                                60r
                              insertRepair(ogp=40b, ngp=40b, path=[20b])

nn=60r
oparent=50r
ogp=40b

               20b                               20b
             /     \                           /     \
          10b       40r       insertCase2   10b       40r
         /  \       /  \         --->      /  \       /  \
                  30b   50b   (do nothing)          30b   50b
                        / \                               / \
                           60r                               60r

ogp=nil
oparent=20b
on=40b
nn=40b

Remove(k=30), removeNodeWithZeroOrOneChild(on=30, term=20b, path[40b, 20b])
===========================================================================

Falls thru deleteCase1,2,3,4,&5 to deleteCase4(). In deleteCase6 it hits the
first 

          20b                                          20b
        /     \            deleteCase6               /     \
     10b       40r            --->                10b       50r
    /  \       /  \       sib.color=parent.color /  \       /  \
             30b   50b    parent.setBlack                 40b   60b
                   / \    sib.rn.setBlack                       / \
                      60r rotateLeft(parent)

12 node tree for iteration
==========================

                                70b
                            /         \
                         40r            90r
                       /    \          /   \
                   20b      50b      80b    110b
                  /   \     / \      / \    /   \
                10r   30r     60r         100r   120r
                / \   / \     / \         / \    / \

cur=10r
path=[20b, 40r, 70b]
endKey=pinf

