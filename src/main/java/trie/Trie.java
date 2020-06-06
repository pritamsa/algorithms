package trie;

class Trie {

    class TrieNode {
        TrieNode[] children;
        boolean isEndOfWord;

        TrieNode() {
            children = new TrieNode[26];

        }

        TrieNode(char c, boolean isEndOfWord) {
            int loc = c - 'a';

            if (children[loc] == null ) {
                children[loc] = new TrieNode();
            }
            this.isEndOfWord = isEndOfWord;

        }
    }

    TrieNode root;
    /** Initialize your data structure here. */
    public Trie() {


    }

    public static void main(String[] args) {
        Trie trie = new Trie();

        trie.insert("apple");
        System.out.println(trie.search("apple"));   // returns true
        System.out.println(trie.search("app"));     // returns false
        System.out.println(trie.startsWith("app")); // returns true
        trie.insert("app");
        System.out.println(trie.search("app"));

    }

    /** Inserts a word into the trie. */
    public void insert(String word) {
        if (word == null || word.trim().length() == 0) {
            return;
        }

        if (root == null) {
            root = new TrieNode();
        }
        TrieNode nd = root;
        for (int i = 0; i < word.length(); i++) {
            int loc = word.charAt(i) - 'a';
            if (nd.children[loc] == null) {
                nd.children[loc] = new TrieNode();
            }

            nd = nd.children[loc];
            if (i == word.length() - 1) {
                nd.isEndOfWord = true;
            }
        }

    }

    /** Returns if the word is in the trie. */
    public boolean search(String word) {
        if (word == null || word.trim().length() == 0 || root == null) {
            return false;
        }

        TrieNode nd = root;

        for (int i = 0; i < word.length() ; i++) {
            int loc = word.charAt(i) - 'a';
            if (nd.children[loc] == null) {
                return false;
            }
            nd = nd.children[loc];

        }
        return (nd != null && nd.isEndOfWord);

    }

    /** Returns if there is any word in the trie that starts with the given prefix. */
    public boolean startsWith(String prefix) {
        if (prefix == null || prefix.trim().length() == 0 || root == null) {
            return false;
        }

        TrieNode nd = root;

        for (int i = 0; i < prefix.length() ; i++) {
            int loc = prefix.charAt(i) - 'a';
            if (nd.children[loc] == null) {
                return false;
            }
            nd = nd.children[loc];

        }
        return nd != null;

    }
}
