package additional;

import java.util.ArrayList;
import java.util.List;

//Given the root of a binary tree, find the maximum value V for which there exists different nodes A and B where V = |A.val - B.val| and A is an ancestor of B.
//
//        (A node A is an ancestor of B if either: any child of A is equal to B, or any child of A is an ancestor of B.)
public class MaxAncestaralDiff {


    int max = Integer.MIN_VALUE;

    public static void main(String[] args) {
        TreeNode nd = new TreeNode(0);
        nd.right = new TreeNode(1);
        int m = (new MaxAncestaralDiff()).maxAncestorDiff(nd);
    }
    public int maxAncestorDiff(TreeNode root) {



        if(root == null || (root.left == null && root.right == null)) {
            return 0;
        }
        maxAncestorDiff(root, new ArrayList());
        return max;


    }

    public void maxAncestorDiff(TreeNode root, List<TreeNode> nds) {
        if(nds == null) nds = new ArrayList<>();
        if(root == null ) {
            return;
        }
        if(root.left == null && root.right == null) {
            nds.add(root);
            return;
        }

        List<TreeNode> leftNds = new ArrayList<>() ;
        List<TreeNode> rightNds = new ArrayList<>();
        maxAncestorDiff(root.left, leftNds);
        maxAncestorDiff(root.right, rightNds);

        nds.add(root);
        nds.addAll(leftNds);
        nds.addAll(rightNds);

        for(TreeNode nd: leftNds) {
            if(nd != null && Math.abs(nd.val-root.val) > max) {
                max = Math.abs(nd.val-root.val);
            }
        }

        for(TreeNode nd: rightNds) {
            if(nd != null && Math.abs(nd.val-root.val) > max) {
                max = Math.abs(nd.val-root.val);
            }
        }


    }

//    Given an array A of positive integers (not necessarily distinct), return the lexicographically largest permutation that is smaller than A, that can be made with one swap (A swap exchanges the positions of two numbers A[i] and A[j]).  If it cannot be done, then return the same array.
//
//
//
//            Example 1:
//
//    Input: [3,2,1]
//    Output: [3,1,2]
//    Explanation: Swapping 2 and 1.
//    Example 2:
//
//    Input: [1,1,5]
//    Output: [1,1,5]
//    Explanation: This is already the smallest permutation.
//    Example 3:
//
//    Input: [1,9,4,6,7]
//    Output: [1,7,4,6,9]
//    Explanation: Swapping 9 and 7.
//    Example 4:
//
//    Input: [3,1,1,3]
//    Output: [1,3,1,3]
//    Explanation: Swapping 1 and 3.
//
//
//    Note:
//
//            1 <= A.length <= 10000
//            1 <= A[i] <= 10000
//
//    Prev
//3 / 3
//
//    Next
//
//            Autocomplete
//
//
//
// != null
//         1

    public int[] prevPermOpt1(int[] A) {
        if(A == null || A.length == 0){
            return A;
        }

        int n = A.length;
        int i = n - 1;
        while(i >= 1 && A[i] >= A[i - 1]){
            i--;
        }

        if(i == 0){
            return A;
        }

        int first = i - 1;
        int second = i;
        for(int j = i + 1; j < n; j++){
            if(A[j] > A[j - 1] && A[j] < A[first]){
                second = j;
            }
        }

        int temp = A[first];
        A[first] = A[second];
        A[second] = temp;
        return A;

    }

//    Return the lexicographically smallest subsequence of text that contains all the distinct characters of text exactly once.
//
//            Example 1:
//
//    Input: "cdadabcc"
//    Output: "adbc"
//    Example 2:
//
//    Input: "abcd"
//    Output: "abcd"
//    Example 3:
//
//    Input: "ecbacba"
//    Output: "eacb"
//    Example 4:
//
//    Input: "leetcode"
//    Output: "letcod"
//
//
//    Constraints:
//
//            1 <= text.length <= 1000
//    text consists of lowercase English letters.
//    Note: This question is the same as 316: https://leetcode.com/problems/remove-duplicate-letters/
//


}