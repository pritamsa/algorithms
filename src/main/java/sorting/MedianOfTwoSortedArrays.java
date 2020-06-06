package sorting;

public class MedianOfTwoSortedArrays {

    public static void main(String[] args) {
        int[] nums1 = {1,2};
        int[] nums2 = {3,4};
        (new MedianOfTwoSortedArrays()).findMedianSortedArrays(nums1, nums2);

    }

    public double findMedianSortedArrays(int[] nums1, int[] nums2) {

        int i = 0;
        int j = 0;
        int k = 0;

        double prevToIdx = 0;
        double valIdx = 0;

        while(i < nums1.length && j < nums2.length) {
            if (nums1[i] <= nums2[j]) {

                if (k == ((nums1.length + nums2.length)/2) - 1 ) {
                    prevToIdx = nums1[i];

                } else if (k == (nums1.length + nums2.length)/2) {
                    valIdx = nums1[i];
                }
                i++;
            } else {

                if (k == ((nums1.length + nums2.length)/2) - 1 ) {
                    prevToIdx = nums2[j];

                } else if (k == (nums1.length + nums2.length)/2) {
                    valIdx = nums2[j];
                }
                j++;
            }

            k++;
        }
        while(i < nums1.length) {

            if (k == ((nums1.length + nums2.length)/2) - 1 ) {
                prevToIdx = nums1[i];

            } else if (k == (nums1.length + nums2.length)/2) {
                valIdx = nums1[i];
            }
            k++;
            i++;
        }
        while(j < nums2.length) {

            if (k == ((nums1.length + nums2.length)/2) - 1 ) {
                prevToIdx = nums2[j];

            } else if (k == (nums1.length + nums2.length)/2) {
                valIdx = nums2[j];
            }
            k++;
            j++;
        }

            if ((nums1.length + nums2.length) %2 == 0) {
                return (valIdx + prevToIdx)/2;
            } else {
                return valIdx;
            }



    }
}
