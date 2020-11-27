package additional;

class Max_level {

  int max_level;
}

//class TreeNode {
//  TreeNode left;
//  TreeNode right;
//  int val;
//  TreeNode(int val) {
//    this.val = val;
//  }
//}

public class ViewFromRight {

  TreeNode root;
  Max_level max = new Max_level();

  // Recursive function to print right view of a binary tree.
  void rightViewUtil(TreeNode node, int level, Max_level max_level) {

    // Base Case
    if (node == null)
      return;

    // If this is the last Node of its level
    if (max_level.max_level < level) {
      System.out.print(node.val + " ");
      max_level.max_level = level;
    }

    // Recur for right subtree first, then left subtree
    rightViewUtil(node.right, level + 1, max_level);
    rightViewUtil(node.left, level + 1, max_level);
  }

  void rightView()
  {
    rightView(root);
  }

  // A wrapper over rightViewUtil()
  void rightView(TreeNode node) {

    rightViewUtil(node, 1, max);
  }

  // Driver program to test the above functions
  public static void main(String args[]) {
//    BinaryTree tree = new BinaryTree();
      TreeNode root = new TreeNode(1);
      root.left = new TreeNode(2);
      root.right = new TreeNode(3);
    root.left.left = new TreeNode(4);
    root.left.right = new TreeNode(5);
    root.right.right = new TreeNode(7);
//    tree.root.left.left = new Node(4);
//    tree.root.left.right = new Node(5);
//    tree.root.right.left = new Node(6);
//    tree.root.right.right = new Node(7);
//    tree.root.right.left.right = new Node(8);
//
//    tree.rightView();
    (new ViewFromRight()).rightView(root);

  }
} 