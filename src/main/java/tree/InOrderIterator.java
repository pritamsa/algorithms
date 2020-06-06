package tree;

class BSTNode {
    BSTNode left;
    BSTNode right;
    BSTNode next;
    int val;
    BSTNode(int val) {
        this.val = val;
    }
}

public class InOrderIterator {

    BSTNode root;
    InOrderIterator(BSTNode root) {
        this.root = root;

    }

    boolean hasNext() {
        return getInOrderSuccessor(root) != null;

    }

    public  int[][] flipAndInvertImage(int[][] A) {
        if (A == null || A.length == 0) {
            return null;
        }

        for (int i = 0; i < A.length; i++) {
            flipInvert(A, i);
        }
        return A;
    }

    private  void flipInvert(int[][] arr, int row) {
        int i = 0;
        int j = arr[row].length - 1;

        while (i <= j) {
            if(arr[row][i] != arr[row][j]) {


            } else {
                if (i == j) {
                    arr[row][i] = ((arr[row][i] == 0) ? 1 : 0);
                } else {
                    arr[row][i] = ((arr[row][i] == 0) ? 1 : 0);
                    arr[row][j] = ((arr[row][j] == 0) ? 1 : 0);
                }
            }
            i++;
            j--;

        }

    }


    BSTNode next() {
        BSTNode nd = getInOrderSuccessor(root);
        if (nd != null) {
            delete(root, nd);
        }

        return nd;
    }

    private BSTNode getInOrderSuccessor(BSTNode rt) {

        if (rt == null) {
            return null;
        }
        if (rt.left == null && rt.right == null) {
            return rt;
        }

        if (rt.left != null) {
            return getInOrderSuccessor(rt.left);
        } else {
            return rt;
        }


    }
    BSTNode delete(BSTNode rt, BSTNode node) {

        if (node == null || rt == null) {
            return null;
        }


        if (node.val < rt.val) {
            rt.left = delete(rt.left, node);
        } else if (node.val > rt.val) {
            rt.right = delete(rt.right, node);
        } else if (rt.val == node.val) {
            if (rt.left == null && rt.right == null) {
                return null;
            }
            BSTNode newRootVal = findMin(rt);
            delete(rt, newRootVal);
            rt.val = newRootVal.val;
            return rt;
        }
        return rt;

    }

    BSTNode findMin(BSTNode rt) {
        if (rt == null) {
            return null;
        }
        if (rt.left == null && rt.right == null) {
            return rt;
        }

        if (rt.left != null) {
            return findMin(rt.left);
        } else if (rt.left == null && rt.right != null) {
            return findMin(rt.right);
        }
        return null;

    }

    public static void main(String[] args) {
        BSTNode root = new BSTNode(7);

        BSTNode n5 = new BSTNode(5);
        BSTNode n10 = new BSTNode(10);
        BSTNode n6 = new BSTNode(6);
        BSTNode n3 = new BSTNode(3);

        n5.left = n3;
        n5.right = n6;
        root.left = n5;
        root.right = n10;

        InOrderIterator iterator = new InOrderIterator(root);

        int[][] arr = new int[3][3];
        arr[0][0] = 1;
                arr[0][1] = 1;
                        arr[0][2] = 0;
        arr[1][0] = 1;
        arr[1][1] = 0;
        arr[1][2] = 1;

        arr[2][0] = 0;
        arr[2][1] = 0;
        arr[2][2] = 0;

        iterator.flipAndInvertImage(arr);

        //iterator.next();
    }
}
