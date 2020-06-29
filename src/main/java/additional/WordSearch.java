package additional;

import java.util.ArrayList;
import java.util.List;

public class WordSearch {

    private List<String> lst = new ArrayList<>();

    public static void main(String[] args) {
        char[][] matrix = {{'A','B','C','E'},
                {'S','F','C','S'},
                {'A','D','E','E'}};

        String[] words = {"oath","pea","eat","rain"};
        boolean ls = (new WordSearch()).exist(matrix, "SEE");

    }

    public  boolean exist(char[][] board, String word) {
        boolean[][] visited = new boolean[board.length][board[0].length];
        boolean ret = false;
        for(int i = 0; i < board.length; i++){
            for (int j = 0; j < board[i].length; j++) {
                ret = ret || dfs(board, visited, word, "", i, j);
            }

        }
        return ret;
    }


    private  boolean dfs(char[][] board, boolean[][] visited, String word, String currWord, int x, int y) {

        if(x < 0 || y < 0 || x >= board.length) {
            return false;
        }
        if(y >= board[x].length) {
            return false;
        }
        if (visited[x][y]) return false;

        currWord += board[x][y];

        if(!word.startsWith(currWord)) {
            return false;
        }
        if(word.compareToIgnoreCase(currWord) == 0) {
            return true;
        }
        visited[x][y] = true;

        boolean top = false;
        boolean bottom = false;
        boolean right = false;
        boolean left = false;


        bottom = dfs(board, visited, word, currWord, x+1, y);

        top = dfs(board, visited, word, currWord, x-1, y);


        right = dfs(board, visited, word, currWord, x, y+1);

        left = dfs(board, visited, word, currWord, x, y-1);


        visited[x][y] = false;
        return left || right || top || bottom;

    }

//    public List<String> findWords(char[][] board, String[] words) {
//
//        Trie trie = new Trie();
//        TrieNode root = new TrieNode();
//        boolean[][] visited = new boolean[board.length][board[0].length];
//
//        for (int i = 0; i < words.length; i++) {
//            trie.insert(words[i], root);
//        }
//
//        for (int i = 0; i < board.length; i++) {
//            for (int j = 0; j < board[0].length; j++) {
//                dfs(visited, board, i, j, "", trie, root);
//            }
//        }
//        return lst;
//    }
//
//    private void dfs(boolean[][] visited, char[][] board, int x, int y, String str, Trie trie, TrieNode root) {
//        if (x < 0 || x >= board.length || y < 0 || y >= board[0].length) return;
//        if(visited[x][y]) return;
//
//        str += board[x][y];
//        if (!trie.startsWith(str,root)) {
//            return;
//        }
//
//        if (trie.find(str, root)) {
//            lst.add(str);
//        }
//        visited[x][y] = true;
//
//        dfs(visited, board, x+1, y, str, trie, root);
//        dfs(visited, board, x-1, y, str, trie, root);
//        dfs(visited, board, x, y+1, str, trie, root);
//        dfs(visited, board, x, y-1, str, trie, root);
//
//        visited[x][y] = false;
//
//
//    }

}

class TrieNode {

    TrieNode[] children;
    boolean isEndOfWord;

    TrieNode() {
        children = new TrieNode[26];
        isEndOfWord = false;
    }

}
class Trie {

    public void insert(String str, TrieNode root) {

        if (root == null) {
            return;
        }
        str = str.toLowerCase().trim();
        TrieNode nd = root;
        for (int i = 0; i < str.length(); i++) {

            int loc = str.charAt(i) - 'a';
            if (nd.children[loc] == null) {
                nd.children[loc] = new TrieNode();
            }
            nd = nd.children[loc];
            if (i == str.length() - 1) {
                nd.isEndOfWord = true;
            }
        }


    }

    public boolean find(String str, TrieNode root) {
        if (root == null) {
            return false;
        }
        str = str.trim().toLowerCase();

        TrieNode nd = root;
        for (int i = 0; i < str.length(); i++) {
            int loc = str.charAt(i) - 'a';
            if (nd.children[loc] == null) {
                return false;
            }
            nd = nd.children[loc];

        }
        return nd != null && nd.isEndOfWord;
    }

    public boolean startsWith(String prefix, TrieNode root) {

        TrieNode nd = root;

        prefix = prefix.trim().toLowerCase();


        for (int i = 0; i < prefix.length(); i++) {
            int j = prefix.charAt(i) - 'a';
            TrieNode[] children = nd.children;

            if (children[j] == null) {
                return false;
            }
            nd = children[j];

        }
        return true;
    }


}
