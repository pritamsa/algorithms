package additional;



import java.util.LinkedList;

class Graph {
    LinkedList<Integer>[] adj = null;

    Graph(int capacity) {
        adj = new LinkedList[capacity];
    }

    public void add(int i, int j) {
        adj[i].add(j);
    }

    public void dfs() {

    }
    public void bfs() {

    }
}

public class Misc {



}
