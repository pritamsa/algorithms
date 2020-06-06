package arrays;
//1:57
//You are given an array of N integers, A1, A2 ,…, AN. Return maximum value of f(i, j) for all 1 ≤ i, j ≤ N.
//        f(i, j) is defined as |A[i] - A[j]| + |i - j|, where |x| denotes absolute value of x.
//
//        For example,
//
//        A=[1, 3, -1]


import java.util.ArrayList;
import java.util.Arrays;

//= (1-3) + (1)
//= 2 + 1 = 3
//3+1 + 1
//= 5
//
//        f(1, 1) = f(2, 2) = f(3, 3) = 0
//        f(1, 2) = f(2, 1) = |0-1|+|1-2|  = |1 - 3| + |1 - 2| = 3
//        f(1, 3) = f(3, 1) = |1 - (-1)| + |1 - 3| = 4
//        f(2, 3) = f(3, 2) = |3 - (-1)| + |2 - 3| = 5
//
//        So, we return 5.
//Solution idea: f(i, j) = |A[i] – A[j]| + |i – j| can be written in 4 ways
//case 1 : if i > j & A[i] > A[j] : A[i] – A[j] + i - j : (A[i] + i) - (A[j] + j)
//case 2 : if i > j & A[i] < A[j] : A[j] – A[i] + i - j : (A[j] - j) - (A[i] - i)
//case 3 : if i < j & A[i] > A[j] : A[i] – A[j] + j - i : (A[i] - i) - (A[j] - j)
//case 4 : if i < j & A[i] < A[j] : A[j] – A[i] + j - i : (A[j] + j) - (A[i] + i)
//So we have groups of A[i] +- i values. We find min and max of those and then find max diff between
public class MaximumAbsoluteDifference {

    public static int maxAbsDiff(Integer[] A) {


        int max1 = Integer.MIN_VALUE;
        int min1 = Integer.MAX_VALUE;
        int max2 = Integer.MIN_VALUE;
        int min2 = Integer.MAX_VALUE;

        for (int i = 0; i < A.length; i++)
        {

            // Updating max and min variables
            // as described in algorithm.
            max1 = Math.max(max1, A[i] + i);
            max2 = Math.max(max2, A[i] - i);

            min1 = Math.min(min1, A[i] + i);
            min2 = Math.min(min2, A[i] - i);
        }

        return Math.max((max1-min1), (max2-min2));
    }

    public static void main(String[] args) {

        Integer[] arr = {1,3,-1};

        ArrayList<Integer> A = new ArrayList<Integer>(Arrays.asList(arr));

        int max = maxAbsDiff(arr);
    }

}
