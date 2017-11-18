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

delete_case1:

                  15b                          15b
              /         \                  /         \
            5r            25r             +            21r
          /  \          /     \                       /    \
         3b  10b      20b      30b                  20b     +
                        \         \                   \
                         21r       35r                 21r


delete_case2: where gp of on is nil; target 5b

               10b
             /     \
           5b        15b
          /  \      /   \
         5b   7b  12b    20b
