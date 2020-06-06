package strings;

//Compare two version numbers version1 and version2.
//
//        If version1 > version2 return 1,
//        If version1 < version2 return -1,
//        otherwise return 0.
//        You may assume that the version strings are non-empty and contain only digits and the . character.
//        The . character does not represent a decimal point and is used to separate number sequences.
//        For instance, 2.5 is not "two and a half" or "half way to version three", it is the fifth second-level revision of the second first-level revision.
//
//        Here is an example of version numbers ordering:
//
//        0.1 < 1.1 < 1.2 < 1.13 < 1.13.4

public class CompareVersions {

    public static void main(String[] args) {
        int comparison = compareVersions("0", "0.0");
        System.out.println(comparison);
    }

    public static int compareVersions(String v1, String v2) {

        int v1Dots = dotCount(v1);
        int v2Dots = dotCount(v2);

        if (v1Dots > v2Dots) {
            v2 = addDots(v2, v1Dots - v2Dots);

        } else if (v2Dots > v1Dots) {
            v1 = addDots(v1, v2Dots - v1Dots);
        }

        String[] v1Vals = v1.split("\\.");

        String[] v2Vals = v2.split("\\.");

        for (int i = 0; i < v1Vals.length ; i++) {
            Integer v1Val = Integer.parseInt(v1Vals[i]);
            Integer v2Val = Integer.parseInt(v2Vals[i]);
            if (v1Val > v2Val) {
                return 1;
            }else if (v2Val > v1Val) {
                return 2;
            }
        }
        return 0;

    }

    private static String addDots(String s1, int target) {

        StringBuffer ret = new StringBuffer(s1);
        for (int i = 0; i < target; i++) {
            ret.append(".").append("0");

        }
        return ret.toString();
    }

    private static int dotCount(String s1) {
        int sum = 0;
        for (int i = 0; i < s1.length(); i++) {
            if (s1.charAt(i) == '.') {
              sum+=1;
            }
        }
        return sum;
    }
}
