package additional;

public class Medians {

    public static void main(String[] args) {
        int[] A = {1,12,15,26,38};
        int[] B = {2, 13, 17, 30, 45};

        int median = findMedianTwoSortedArraysSameSize(A, B);

    }

    public static int findMedianTwoSortedArraysSameSize(int[] A, int[] B) {

        if (A == null || B == null || A.length == 0 || B.length == 0) {
            return -1;
        }
        int count = 0;

        int i = 0;
        int j = 0;
        int m = A.length;
        int n = B.length;
        int prev = 0;
        int curr = 0;

        int targetCount = 1+((m+n)/2);
        while (i < m && j < n) {
            if (A[i] <= B[j]) {
                prev = curr;
                curr = A[i];

                count++;
                i++;

            } else if (B[j] < A[i]) {
                prev = curr;
                curr = B[j];

                count++;
                j++;

            }

            if (count == targetCount) {
                return (prev + curr)/2;
            }

        }

        while (i < m) {
            prev = curr;
            curr = A[i];
            count++;
            i++;

            if (count == targetCount) {
                return (prev + curr)/2;
            }

        }

        while (j < n) {
            prev = curr;
            curr = B[j];
            count++;
            j++;

            if (count == targetCount) {
                return (prev + curr)/2;
            }


        }

        return -1;
    }

    public static double findMedianTwoSortedArraysDiffSize(int[] A, int[] B) {

        int m = A.length;
        int n = B.length;

        //swap if needed
        if (m > n) {
            int[] temp = A;
            A = B;
            B = temp;

            int tmp = m;
            m = n;
            n = tmp;

        }

        int halfLen = (m + n + 1)/2;

        int iMin = 0;
        int iMax = m;

        while (iMin <= iMax) {
            int i = (iMin + iMax)/2;

            int j = halfLen - i;

            if (i > iMin && A[i-1] > B[j]) {
                //i is too big, go lower half
                iMax = i-1;
            } else if (i < iMax && A[i] < B[j - 1]) {
                //i is too small, go in upper half
                iMin = i + 1;
            } else {
                int lowerMax = 0;
                int higherMin = 0;

                if (i == 0) {
                    lowerMax = B[j-1];
                } else if (j == 0 ) {
                    lowerMax = A[i-1];
                }else {
                    lowerMax = Math.max(A[i-1], B[j-1]);
                }
                if ((m+n) % 2 == 1) {
                    return lowerMax;
                }

                if (i == m) {
                    higherMin = B[j];
                } else if (j == n) {
                    higherMin = A[i];
                } else {
                    higherMin = Math.min(A[i], B[j]);
                }
                return (lowerMax + higherMin)/2;
            }

        }

        return -1;

    }

    public static int findMedianRowWiseSortedMatrix() {

        return -1;

    }

    public static int findMedianUnsortedArray(int[] A, int[] B) {

        return -1;

    }
}
