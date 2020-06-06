package sorting;


public class MergeSort {

    public static void main(String[] args) {

        int[] arr = {8,4,7,3,1,2,6,5,10, 15};

        sort(arr, 0, arr.length-1);


    }

    public static void sort(int[] arr, int l, int r) {
        if (l < r) {
            int m = (l + r) / 2;


            sort(arr, l, m );
            sort(arr, m + 1, r);

            mgr(arr, l, r, m);
        }

    }

    public static void mgr(int[] arr, int l, int r, int m) {

        int[] L = new int[m - l + 1];
        int[] R = new int[r - m];

        int j = 0;
        for (int i = 0; i < L.length ; ++i) {
            L[j++] = arr[l+i];
        }
        int k = 0;
        for (int i = 0; i < R.length ; i++) {
            R[k++] = arr[i + m + 1];
        }
        j = 0;
       k = 0;
       int t = l;
        while (j < L.length && k < R.length) {
            if (L[j] > R[k]) {
                arr[t] = R[k];
                k++;
            } else {
                arr[t] = L[j];
                j++;
            }
            t++;

        }

        while (j < L.length) {
            arr[t] = L[j];
            j++;
            t++;
        }

        while (k < R.length) {
            arr[t] = R[k];
            t++;
            k++;
        }

    }



    public static void merge(int[] arr, int l, int r, int m) {

        int[] L = new int[m - l + 1];
        int[] R = new int[r-m];

        int j = 0;
        for (int i = 0; i < L.length ; ++i) {
            L[j++] = arr[l+i];
        }
        int k = 0;
        for (int i = 0; i < R.length ; i++) {
            R[k++] = arr[i + m + 1];
        }
        j = 0; k = 0;
        int t = l;
        while (j < L.length && k < R.length) {
            if (L[j] < R[k])
            { arr[t++] = L[j++];} else {
                arr[t++] = R[k++];
            }
        }

        while (j < L.length) {
            arr[t++] = L[j++];
        }
        while (k < R.length) {
            arr[t++] = R[k++];
        }
    }
}
