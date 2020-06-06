package sorting;

public class SelectionSort {

    public static void main(String[] args) {
        int[] arr = {30, 22, 27, 28, 12, 19, 17, 8, 3, 14, 5, 0};
        selectionSort(arr);
    }

    public static void selectionSort(int[] arr) {
        selectionSortUtil(arr, 0, arr.length-1);
    }

    private static void selectionSortUtil(int[] arr, int unsortedSt, int unsortedEn) {

        while(unsortedSt < unsortedEn) {
            int min = Integer.MAX_VALUE;
            int minIndex = -1;
            for (int i = unsortedSt; i <= unsortedEn; i++ ) {
                if (min > arr[i]) {
                    min = arr[i];
                    minIndex = i;
                }
            }
            int temp = arr[unsortedSt];
            arr[unsortedSt] = arr[minIndex];
            arr[minIndex] = temp;
            unsortedSt++;

        }
    }
}
