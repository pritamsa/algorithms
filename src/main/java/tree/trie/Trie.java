package tree.trie;

import javax.swing.tree.TreeNode;
import java.util.HashMap;
import java.util.Map;

class TrieNode {

    TrieNode[] children;
    boolean isEndOfWord;

    TrieNode() {
        children = new TrieNode[26];
        isEndOfWord = false;
    }

}
public class Trie {

    public static void main(String[] args) {
        TrieNode root = new TrieNode();
        Trie t = new Trie();

        t.insert("auth", root);
        t.insert("capture", root);
        t.insert("someword", root);

        boolean word1 = t.find("auth", root);
        boolean word2 = t.find("capture", root);
        boolean word3 = t.find("someword", root);
    }

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
