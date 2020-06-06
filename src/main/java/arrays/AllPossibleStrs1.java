package arrays;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.LinkedBlockingQueue;

//Generate a list of all strings from groups of characters, one character per group.
public class AllPossibleStrs1 {

    public static void main(String[] args) {
        System.out.println("First Test");
        String[] arr = {"abc", "pqr", "stu"};
        LinkedList<String> set1 = new LinkedList(Arrays.asList(arr));
        allPossibleStrs(set1);
        int k = 3;
        //printAllKLength(set1, k);

//        System.out.println("\nSecond Test");
//        char[] set2 = {'a', 'b', 'c', 'd'};
//        k = 1;
//        printAllKLength(set2, k);
    }



    public static LinkedList<String> allPossibleStrs(LinkedList<String> sets) {

        Queue<String> q = new LinkedBlockingQueue<String>();

        String s1 = sets.get(0);
        StringBuffer buff = new StringBuffer();

        for (int i = 0; i < s1.length(); i++) {
            q.add(buff.append(s1.charAt(i)).toString());
            buff.delete(0, buff.length());
        }

        int j = 1;
        int currSize = q.size();

        while (!q.isEmpty() && j < sets.size()) {
            currSize = q.size();

            for (int i = 0; i < currSize; i++) {
                String prefix = q.remove();
                String str = sets.get(j);
                for (int k = 0; k < str.length() ; k++) {
                    q.add(prefix+ str.charAt(k));
                }
            }

            j++;

        }
        while (!q.isEmpty()) {
            System.out.println(q.remove());
        }

        return null;

    }



}
