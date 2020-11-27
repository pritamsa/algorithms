package tree;


import java.util.Arrays;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.LinkedBlockingQueue;

public class SerializeDeserializeTree {

    public static void main(String[] args) {
        TreeNode nd1 = new TreeNode(15);
        TreeNode nd2 = new TreeNode(10);
        TreeNode nd3 = new TreeNode(5);
        TreeNode nd4 = new TreeNode(12);

        TreeNode nd6 = new TreeNode(18);
        TreeNode nd7 = new TreeNode(20);
        TreeNode nd8 = new TreeNode(16);

        nd1.left = nd2;
        nd2.left = nd3;
        nd2.right = nd4;

        nd6.right = nd7;
        nd6.left = nd8;
        nd1.right = nd6;
        int c = 'b' - 'a';

        String str = (new SerializeDeserializeTree()).serialize(nd1);
        System.out.println(str);
        LinkedBlockingQueue<String> lst = new LinkedBlockingQueue<String>(Arrays.asList(str.split("\\|")));
        TreeNode newRoot = (new SerializeDeserializeTree()).deSerialize1(lst, Integer.MIN_VALUE, Integer.MAX_VALUE);
    }
    public String serialize(TreeNode node) {
        StringBuilder builder = new StringBuilder("");
        serialize(node, builder);
        return builder.toString();
    }

    private void serialize(TreeNode root, StringBuilder builder) {
        if (root != null) {
                builder.append(root.val);
                builder.append("|");

                String left = null;
                String right = null;


                serialize(root.left, builder);


                serialize(root.right, builder);

            } else {
//                builder.append("-1");
//                builder.append("|");
            }

    }

    private TreeNode deSerialize1(LinkedBlockingQueue<String> strs, int lower, int upper) {
        if (strs == null || strs.isEmpty() || strs.size() == 0) {
            return null;
        }

        int val = Integer.parseInt(strs.peek());
        if (val < lower || val > upper) {
            return null;
        }
        strs.remove();

        TreeNode nd = null;

        if (val == -1) {

            return null;
        }

        nd = new TreeNode(val);

        nd.left = deSerialize1(strs, lower, val);

        nd.right = deSerialize1(strs, val, upper);


        return nd;
    }

    private TreeNode deSerialize(LinkedBlockingQueue<String> strs) {
        if (strs == null || strs.isEmpty() || strs.size() == 0) {
            return null;
        }

        String num = strs.remove();
        int val = Integer.parseInt(num);
        TreeNode nd = null;

        if (val == -1) {

            return null;
        }

            nd = new TreeNode(val);

            nd.left = deSerialize(strs);

            nd.right = deSerialize(strs);


        return nd;
    }
}
