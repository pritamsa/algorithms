package strings;

public class ExcelColNameToNumber {

    public String numToTitle(int n) {
        if (n <=0 ) {
            return "";
        }
        StringBuilder columnName = new StringBuilder();

        while (n > 0) {
            // Find remainder
            int rem = n % 26;

            // If remainder is 0, then a
            // 'Z' must be there in output
            if (rem == 0) {
                columnName.append("Z");
                n = (n / 26) - 1;
            }
            else // If remainder is non-zero
            {
                columnName.append((char)((rem - 1) + 'A'));
                n = n / 26;
            }
        }

        return columnName.reverse().toString();
    }
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
