package strings;

import java.util.Collections;
import java.util.HashMap;
import java.util.Set;
import java.util.TreeMap;

//The only argument given is integer A.
//        Output Format
//
//        Return a string denoting roman numeral version of A.
//        Constraints
// 1. I, 2 II, 3 III IV V VI VII VIII IX X
//  10 : X
//        1 <= A <= 3999
//        For Example
//
//        Input 1:
//        A = 5
//        Output 1:
//        "V"
//
//        Input 2:
//        A = 14
//        Output 2:
//        "XIV"
public class IntegerToRoman {


    public static void main(String[] args) {
        System.out.println(intToRoman("5"));
    }
    public static String intToRoman(String numStr) {
        if (numStr != null) {
            Integer num = null;
            try {
              num = Integer.parseInt(numStr);
            } catch(NumberFormatException ex) {
                return "";
            }
            return getRomanNum(num);
        }
        return null;

    }


    private static String getRomanNum(int num) {
        int[] keys = {1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 3, 2, 1};
        String[] vals = {"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "III", "II", "I"};


        StringBuffer ret = new StringBuffer("");

        int i = 0;
        while (num > 0) {
            if (num >= keys[i]) {
                num -= keys[i];
                ret.append(vals[i]);
            } else {
                while (num < keys[i] && i < keys.length) {
                    i++;
                }
            }

        }

        return ret.toString();
    }
}
