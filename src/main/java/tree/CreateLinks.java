package tree;

import java.util.ArrayList;
import java.util.Queue;
import java.util.concurrent.LinkedBlockingQueue;

public class CreateLinks {

    public static void main(String[] args) {
        BSTNode root = new BSTNode(1);

        BSTNode l = new BSTNode(2);
        BSTNode r = new BSTNode(3);

        BSTNode ll = new BSTNode(4);
        BSTNode lr = new BSTNode(5);

        BSTNode rl = new BSTNode(6);
        BSTNode rr = new BSTNode(7);
        l.left = ll;
        l.right = lr;

        r.left = rl;
        r.right = rr;

        root.left = l;
        root.right = r;

        (new CreateLinks()).connect(root);


    }
    public BSTNode connect(BSTNode root) {

        if (root == null) {
            return null;
        }

        Queue<BSTNode> q = new LinkedBlockingQueue<>();

        q.add(root);
        ArrayList<BSTNode> lst = new ArrayList<>();
        while (!q.isEmpty()) {
            int currSize = q.size();

            for (int i = 0; i < currSize; i++) {
                BSTNode nd = q.remove();

                if (nd != null) {
                    if (nd.left != null) {
                        q.add(nd.left);
                    }
                    if (nd.right != null) {
                        q.add(nd.right);
                    }
                }

                if (nd != null) {
                    lst.add(nd);
                }


            }
            if (lst.size() > 0) {
                for(int j = 0; j < lst.size(); j++) {
                    lst.get(j).next = (j == lst.size() - 1) ? null : lst.get(j+1);
                }
                lst.clear();
            }

        }

        return root;

    }
}
