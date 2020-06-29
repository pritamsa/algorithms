package additional;

//Given a string which contains only lowercase letters, remove duplicate letters so that every letter
//        appears once and only once.
//        You must make sure your result is the smallest in lexicographical order among all possible results.
//Idea is simple. First, get count array to find count for each char.
//Then start iterating the string. start pos = 0
//For each string char loc, if the char is less than the charAt(pos), pos = string char loc
//decrement cnt for string char loc. If it reaches zero, break out of the loop
//recursively call: str.charAt(pos) + samefunction(string after pos with char at 'pos' removed from it);
public class RemoveDuplicates {

    public static void main(String[] args) {
        System.out.println(removeDuplicates1("llpqzqtatm"));

    }

    public static String removeDuplicates1(String str) {

        int pos = 0;

        int[] cnt = new int[26];

        for (int i = 0; i < str.length(); i++) {
            int j = str.charAt(i) - 'a';
            cnt[j]++;
        }

        //"llpqzqtatm" : cnt[l] = 2, cnt[p] = 1, cnt[q] = 2 cnt[z] = 1 cnt[t] : 2 cnt[a] : 1 cnt[m] :1
        for (int i = 0; i < str.length(); i++) {
            int j = str.charAt(i) - 'a';
            //this is for lexicographical order
            if (str.charAt(pos) > str.charAt(i)) {
                pos = i;
            }
            cnt[j]--;
            if(cnt[j] == 0) break;
        }
        if (pos == str.length() - 1) return str.charAt(pos) + "";
        return str.charAt(pos) + removeDuplicates1(str.substring(pos+1).replaceAll(str.charAt(pos) +"",
                ""));
    }
    public static String removeDuplicates(String str) {

        int[] cnt = new int[26];

        for (int i = 0; i < str.length(); i++) {
            int loc = str.charAt(i) - 'a';
            cnt[loc]++;
        }

        int pos = 0;
        for (int i = 0; i < str.length(); i++) {
            int loc = str.charAt(i) - 'a';
            if(str.charAt(pos) > str.charAt(i)) {
                pos = i;
            }
            cnt[loc]--;
            if(cnt[loc] == 0) break;

        }
        if (pos == str.length() - 1) {
            return "" +  str.charAt(pos);
        }
        String newStr =  str.charAt(pos)
                + removeDuplicates(str.substring(pos + 1).replaceAll(str.charAt(pos) +"", ""));

        return newStr;


    }

}
