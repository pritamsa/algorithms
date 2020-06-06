package tree;

public class ConstructBinTreeFromPreAndIn {

    public TreeNode buildTree(int[] preorder, int[] inorder) {
        return buildUtil(0, 0, inorder.length - 1, preorder, inorder);
    }

    public TreeNode buildUtil(int preIdx, int inSt, int inEn, int[] pre, int[] in) {

        if (inSt == inEn) {
            if (pre[preIdx] == in[inSt]) {
                boolean i = 2 >= 4;
                return new TreeNode(pre[preIdx]);
            }
        }

        TreeNode root = new TreeNode(pre[preIdx]);

        int inIdxOfRoot = getIdxIn(in, pre[preIdx]);
        if (inIdxOfRoot <= 0) {
            root.left = null;
        } else {
            root.left = buildUtil(preIdx++, inSt, inIdxOfRoot-1, pre, in );
        }
        if (inIdxOfRoot >= in.length - 1) {
            root.right = null;
        } else {
            root.right = buildUtil(preIdx++, inIdxOfRoot+1, inEn, pre, in );
        }

        return root;

    }

    private int getIdxIn(int[] in, int val) {
        for (int i = 0; i < in.length; i++) {
            if (val == in[i]) {
                return i;
            }
        }
        return -1;
    }

}
