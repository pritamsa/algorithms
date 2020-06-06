package tree.trie;

import java.util.*;

public class WordSearch {
//    static List<String> topLst = new ArrayList<>();

    static Set<String> res = new HashSet<String>();

    public List<String> find(char[][] board, String[] words) {
        TrieNode root = new TrieNode();
        Trie trie = new Trie();
        for (String word : words) {
            trie.insert(word, root);
        }

        int m = board.length;
        int n = board[0].length;
        boolean[][] visited = new boolean[m][n];
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                dfs(board, visited, "", i, j, trie, root);
            }
        }

        return new ArrayList<String>(res);
    }

    public void dfs(char[][] board, boolean[][] visited, String str, int x, int y, Trie trie, TrieNode root) {
        if (x < 0 || x >= board.length || y < 0 || y >= board[0].length) return;
        if (visited[x][y]) return;

        str += board[x][y];
        if (!trie.startsWith(str, root)) return;

        if (trie.find(str, root)) {
            res.add(str);
        }

        visited[x][y] = true;
        dfs(board, visited, str, x - 1, y, trie, root);
        dfs(board, visited, str, x + 1, y, trie, root);
        dfs(board, visited, str, x, y - 1, trie, root);
        dfs(board, visited, str, x, y + 1, trie, root);
        visited[x][y] = false;
    }
    public static void main(String[] args) {

        char[][] matrix = {{'o','a','a','n'}, {'e','t','a','e'}, {'i','h','k','r'}, {'i','f','l','v'}};

        String[] words = {"oath","pea","eat","rain"};
        List<String> ls = findWords(matrix, words);//(new WordSearch()).find(matrix, words);

    }

    public static List<String> findWords(char[][] board, String[] words) {
        Set<String> ls = new HashSet<String>();
        boolean[][] visited = new boolean[board.length][board[0].length];
        if (board == null || words == null || board.length == 0 || board[0].length == 0 || words.length == 0) {

            String[] arr = (String[]) ls.toArray();
            return Arrays.asList(arr);
        }

        Trie trie = new Trie();
        TrieNode root = new TrieNode();
        for(int i = 0; i < words.length; i++) {
            trie.insert(words[i], root);

        }

        for (int i = 0; i < board.length; i++) {
            for (int j = 0; j < board[i].length; j++) {
                dfs(visited, board, i, j, trie, "", root, ls );

            }
        }
        //String[] arr = (String[]) ls.toArray();
        return null;


    }

    public static void dfs(boolean[][] visited, final char[][] board, int x, int y, Trie trie, String word, TrieNode root,
                           Set<String> ls) {


        if (x < 0 || x >= board.length || y < 0 || y >= board[0].length) return;
        if (visited[x][y]) return;

        word += board[x][y];


        if (!trie.startsWith(word, root)) {
            return;
        }

        if (trie.find(word, root)) {

            res.add(word);

        }
        visited[x][y] = true;


            dfs(visited,board, x - 1, y, trie, word, root, ls);

            dfs(visited,board, x + 1, y, trie, word, root, ls);

            dfs(visited,board, x, y - 1, trie, word, root, ls);

            dfs(visited,board, x, y + 1, trie, word, root, ls);

        visited[x][y] = false;

    }

    private static boolean isSafe(final int i, final int j, char[][] board) {
        return i < board.length && j < board[0].length && i >= 0 && j >= 0;
    }


}
