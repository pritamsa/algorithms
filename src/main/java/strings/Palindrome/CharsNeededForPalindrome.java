package strings.Palindrome;

public class CharsNeededForPalindrome {
    public static int solve(String A) {

        int j = A.length() - 1;

        while(true) {
            while (A.charAt(j) != A.charAt(0)) {
                j--;
            }
            if (j == 0) {
                return A.length() - 1;
            }
            System.out.println(A.substring(1, j) + " " + j);
            if (isPalindrome(A.substring(1, j))) {
                break;
            } else {
                j--;
            }
        }
        return A.length() - 1 - j;
    }

    private static boolean isPalindrome(String A) {
        if (A == null || A.trim().length() == 0) {
            return false;

        }
        if (A.length() == 1) {
            return true;
        }
        if (A.length() == 2) return (A.charAt(0) == A.charAt(1));
        return (A.charAt(0) == A.charAt(A.length() - 1)) &&
                isPalindrome(A.substring(1,A.length() - 1));
    }

    public static void main(String[] args) {
        int val = getMinCharToAddedToMakeStringPalin("AACECAAAA");
    }

    static int getMinCharToAddedToMakeStringPalin(String str)
    {
        StringBuilder s = new StringBuilder();
        s.append(str);

        // Get concatenation of string, special character
        // and reverse string
        String rev = s.reverse().toString();
        s.reverse().append("$").append(rev);

        // Get LPS array of this concatenated string
        int lps[] = computeLPSArray(s.toString());
        return str.length() - lps[s.length() - 1];
    }
    // returns vector lps for given string str
    public static int[] computeLPSArray(String str)
    {
        int n = str.length();
        int lps[] = new int[n];
        int i = 1, len = 0;
        lps[0] = 0; // lps[0] is always 0

        while (i < n)
        {
            if (str.charAt(i) == str.charAt(len))
            {
                len++;
                lps[i] = len;
                i++;
            }
            else
            {
                // This is tricky. Consider the example.
                // AAACAAAA and i = 7. The idea is similar
                // to search step.
                if (len != 0)
                {
                    len = lps[len - 1];

                    // Also, note that we do not increment
                    // i here
                }
                else
                {
                    lps[i] = 0;
                    i++;
                }
            }
        }
        return lps;
    }
}
