package additional;

import java.util.Comparator;
import java.util.PriorityQueue;

//Given an array of point co-ordinates, find the points that are closest to 0.0
//Build minheap with a camparator that compares distance of any point from 0,0
class PointsComp implements Comparator<Integer[]> {

    @Override
    public int compare(Integer[] o1, Integer[] o2) {

        double distO1 = Math.sqrt(o1[0]*o1[0] + o1[1]*o1[1]);
        double distO2 = Math.sqrt(o2[0]*o2[0] + o2[1]*o2[1]);
        return (distO1 > distO2) ? 1 : (distO2 == distO1) ? 0 : -1 ;
    }
}

public class ClosestPoints {

    static PriorityQueue<Integer[]> minHeap = new PriorityQueue<>(new PointsComp());

    public static void main(String[] args) {
        Integer[][] arr = {{3,3},{5,-1},{-2,4}};
        getClosestPoints(arr, 2);
    }

    public static Integer[][] getClosestPoints(Integer[][] points, int k) {

        Integer[][] ret = new Integer[points.length][k];
        for (int i = 0; i < points.length; i++) {
            Integer[] point = points[i];

            minHeap.add(point);

        }

        for (int i = 0; i < k ; i++) {
            ret[i] = minHeap.remove();
        }
        return ret;


    }

}
