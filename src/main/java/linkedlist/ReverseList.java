package linkedlist;

import javax.swing.tree.TreeNode;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
class ListNode {
    int val;
    ListNode next;
    ListNode(int val) {
        this.val = val;
    }
}

public class ReverseList {

    public static ListNode reverseListRecur(ListNode root) {
        if (root == null) {
            return null;
        }

        if (root.next == null) {
            return root;
        }
        ListNode curr = root;

        ListNode rootOfReversed = reverseListRecur(root.next);

        ListNode nd = rootOfReversed;
        while (nd.next != null) {
            nd = nd.next;
        }
        nd.next = curr;
        curr.next = null;
        return rootOfReversed;

    }

    public static ListNode reverseListIter(ListNode root) {

        if (root == null) {
            return null;
        }

        if (root.next == null) {
            return root;
        }

        ListNode nd = root;
        ListNode next = nd.next;
        ListNode nNext = next.next;
        nd.next = null;

        while (next != null) {
            next.next = nd;
            nd = next;
            next = nNext;
            if (nNext != null) {
                nNext = nNext.next;
            }

        }
        return nd;
    }

    public static void main(String[] args) {
        ListNode l1 = new ListNode(7);

        ListNode l2 = new ListNode(8);

        ListNode l3 = new ListNode(10);

        ListNode l4 = new ListNode(12);

        ListNode l5 = new ListNode(14);

        l1.next = l2;
        l2.next = l3;
        l3.next = l4;
        l4.next = l5;

        ListNode l = reverseListRecur(l1);

    }
}
