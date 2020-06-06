package tree;

public class IsValidBST {

    public static void main(String[] args) {
        TreeNode nd3 = new TreeNode(3);
        TreeNode nd1 = new TreeNode(1);
        TreeNode nd5 = new TreeNode(5);

        TreeNode nd0 = new TreeNode(0);
        TreeNode nd2 = new TreeNode(2);
        TreeNode nd4 = new TreeNode(4);
        TreeNode nd6 = new TreeNode(6);

        nd3.left = nd1;
        nd3.right = nd5;

        nd1.left = nd0;
        nd1.right = nd2;

        nd5.left = nd4;
        nd5.right = nd6;
        (new IsValidBST()).isValidBST(nd3);

    }
    public boolean isValidBST(TreeNode root) {
        if (root == null) {
            return true;
        }
        return isValid(root, null, null);

    }

    public boolean isValid(TreeNode root, Integer upperBound, Integer lowerBound) {
        if (root == null) {
            return true;
        }
        boolean rootValid = false;

        if (upperBound == null && lowerBound != null) {
            rootValid = root.val > lowerBound;
        } else if (lowerBound == null && upperBound != null) {
            rootValid = root.val < upperBound;
        } else if (upperBound == null && lowerBound == null) {
            rootValid = true;
        } else if (upperBound != null && lowerBound != null) {
            rootValid = (upperBound > root.val) && (root.val > lowerBound);
        }
        if (root.left == null && root.right == null) {
            return rootValid;
        }

        boolean leftValid = isValid(root.left, root.val, lowerBound);
        boolean rightValid = isValid(root.right, upperBound, root.val);
        return rootValid && leftValid && rightValid;
    }
}
