package additional;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

//Given an array of 4 digits, return the largest 24 hour time that can be made.
//
//        The smallest 24 hour time is 00:00, and the largest is 23:59.  Starting from 00:00, a time is larger if more time has elapsed since midnight.
//
//        Return the answer as a string of length 5.  If no valid time can be made, return an empty string.
//
//
//
//        Example 1:
//
//        Input: [1,2,3,4]
//        Output: "23:41"
//        Example 2:
//
//        Input: [5,5,5,5]
//        Output: ""
//Brute force: fine all perms of 4 numbers and select valid ones
public class LargestClockTime {

    public static void main(String[] args) {
        int[] arr = {1,2, 3};
        largestTimeFromDigits1(arr);
//        largestTimeFromDigits(arr);
//        Arrays.stream(arr).sum();
    }

    public boolean isValidTime(String str) {
        return true;
    }
    public static String largestTimeFromDigits1(int[] A) {
        List<String> lst = getPerms(A, 0);
        return "";
    }

    public static List<String> getPerms(int[] arr, int st) {
        List<String> perms = new ArrayList<>();

        if (arr == null) {
            return perms;
        }
        if (st == arr.length - 1) {
            perms.add(arr[st] + "");
            return perms;
        }

        int num = arr[st];

        List<String> vals = getPerms(arr, st+1);
        for (String val: vals ) {
            if (val != null && val.trim().length() > 0 ) {
                for (int i = 0; i <= val.length(); i++) {
                    String newWord = insertNum(val, i, num);
                    perms.add(newWord);
                }
            }
        }
        return perms;

    }

    private static String insertNum(String str, int loc, int num) {
        if (loc == 0) {
            return num + str;
        } else if (loc == str.length()){
            return str+num;
        } else {
            return str.substring(0, loc) + num + str.substring(loc);
        }
    }

    public static String largestTimeFromDigits(int[] A) {
        int largestHr = getLargestHr(A);//9077
        int largestMin = getLargestMin(A);

        if (largestHr == -1 || largestMin == -1) {
            return "";
        } else return largestHr+":"+ largestMin;
    }

    private static int getLargestHr(int[] A) {
        int msb = -1;
        int lsb = -1;
        int idx1 = -1;
        int idx2 = -1;


        for(int i = 0; i < A.length; i++) {
            if(msb < A[i] && A[i] < 3) {
                msb = A[i];
                idx1 = i;
            }
        }

        for(int i = 0; i < A.length; i++) {
            if (idx1 != i) {
                if (msb == 2) {
                    if (A[i] <= 3 && lsb < A[i]) {lsb = A[i]; idx2 = i;}
                } else if (msb == 1 || msb == 0) {
                    if (A[i] <= 9 && lsb < A[i]) {lsb = A[i]; idx2 = i;}
                }

            }
        }

        if (idx1 > -1) A[idx1] = -1;
        if (idx2 > -1) A[idx2] = -1;
        if (msb == -1 || lsb == -1) {
            return -1;
        }
        return msb*10+lsb;
    }

    private static int getLargestMin(int[] A) {
        int msb = -1;
        int lsb = -1;
        for(int i = 0; i < A.length; i++) {
            if(A[i] >= 0)
            { if (msb == -1) {
                msb = A[i];
            } else {
                lsb = A[i];
            }


            }
        }

        if (msb == -1 || lsb == -1) {
            return -1;
        }

        if (msb*10+lsb >= 60){
           if (lsb*10+msb >= 60) {
               return -1;
           } else {
               return lsb*10+msb;
           }
        } else {
            if (lsb*10+msb >= 60) {
                return msb*10+lsb;
            } else {
                return Math.max(msb*10+lsb, lsb*10+msb);
            }
        }

    }


}

//    Given an integer array nums, find the sum of the elements between indices i and j (i ≤ j), inclusive.
//
//        The update(i, val) function modifies nums by updating the element at index i to val.
//
//        Example:
//
//        Given nums = [1, 3, 5]
//
//        sumRange(0, 2) -> 9
//        update(1, 2)
//        sumRange(0, 2) -> 8
//        Note:
//
//        The array is only modifiable by the update function.
//        You may assume the number of calls to update and sumRange function is distributed evenly.