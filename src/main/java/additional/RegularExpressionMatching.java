package additional;
//find isFirstMatch for 1st char
//and then is length of pattern 2 or more and 2nd char is *, isFirstMatch && match(pattern, str.substring(1)) ||
//match(str, pattern.substring(2))
//because char* means match 0 or many chars.
public class RegularExpressionMatching {

    public static void main(String[] args) {
        System.out.println(isMatch("aa", "a*"));

    }

    //Correct one
    public static boolean isMatch(String text, String pattern) {
        if (pattern.isEmpty()) return text.isEmpty();
        boolean first_match = (!text.isEmpty() &&
                (pattern.charAt(0) == text.charAt(0) || pattern.charAt(0) == '.'));

        if (pattern.length() >= 2 && pattern.charAt(1) == '*'){
            return (isMatch(text, pattern.substring(2)) ||
                    (first_match && isMatch(text.substring(1), pattern)));
        } else {
            return first_match && isMatch(text.substring(1), pattern.substring(1));
        }
    }


}

