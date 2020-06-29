package additional;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
//1->2->3->4->5->6 with len = 3
//ans: 3->2->1->6->5->4
//If len = 1 return as is
//else nd = root, netxt = nd.next & nNext = next.next.
//keep i
public class Reverse {

    public static void main(String[] args) {
        ListNode nd1 = new ListNode(1);
        ListNode nd2 = new ListNode(2);
//        ListNode nd3 = new ListNode(3);
//        ListNode nd4 = new ListNode(4);
//        ListNode nd5 = new ListNode(5);
//        ListNode nd6 = new ListNode(6);
//        ListNode nd7 = new ListNode(7);
//        ListNode nd8 = new ListNode(8);
        //ListNode nd9 = new ListNode(9);
        nd1.next = nd2;
//        nd2.next = nd3;
//        nd3.next = nd4;
//        nd4.next = nd5;
//        nd5.next = nd6;
//        nd6.next = nd7;
//        nd7.next = nd8;
        //nd8.next = nd9;
        ListNode hd = (new Reverse()).reverseKGroup(nd1,3);
    }

    public ListNode reverseKGroup(ListNode head, int k) {

        if (head == null) return null;
        if (head.next == null) return head;

        int i = 0;

        ListNode nd = head;
        ListNode newTail = head;
        ListNode next = nd.next;
        nd.next = null;
        ListNode nNext = (next != null) ? next.next : null;
        ListNode prevNd = null;

        while (nd != null && i < k-1) {
            if (next != null) next.next = nd;
            prevNd = nd;
            nd = next;
            next = nNext;
            i++;
            nNext = (next != null) ? next.next : null;

        }

        if (i != k && nd == null) {
            return reverseKGroup(prevNd, i);
        }

        newTail.next = reverseKGroup(next, k);
        return (nd != null) ? nd : prevNd;
    }
//TreeMap : 2:1, 3:1, 5:1,   7:-1 (2) 8:1 (3) -------
//    {3,7},  3:1 ,  7:0
//    {5,7},  5:1 +
//    {7,15},
//    {8,12},
//    {8,20},
//    {14,20}
//}
//3,1
//        5,2
//        7,1
//        8,3
//        12,2
//        14,3
//        15,2
//        20,0

    // arrival : 5, 4 , 2, 7, 8, 14 && dept : 7,7, 12, 12, 15, 20,20
    //Time :
    public Map<Integer,Integer> arrivalDeparture(List<List<Integer>> arr, int stop) {

        int[] arrivals = new int[arr.size()]; //O(n)
        int[] depts = new int[arr.size()];
        int count = 0;
        Map<Integer, Integer> map = new HashMap<>(); // O(m)

        //Construct proper arrays (O(n) the number of entries in input listn) where n i)
        for (int i = 0; i < arr.size(); i++) {// 3, 5, 7, 8, 14 && dept : 7,7, 12, 12, 15, 20,20
            arrivals[i] = arr.get(i).get(0);
            depts[i] = arr.get(i).get(1);
        }

        //O(nlogn)
        Arrays.sort(arrivals);
        Arrays.sort(depts);

        int i = 0;
        int j= 0;

        //Keep adding and subtracting counts O(n)
        while (i < arrivals.length && j < depts.length) {// 3, 5, 7, 8, 14 && dept : 7,7, 12, 12, *15, 20,20
            if(arrivals[i] < depts[j]) {
                count++;
                i++;
                map.put(arrivals[i], count);//3:1, 5:2, 7:1, 8:2,12:0,  14:1
            } else if (depts[j] < arrivals[i]) {
                count--;
                j++;
                map.put(depts[j], count);
            } else {
                map.put(depts[j], count);
                i++;
            }

        }

        //Decrement remaining O(k) where k is number of remaining operations
        while (j < depts.length) {
            count--;
            map.put(depts[j], count);
        }

        return map;
    }

}
