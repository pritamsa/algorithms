package additional;

public class WildCardMatching {

    public static void main(String[] args) {
        System.out.println(match("mississippi", "m??*ss*?i*pi"));
       /* System.out.println(strmatch("mis sissippi".toCharArray(), "m?? *ss*?i*pi".toCharArray(),
                "mississippi".length(), "m??*ss*?i*pi".length()));*/
    }


    public static boolean match(String s, String p) {
        int i = 0;
        int j = 0;
        int sIdx = -1;
        int starIdx = -1;

        while(i < s.length()) {
            if (j < p.length() && s.charAt(i) == p.charAt(j)) {
                i++;
                j++;
            } else if (j < p.length() && p.charAt(j) == '?') {
                i++;
                j++;
            } else if (j < p.length() && p.charAt(j) == '*') {
                sIdx = i;
                starIdx = j;
                j++;
            } else if (starIdx != -1) {
                i = sIdx + 1;
                sIdx = i;
                j = starIdx+1;
            } else {
                return false;
            }



        }

        while (j < p.length()) {
            if (p.charAt(j) != '*') {
                return false;
            }
            j++;

        }
        return true;

    }

    static boolean strmatch(char txt[], char pat[],
                     int n, int m) {
        // empty pattern can only
        // match with empty string.
        // Base Case :
        if (m == 0)
            return (n == 0);

        // step-1 :
        // initailze markers :
        int i = 0, j = 0, index_txt = -1,
                index_pat = -1;

        while (i < n) {

            // For step - (2, 5)
            if (j < m && txt[i] == pat[j]) {
                i++;
                j++;
            }

            // For step - (3)
            else if (j < m && pat[j] == '?') {
                i++;
                j++;
            }

            // For step - (4)
            else if (j < m && pat[j] == '*') {
                index_txt = i;
                index_pat = j;
                j++;
            }

            // For step - (5)
            else if (index_pat != -1) {
                j = index_pat + 1;
                i = index_txt + 1;
                index_txt++;
            }

            // For step - (6)
            else {
                return false;
            }
        }

        // For step - (7)
        while (j < m && pat[j] == '*') {
            j++;
        }

        // Final Check
        if (j == m) {
            return true;
        }

        return false;
    }
}
