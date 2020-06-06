package backtracking;

import org.apache.commons.lang3.StringUtils;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class Solution {

    public static void main(String[] args) {

        //boolean valid = (new Solution()).validTime("01");
        //System.out.println(valid);
        String[] words = {"a"};
        String word = "AAAAA";
        boolean sep = (new Solution()).camelCaseSeparation(words, word);
    }

    boolean validTime(String time) {

        if (time == null || time.trim().length() == 0) {
            return false;
        }

        String[] timeRepresentation = time.split(":");

        if (timeRepresentation.length == 2) {
            return isValidHourAndMinutes(timeRepresentation[0], timeRepresentation[1]);
        }
        return false;


    }


    boolean isValidHourAndMinutes(String hr, String minutes) {

        Integer numHr = Integer.parseInt(hr);
        Integer numMinutes = Integer.parseInt(minutes);

        if (numHr == 24) {
            return numMinutes == 0;
        }

        return numHr >= 0 && numHr <= 24 && numMinutes >= 0 && numMinutes < 60;


    }

//    boolean camelCaseSeparation(String[] words, String variableName) {
//
//        if (variableName == null || variableName.trim().length() == 0 || words == null || words.length == 0) {
//            return false;
//        }
//
//        Set<String> set = new HashSet<>();
//
//        for (int i = 0; i < words.length; i++) {
//            set.add(words[i]);
//        }
//
//        //String[] values = //splitCamelCase(variableName).split(" ");//StringUtils.splitByCharacterTypeCamelCase(variableName);
//
//        List<String> lst = getWords(variableName);
////        if (values == null || values.length == 0) {
////            return false;
////        }
//        if (lst == null || lst.size() == 0) {
//            return false;
//        }
//        for (String w : lst) {
//            if (!set.contains(w.toLowerCase())) {
//                return false;
//           }
//            }
////        for (int i = 0; i < lst.size(); i++) {
////            if (!set.contains(values[i].toLowerCase())) {
////                return false;
////            }
////        }
//            return true;
//
//        }


    boolean camelCaseSeparation(String[] words, String variableName) {
        if (variableName == null || variableName.trim().length() == 0 || words == null || words.length == 0) {
            return false;
        }

        Set<String> set = new HashSet<>();

        for (int i = 0; i < words.length; i++) {
            set.add(words[i]);
        }

        String[] values = splitCamelCase(variableName).split(" ");//StringUtils.splitByCharacterTypeCamelCase(variableName);

        if (values == null || values.length == 0) {
            return false;
        }
        for (int i = 0; i < values.length; i++) {
            if (!set.contains(values[i].toLowerCase())) {
                return false;
            }
        }
        return true;

    }


    List<String> getWords(String str) {

        int prev = 0;
        List<String> lst = new ArrayList<>();

        for (int i = 0; i < str.length(); i++) {
            int v = str.charAt(i) - '@';
            if ((v >= 1 && v <= 26) || i == str.length() - 1) {
                lst.add(str.substring(prev, i));
                prev = i;
            }

        }

        return lst;
    }

    String splitCamelCase(String s) {
        return s.replaceAll(
                String.format("%s|%s|%s",
                        "(?<=[A-Z])(?=[A-Z][a-z])",
                        "(?<=[^A-Z])(?=[A-Z])",
                        "(?<=[A-Za-z])(?=[^A-Za-z])"
                ),
                " "
        );
    }
}

