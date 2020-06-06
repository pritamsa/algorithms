package linkedlist;

public class MergeTwoLists {
    public static void main(String[] args) {
        ListNode l1 = new ListNode(1);
        ListNode l2 = new ListNode(2);
        ListNode l4 = new ListNode(4);

        ListNode l11 = new ListNode(1);
        ListNode l12 = new ListNode(3);
        ListNode l14 = new ListNode(4);
        (new MergeTwoLists()).mergeTwoLists(l1, l11);
    }

    public ListNode mergeTwoLists(ListNode l1, ListNode l2) {
        ListNode lstHead = null;
        ListNode lst = null;

        ListNode nd1 = l1;
        ListNode nd2 = l2;

        while (nd1 !=null && nd2 != null) {
            if (nd1.val <= nd2.val) {
                if (nd1 == l1 && lst == null) {
                    lst = nd1;
                    lstHead = lst;
                } else {
                    lst.next = nd1;
                    lst = lst.next;
                    nd1 = nd1.next;
                }

            } else {
                if (nd2 == l2 && lst == null) {
                    lst = nd2;
                    lstHead = lst;
                } else {
                    lst.next = nd2;
                    lst = lst.next;
                    nd2 = nd2.next;
                }
            }
        }
        return lstHead;

    }}
