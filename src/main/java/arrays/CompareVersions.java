package arrays;

public class CompareVersions {

    public static void main(String[] args) {
        (new CompareVersions()).compareVersion("1", "1.1");
    }

    public int compareVersion(String version1, String version2) {

        String[] arr1 = version1.split("\\.");
        String[] arr2 = version2.split("\\.");

        if (arr1.length > arr2.length) {
            arr2 = getLargeArr(arr2, arr1.length);
        } else if (arr1.length < arr2.length) {
            arr1 = getLargeArr(arr1, arr2.length);
        }


        for (int j = 0; j < arr1.length; j++) {
            if (compare(arr1[j], arr2[j]) == 1) {
                return 1;
            } else if (compare(arr1[j], arr2[j]) == -1) {
                return -1;
            }
        }
        return 0;
    }

    private String[] getLargeArr(String[] arr, int target) {
        String[] ret = new String[target];
        int i = 0;
        while (i < arr.length) {
            ret[i] = arr[i];
            i++;
        }
        while (i < target) {
            ret[i] = "0";
            i++;
        }
        return ret;
    }

    private int compare(String s1, String s2) {
        return Integer.parseInt(s1) > Integer.parseInt(s2) ? 1
                : Integer.parseInt(s1) == Integer.parseInt(s2) ? 0 : -1;

    }
}