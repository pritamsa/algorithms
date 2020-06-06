package linkedlist;


class Node {
    int val;
    Node next;
    Node(int val) {
        this.val = val;
    }

}
//1-> 2 -> 3-> 4->5 - null
public class RemoveKthElement {

    public static void main(String[] args) {
        Node n1 = new Node(1);
        Node n2 = new Node(2);
        Node n3 = new Node(3);
        Node n4 = new Node(4);
        Node n5 = new Node(5);
        n1.next = n2;
        n2.next = n3;
        n3.next = n4;
        n4.next = n5;
        removeKthElement(n1, 2);

    }

    public static Node removeKthElement(Node root, int k) {
        int i = 0;
        Node nd = root;
        Node slow = null;

        while (nd != null) {
            nd = nd.next;
            if (slow != null) {
                slow = slow.next;
            }
            i++;
            if (i == k + 1) {
                slow = root;
            }

        }

        if (slow != null) {
            if (slow.next != null) {
                Node nNext = slow.next.next;
                slow.next = nNext;
            }

        }
        return root;

    }

}
