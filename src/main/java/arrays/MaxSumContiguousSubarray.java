package arrays;

import java.util.ArrayList;
import java.util.Arrays;

//Find the contiguous subarray within an array, A of length N which has the largest sum.
//        Input Format:
//
//        The first and the only argument contains an integer array, A.
//        Output Format:
//
//        Return an integer representing the maximum possible sum of the contiguous subarray.
//        Constraints:
//
//        1 <= N <= 1e6
//        -1000 <= A[i] <= 1000
//        For example:
//
//        Input 1:
//        A = [1, 2, 3, 4, -10]
//
//        Output 1:
//        10
//
//        Explanation 1:
//        The subarray [1, 2, 3, 4] has the maximum possible sum of 10.
//
//        Input 2:
//        A = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
//
//        Output 2:
//        6
//
//        Explanation 2:
//        The subarray [4,-1,2,1] has the maximum possible sum of 6.
//11:01
// give max st max end as well as max sum
//Correct algorithm
public class MaxSumContiguousSubarray {
    public static int maxSubArray(final int[] A, int[] ret) {


        int max_sum = Integer.MIN_VALUE;
        int maxSt = -1;
        int maxEn = -1;

        int sum_here = Integer.MIN_VALUE;
        int st = -1;
        int en = -1;

        for (int i = 0; i < A.length; i++) {
            int localSum = 0;

            if (sum_here != Integer.MIN_VALUE) {
                localSum = sum_here;
            } else {
                st = 0;
            }
            localSum += A[i];

            if (A[i] > localSum) {
                st = i;
                en = i;
                sum_here = A[i];
            } else {
                sum_here = localSum;
                en = i;
            }

            if (max_sum < sum_here) {
                max_sum = sum_here;
                maxSt = st;
                maxEn = en;
            }


        }

        int j = 0;
        for (int i = maxSt; i <= maxEn ; i++) {
            ret[j++] = A[i];
        }


        return max_sum;}

    public static void main(String[] args) {

        int[] arr = {1, 0, 3, 4, -10};
        int[] ret = new int[arr.length];




        int max = maxSubArray(arr,ret);
    }
}
