package additional;

import java.util.ArrayList;
import java.util.List;

public class CountSmaller {
//    You are given an integer array nums and you have to return a new counts array. The counts array has the property where counts[i] is the number of smaller elements to the right of nums[i].
//
//    Example:
//
//    Input: [5,2,6,1]
//    Output: [2,1,1,0]
//    Explanation:
//    To the right of 5 there are 2 smaller elements (2 and 1).
//    To the right of 2 there is only 1 smaller element (1).
//    To the right of 6 there is 1 smaller element (1).
//    To the right of 1 there is 0 smaller element.
//(nlogn) with additional space of 2n
//Another approach is to build a binary search tree. 3n additional space(for each node, its value, left & right)
    //Time complexity of tree will be O(n)
    public List<Integer> countSmaller(int[] nums) {
        List<Integer> ans = new ArrayList<>();

        int n = nums.length;

        int[][] arr = new int[n][2];

        //In each array loc, save original loc and the value
        for(int i=0;i<n;i++)
            arr[i] = new int[]{nums[i],i};

        //counts array
        int[] count = new int[n];

        //call merge sort to sort that array that we just created and while soring add count
        sort(0,n-1,count,new int[n][2],arr);

        for(int c : count)
            ans.add(c);

        return ans;
    }

    private void sort(int start,int end,int[] count,int[][] temp,int[][] nums){
        if(start >= end)
            return;

        int mid = start + (end - start)/2;

        sort(start,mid,count,temp,nums);
        sort(mid+1,end,count,temp,nums);
        merge(start,end,count,temp,nums);
    }

    private void merge(int start,int end,int[] count,int[][] temp,int[][] nums){
        int mid = (start+end)/2;

        int left = start;
        int right = mid+1;
        int k = start;

        while(left <= mid && right <= end){
            if(nums[left][0] <= nums[right][0]){
                count[nums[left][1]] += right - (mid + 1);
                temp[k++] = nums[left++];
            }else{
                temp[k++] = nums[right++];
            }
        }

        while(left <= mid){
            count[nums[left][1]] += right - (mid + 1);
            temp[k++] = nums[left++];
        }

        while(right <= end)
            temp[k++] = nums[right++];

        for(int i=start;i<=end;i++)
            nums[i] = temp[i];
    }
}
