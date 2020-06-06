package arrays;

public class SearchInMatrix {

    public static void main(String[] args) {
        int[][] matrix = {{1},{3}};
        (new SearchInMatrix()).searchMatrix(matrix, 3);
    }


    public boolean searchMatrix(int[][] matrix, int target) {
        if (matrix == null || matrix.length == 0 || matrix[0].length == 0) {
            return false;
        }
        if (matrix.length == 1 && matrix[0].length == 1) {
            return (matrix[0][0] == target);
        }
        int i = 0;

        while (i < matrix.length && target > matrix[i][0]) {
            i++;
        }

        if (i < matrix.length) {
            if(matrix[i][0] == target) {
                return true;
            }
        }
        i--;
        if (i >= 0) {
            return searchRow(matrix, target, i);
        }
        return false;

    }

    private boolean searchRow(int[][] matrix, int target, int row) {
        for (int i = 0; i < matrix[row].length; i++) {
            if (matrix[row][i] == target) {
                return true;
            }
        }
        return false;
    }

}
