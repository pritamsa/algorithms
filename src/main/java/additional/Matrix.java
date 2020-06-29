package additional;

//[[-4,-2147483648,6,-7,0],[-8,6,-8,-6,0],[2147483647,2,-9,-6,-10]]
//        Output
//        [[0,0,0,0,0],[0,0,0,0,0],[0,2,-9,-6,0]]
//        Expected
//        [[0,0,0,0,0],[0,0,0,0,0],[2147483647,2,-9,-6,0]]
// Set zeros to the row and col of a matrix when an element has a value 0
public class Matrix {

    public static void main(String[] args) {
        int[][] arr = {{-4,-2147483648,6,-7,0},{-8,6,-8,-6,0},{2147483647,2,-9,-6,-10}};
        (new Matrix()).setZeroes(arr);
    }

    public void setZeroes(int[][] matrix) {

        boolean isFirstCol = false;
        boolean isFirstRow = false;


        for(int i = 0; i < matrix.length; i++) {
            if(matrix[i][0] == 0) {
                isFirstCol = true;
            }
        }

        for(int i = 0; i < matrix[0].length; i++) {
            if(matrix[0][i] == 0) {
                isFirstRow = true;
            }
        }

        for(int i = 1; i < matrix.length; i++) {

            for(int j = 1; j < matrix[i].length; j++) {

                if (matrix[i][j] == 0) {
                    matrix[i][0] = 0;
                    matrix[0][j] = 0;
                }

            }
        }

        int t = 0;
        //col
        for(int j = 1; j < matrix[t].length; j++) {


                if(matrix[t][j] == 0) {
                    for(int k = 0; k < matrix.length; k++) {
                        matrix[k][j] = 0;
                    }

                }



        }
        //row
        int d = 0;
        for(int i = 1; i < matrix.length; i++) {

                if(matrix[i][d] == 0) {
                    for(int k = 0; k < matrix[i].length; k++) {
                        matrix[i][k] = 0;
                    }
                }


        }
        //firstrow
        if (isFirstRow) {
            for(int j = 0; j < matrix[0].length; j++) matrix[0][j] = 0;

        }

        //first col
        if (isFirstCol) {
            for(int i = 0; i < matrix.length; i++) matrix[i][0] = 0;
        }


    }
}
