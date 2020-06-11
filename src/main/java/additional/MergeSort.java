package additional;

import linkedlist.MergeTwoLists;

public class MergeSort {

    public static void main(String[] args) {
        int[] arr = {17, 22, 7, 23, 5, 2, 8, 9};
        (new MergeSort()).sort(arr);
    }

    public int[] sort(int[] arr) {
        sort(arr, 0, arr.length - 1);
        return arr;
    }

    private void sort(int[] arr, int st, int en) {
        if (st == en) {
            return;
        }
        int mid = (st + en)/2;
        sort(arr, st, mid);
        sort(arr, mid+1, en);

        merge(arr, st, en, mid);


    }

    private void merge(int[] arr, int l, int r, int m) {
        if (l == r) {
            return;
        }
        int[] L = new int[m-l+1];
        int[] R = new int[r-m];

        //Copy left
        int k=0;
        for (int i = l; i <= m ; i++) {
            L[k++] = arr[i];
        }

        //Copy right
        k =0;
        for (int i = m+1; i <= r ; i++) {
            R[k++] = arr[i];
        }

        k = l;
        int i = 0;
        int j = 0;

        while(i < L.length && j < R.length) {
            if (L[i] <= R[j]) {
                arr[k] = L[i];
                k++;
                i++;
            } else {
                arr[k] = R[j];
                k++;
                j++;
            }
        }
        //Copy rest
        while(i < L.length) {
            arr[k] = L[i];
            k++;
            i++;
        }
        //Copy rest
        while(j < L.length) {
            arr[k] = R[j];
            k++;
            j++;
        }

    }
}
