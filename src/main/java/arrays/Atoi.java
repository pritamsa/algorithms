package arrays;

public class Atoi {

    public static void main(String[] args) {
        myAtoi("-91283472332");

    }
    public static int myAtoi(String str) {

        str = str.trim();

        if (!Character.isDigit(str.charAt(0)) && str.charAt(0) != '+' && str.charAt(0) != '-') {
            return 0;
        }
        boolean negative = false;
        if (str.charAt(0) == '-') {
            negative = true;
        }
        int st = 0;
        if (!Character.isDigit(str.charAt(0))) {
            st = 1;
        }
        int en = st;
        while (en < str.length() && Character.isDigit(str.charAt(en))) {
            en++;
        }
        en--;
        int sum = 0;
        int p = 0;
        while(en >= st) {
            int v = str.charAt(en) - '0';
            sum += v*Math.pow(10,p);
            en--;
            p++;
            if (negative && (-sum) <= Integer.MIN_VALUE) {
                return Integer.MIN_VALUE;
            }
            if (!negative && sum >= Integer.MAX_VALUE) {
                return Integer.MAX_VALUE;
            }

        }
        if (negative) {
            if (sum >= Integer.MAX_VALUE) {
                return Integer.MIN_VALUE;
            }
            sum = -sum;
        }
        return sum;


    }
}
