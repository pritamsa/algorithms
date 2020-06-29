package additional;

import java.util.HashSet;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.LinkedBlockingQueue;

//There are N network nodes, labelled 1 to N.
//
//        Given times, a list of travel times as directed edges times[i] = (u, v, w), where u is
// the source node, v is the target node, and w is the time it takes for a signal to travel from source to target.
//
//        Now, we send a signal from a certain node K. How long will it take for all nodes to
// receive the signal? If it is impossible, return -1.
//        Input: times = [[2,1,1],[2,3,1],[3,4,1]], N = 4, K = 2
//        Output: 2

//Dijkstra’s shortest path algorithm from one source to multiple destinations
//Idea is to maintain min spanning tree.
//Keep a min = MAX_VALUE and minIdx = -1
// - From the times[] construct edge matrix first. Edge value = Max_value if the 2 verices are not connected
//- Keep a set and keep a distance array.
// - current source = source
//- Add the current source to the set first
//- Initialize dist[source] = 0. Rest of dist[] = Max value
// for that source fine all adj vertices,
//  for each adj vertex (0 to n) : if the vertex is connected and is not in the set yet ,
// dist[adj vertes] = min(dist[adj_vertex], dist[current source] + matrix[current source][adj vertex]
// if dist[adj_vertex] < min , min = dist[adj_vertex], minIdx = adj_vertex,
//For each adjecency loop, initialize minIdx = -1 and min as MAX_VALUE. So you can find a min distance vertex at
// the end of each loop. add min distance adj_vertex to the set.
// After all vertices are exausted, return dist[] array.
public class NwNodes {

    public static void main(String[] args) {
        int[][] times = {{1,0,1},{1,2,1},{2,3,1}};
        int[] nwTimes = calcTimes(times, 4, 0);

    }

    public static int[] calcTimes(int[][] times, int n, int k) {

        int[][] matrix = new int[n][n];

        int[] dist = new int[n];

        Set<Integer> set = new HashSet<>();

        for (int i = 0; i < matrix.length; i++) {
            for (int j = 0; j < matrix[i].length; j++) {
                matrix[i][j] = Integer.MAX_VALUE;
            }
        }
        //Create edge matrix
        for (int i = 0; i < times.length ; i++) {
            int[] arr = times[i];

            matrix[arr[0]][arr[1]] = arr[2];
            matrix[arr[1]][arr[0]] = arr[2];
        }

        //Min distances from source vertex
        for (int i = 0; i < n; i++) {
            dist[i] = Integer.MAX_VALUE;
        }


        dist[k] = 0;
        int min = Integer.MAX_VALUE;
        int minIdx = k;
        int v = k;
        set.add(k);

            //Dist keeps distances from the source, so for each vertex 'j' adj to 'v', we find the min of dist[j] and
            // matrix[v][j] + dist[v]
        while (v != -1) {
            for (int j = 0; j < matrix[v].length; j++) {
                if ( v != j) {
                    dist[j] = matrix[v][j] != Integer.MAX_VALUE ?
                            Math.min(dist[j], matrix[v][j] + dist[v]) : dist[j];
                    if (min > dist[j] && !set.contains(j)) {
                        min = dist[j];
                        minIdx = j;
                    }
                }

            }
            v = minIdx;
            set.add(minIdx);
            min = Integer.MAX_VALUE;
            minIdx = -1;

        }

        return dist;
    }
}
