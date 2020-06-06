package arrays;

//You are in an infinite 2D grid where you can move in any of the 8 directions :
//
//        (x,y) to
//        (x+1, y),
//        (x - 1, y),
//        (x, y+1),
//        (x, y-1),
//        (x-1, y-1),
//        (x+1,y+1),
//        (x-1,y+1),
//        (x+1,y-1)
//        You are given a sequence of points and the order in which you need to cover the points. Give the minimum number of steps in which you can achieve it. You start from the first point.
//
//        Input :
//11
public class MinStepsInfiniteGrid {

//    class A {}
//    class B extends A {}
//    class C extends B {}

    public static void main(String[] args) {
        class A {}
        class B extends A {}
        class C extends B {}

        B b = new B();
        boolean t = (b instanceof A);
    }

    private int getDist(int[] pt1, int[] pt2) {

        getNormalDist(pt1[0], pt1[1], pt2[0], pt2[1]);
        if (isInDiagonal(pt1[0], pt1[1], pt2[0], pt2[1])) {
            getDiagonalDist( pt1[1],  pt2[1]);
        } else {
            //get diagonal pt in col and find dist through it
            //get diagonal pt in row and find dist through it.
        }
        //return min
        return 1;
    }


    private int[] getDiagonalPtInColumn(int i1, int j1, int col, int maxRows, boolean fwd) {
        int offset = Math.abs(col - j1);
        int[] ret = new int[2];
        if (fwd) {
            ret[0] = i1 + offset;
        } else {
            ret[0] = i1 - offset;
        }
        ret[1] = col;
        if (ret[0] < 0 || ret[0] > maxRows) {
            return null;
        }
        return ret;
    }

    private int[] getDiagonalPtInRowFwd(int i1, int j1, int row, int maxCols, boolean fwd) {
        int offset = Math.abs(row - i1);
        int[] ret = new int[2];
        ret[0] = row;
        if (fwd) {
            ret[1] = j1 + offset;
        } else {
            ret[1] = j1 - offset;
        }
        if (ret[1] < 0 || ret[1] > maxCols) {
            return null;
        }
        return ret;

    }

    private int getNormalDist(int i1, int j1, int i2, int j2) {
        return Math.abs(j2-j1) + Math.abs(i2-i1);

    }

    private boolean isInDiagonal(int i1, int j1, int i2, int j2) {
        return Math.abs(i2-i1) == Math.abs(j2-j1);
    }
    private int getDiagonalDist(int j1, int j2) {

        return Math.abs(j2-j1);
    }

}
