package additional;

public class Flip {

    public static void main(String[] args) {
        (new Flip()).flip("010");
    }
    public int[] flip(String A) {
        int[] ret = new int[2];
        int[][] dp = new int[A.length()][A.length()];

        //"1001110011"
        for (int i = 0; i < A.length(); i++) {
            dp[i][i] = A.charAt(i) == '0' ? 1 : 0;

        }
        int max = 0;

        for (int i = 0; i < A.length(); i++) {
            for (int j = i; j < A.length(); j++) {
                if (j > 0)
                { dp[i][j] = dp[i][j-1] + (A.charAt(j) == '0' ? 1 : 0 );}

            }

        }

        for (int i = 0; i < A.length(); i++) {
            for (int j = i; j < A.length(); j++) {
                if (max < dp[i][j]) {
                    max = dp[i][j];
                    ret[0] = i+1;
                    ret[1] = j+1;
                }
            }

        }
        return ret;

    }
}
