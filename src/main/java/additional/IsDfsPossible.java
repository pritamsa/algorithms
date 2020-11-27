package additional;

import java.util.HashSet;
import java.util.Set;

public class IsDfsPossible {

  public static void main(String[] args) {
    //int[][] pre = {{1,0},{2,0},{3,1},{3,2}};
    int[][] pre = {{0,1},{1,0}};
    (new FindOrder()).findOrder(2, pre);
  }

  public boolean canFinish(int numCourses, int[][] prerequisites) {
    int[][] board = new int[numCourses][numCourses];
    boolean[] visited = new boolean[numCourses];
    Set<Integer> cou = new HashSet<>();

    for (int i = 0; i < prerequisites.length; i++) {
      addEdge(prerequisites[i][0], prerequisites[i][1], board );
    }
    for (int i = 0; i < numCourses; i++) {
      boolean val = dfs(i, cou, board, visited);
      if (!val) return false;
    }
    return true;

  }

  private void addEdge(int row, int col, int[][] board) {
    board[row][col] = 1;

  }

  public boolean dfs(int v, Set<Integer> cou, int[][] board, boolean[] visited ) {
    if (visited[v]) {
      return false;
    }

    cou.add(v);
    visited[v] = true;
    for (int i = 0; i < board.length; i++) {
      if (board[v][i] == 1) {
        boolean val = dfs(i, cou, board, visited);
        if (!val) {
          return false;
        }
      }
    }
    visited[v] = false;
    return true;
  }

}
