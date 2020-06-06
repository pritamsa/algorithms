package bits;

public class hammingWeight {

    public static void main(String[] args) {
        Integer k = 11;
        hammingWeight(11);
    }

    public static int hammingWeight(int n) {

        System.out.println(Integer.toBinaryString(n));
        int sum = 0;

        for (int i = 0; i < 32; i++) {
            sum += (n & (1<<i));
        }
        return sum;

    }
}
