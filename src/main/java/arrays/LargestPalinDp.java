package arrays;

import java.util.*;

public class LargestPalinDp {
    public static void main(String[] args) {
        new LargestPalinDp().longestPalindrome("babad");
    }
    private boolean isPalin(String s) {
        HashMap<String, String> map = new HashMap<>();
        Set<Map.Entry<String,String>> set =  map.entrySet();
        ArrayList<ArrayList<String>> ret = new ArrayList<>();
        for (Map.Entry<String,String> entry: set) {

        }
        int st = 0;
        int en = s.length() - 1;
        while(st <= en) {
            if (s.charAt(st) != s.charAt(en)) {
                return false;
            }
            en--;
            st++;
        }
        return true;
    }
    public String longestPalindrome(String s) {

        int n = s.length();
        boolean[][] dp = new boolean[s.length()][s.length()];

        for (int i = 0; i < s.length(); i++) {
            dp[i][i] = true;

        }

        s = s+ " ";
        for (int L = 2; L <= n; L++) {
            for (int i = 0; i <= n - L; i++) {
                int j = i+L - 1;

                if (i != j) {
                    if (j == i+1) {
                        if (s.charAt(i) == s.charAt(j)) {
                            dp[i][j] = true;
                        }
                    } else {
                        dp[i][j] = (s.charAt(i) == s.charAt(j) && dp[i+1][j-1]);

                    }

                }


            }
        }

        int maxLen = 0;
        int st = 0;
        int en = 0;

        for (int L = 1; L <= n; L++) {
            for (int i = 0; i <= n - L; i++) {
                int j = i + L - 1;
                if (dp[i][j]) {
                    if (maxLen < L) {
                        maxLen = L;
                        st = i;
                        en = j + 1;
                    }
                }
            }
        }



        return s.substring(st,en);

    }

}
