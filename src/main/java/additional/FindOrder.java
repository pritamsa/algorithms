package additional;

import java.util.Stack;

public class FindOrder {

  public static void main(String[] args) {
    int[][] pre = {{0,1},{1,0}};
    int[] arr = (new FindOrder()).findOrder(2, pre);
  }

  public int[] findOrder(int numCourses, int[][] prerequisites) {
    int[][] board = new int[numCourses][numCourses];
    boolean[] visited = new boolean[numCourses];
    Stack<Integer> cou = new Stack<>();


    for (int i = 0; i < prerequisites.length; i++) {
      addEdge(prerequisites[i][0], prerequisites[i][1], board );
    }
    for (int i = 0; i < numCourses; i++) {
      if (!visited[i]) {

        boolean val = dfs(i, cou, board, visited);
        if (! val) return null;
      }

    }


    int[] arr = new int[cou.size()];
    int k = 0;
    while (!cou.isEmpty()) {
      arr[arr.length - 1 - k] = cou.pop();
      k++;
    }

    return arr;
  }

  private void addEdge(int row, int col, int[][] board) {
    board[row][col] = 1;

  }

  public boolean dfs(int v, Stack<Integer> cou, int[][] board, boolean[] visited ) {
    if (visited[v]) {
      return false;
    }


    visited[v] = true;
    for (int i = 0; i < board.length; i++) {
      if (board[v][i] == 1) {
        //if (!visited[i]) {
          boolean val = dfs(i, cou, board, visited);
          if (!val) {
            return false;
          }
        //}

      }
    }
  visited[v] = false;
    if (!cou.contains(v)) cou.push(v);
    return true;
  }


}
