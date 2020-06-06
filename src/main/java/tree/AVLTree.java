package tree;

class AVLNode{
    public int val;
    public int height;
    AVLNode left;
    AVLNode right;
    AVLNode(int val, int height) {
        this.val = val;
        this.height = height;
    }
        }
public class AVLTree {

    private static int height(AVLNode nd) {
        if (nd == null) return 0;
        return nd.height;
    }

    private static int getBalanceFactor(AVLNode root){
        if (root == null) return 0;
        return height(root.left) - height(root.right);

    }

    public static AVLNode rightRotation(AVLNode nd) {
        if (nd != null) {
            AVLNode ndL = nd.left;
            AVLNode ndLR = ndL.right;

            ndL.right = nd;
            nd.left = ndLR;

            nd.height = 1 + Math.max( height(nd.left), height(nd.right));
            ndL.height = 1 + Math.max( height(ndL.left), height(ndL.right));

            return ndL;

        }
        return null;
    }

    public static AVLNode leftRotation(AVLNode nd) {
        if (nd != null) {
            AVLNode ndR = nd.right;
            AVLNode ndRL = ndR.left;

            ndR.left = nd;
            nd.right = ndRL;

            nd.height = 1 + Math.max( height(nd.left), height(nd.right));
            ndR.height = 1 + Math.max( height(ndR.left), height(ndR.right));



            return ndR;

        }
        return null;

    }

    public static AVLNode insert(AVLNode root, int val) {

        if (root == null) {
            return new AVLNode(val, 1);
        }
        if (val < root.val) {
            root.left = insert(root.left, val);

        } else if (val > root.val) {
            root.right = insert(root.right, val);

        }

        root.height =  1 + Math.max( root.left.height, root.right.height);
        int balance = getBalanceFactor(root);

        //case 1 : left left : right rotation
        if (balance > 1 && val < root.left.val) {

        }

        //case 2 : left right : left rotate and then right rotate
        if (balance > 1 && val > root.left.val) {

        }

        //case 3 : right right : left rotation
        if (balance < 1 && val > root.right.val) {

        }

        //case 4 : right left : right rotate and then left rotate
        if (balance < 1 && val < root.right.val) {

        }

        return null;
    }

}
