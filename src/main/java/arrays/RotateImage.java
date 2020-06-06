package arrays;

public class RotateImage {

    public static void main(String[] args) {
        int[][] matrix = {{1,2,3},{4,5,6},{7,8,9}};
        (new RotateImage()).rotate(matrix);
    }

    public void rotate(int[][] matrix) {
        if (matrix != null) {
            matrix = transpose(matrix);
            matrix = reverse(matrix);
        }
    }

    private int[][] transpose(int[][] matrix) {

        for (int i = 0; i < matrix.length; i++) {
            for (int j = i; j < matrix[0].length; j++) {
                int temp = matrix[i][j];
                matrix[i][j] = matrix[j][i];
                matrix[j][i] = temp;

            }

        }
        return matrix;
    }

    private int[][] reverse(int[][] matrix) {

        for (int i = 0; i < matrix.length; i++) {
            int st = 0;
            int en = matrix[i].length - 1;
            while(en >= st) {
                int temp = matrix[i][st];
                matrix[i][st] = matrix[i][en];
                matrix[i][en] = temp;
                st++;
                en--;
            }

        }
        return matrix;
    }

}
