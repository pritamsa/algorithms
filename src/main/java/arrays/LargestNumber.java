package arrays;

import java.util.*;

//Given a list of non negative integers, arrange them such that they form the largest number.
//
//        For example:
//
//        Given [3, 30, 34, 5, 9], the largest formed number is 9534330.
//
//        Note: The result may be very large, so you need to return a string instead of an integer.
class Comp implements Comparator<Integer> {

    int maxDigits;

    Comp(int maxDigits) {
        this.maxDigits = maxDigits;
    }
    public int compare(Integer o1, Integer o2) {
        Integer modifiedO1 = increaseDigits(o1, maxDigits);
        Integer modifiedO2 = increaseDigits(o2, maxDigits);
        return modifiedO2.compareTo(modifiedO1) > 1 ? 1 : modifiedO2.compareTo(modifiedO1) < 0 ? -1 : 0;
    }

    int getNumDigits(Integer num) {
        return String.valueOf(num).trim().length();
    }

    int increaseDigits(int num, int targetDigits) {
        int currDigits = getNumDigits(num);
        if (currDigits == targetDigits) {
            return num;
        } else if (currDigits < targetDigits){
            int digits = targetDigits - currDigits;
            int mul = (int)Math.pow(10, digits);
            return num*mul;
        }
        return -1;
    }
}
public class LargestNumber {

    public static void main(String[] args) {

        Integer[] arr = {3, 30, 34, 5, 9};

        (new LargestNumber()).largestNumber(Arrays.asList(arr));
    }


    //11:47
    public String largestNumber(final List<Integer> A) {

        List<Integer> numsWithIncreasedDigits = new LinkedList<Integer>();
        int maxNumDigits = getMaxNumDigitAndNumDigits(A);

        for (int i = 0; i < A.size() ; i++) {
            numsWithIncreasedDigits.add(increaseDigits(A.get(i), maxNumDigits));

        }
        Collections.sort(A, new Comp(maxNumDigits));
        StringBuffer ret = new StringBuffer();
        for (int i = 0; i < A.size() ; i++) {
            ret.append(A.get(i));
        }
        return ret.toString();
    }

    int getMaxNumDigitAndNumDigits( List<Integer> arr) {
        int maxNumDigits = Integer.MIN_VALUE;
        if (arr == null || arr.size() == 0) {
            return 0;
        }
        for (int i = 0; i < arr.size(); i++) {
            int digits = getNumDigits(arr.get(i));

            if (digits > maxNumDigits) {
                maxNumDigits = digits;
            }

        }
        return maxNumDigits;
    }
    int getNumDigits(Integer num) {
        return String.valueOf(num).trim().length();
    }

    int increaseDigits(int num, int targetDigits) {
       int currDigits = getNumDigits(num);
       if (currDigits == targetDigits) {
           return num;
       } else if (currDigits < targetDigits){
           int digits = targetDigits - currDigits;
           int mul = (int)Math.pow(10, digits);
           return num*mul;
       }
       return -1;
    }

    int reduceDigits(int num, int targetDigits) {
        int currDigits = getNumDigits(num);
        if (currDigits == targetDigits) {
            return num;
        } else if (currDigits > targetDigits){
            int digits = currDigits - targetDigits;
            int div = (int)Math.pow(10, digits);
            return num/div;
        }
        return -1;
    }


}
