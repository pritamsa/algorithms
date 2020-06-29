package arrays;
//Spiral Order Matrix II
//Given an integer A, generate a square matrix filled with elements from 1 to A2 in spiral order.
//
//        Input Format:
//
//        The first and the only argument contains an integer, A.
//        Output Format:
//
//        Return a 2-d matrix of size A x A satisfying the spiral order.
//        Constraints:
//
//        1 <= A <= 1000
//        Examples:
//
//        Input 1:
//        A = 3
//
//        Output 1:
//        [ 1, 2, 3 ],
//        [ 8, 9, 4 ],
//        [ 7, 6, 5 ]   ]
//
//        Input 2:
//        4
//
//        Output 2:
//        [1, 2, 3, 4],
//        [12, 13, 14, 5],
//        [11, 16, 15, 6],
//        [10, 9, 8, 7]   ]

import java.util.ArrayList;
import java.util.List;

public class SpiralOrderMatrix {

    public static void main(String[] arrs) {
        int[][] board = {{1, 2, 3},
                {4, 5, 6},
                {7, 8, 9}
        };
        List<Integer> ret = spiralOrder(board);
    }

    public static List<Integer> spiralOrder(int[][] matrix) {

        List<Integer>  ret = new ArrayList<>();
        if (matrix == null || matrix.length == 0) {
            return ret;
        }

        int rowMax = matrix.length - 1;
        int colMax = matrix[0].length - 1;

        int rowMin = 0;
        int colMin = 0;


        while (rowMin <= rowMax && colMin <= colMax) {

            for (int i = colMin; i <= colMax; i++) {
                ret.add(matrix[rowMin][i]);

            }

            for (int i = rowMin+1; i <= rowMax; i++) {
                ret.add(matrix[i][colMax]);

            }

            if(rowMin < rowMax && colMin < colMax) {
                for (int i = colMax - 1; i >= colMin; i--) {
                    ret.add(matrix[rowMax][i]);
                }

                for (int i = rowMax-1; i > rowMin; i--) {
                    ret.add(matrix[i][colMin]);

                }

            }

            rowMin++;
            rowMax--;
            colMin++;
            colMax--;


        }
        return ret;

    }

}
