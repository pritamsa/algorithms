package strings;

import java.util.ArrayList;
import java.util.List;

public class RotateMatrix {

    public static void main(String[] args) {

        List<List<Integer>> sourceMatrix = new ArrayList<>(2);
        sourceMatrix.add( new ArrayList<>());
        sourceMatrix.add( new ArrayList<>());
        sourceMatrix.get(0).add( 50);
        sourceMatrix.get(1).add( 90);
        RotateMatrix(sourceMatrix, 90);

    }

    static List<List<Integer>> RotateMatrix(List<List<Integer>> sourceMatrix, int rotationDegree) {
        List<List<Integer>> rotated = new ArrayList<>();

        if (sourceMatrix == null || sourceMatrix.size() == 0
                || sourceMatrix.get(0).size() == 0 || rotationDegree % 90 != 0) {
            return rotated;

        }


        int numRotations = rotationDegree/90;
        return rotate(sourceMatrix, numRotations);

    }

    private static List<List<Integer>> rotate(List<List<Integer>> matrix,
                                              int numRotations) {
        List<List<Integer>> rotated = matrix;
        int cols = 0;
        for (int j = 0; j < matrix.size(); j++) {
            if (matrix.get(j) != null && matrix.get(j).size() > cols) {
                cols = matrix.get(j).size();
            }
        }
        for (int i = 0; i < numRotations; i++) {
            rotated = transpose(matrix, cols);
            rotated = reverse(rotated);
        }
        return rotated;

    }


    private static List<List<Integer>> transpose(List<List<Integer>> matrix, int rows) {
        List<List<Integer>> transposed = new ArrayList<>();
//        for (int k = 0; k < rows; k++) {
//            transposed.add(new ArrayList<>());
//
//        }
        for (int i = 0; i < matrix.size(); i++) {
            for (int j = 0; j < matrix.get(i).size(); j++) {
                int temp = matrix.get(i).get(j);

                if (transposed.get(j) == null) {
                    transposed.add(j,new ArrayList<>());
                }
                transposed.get(j).add(temp);
            }

        }
        return transposed;

    }

    private static List<List<Integer>> reverse(List<List<Integer>> matrix) {

        for (int i = 0; i < matrix.size(); i++) {
            int st = 0;
            int en = matrix.get(i).size() - 1;
            while(en >= st) {
                int temp = matrix.get(i).get(st);
                matrix.get(i).set(st, matrix.get(i).get(en));
                matrix.get(i).set(en, temp);

                st++;
                en--;
            }

        }
        return matrix;

    }

}
