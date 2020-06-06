package tree;

import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.LinkedBlockingQueue;

public class LeafSimilarTree {
    /**
     * Definition for a binary tree node.
     * public class TreeNode {
     *     int val;
     *     TreeNode left;
     *     TreeNode right;
     *     TreeNode() {}
     *     TreeNode(int val) { this.val = val; }
     *     TreeNode(int val, TreeNode left, TreeNode right) {
     *         this.val = val;
     *         this.left = left;
     *         this.right = right;
     *     }
     * }
     */

    public static void main(String[] args) {
        LeafSimilarTree l = new LeafSimilarTree();
        TreeNode root1 = new TreeNode(1);
        TreeNode rl1 = new TreeNode(2);
        root1.right = rl1;
//        TreeNode rl2 = new TreeNode(1);

//        rl2.left = new TreeNode(9);
//        rl2.right = new TreeNode(8);

//        rl1.left = new TreeNode(6);
//
//        TreeNode rl3 = new TreeNode(2);
//
//        rl3.left = new TreeNode(7);
//        rl3.right = new TreeNode(3);
//        rl1.right = rl3;
//
//
//        root1.left = rl1;
//        root1.right = rl2;

        TreeNode root2 = new TreeNode(7);
        TreeNode rp1 = new TreeNode(5);
//        TreeNode rp2 = new TreeNode(3);
//        root2.left = rp1;
        root2.right = rp1;

        boolean leafSim = l.leafSimilar(root1, root2);
    }
        public boolean leafSimilar(TreeNode root1, TreeNode root2) {

            if (root1 == null || root2 == null) {
                return false;
            }
            LinkedList<TreeNode> r1LeafSeq = getLeaves(root1);
            LinkedList<TreeNode> r2LeafSeq = getLeaves(root2);

            if (r1LeafSeq == null || r1LeafSeq == null || r1LeafSeq.size() != r2LeafSeq.size()) {
                return false;
            }
            for (int i = 0; i < r1LeafSeq.size(); i++) {
                if (r1LeafSeq.get(i).val != r2LeafSeq.get(i).val) {
                    return false;
                }
            }
            return true;

        }

        private LinkedList<TreeNode> getLeafSeq(TreeNode rt) {

            Queue<TreeNode> q = new LinkedBlockingQueue();
            LinkedList<TreeNode> lst = new LinkedList<>();

            q.add(rt);

            boolean added = false;
            int currSize = 0;
            while(!q.isEmpty()) {
                currSize = q.size();


                for (int i = 0; i < currSize; i++ ) {
                    TreeNode nd = q.remove();
                    if (nd.left == null && nd.right == null) {
                        lst.add(nd);
                    }
                    if (nd.left != null) {
                        q.add(nd.left);

                    }
                    if (nd.right != null) {
                        q.add(nd.right);

                    }

                }

            }
            return lst;
        }

        private LinkedList<TreeNode> getLeaves(TreeNode root) {
            if (root == null) {
                return null;

            }
            LinkedList<TreeNode> ret = new LinkedList<TreeNode>();
            if (root.left == null && root.right == null) {
                ret.add(root);

            } else {
                LinkedList<TreeNode> leftLeaves = getLeaves(root.left);
                LinkedList<TreeNode> rightLeaves = getLeaves(root.right);
                if (leftLeaves != null) {
                    ret.addAll(leftLeaves);
                }
                if (rightLeaves != null) {
                    ret.addAll(rightLeaves);
                }
            }
            return ret;
        }


}
