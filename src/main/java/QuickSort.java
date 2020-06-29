import java.util.*;

public class QuickSort {

    public static void main(String[] args) {
        int[] arr = {2,3,2};
        //quickSort(arr, 0, arr.length-1 );
        //isAnagram("a", "b");

        //majorityElement(arr);

    }


        static int appearsNBy3(int arr[], int n)
        {
            int count1 = 0, count2 = 0;

            // take the integers as the maximum
            // value of integer hoping the integer
            // would not be present in the array
            int first =  Integer.MIN_VALUE;;
            int second = Integer.MAX_VALUE;

            for (int i = 0; i < n; i++) {

                // if this element is previously
                // seen, increment count1.
                if (first == arr[i])
                    count1++;

                    // if this element is previously
                    // seen, increment count2.
                else if (second == arr[i])
                    count2++;

                else if (count1 == 0) {
                    count1++;
                    first = arr[i];
                }

                else if (count2 == 0) {
                    count2++;
                    second = arr[i];
                }

                // if current element is different
                // from both the previously seen
                // variables, decrement both the
                // counts.
                else {
                    count1--;
                    count2--;
                }
            }

            count1 = 0;
            count2 = 0;

            // Again traverse the array and
            // find the actual counts.
            for (int i = 0; i < n; i++) {
                if (arr[i] == first)
                    count1++;

                else if (arr[i] == second)
                    count2++;
            }

            if (count1 > n / 3)
                return first;

            if (count2 > n / 3)
                return second;

            return -1;
        }

    public static List<Integer> majorityElement(int[] nums) {
        if (nums == null || nums.length == 0) {
            return null;
        }
        List<Integer> ret = new LinkedList();
        if (nums == null || nums.length == 0) {
            return null;
        }

        Map<Integer, Integer> map = new HashMap<>();

        for(int i = 0; i < nums.length; i++) {

            Integer val = map.get(nums[i]);
            if (val == null || val == 0) {
                val = 1;
            }else {
                val++;
            }
            map.put(nums[i], val);
        }

        Set<Integer> set = map.keySet();
        int maxCount = 0;
        int maxVal = 0;

        for(Integer i: set) {
            if (maxCount < map.get(i) ) {
                maxCount = map.get(i);
                maxVal = i;
            }
            if(map.get(i) >= nums.length/3) {
                if(map.get(i) == maxCount) {
                    ret.add(i);
                }
            }
        }
        return ret;
    }
        // Driver c



    public static boolean rotateString(String A, String B) {

        if (A == null || B == null) {
            return false;
        }
        A = A.trim();
        B = B.trim();
        if (A.compareTo(B) ==0) {
            return true;
        }
        StringBuilder temp = new StringBuilder();

        for (int i = 0; i < A.length() - 1; i++) {
            temp.delete(0, temp.length());
            temp.append(A.substring(i+1));
            temp.append(A.substring(0,i+1));
            if (temp.toString().compareTo(B) == 0) {
                return true;
            }
        }
        return false;
    }
    public static boolean isAnagram(String s, String t) {
        Map<Character, Integer> map = new HashMap<>();

        for (int i = 0; i < s.length(); i++) {

            Integer val = map.get(s.charAt(i));
            if (val != null) {
                val++;
            } else {
                val = 1;
            }
            map.put(s.charAt(i), val);

        }
        for (int i = 0; i < t.length(); i++) {
            Integer val = map.get(t.charAt(i));
            if (val == null || val == 0) {
                return false;
            }
            val--;
            if (val == 0) {
                map.remove(s.charAt(i));
            } else {
                map.put(s.charAt(i), val);
            }

        }
        return map.isEmpty();
    }

    public static void quickSort(int[] arr, int st, int en) {
        if (st >= en) {
            return;
        }
        int idx = partition(arr, st, en);
        quickSort(arr, st, idx-1 );
        quickSort(arr, idx+1, en );


    }


    public static int partition(int[] arr, int st, int en ) {

        int pivot = en;

        int i = st-1;
        int j = st;

        while (j < en) {
            if (arr[j] <= arr[pivot]) {
                i++;
                int temp = arr[j];
                arr[j] = arr[i];
                arr[i] = temp;

            }
            j++;
        }
        i++;
        //exchange with pivot
        int temp = arr[pivot];
        arr[pivot] = arr[i];
        arr[i] = temp;
        return i;

    }

}

