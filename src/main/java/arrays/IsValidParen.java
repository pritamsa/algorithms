package arrays;

import java.util.Stack;

public class IsValidParen {

    public static void main(String[] args) {
        String s = "()";
        boolean valid = (new IsValidParen()).isValid(s);
    }

    public boolean isValid(String s) {

        if (s == null || s.trim().length() == 0) {
            return false;

        }

        Stack<Character> st = new Stack<>();

        for (int i = 0; i < s.length(); i++) {
            if (isOpeningBrace(s.charAt(i))) {
                st.push(s.charAt(i));
            }

            if (isClosingBrace(s.charAt(i))) {
                if (!st.isEmpty()) {
                    char t = st.pop();
                    if (t != s.charAt(i)) {
                        return false;
                    }
                } else {
                    return false;
                }

            }
        }
        if (!st.isEmpty()) {
            return false;
        }
        return true;

    }



    private boolean isOpeningBrace(char c) {
        return (c == '(' || c == '{' || c == '[');
    }

    private boolean isClosingBrace(char c) {
        return (c == ')' || c == '}' || c == ']');
    }

}
