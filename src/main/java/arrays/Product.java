package arrays;

public class Product {

    public static void main(String args[]) {
        (new Product()).multiply("123456789",
                "987654321");
    }

    public String multiply(String num1, String num2) {
        long sum = 0;
        int t = 0;

        for (int i = num2.length() - 1; i >=0; i--) {
            sum += prod(num1, num2.charAt(i)) * Math.pow(10, t);
            t++;
        }
        return Long.toString(sum);

    }

    private long prod(String num1, char c) {
        int i = (int)(c-'0');

        long carry = 0;
        long sum = 0;
        int t = 0;
        for (int j = num1.length() - 1; j >=0; j--) {

            int k = (int)(num1.charAt(j)-'0');

            long add = (i*k)%10 + carry;
            sum += add*Math.pow(10,t);
            carry = (i*k)/10;
            t++;
        }
        return (long)(carry*Math.pow(10,t)) + sum;

    }
}
