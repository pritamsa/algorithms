package arrays;

public class TrappedWater {

    public static void main(String[] args) {
//        int[] height = {0,1,0,2,1,0,1,3,2,1,2,1};
//        (new TrappedWater()).trap(height);

        int[][] matrix = {{0,1,2,0}, {3,4,5,2}, {1,3,1,5}};
        (new TrappedWater()).setZeroes(matrix);

    }

    public int trap(int[] height) {

        int max_on_left = height[0];
        int max_on_Right = findMax(height, 2);
        int max = -1;
        int sum = 0;

        for (int i = 1; i < height.length - 1; i++) {
            if (i > 1) {
                max_on_left = Math.max(height[i-1], max_on_left);
                if (max_on_Right == height[i]) {
                    max_on_Right = findMax(height, i+1);
                }
            }
            int lower = Math.min(max_on_left, max_on_Right);
            int water = 0;
            if (lower > height[i]) {
                water = lower - height[i];
            }
            sum += water;
        }
        return sum;

    }

    private int findMax(int[] height, int st) {
        int max = 0;

        for (int i = st; i < height.length; i++) {
            if (height[i] > max) {
                max = height[i];
            }
        }
        return max;
    }

    public void setZeroes(int[][] matrix) {

        for(int i = 0; i < matrix.length; i++) {
            for(int j = 0; j < matrix[i].length; j++) {
                if(matrix[i][j] == 0) {
                    makeNegative(matrix, i, j);
                }

            }

        }

        for(int i = 0; i < matrix.length; i++) {
            for(int j = 0; j < matrix[i].length; j++) {
                if(matrix[i][j] < 0) {
                    matrix[i][j] = 0;
                }
            }
        }
    }

    private void makeNegative(int[][] matrix, int row, int col) {

        for(int j = 0; j < matrix[row].length; j++) {
            if (matrix[row][j] > 0) {
                matrix[row][j] = -matrix[row][j];
            }
        }

        for(int j = 0; j < matrix.length; j++) {
            if (matrix[j][col] > 0) {
                matrix[j][col] = -matrix[j][col];
            }
        }

    }

}
