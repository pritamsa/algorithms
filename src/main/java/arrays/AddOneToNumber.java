package arrays;

import java.lang.reflect.Array;
import java.util.ArrayList;
import java.util.Arrays;


//Given a non-negative number represented as an array of digits,
//
//        add 1 to the number ( increment the number represented by the digits ).
//
//        The digits are stored such that the most significant digit is at the head of the list.
//
//        Example:
//
//        If the vector has [1, 2, 3]
//
//        the returned vector should be [1, 2, 4]
//
//        as 123 + 1 = 124.
public class AddOneToNumber {

    public static void main(String[] args) {

        Integer[] arr = {9,9,9,9};

        ArrayList<Integer> A = new ArrayList<Integer>(Arrays.asList(arr));

        ArrayList<Integer> aPlusOne = plusOne(A);
    }

    public static ArrayList<Integer> plusOne(ArrayList<Integer> A) {
        ArrayList<Integer> ret = new ArrayList<Integer>();
        if (A == null || A.size() == 0) {
            return null;
        }


        int carry = 0;

        for (int i = A.size() - 1; i >=0; i-- ) {
            int val = 0;
            if (i ==  A.size() - 1) {
                val = A.get(i) + 1;

            } else {
                val = A.get(i) + carry;
            }
                carry = val/10;
                if (carry != 0) {
                    val = val % 10;
                }
                ret.add(0, val);

        }

        if (carry > 0) {
            ret.add(0,carry);
        }
        return ret;



    }
}

