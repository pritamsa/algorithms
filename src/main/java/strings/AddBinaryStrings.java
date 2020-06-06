package strings;

//Given two binary strings, return their sum (also a binary string).
//
//        Example:
//
//        a = "100"
//
//        b = "11"
//        Return a + b = “111”.
//1111+ 111 =    10110
public class AddBinaryStrings {
    public static void main (String[] args) {
        String str1 = "1111";
        String str2 ="11111";
        String str3 = addBinary(str1, str2);
    }

    public static String addBinary( String str1,  String str2 ) {

        StringBuffer result = new StringBuffer("");

        if (str1.length() > str2.length()) {
            str2 = getPaddedStr(str2, str1.length());

        } else if (str2.length() > str1.length()) {
            str1 = getPaddedStr(str1, str2.length());

        }
        int carry = 0;
        for (int i = str2.length() - 1; i >= 0 ; i--) {
            int sum = Character.getNumericValue(str2.charAt(i))
                    + Character.getNumericValue(str1.charAt(i)) + carry;
            int rem = sum % 2;
            carry = sum/2;
            result.append(rem);

        }
        if (carry == 0) return result.reverse().toString();
        else return (carry + result.reverse().toString());
    }

    public static String getPaddedStr(String str, int targetLen) {
        int addedZeros = targetLen - str.length();
        for (int i = 0; i < addedZeros; i++) {
            str = '0' + str;

        }
        return str;
    }
}
