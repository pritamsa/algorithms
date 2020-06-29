package additional;

import java.util.LinkedList;
import java.util.List;
//valid ip address ix in the form w.x.y.z where w and x and y and z are > 0 & < 255
public class GenerateAllIValidp {

    public List<String> generateAllIps(String str) {
        String sNew = str;
        List<String> ret = new LinkedList<>();

        for (int i = 1; i < str.length() - 2; i++) {
            for (int j = i+1; j < str.length() - 1; j++) {
                for (int k = j+1; k < str.length();k++ ) {
                    sNew += sNew.substring(0,k) + "." + sNew.substring(k);
                    sNew += sNew.substring(0,j) + "." + sNew.substring(j);
                    sNew += sNew.substring(0,i) + "." + sNew.substring(i);
                    if (isValidIp(sNew)) {
                        ret.add(sNew);
                    }
                }
            }

        }
        return ret;
    }

    private boolean isValidIp(String str) {
        String[] nums = str.split("\\.");
        if(str.contains(".0")) {
            return false;
        }
        for (String num: nums) {
            Integer val = Integer.parseInt(num);
            if (val < 0 || val > 255) {
                return false;
            }
        }
        return true;
    }
}
