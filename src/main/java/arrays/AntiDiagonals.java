package arrays;

//Give a N*N square matrix, return an array of its anti-diagonals. Look at the example for more details.
//
//        Example:
//
//
//        Input:
//
//        1 2 3 5
//        4 5 6 3
//        7 8 9 2
//        6 5 4 2
//
//        Return the following :
//
//        [
//        [1],
//        [2, 4],
//        [3, 5, 7],
//        [5, 6, 8, 6],
//        [3, 9, 5]
//        [2, 4],
//        [2]
//        ]
//
//
//        Input :
//        1 2
//        3 4
//
//        Return the following  :
//
//        [
//        [1],
//        [2, 3],
//        [4]
//        ]
//2:47
import java.util.ArrayList;

public class AntiDiagonals {
    public static ArrayList<ArrayList<Integer>> antiDiagonal(ArrayList<ArrayList<Integer>> A) {

        if (A == null || A.size() == 0) {
            return null;
        }
        ArrayList<ArrayList<Integer>> ret = new  ArrayList<ArrayList<Integer>>();

        int i = 0;
        for (int j = 0; j <= A.get(0).size()-2;j++) {
            ArrayList<Integer> lst = new ArrayList<Integer>();
            int k = j;
            while (k >= 0 && i < A.size()) {
                lst.add(A.get(i).get(k));
                i++;
                k--;
            }
            i = 0;
            ret.add(lst);

        }

        int t = A.get(0).size()-1;
        for (int a = 0; a < A.size();a++) {
            ArrayList<Integer> lst = new ArrayList<Integer>();
            int k = a;
            while (k < A.size() && t >= 0) {
                lst.add(A.get(k).get(t));
                k++;
                t--;
            }
            t = A.get(0).size()-1;
            ret.add(lst);
        }

        return ret;
    }

    public static void main(String[] args) {
        ArrayList<ArrayList<Integer>> A = new ArrayList<ArrayList<Integer>>();

        ArrayList<Integer> lst1 = new ArrayList<Integer>();
        ArrayList<Integer> lst2 = new ArrayList<Integer>();
        ArrayList<Integer> lst3 = new ArrayList<Integer>();
        ArrayList<Integer> lst4 = new ArrayList<Integer>();


        lst1.add(1);
       lst1.add(2);
       lst1.add(3);
       lst1.add(5);

       lst2.add(4);
       lst2.add(5);
       lst2.add(6);
       lst2.add(3);

        lst3.add(7);
        lst3.add(8);
        lst3.add(9);
        lst3.add(2);

        lst4.add(6);
        lst4.add(5);
        lst4.add(4);
        lst4.add(2);

        A.add(lst1);
        A.add(lst2);
        A.add(lst3);
        A.add(lst4);

        antiDiagonal(A);


    }
}
