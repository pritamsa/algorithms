package arrays;

public class LIS {

    public static void main(String[] args) {
        int[] arr = {10,9,2,5,3,7,101,18};
        lengthOfLIS(arr);

    }

    public static int lengthOfLIS(int[] nums) {

        int[] dp = new int[nums.length + 1];

        dp[0] = 1;
        int max = 0;

        for (int i = 1; i < nums.length; i++) {
            dp[i] = 1;
        }

        for (int i = 1; i < nums.length; i++) {

            for (int j = 0; j < i; j++) {
                if (nums[j] < nums[i] && dp[i] < dp[j] + 1) {
                    dp[i] = dp[j] + 1;

                }

                if (max < dp[i]) {
                    max = dp[i];
                }

            }
        }
        return max;

    }
}

