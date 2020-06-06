package arrays.missingnumbers;

//fail safe vsfail fast
//collections does not implement serializable. why?
//
public class FindDuplicatesInImmutableArr {

    public static Integer findDuplicate(final Integer[] arr) {

        Integer slow = arr[0];
        Integer fast = arr[arr[0]];

        while(!slow.equals(fast) && slow < arr.length && fast < arr.length) {
            slow = arr[slow];
            fast = arr[arr[fast]];
        }
        fast = 0;
        while(!fast.equals(slow)) {
            slow = arr[slow];
            fast = arr[fast];
        }
        return slow;

    }

    public static void main(String[] args) {
        Integer[] arr = {1, 2, 3, 4, 5, 6, 3};
        Integer dup = findDuplicate(arr);
    }

}
