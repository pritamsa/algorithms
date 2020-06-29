package additional;

import java.util.*;
//Given an unsorted array, find the maximum difference between the successive elements in its sorted form.
//
//        Return 0 if the array contains less than 2 elements.
public class MaximumGap {

    static int maxSortedAdjacentDiff(int[] arr, int n) {
        // Find maximum and minimum in arr[]
        int maxVal = arr[0];
        int minVal = arr[0];
        for (int i = 1; i < n; i++) {
            maxVal = Math.max(maxVal, arr[i]);
            minVal = Math.min(minVal, arr[i]);
        }

        // Arrays to store maximum and minimum values
        // in n-1 buckets of differences.
        int maxBucket[] = new int[n - 1];
        int minBucket[] = new int[n - 1];
        Arrays.fill(maxBucket, 0, n - 1, Integer.MIN_VALUE);
        Arrays.fill(minBucket, 0, n - 1, Integer.MAX_VALUE);

        // Expected gap for every bucket.
        float delta = (float) (maxVal - minVal) / (float) (n - 1);

        // Traversing through array elements and
        // filling in appropriate bucket if bucket
        // is empty. Else updating bucket values.
        for (int i = 0; i < n; i++) {
            if (arr[i] == maxVal || arr[i] == minVal) {
                continue;
            }

            // Finding index of bucket.
            int index = (int) (Math.round((arr[i] - minVal) / delta));

            // Filling/Updating maximum value of bucket
            if (maxBucket[index] == Integer.MIN_VALUE) {
                maxBucket[index] = arr[i];
            } else {
                maxBucket[index] = Math.max(maxBucket[index], arr[i]);
            }

            // Filling/Updating minimum value of bucket
            if (minBucket[index] == Integer.MAX_VALUE) {
                minBucket[index] = arr[i];
            } else {
                minBucket[index] = Math.min(minBucket[index], arr[i]);
            }
        }

        // Finding maximum difference between maximum value
        // of previous bucket minus minimum of current bucket.
        int prev_val = minVal;
        int max_gap = 0;
        for (int i = 0; i < n - 1; i++) {
            if (minBucket[i] == Integer.MAX_VALUE) {
                continue;
            }
            max_gap = Math.max(max_gap, minBucket[i] - prev_val);
            prev_val = maxBucket[i];
        }
        max_gap = Math.max(max_gap, maxVal - prev_val);

        return max_gap;
    }

    // Driver program to run the case
    public static void main(String[] args) {

        int arr[] = {3,6,9,1, 10, 14, 7, 8, 0, 17, 20}; //0 1 3 6 7 8 9 10 14 17 20 //2,3, 4, 3,3
        int n = arr.length;
        System.out.println(maxSortedAdjacentDiff(arr, n));
    }
}
