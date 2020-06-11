package strings;

public class ExcelColNameToNumber {
    public int titleToNumber(String s) {

        final int base = 26;

        int sum = 0;

        int i = 0;

        for (int j = s.length() - 1; j >=0; j--) {

            sum += Math.pow(26, i)*getValue(s.charAt(j));
            i++;

        }
        return sum;

    }

    private int getValue(char c) {
        final int val = c - 'A' + 1;
        return val;
    }
}
