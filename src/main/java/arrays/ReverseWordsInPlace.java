package arrays;

import java.util.Stack;

public class ReverseWordsInPlace {

    public static void main(String[] args) {
        char[] arr = {'t','h','e',' ','s','k','y',' ','i','s',' ','b','l','u','e'};
        ReverseWordsInPlace w = new ReverseWordsInPlace();
        w.reverseWords(arr);
    }

    public void reverseWords(char[] s) {
        reverse(s, 0, s.length - 1);
        int st = 0;
        int en = 0;


        while(en != s.length -1 ) {
            en = findEndIndexAfter(s, st);
            reverse(s, st, en);
            st = en+2;
        }
    }


    private int findEndIndexAfter(char[] s, int st) {

        int idx = st;
        while (idx < s.length && s[idx] != ' ') {
            idx++;
        }
        if (idx == s.length - 1) {
            return idx;
        }
        return idx - 1;
    }
    private void reverse(char[] s, int st, int en) {
        if (en < st) {
            return;
        }

        while(en >= st) {
            if (s[st] != s[en]) {
                char temp = s[en];
                s[en] = s[st];
                s[st] = temp;

            }
            st++;
            en--;

        }

    }

}
