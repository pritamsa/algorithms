package strings;

//Given a string A representing a roman numeral.
//        Convert A into integer.
//
//        A is guaranteed to be within the range from 1 to 3999.
//
//        NOTE: Read more
//        details about roman numerals at Roman Numeric System
//
//
//
//        Input Format
//
//        The only argument given is string A.
//        Output Format
//
//        Return an integer which is the integer verison of roman numeral string.
//        For Example
//
//        Input 1:
//        A = "XIV"
//        Output 1:
//        14
//
//        Input 2:
//        A = "XX"
//        Output 2:
//        20
//
public class RomanToInteger {

    public static void main(String[] args) {
        System.out.println((new RomanToInteger()).convertRomanToInteger("MCMXCIV"));
    }

    public int convertRomanToInteger(String romanStr) {

        int sum = 0;
        romanStr = romanStr.trim();
        if (romanStr.length() < 2) {
            return getValue(romanStr.charAt(0));
        }

        for (int i = 0; i < romanStr.length(); i++) {
            if(i < romanStr.length() - 1) {
                if(getValue(romanStr.charAt(i)) < getValue(romanStr.charAt(i+1))) {
                    sum -= getValue(romanStr.charAt(i));
                } else {
                    sum += getValue(romanStr.charAt(i));
                }
            } else {
                sum += getValue(romanStr.charAt(i));
            }
        }
        return sum;
    }



    public int getValue(final char romanLet) {
        switch (romanLet) {
            case 'I':
                return 1;
            case 'V':
                return 5;
            case 'X':
                return 10;
            case 'L':
                return 50;
            case 'C':
                return 100;
            case 'D':
                return 500;
            case 'M':
                return 1000;
        }
        return 0;
    }
}
