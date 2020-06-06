package graph;

import java.util.Stack;

public class NumberOfIslands {

    public static void main(String[] args) {
        char[][] matrix = new char[5][5];
        matrix[0][0] = '1';
        matrix[0][1] = '1';
        matrix[0][2] = '1';
        matrix[0][3] = '1';
        matrix[0][4] = '0';

        matrix[1][0] = '1';
        matrix[1][1] = '1';
        matrix[1][2] = '0';
        matrix[1][3] = '1';
        matrix[1][4] = '0';

        matrix[2][0] = '1';
        matrix[2][1] = '1';
        matrix[2][2] = '0';
        matrix[2][3] = '0';
        matrix[2][4] = '0';

        matrix[3][0] = '0';
        matrix[3][1] = '0';
        matrix[3][2] = '0';
        matrix[3][3] = '0';
        matrix[3][4] = '0';


        (new NumberOfIslands()).findNumOfIslands(matrix);
    }

    private int findNumOfIslands(char[][] matrix) {

        boolean[][] visited = new boolean[matrix.length][matrix[0].length];

        int count = 0;

        for (int i = 0; i < matrix.length; i++) {
            for (int j = 0; j < matrix[i].length; j++) {
                if (!visited[i][j] && matrix[i][j] == '1') {
                    islandUtil(matrix, visited, i, j);

                    count++;
                }

            }
        }
        return count;
    }

    private void islandUtil(char[][] matrix, boolean[][] visited, int row, int col) {
        visited[row][col] = true;

        int[] offset = {1, -1};

        for (int i = 0; i < offset.length; i++) {
            int neRow = row + offset[i];

            if (isSafe(matrix,neRow,col)) {
                if (matrix[neRow][col] == '0') {
                    visited[neRow][col] = true;
                } else {
                    if(!visited[neRow][col]) {
                        islandUtil(matrix, visited, neRow, col);
                    }
                }

            }
        }

        for (int i = 0; i < offset.length; i++) {
            int neCol = col + offset[i];

            if (isSafe(matrix,row,neCol)) {
                if (matrix[row][neCol] == '0') {
                    visited[row][neCol] = true;
                } else {
                    if(!visited[row][neCol]) {
                        islandUtil(matrix, visited, row, neCol);
                    }
                }

            }
        }
    }

    private boolean isSafe(final char[][] matrix, int row, int col) {
        return row < matrix.length && row >=0 && col < matrix[row].length && col >=0;
    }
}
