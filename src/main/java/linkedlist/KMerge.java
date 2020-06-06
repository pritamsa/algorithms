package linkedlist;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedList;

public class KMerge {
    //1->4->5,
    //  1->3->4,
    //  2->6
    public static void main(String[] args) {
        ListNode nd1 = new ListNode(1);
        nd1.next = new ListNode(4);
        nd1.next.next = new ListNode(5);

        ListNode nd2 = new ListNode(1);
        nd2.next = new ListNode(3);
        nd2.next.next = new ListNode(4);

        ListNode nd3 = new ListNode(2);
        nd3.next = new ListNode(6);

        ListNode[] lists = {nd1, nd2, nd3};

        (new KMerge()).mergeKLists(lists);

    }
    public ListNode mergeKLists(ListNode[] lists) {

        if (lists == null && lists.length == 0) {
            return null;
        }

        if (lists.length == 1) {
            return lists[0];
        }

        ArrayList<ListNode> merged = new ArrayList<>(Arrays.asList(lists));

        ArrayList<ListNode> merged1 = new ArrayList<>();

        while (merged.size() != 1) {
            for (int i = 0; i < merged.size(); i = i+2) {

                ListNode mergedLs = null;
                if (i == merged.size() - 1) {
                    mergedLs = merged.get(i);
                } else {
                    mergedLs = mergeLists(merged.get(i), merged.get(i + 1));
                }
                merged1.add(mergedLs);
            }
            merged = (ArrayList<ListNode>) merged1.clone();
            merged1.clear();

        }

        return merged.get(0);

    }

    public ListNode mergeLists(ListNode head1, ListNode head2) {

        ListNode nd1 = head1;
        ListNode nd2 = head2;

        ListNode newList = null;
        ListNode newListRoot = null;

        while (nd1 != null && nd2 != null) {
            ListNode newNd = null;
            if (nd1.val <= nd2.val ) {
                newNd = new ListNode(nd1.val);
                nd1 = nd1.next;
            } else {
                newNd = new ListNode(nd2.val);
                nd2 = nd2.next;
            }

            if (newList == null) {
                newList = newNd;
                newListRoot = newList;
            } else {
                newList.next = newNd;
                newList = newList.next;
            }
        }

        if (nd1 != null) {
            while (nd1 != null) {
                ListNode newNd = new ListNode(nd1.val);
                newList.next = newNd;
                nd1 = nd1.next;
            }
        }

        if (nd2 != null) {
            while (nd2 != null) {
                ListNode newNd = new ListNode(nd2.val);
                newList.next = newNd;
                nd2 = nd2.next;
            }
        }
        return newListRoot;
    }
}
