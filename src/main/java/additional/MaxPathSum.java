package additional;


import java.util.ArrayList;
import java.util.List;

//Given a non-empty binary tree, find the maximum path sum.
//
//        For this problem, a path is defined as any sequence of nodes from some starting node to any node in the tree along the parent-child connections. The path must contain at least one node and does not need to go through the root.
//
//        Example 1:
//
//        Input: [1,2,3]
//
//        1
//        / \
//        2   3
//
//        Output: 6
//        Example 2:
//
//        Input: [-10,9,20,null,null,15,7]
//
//        -10
//        / \
//        9  20
//        /  \
//        15   7
//
//        Output: 42
//Idea is to keep pathSums globally where you add the path of left + root + right and also root.
//And each node finds left and right and ignores ones that are less than zero.
//Then it forwards max(left+root, right+ root) to its parent.
class TreeNode {
      int val;
      TreeNode left;
      TreeNode right;
      TreeNode() {}
      TreeNode(int val) { this.val = val; }
      TreeNode(int val, TreeNode left, TreeNode right) {
          this.val = val;
          this.left = left;
          this.right = right;
      }
  }
public class MaxPathSum {

    private static List<Integer> pathSums = new ArrayList<>();

    public static void main(String[] args) {
        TreeNode nd = new TreeNode(-2);
        TreeNode nd1 = new TreeNode(-1);
        TreeNode nd2 = new TreeNode(20);
        TreeNode nd3 = new TreeNode(15);
        TreeNode nd4 = new TreeNode(7);

//        nd.left = nd1;
//
//        nd2.left = nd3;
//        nd2.right = nd4;
        nd.left = nd1;

        System.out.println(maxPathSum(nd));

    }

    public static int maxPathSum(TreeNode root) {
        int pathSum = maxPathSumUtil(root);
        Integer max = pathSums.stream().max(Integer::compare).get();
        return Math.max(max, pathSum);
    }

    public static int maxPathSumUtil(TreeNode root) {
        if (root == null) {
            return 0;

        }

        if (root.left == null && root.right == null) {
            pathSums.add(root.val);
            return root.val;
        }

        int leftSum = Math.max(maxPathSumUtil(root.left), 0);
        int rightSum =  Math.max(maxPathSumUtil(root.right), 0);

        int maxVal = root.val + Math.max(leftSum, rightSum);
        int otherPath =root.val + rightSum + leftSum;
        pathSums.add(Math.max(otherPath, root.val));
        return maxVal;

    }
}
