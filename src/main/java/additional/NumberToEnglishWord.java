package additional;

//Convert a non-negative integer to its english words representation. Given input is guaranteed to be less than 231 - 1.
//
//        Example 1:
//
//        Input: 123
//        Output: "One Hundred Twenty Three"
//        Example 2:
//
//        Input: 12345
//        Output: "Twelve Thousand Three Hundred Forty Five"
//        Example 3:
//
//        Input: 1234567
//        Output: "One Million Two Hundred Thirty Four Thousand Five Hundred Sixty Seven"
//        Example 4:
//
//        Input: 1234567891
//        Output: "One Billion Two Hundred Thirty Four Million Five Hundred Sixty Seven Thousand Eight Hundred Ninety One"


public class NumberToEnglishWord {

    public static void main(String[] args) {
        System.out.println((new NumberToEnglishWord()).numberToWords(200));

    }

    public String numberToWords(int num) {

        int bil = numBillions(num);
        num = num % 1000000000;

        int mil = numMillions(num);
        num = num%1000000;

        int thousands = numThousands(num);
        num = num%1000;

        int hundreds = numHundreds(num);
        num = num%100;

        String st = "";
        int t = num/10;
        if (t > 1) {
            st += getStr2(t) + " ";
            num = num%10;
        }
        st += getStr(num);

        String ret = "";
        if (bil > 0) {
            ret += numberToWords(bil) + " billion ";
        }

        if (mil > 0) {
            ret += numberToWords(mil) + " million ";
        }

        if (thousands > 0) {
            ret += numberToWords(thousands) + " thousand ";
        }

        if (hundreds > 0) {
            ret += numberToWords(hundreds) + " hundred ";
        }

        ret += st;
        return ret;
    }

    private int numBillions(int val) {
        return val/1000000000;
    }
    private int numMillions(int val) {
        return val/1000000;
    }
    private int numThousands(int val) {
        return val/1000;
    }
    private int numHundreds(int val) {
        return val/100;
    }

    private String getStr2(int val) {
        switch (val) {

            case 2: return "twenty";
            case 3: return "thirty";
            case 4: return "forty";
            case 5: return "fifty";
            case 6: return "sixty";
            case 7: return "seventy";
            case 8: return "eighty";
            case 9: return "ninety";

        }
        return "";
    }

    public String getStr(int val) {
        switch (val) {
            case 1: return "one";
            case 2: return "two";
            case 3: return "three";
            case 4: return "four";
            case 5: return "five";
            case 6: return "six";
            case 7: return "seven";
            case 8: return "eight";
            case 9: return "nine";
            case 10: return "ten";
            case 11: return "eleven";
            case 12: return "twelve";
            case 13: return "thirteen";
            case 14: return "fourteen";
            case 15: return "fifteen";
            case 16: return "sixteen";
            case 17: return "seventeen";
            case 18: return "eighteen";
            case 19: return "nineteen";
        }
        return "";
    }







}
