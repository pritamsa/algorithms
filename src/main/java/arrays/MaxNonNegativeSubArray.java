package arrays;

//Max Non Negative SubArray: 4:05

//Given an array of integers, A of length N, find out the maximum sum sub-array of non negative numbers from A.
//
//        The sub-array should be contiguous i.e., a sub-array created by choosing the second and fourth
// element and skipping the third element is invalid.
//
//        Maximum sub-array is defined in terms of the sum of the elements in the sub-array.
//
//        Find and return the required subarray.
//
//        NOTE:
//
//        1. If there is a tie, then compare with segment's length and return segment which has maximum length.
//        2. If there is still a tie, then return the segment with minimum starting index.
//
//
//        Input Format:
//
//        The first and the only argument of input contains an integer array A, of length N.
//        Output Format:
//
//        Return an array of integers, that is a subarray of A that satisfies the given conditions.
//        Constraints:
//
//        1 <= N <= 1e5
//        1 <= A[i] <= 1e5
//        Examples:
//
//        Input 1:
//        A = [1, 2, 5, -7, 2, 3]
//
//        Output 1:
//        [1, 2, 5]
//
//        Explanation 1:
//        The two sub-arrays are [1, 2, 5] [2, 3].
//        The answer is [1, 2, 5] as its sum is larger than [2, 3].
//
//        Input 2:
//        A = [10, -1, 2, 3, -4, 100]
//
//        Output 2:
//        [100]
//
//        Explanation 2:
//        The three sub-arrays are [10], [2, 3], [100].
//        The answer is [100] as its sum is larger than the other two.
// 0,0,-1,0

import java.util.ArrayList;
import java.util.Arrays;

//1, 2, 5, -7, 2, 3
//10, -1, 2, 3, -4, 100
public class MaxNonNegativeSubArray {

    public static ArrayList<Integer> maxNonNegativeSubArray(ArrayList<Integer> A) {
        if (A == null || A.size() == 0) {
            return A;
        }

        ArrayList<Integer> ret = new ArrayList(A.size());

        int st = -1;
        int en = -1;
        int sum_here = 0;

        int max_so_far = Integer.MIN_VALUE;
        int maxSt = -1;
        int maxEn = -1;

        for (int i = 0; i < A.size() ; i++) {
            if (sum_here == 0) {
                st = i;

            }

            if (A.get(i) >= 0) {
                sum_here += A.get(i);
                if (en == -1) {
                    en = st;
                } else {
                    en++;
                }
            } else {
                sum_here = 0;
                st = -1;
                en = -1;
            }

            if (sum_here > max_so_far) {
                maxSt = st;
                maxEn = en;
                max_so_far = sum_here;
            }
            if (i - 1>= 0 && A.get(i) == 0 && sum_here == max_so_far) {
                maxSt = st;
                maxEn = en;
            }

        }

        for (int k = maxSt; k <=maxEn ; k++) {
            ret.add(A.get(k));
        }


        return ret;

    }

    public static void main(String[] args) {

        Integer[] arr = {0,0,-1,0 };

        //ArrayList<Integer> A = new ArrayList<Integer>(Arrays.asList(arr));

        ArrayList<Integer> max = maxNonNegativeSubArray(new ArrayList<Integer>(Arrays.asList(arr)));
    }


}
