package backtracking;

public class WildCardMatch {

    public static void main(String[] args) {
        String str = "adceb";
        String pat = "*a*b";
        boolean ret = (new WildCardMatch()).strmatch(str.toCharArray(), pat.toCharArray(),str.length(), pat.length() );
        System.out.println(ret);
    }

    boolean strmatch(char txt[], char pat[],
                     int n, int m)
    {
        // empty pattern can only
        // match with empty string.
        // Base Case :
        if (m == 0)
            return (n == 0);

        // step-1 :
        // initailze markers :
        int i = 0, j = 0, index_txt = -1,
                index_pat = -1;

        while (i < n)
        {

            // For step - (2, 5)
            if (j < m && txt[i] == pat[j])
            {
                i++;
                j++;
            }

            // For step - (3)
            else if (j < m && pat[j] == '?')
            {
                i++;
                j++;
            }

            // For step - (4)
            else if (j < m && pat[j] == '*')
            {
                index_txt = i;
                index_pat = j;
                j++;
            }

            // For step - (5)
            else if (index_pat != -1)
            {
                j = index_pat + 1;
                i = index_txt + 1;
                index_txt++;
            }

            // For step - (6)
            else
            {
                return false;
            }
        }

        // For step - (7)
        while (j < m && pat[j] == '*')
        {
            j++;
        }

        // Final Check
        if (j == m)
        {
            return true;
        }

        return false;
    }

    public boolean isMatch(String s, String p) {

        if (s == null || p == null) {
            return false;
        }
        s = s.trim();
        p = p.trim();

        if (p.equals("*") && s != null) {
            return true;
        }

        if ((p.length() == 0 && s.length() != 0) || (p.length() != 0 && s.length() == 0 && !allStarPattern(p))) {
            return false;
        }
        if(p.equals("?") && s.length() == 1) {
            return true;
        }

        if (p.equals("s")) {
            return true;
        }

        int i = 0;
        int j = 0;

        if (s.length() > 0) {
            if (s.charAt(i) == p.charAt(j) || (p.charAt(j) == '?')) {
                if (i == s.length() - 1 && j == (p.length() - 1)) {
                    return true;
                }
                String newStr = s.substring(i + 1);
                String newPat = p.substring(j + 1);
                return isMatch(newStr, newPat);

            } else if (p.charAt(j) == '*') {
                if (j == p.length() - 1) {
                    return true;
                }
                int st = i;
                String newStr = s.substring(st);
                String newPat = p.substring(j + 1);
                boolean remainingMatch = isMatch(newStr, newPat);
                while (!remainingMatch && st < s.length()) {
                    st++;
                    newStr = s.substring(st);
                    remainingMatch = isMatch(newStr, newPat);
                }
                if (remainingMatch) {
                    return true;
                } else {
                    return allStarPattern(newPat);
                }


            }
        } else {
            return allStarPattern(p);
        }

        return false;

    }

    private boolean allStarPattern(String p) {
        for (int k = 0; k < p.length(); k++) {
            if (p.charAt(k) != '*') {
                return false;
            }
        }
        return true;
    }
}
