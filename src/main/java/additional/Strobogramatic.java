package additional;

import java.util.Arrays;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;

public class Strobogramatic {

    public static void main(String[] args) {
        List<String> lst = findStrobogrammatic(4);
    }

    private static List<String> helper(int n, int end){
        if(n<1)
            return new LinkedList<>(Arrays.asList(""));

        if(n==1)
            return new LinkedList<String>(Arrays.asList("0","1","8"));;

        HashMap<Character, Character> pair = new HashMap<Character, Character>();
        pair.put('1','1'); pair.put('8','8'); pair.put('6','9'); pair.put('9','6');

        List<String> res = helper(n-2, end);

        int fixLen = res.size();

        for(int i = 0; i < fixLen; ++i){
            String curr = res.remove(0);
            if(n != end)
                res.add("0"+curr+"0");

            for(Character ch: pair.keySet()){
                res.add(ch+curr+pair.get(ch));
            }
        }

        return res;
    }
    public static List<String> findStrobogrammatic(int n) {

        //0 is the main case to handle here because caller fxn may need to return 00 which is needed by calling fxn

        return helper(n, n);

    }
}
