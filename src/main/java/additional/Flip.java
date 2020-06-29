package additional;

import java.util.List;

//A flip makes 0s as 1s and 1s as 0s. Give a max length of flips yielding optimal number of ones.
public class Flip {

    public static void main(String[] args) {
        (new Flip()).flip("010");
    }
    public int maximumGap(final List<Integer> A) {

        if (A == null || A.size() == 1) {
            return 0;
        }
        int max = 0;

        for (int i = 0; i < A.size(); i++) {
            if (max < A.get(i)) {
                max = A.get(i);
            }
        }

        int[] arr = new int[max+1];

        for (int i = 0; i < max+1; i++) {
            arr[A.get(i)] = 1;
        }

        int maxM = 0;
        int m = 0;
        int k = 0;
        while(A.get(k) != 1) {
            k++;

        }
        for (int i = k; i < arr.length; i++) {
            if (arr[i] == 0) {
                m++;
            } else {
                if (maxM < m) {
                    maxM = m;

                }
                m = 0;
            }
        }
        return maxM;

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
