package additional;

import java.util.ArrayList;
import java.util.List;

public class MaxPaths {
  public static void main(String[] args) {
    TreeNode nd = new TreeNode(-2);
    //TreeNode ndL = new TreeNode(2);
    TreeNode ndR = new TreeNode(-1);
    //nd.left = ndL;
    nd.right = ndR;

    int m = (new MaxPaths()).maxPathSum(nd);
  }

  List<Integer> altPaths = new ArrayList<>();
  public int maxPathSum(TreeNode root) {


    int max = maxPaths(root);
    int max1 = (altPaths != null && altPaths.size() > 0) ? altPaths.stream().max(Integer::compare)
        .get() : Integer.MIN_VALUE;

    return Math.max(max1,max);

  }

  public int maxPaths(TreeNode root) {
    if (root == null) {
      return Integer.MIN_VALUE;
    }
    if (root.left == null && root.right == null) {
      altPaths.add(root.val);
      return root.val;
    }

    int leftVal = Math.max(maxPaths(root.left), 0);
    int rightVal = Math.max(maxPaths(root.right), 0) ;

    int altPath = leftVal + root.val + rightVal;
    altPaths.add(altPath);

    int mx = Math.max(leftVal,rightVal);

    return (mx + root.val);

  }



}
