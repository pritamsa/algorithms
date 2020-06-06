package backtracking;

public class RegexMatch {

    public boolean isMatch(String s, String p) {

        if (s == null || p == null) {
            return false;
        }

        s = s.trim();
        p = p.trim();

        if (allStarPattern(p)) {
            return true;
        }
        int stIdx = 0;
        int patIdx = 0;
        int str_idx_prev = -1;
        int pat_idx_prev = -1;

        while(stIdx < s.length() && patIdx < p.length()) {

            if (s.charAt(stIdx) == p.charAt(patIdx)) {
                stIdx++;
                patIdx++;
            } else if (p.charAt(patIdx) == '.') {
                stIdx++;
                patIdx++;

            } else if (p.charAt(patIdx) == '*') {
                str_idx_prev = stIdx;
                patIdx++;
                pat_idx_prev = patIdx;

            } else if (pat_idx_prev != -1) {
                patIdx = pat_idx_prev;
                stIdx++;
                str_idx_prev = stIdx;

            } else {
                return false;
            }

        }

        if (patIdx == p.length() && stIdx < s.length()) {
            return false;
        }

        if (patIdx < p.length() && stIdx == s.length()) {

            return allStarPattern(p.substring(patIdx));
        }

        return false;



    }
    private boolean allStarPattern(String p) {
        for (int k = 0; k < p.length(); k++) {
            if (p.charAt(k) != '*') {
                return false;
            }
        }
        return true;
    }
}
