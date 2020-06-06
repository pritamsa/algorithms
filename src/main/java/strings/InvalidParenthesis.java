package strings;

import java.util.HashSet;
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.LinkedBlockingQueue;

public class InvalidParenthesis {

    private static boolean isValidString(final String str) {
        int ct = 0;
        if (str != null && str.trim().length() > 0) {
            for (int i = 0; i < str.length(); i++) {
                if (str.charAt(i) == '(') {
                    ct++;
                } else if (str.charAt(i) == ')') {
                    ct--;
                }
                if (ct < 0) {
                    return false;
                }

            }

        }
        return ct == 0;
    }

    private static boolean isParenthesis(char c) {
        return (c == '(') || (c == ')');
    }

    public static LinkedList<String> removeInvalidParenthesis(String str) {

        if (str == null || str.trim().length() == 0 ) {
            return null;
        }
        HashSet<String> st = new HashSet<>();
        Queue<String> q = new LinkedBlockingQueue<>();
        LinkedList<String> validStrs = new LinkedList<>();

        q.add(str);
        st.add(str);

        boolean stringsAdded = false;
        while (!q.isEmpty()) {

            String s = q.remove();
            if (isValidString(s)) {
                validStrs.add(s);

            }

            if (!stringsAdded) {
                //Remove each char to get a new string. Keep track of visited strings in a hashset.
                for (int i = 1; i < s.length(); i++) {

                    String newSt = str.substring(i - 1, i) + str.substring(i + 1);
                    if (!st.contains(newSt)) {
                        q.add(newSt);
                        st.add(newSt);
                    }

                }
                stringsAdded = true;
            }
        }
        return validStrs;


    }

}
