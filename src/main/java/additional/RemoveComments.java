package additional;
//Given a C++ program, remove comments from it. The program source is an array where source[i] is the i-th line of the source code. This represents the result of splitting the original source code string by the newline character \n.
//
//        In C++, there are two types of comments, line comments, and block comments.
//
//        The string // denotes a line comment, which represents that it and rest of the characters to the right of it in the same line should be ignored.
//
//        The string /* denotes a block comment, which represents that all characters until the next (non-overlapping) occurrence of */ should be ignored. (Here, occurrences happen in reading order: line by line from left to right.) To be clear, the string /*/ does not yet end the block comment, as the ending would be overlapping the beginning.
//
//The first effective comment takes precedence over others: if the string // occurs in a block comment, it is ignored. Similarly, if the string /* occurs in a line or block comment, it is also ignored.
//
//If a certain line of code is empty after removing comments, you must not output that line: each string in the answer list will be non-empty.
//
//There will be no control characters, single quote, or double quote characters. For example, source = "string s = "/* Not a comment. */";" will not be a test case. (Also, nothing else such as defines or macros will interfere with the comments.)
//
//        It is guaranteed that every open block comment will eventually be closed, so /* outside of a line or block comment always starts a new comment.
//
//Finally, implicit newline characters can be deleted by block comments. Please see the examples below for details.
//
//After removing the comments from the source code, return the source code in the same format.
//
//Example 1:
//Input:
//source = ["/*Test program */", "int main()", "{ ", "  // variable declaration ", "int a, b, c;", "/* This is a test", "   multiline  ", "   comment for ", "   testing */", "a = b + c;", "}"]
//
//        The line by line code is visualized as below:
//        /*Test program */
//        int main()
//        {
//        // variable declaration
//        int a, b, c;
///* This is a test
//   multiline
//   comment for
//   testing */
//        a = b + c;
//        }
//
//        Output: ["int main()","{ ","  ","int a, b, c;","a = b + c;","}"]
//
//        The line by line code is visualized as below:
//        int main()
//        {
//
//        int a, b, c;
//        a = b + c;
//        }
//
//        Explanation:
//        The string /* denotes a block comment, including line 1 and lines 6-9. The string // denotes line 4 as comments.
//Example 2:
//Input:
//source = ["a/*comment", "line", "more_comment*/b"]
//        Output: ["ab"]
//        Explanation: The original source string is "a/*comment\nline\nmore_comment*/b", where we have bolded the newline characters.  After deletion, the implicit newline characters are deleted, leaving the string "ab", which when delimited by newline characters becomes ["ab"].
//        Note:
//
//        The length of source is in the range [1, 100].
//        The length of source[i] is in the range [0, 80].
//        Every open block comment is eventually closed.
//        There are no single-quote, double-quote, or control characters in the source code.
//
//        Prev
//        3 / 4
//
//        Next
//
//        Autocomplete
//
//
//
//        )
//        1

import java.util.ArrayList;
import java.util.List;

public class RemoveComments {

    public static void main(String[] args) {
        String[] source = {
                "void func(int k) {", "// this function does nothing /*", "   k = k*2/4;", "   k = k/2;*/", "}"
        };
        List<String> ret = new RemoveComments().removeComments(source);

    }
    public double findMedianSortedArrays(int[] A, int[] B) {
        int m = A.length;
        int n = B.length;
        if (m > n) { // to ensure m<=n
            int[] temp = A; A = B; B = temp;
            int tmp = m; m = n; n = tmp;
        }
        int iMin = 0, iMax = m, halfLen = (m + n + 1) / 2;
        while (iMin <= iMax) {
            int i = (iMin + iMax) / 2;
            int j = halfLen - i;
            if (i < iMax && B[j-1] > A[i]){
                iMin = i + 1; // i is too small
            }
            else if (i > iMin && A[i-1] > B[j]) {
                iMax = i - 1; // i is too big
            }
            else { // i is perfect
                int maxLeft = 0;
                if (i == 0) { maxLeft = B[j-1]; }
                else if (j == 0) { maxLeft = A[i-1]; }
                else { maxLeft = Math.max(A[i-1], B[j-1]); }
                if ( (m + n) % 2 == 1 ) { return maxLeft; }

                int minRight = 0;
                if (i == m) { minRight = B[j]; }
                else if (j == n) { minRight = A[i]; }
                else { minRight = Math.min(B[j], A[i]); }

                return (maxLeft + minRight) / 2.0;
            }
        }
        return 0.0;
    }


    public List<String> removeComments(String[] source) {

        boolean inBlock = false;
        StringBuilder newline = new StringBuilder();
        List<String> ans = new ArrayList();
        for (String line: source) {
            int i = 0;
            char[] chars = line.toCharArray();
            if (!inBlock) newline = new StringBuilder();
            while (i < line.length()) {
                if (!inBlock && i+1 < line.length() && chars[i] == '/' && chars[i+1] == '*') {
                    inBlock = true;
                    i++;
                } else if (inBlock && i+1 < line.length() && chars[i] == '*' && chars[i+1] == '/') {
                    inBlock = false;
                    i++;
                } else if (!inBlock && i+1 < line.length() && chars[i] == '/' && chars[i+1] == '/') {
                    break;
                } else if (!inBlock) {
                    newline.append(chars[i]);
                }
                i++;
            }
            if (!inBlock && newline.length() > 0) {
                ans.add(new String(newline));
            }
        }
        return ans;
    }
     /*/ declare members/**/
    private int findEndIdx(String str, int st) {
        while (st < (str.length()-1)) {
            if(str.charAt(st) == '*' && str.charAt(st+1) == '/') {
                return st+2;
            }
        }
        return str.length()-1;
    }
}
