package strings;
//Given an expression string exp , write a program to examine whether the pairs and the orders of “{“,”}”,”(“,”)”,”[“,”]” are correct in exp.
public class MatchingParenthesis {

    public static boolean isBalanced(String str) {
        int ct = 0;
        if (str != null && !str.trim().isEmpty()) {
            for (int i = 0; i < str.length(); i++) {
                if (str.charAt(i) == '(' ) {
                    ct++;
                } else if (str.charAt(i) == ')') {
                    ct--;
                }
                if (ct < 0) {
                    return false;
                }
            }
            return ct == 0;
        }
        return false;
    }

    public static void main(String[] args) {
        boolean bal = isBalanced("()");
    }

}
