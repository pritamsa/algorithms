package arrays;

//Pascal Triangle

import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

//Given numRows, generate the first numRows of Pascal’s triangle.
//
//        Pascal’s triangle : To generate A[C] in row R, sum up A’[C] and A’[C-1] from previous row R - 1.
//
//        Example:
//
//        Given numRows = 5,
//
//        Return
//
//        [
//        [1],
//        [1,1],
//        [1,2,1],
//        [1,3,3,1],
//        [1,4,6,4,1]
//        ]
//11:10 to 11:35
public class PascalTriangle {
    public static void main(String[] args) {


        ArrayList<ArrayList<Integer>> solve  = solve(5);
    }
    public static ArrayList<ArrayList<Integer>> solve(int A) {

        ArrayList<Integer> lst = new ArrayList<Integer>();
        ArrayList<Integer> lst2 = new ArrayList<Integer>();
        ArrayList<ArrayList<Integer>> ret = new ArrayList<ArrayList<Integer>>(A);


        lst.add(1);

        ret.add((ArrayList<Integer>) lst.clone());

        int j = 0;
        while(j < A - 1) {
            lst.add(0, 0);
            lst.add(0);

            for (int i = 0; i < lst.size() - 1;i++) {
               lst2.add(lst.get(i) + lst.get(i+1));
            }
            ret.add((ArrayList<Integer>) ((ArrayList<Integer>) lst2).clone());
            lst = (ArrayList<Integer>) lst2.clone();
            lst2.clear();
            j++;

        }
        return ret;

    }
}
