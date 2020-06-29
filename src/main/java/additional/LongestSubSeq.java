package additional;

import java.util.Arrays;
import java.util.List;

public class LongestSubSeq {

    public static void main(String[] args) {
        Integer[] arr = {1, 11, 2, 10, 4, 5, 2, 1};
        int l = longestSubsequenceLength(Arrays.asList(arr));
    }

    public int longestSubsequenceLength1(final List<Integer> A) {
        if (A == null || A.size() == 0) {
            return 0;
        }
        if (A.size() == 1) {
            return 1;
        }

        int[] lcs = new int[A.size()];
        int[] lds = new int[A.size()];

        for (int i = 0; i < A.size(); i++) {
            lcs[i] = 1;
        }

        for (int i = 0; i < A.size(); i++) {
            lds[i] = 1;
        }

        int longestLcsLen = 1;
        int longestLdsLen = 1;
        int longestSuLen = 1;

        for (int i = 1; i < A.size(); i++) {
            for (int j = 0; j < i; j ++) {
                if (A.get(i) > A.get(j)) {
                    lcs[i] = Math.max(lcs[j] + 1, lcs[i]);
                }

            }

        }

        for (int i = A.size() - 2; i >= 0; i--) {
            for (int j = A.size() - 1; j > i; j--) {
                if (A.get(i) > A.get(j)) {
                    lds[i] = Math.max(lds[j] + 1, lds[i]);
                }

            }

        }

        for (int i = 0; i < lcs.length; i++) {
            if (longestSuLen > lcs[i] + lds[i] - 1) {
                longestSuLen = lcs[i] + lds[i] - 1;
            }

        }

        return longestSuLen;

    }

    //longest continuous sunse that is 1st increase and then decreasing.
    public static int longestSubsequenceLength(final List<Integer> A) {
        if (A == null || A.size() == 0) {
            return 0;
        }
        if (A.size() == 1) {
            return 1;
        }

        int c1 = 0;
        int c2 = 0;
        int max_ct = 0;

        for(int i = 1; i < A.size(); i++) {
            if(A.get(i) > A.get(i-1)) {
                if(c2 != 0 && c1 != 0) {
                    c1=0; c2=0;
                }
                c1++;
                if(max_ct < c1+c2) max_ct = c1+c2;

            } else if (A.get(i) < A.get(i-1)) {

                if(c1 > 0) {
                    c2++;
                    if(max_ct < c1+c2) max_ct = c1+c2;
                }
            }
        }
        return max_ct + 1;

    }
}
