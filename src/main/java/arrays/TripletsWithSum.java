package arrays;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;

//Given an array of real numbers greater than zero in form of strings.
//        Find if there exists a triplet (a,b,c) such that 1 < a+b+c < 2 .
//        Return 1 for true or 0 for false.
//
//        Example:
//
//        Given [0.6, 0.7, 0.8, 1.2, 0.4] ,
//
//        You should return 1
//
//        as
//
//        0.6+0.7+0.4=1.7
//
//        1<1.7<2
//
//        Hence, the output is 1.
//
//        O(n) solution is expected.
//
//        Note: You can assume the numbers in strings don’t overflow the primitive data type and there are no leading zeroes in numbers. Extra memory usage is allowed.
//10:27
public class TripletsWithSum {
    public static void main(String[] args) {
        String[] arr = {"0.6", "0.7", "0.8", "1.2", "0.4"};

        int ret = solve(new ArrayList<String>(Arrays.asList(arr)));


    }

    public static int solve(ArrayList<String> A) {

        if (A == null) {
            return -1;
        }

        ArrayList<Float> nums = new ArrayList<Float>(A.size());
        for (int i = 0; i < A.size(); i++) {
            nums.add(Float.valueOf(A.get(i)));
        }

        Collections.sort(nums);

        for (int i = 0; i < nums.size() - 2; i++) {
            int j = i + 1;
            int k = nums.size() - 1;
            Float sum = 0F;
            while (j < k) {
                sum = nums.get(i) + nums.get(j) + nums.get(k);

                if (sum >= 2) {
                    k--;
                } else if (sum <= 1) {
                    j++;
                } else return 1;
            }
        }
        return -1;
    }
}
