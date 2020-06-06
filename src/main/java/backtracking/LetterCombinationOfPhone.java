package backtracking;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class LetterCombinationOfPhone {

    HashMap<Integer, String> map = new HashMap<>();

    public static void main(String[] args) {
        (new LetterCombinationOfPhone()).letterCombinations("74");
    }
    public List<String> letterCombinations(String digits) {
        List<String> lst = new ArrayList<>();
        if (digits == null || digits.trim().length() == 0) {
            return lst;
        }
        fillHashMap();

        getCombinations(0, new StringBuilder(""), digits, lst);
        return lst;


    }

    private void getCombinations(int curr, StringBuilder builder, String digits, List<String> lst) {


        if (curr == digits.length()) {
            lst.add(builder.toString());

            return;
        }
        int num = (int)digits.charAt(curr) - '0';
        String val = map.get(num);

        for (int i = 0; i < 3; i++) {

            builder.append(val.charAt(i));
            getCombinations(curr+1, builder, digits, lst);
            builder.delete(builder.length()-1,builder.length());

        }


    }

    private void fillHashMap() {
        map.put(1, "");
        map.put(2, "abc");
        map.put(3, "def");
        map.put(4, "ghi");
        map.put(5, "jkl");
        map.put(6, "mno");
        map.put(7, "pqrs");
        map.put(8, "tuv");
        map.put(9, "wxyz");
        map.put(0, " ");

    }
}
