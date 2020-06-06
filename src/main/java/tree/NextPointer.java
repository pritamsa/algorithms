package tree;

import java.util.ArrayList;
import java.util.List;

class Node {
    public int val;
    public Node left;
    public Node right;
    public Node next;

    public Node() {}

    public Node(int _val) {
        val = _val;
    }

    public Node(int _val, Node _left, Node _right, Node _next) {
        val = _val;
        left = _left;
        right = _right;
        next = _next;
    }
};
public class NextPointer {

    public static void main(String[] args) {
        Node n3 = new Node (3);
        Node n2 = new Node (2);

        Node n5 = new Node (5);
        Node n6 = new Node (6);
        //Node n7 = new Node (7);

        n2.left = n5;
        n2.right = n6;
        //n3.right = n7;
        Node n1 = new Node (1, n2, n3, null);

        (new NextPointer()).connect(n1);
    }

    public Node connect(Node root) {


        getConnections(root);
        return root;
    }

    public List<List<Node>> getConnections(Node root) {
        List<List<Node>> ret = new ArrayList<>();
        List<Node> lst = new ArrayList<Node>();

        if (root == null) {
            return ret;
        }
        root.next = null;
        lst.add(root);
        ret.add(lst);

        if (root.left == null && root.right == null) {
            return ret;
        }

        List<List<Node>> left = getConnections(root.left);
        List<List<Node>> right = getConnections(root.right);

        int maxLevel = Math.max(left.size(), right.size());

        int level = 0;
        //level by level merge
        while(level < maxLevel) {
            List<Node> leftLs = null;
            List<Node> rightLs = null;

            if (left != null && left.size() > 0 && level < left.size()) {
                if (level < left.size()) {
                    leftLs = left.get(level);
                }

                for (int i = 0; i < leftLs.size() - 1; i++) {
                    if (leftLs.get(i) != null ) {
                        leftLs.get(i).next = leftLs.get(i+1);
                    }
                }
            }
            if (right != null && right.size() > 0 && level < right.size()) {
                rightLs = right.get(level);

                if (leftLs != null && left.size() > 0 && rightLs != null) {
                    leftLs.get(leftLs.size() - 1).next = rightLs.get(0);
                }
                for (int i = 0; i < rightLs.size() - 1; i++) {
                    if (rightLs.get(i) != null ) {
                        rightLs.get(i).next = rightLs.get(i+1);
                    }
                }

            }
            List<Node> lsN = new ArrayList<Node>();
            if (leftLs != null) lsN.addAll(leftLs);
            if (rightLs != null) lsN.addAll(rightLs);
            ret.add(lsN);
            level++;
        }
        return ret;

    }
}
