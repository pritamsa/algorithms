package additional;

//Given two strings S and T, return if they are equal when both are typed into empty text editors. # means a backspace character.
//
//        Note that after backspacing an empty text, the text will continue empty.
//
//        Example 1:
//
//        Input: S = "ab#c", T = "ad#c"
//        Output: true
//        Explanation: Both S and T become "ac".
//        Example 2:
//
//        Input: S = "ab##", T = "c#d#"
//        Output: true
//        Explanation: Both S and T become "".
//        Example 3:
//
//        Input: S = "a##c", T = "#a#c"
//        Output: true
//        Explanation: Both S and T become "c".
//        Example 4:
//
//        Input: S = "a#c", T = "b"
//        Output: false
//        Explanation: S becomes "c" while T becomes "b".

public class BackspaceStrCompare {

    public static void main(String[] args) {
        System.out.println(backSpaceStrComp("", "#"));
    }

    public static boolean backSpaceStrComp(String str1, String str2) {
        int c = 0;
        int c1 = 0;

        int i = str1.length() - 1;
        int j = str2.length() - 1;

        while (i >= 0 && j >=0 ) {
            if (c1==0 && c == 0 && str1.charAt(i) == str2.charAt(j)  && str1.charAt(i) != '#' && str2.charAt(j) != '#') {
                i--;
                j--;
            } else if(str1.charAt(i) == '#' || str2.charAt(j) == '#') {
                if (str1.charAt(i) == '#') {
                    i--;
                    c++;
                }
                if (str2.charAt(j) == '#') {
                    j--;
                    c1++;
                }
            } else if (c > 0 || c1 > 0) {
                if (c > 0) {
                    i--;
                    c--;
                }
                if (c1 > 0) {
                    j--;
                    c1--;
                }
            }
            else {
                return false;
            }

        }

        while (i>=0) {
            if(c == 0 && str1.charAt(i) != '#') {
                return false;
            } else if (str1.charAt(i) == '#') {
                c++;
                i--;
            } else if (str1.charAt(i) != '#') {
                c--;
                i--;
            }
        }

        while (j>=0) {
            if(c1 == 0 && str2.charAt(j) != '#') {
                return false;
            } else if (str2.charAt(j) == '#') {
                c1++;
                j--;
            } else if (str2.charAt(j) != '#') {
                c1--;
                j--;
            }
        }

        if (i == -1 && j == -1) {
            return true;
        }


        return false;

    }

//    private int getNewPtr(String str, int i) {
//        int c = 0;
//        String suf = "";
//        if (i < str.length() - 1) {
//            suf = str.substring(i+1);
//        }
//        while(str.charAt(i) == '#') {
//            c++;
//            i--;
//            if (i == 0) {
//                return 0;
//            }
//
//        }
//        if (i+1 <= c) {
//            return 0;
//        }
//        i = i - c;
//        return i;
//
//    }
}
