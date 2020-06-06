package tree;

import sun.reflect.generics.tree.Tree;

import java.util.ArrayList;
import java.util.List;
import java.util.Stack;

class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode(int val) {
        this.val = val;
    }
}
public class InPrePostOrder {

    public static void main(String[] args) {

        TreeNode root = new TreeNode(1);

        TreeNode l = new TreeNode(2);
        TreeNode r = new TreeNode(3);

        TreeNode ll = new TreeNode(4);
        TreeNode lr = new TreeNode(5);

        TreeNode rl = new TreeNode(6);
        TreeNode rr = new TreeNode(7);
        l.left = ll;
        l.right = lr;

        r.left = rl;
        r.right = rr;

        root.left = l;
        root.right = r;
        (new InPrePostOrder()).inorderTraversal(root);

    }
    public List<Integer> inorderTraversal(TreeNode root) {

        if(root == null) {
            return null;
        }

        Stack<TreeNode> s = new Stack<>();
        List<Integer> ret = new ArrayList<>();

        TreeNode curr = root;
        s.push(curr);


        while (!s.isEmpty() || curr != null) {
            if (curr != null) {

                while(curr.left != null) {
                    curr = curr.left;
                    s.push(curr);
                }
            }
            if (!s.isEmpty()) {
                TreeNode nd = s.pop();
                ret.add(nd.val);
                curr = nd.right;
                if(curr != null) {
                    s.push(curr);
                }
            }
        }
        return ret;
    }
    public void inOrder(TreeNode root) {
        if(root != null) {
            Stack<TreeNode> s = new Stack<TreeNode>();
            TreeNode temp = root;
            s.push(temp);

            while (temp != null || !s.isEmpty()) {
                while (temp != null) {
                    s.push(temp.left);
                    temp = temp.left;
                }

                if (temp == null && !s.isEmpty()) {
                    TreeNode nd = s.pop();
                    if (nd != null) {
                        System.out.println(nd.val);
                    }

                    temp = nd.right;
                    if (temp != null) {
                        s.push(temp);
                    }
                }
            }

        }
    }

    public void preOrder(TreeNode root) {
        Stack<TreeNode> s = new Stack<TreeNode>();

        s.push(root);
        TreeNode curr = root;



            while (!s.isEmpty()) {
                TreeNode nd = s.pop();
                System.out.println(nd.val);
                s.push(nd.right);
                s.push(nd.left);
            }



    }

    public void postOrder(TreeNode root) {

        TreeNode curr = root;
        Stack<TreeNode> s = new Stack<TreeNode>();

        while (curr != null || !s.isEmpty()) {
            while(curr != null) {
                if (curr.right != null) {
                    s.push(curr.right);
                }
                s.push(curr);
                curr = curr.left;
            }

            while (curr == null && !s.isEmpty()) {
                TreeNode nd = s.pop();

                if (!s.isEmpty() && s.peek() != null && s.peek().equals(nd.right)) {
                    curr = s.pop();
                    s.push(nd);
                } else {
                    System.out.println(nd.val);
                }
            }
        }
    }

}
