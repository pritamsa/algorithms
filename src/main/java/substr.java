import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Stack;

public class substr {

  public static void main(String[] args) {
    isPalindrome("A man, a plan, a canal: Panama");
  }

  public static boolean isPalindrome(String s) {

    int st = 0;
    int en = s.length() - 1;
    s = s.toLowerCase();

    while (st <= en ) {
      if ( !Character.isAlphabetic(s.charAt(st))  && !Character.isDigit(s.charAt(st))) {
        st++;
      } else if ( !Character.isAlphabetic(s.charAt(en))  && !Character.isDigit(s.charAt(en))) {
        en--;
      } else if(s.charAt(en) != s.charAt(st)) {
        return false;
      } else {
        en--;
        st++;
      }

    }
    return true;
  }
  //a b c d e fghtr cab : target c ab //kmp algorithm
  public boolean isSubStr(String source, String target) {

    int i = 0;
    int j = 0;
    boolean previousMatch = false;
    Stack<Integer> cou = new Stack<>();


    Map<Integer, Integer> arr;

    List<Integer> altPaths = new ArrayList<>();

    altPaths.stream().max(Integer::compare).get();
    //O(n)
    while (i < source.length()) {

      if (source.charAt(i) == source.charAt(j)) {
        //
        if (j == 0) {
          previousMatch = true;
        }
        i++;
        if (previousMatch) {
          if(j == target.length() - 1) {
            return true;
          }
          j++;
        }

      } else {
        if (previousMatch) {
          previousMatch = false;
        }
        i++;
      }


    }

    return false;

  }

}
