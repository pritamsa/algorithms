package additional;

import java.util.ArrayList;
import java.util.List;

class ListNode {
      int val;
      ListNode next;
      ListNode() {}
      ListNode(int val) { this.val = val; }
      ListNode(int val, ListNode next) { this.val = val; this.next = next; }
  }
public class NthFromEnd {

    public static void main(String[] args ) {
        ListNode hd = new ListNode(3);
        ListNode hd1 = new ListNode(7);

        ListNode hd2 = new ListNode(9);
        ListNode hd3 = new ListNode(3);
        ListNode hd4 = new ListNode(5);
        ListNode hd5 = new ListNode(8);
        ListNode hd6 = new ListNode(0);

        hd.next = hd1;
        hd1.next = hd2;
        hd2.next = hd3;
        hd3.next = hd4;
        hd4.next = hd5;
        hd5.next = hd6;

        //ListNode res = removeNthFromEnd(hd,8);

        ListNode nd = new ListNode(2);
        ListNode nd1 = new ListNode(4);

        ListNode nd2 = new ListNode(3);

        ListNode pd = new ListNode(5);
        ListNode pd1 = new ListNode(6);
        ListNode pd2 = new ListNode(4);
        nd.next = nd1;
        nd1.next = nd2;

        pd.next = pd1;
        pd1.next = pd2;
        addTwoNumbers(nd, pd);
    }

    public static ListNode addTwoNumbers(ListNode l1, ListNode l2) {

        if (l1 == null && l2 == null) return null;
        if (l1 == null && l2 != null) return l2;
        if (l1 != null && l2 == null) return l1;

        ListNode nd1 = l1;
        ListNode nd2 = l2;
        ListNode res = null;
        int carry = 0;

        while(true) {

            if (nd1 == null && nd2 == null) break;
            if (nd1 == null) nd1 = new ListNode(0);
            if (nd2 == null) nd2 = new ListNode(0);

            int val = carry + nd1.val + nd2.val;
            carry = val/10;
            val = val%10;

            ListNode ndSum = new ListNode(val);
            if(res == null) {
                res = ndSum;
            } else {
                res.next = ndSum;
            }
            nd1 = nd1.next;
            nd2 = nd2.next;

        }
        return res;


    }
    public static ListNode removeNthFromEnd(ListNode head, int n) {

        int i = 0;
        ListNode nd = head;
        ListNode pt = null;

        if (head == null || n < 1) {
            return head;
        }


        while (nd != null) {
            if (pt != null) {
                pt = pt.next;
            }
            nd = nd.next;
            i++;
            if (i == n + 1 && pt == null) {
                pt = head;
                i = 0;
            }

        }

        if (pt != null) {
            ListNode temp = pt.next.next;
            pt.next.next = null;
            pt.next = temp;
        } else if (i == n) {
            return head.next;
        }
        return head;
    }
}
