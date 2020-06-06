package linkedlist;

public class AddNumbers {

    public static void main(String[] args) {
        ListNode n1 = new ListNode(7);
        ListNode n2 = new ListNode(2);
        ListNode n3 = new ListNode(4);
        ListNode n4 = new ListNode(3);
        n1.next = n2;
        n2.next = n3;
        n3.next = n4;

        ListNode n5 = new ListNode(5);
        ListNode n6 = new ListNode(6);
        ListNode n7 = new ListNode(4);

        n5.next = n6;
        n6.next = n7;
        (new AddNumbers()).addTwoNumbers(n1,n5);
    }
    public ListNode addTwoNumbers(ListNode l1, ListNode l2) {

        int carry = 0;

        ListNode n1 = reverse(l1);
        ListNode n2 = reverse(l2);

        ListNode ret = null;
        ListNode ansHead = null;
        ListNode ans = null;
        while(n1 != null || n2 != null ) {
            int sum = 0;

            if (n1 != null) {
                sum += n1.val;
            }

            if (n2 != null) {
                sum += n2.val;
            }
            sum += carry;

            int val = sum % 10;
            carry = sum/10;
            ans = new ListNode(val);

            if (ansHead == null) {
                ansHead = ans;
                ret = ansHead;
            } else {
                ansHead.next = ans;
                ansHead = ansHead.next;
            }
            if(n1 != null) n1 = n1.next;
            if(n2 != null) n2=n2.next;


        }
        return reverse(ret);

    }

    ListNode reverse(ListNode l1) {

        ListNode curr = l1;
        ListNode next = curr.next;
        ListNode nNext = next.next;

        curr.next = null;

        while (next != null) {
            if (next != null) {
                next.next = curr;
            }
            curr = next;
            //if (nNext != null) {
                next = nNext;
            //}
            if (next != null) {
                nNext = next.next;
            }
        }
        return curr;

    }
}
