package arrays;

//Repeat and Missing Number Array:
//        Please Note:
//        There are certain problems which are asked in the interview to also check how you take care of overflows in your problem.
//        This is one of those problems.
//        Please take extra care to make sure that you are type-casting your ints to long properly and at all places. Try to verify if your solution works if number of elements is as large as 105
//
//        Food for thought :
//        Even though it might not be required in this problem, in some cases, you might be required to order the operations cleverly so that the numbers do not overflow.
//        For example, if you need to calculate n! / k! where n! is factorial(n), one approach is to calculate factorial(n), factorial(k) and then divide them.
//        Another approach is to only multiple numbers from k + 1 ... n to calculate the result.
//        Obviously approach 1 is more susceptible to overflows.
//        You are given a read only array of n integers from 1 to n.
//
//        Each integer appears exactly once except A which appears twice and B which is missing.
//
//        Return A and B.
//
//        Note: Your algorithm should have a linear runtime complexity. Could you implement it without using extra memory?
//
//        Note that in your output A should precede B.
//
//        Example: 1 2 2 4 5
//
//        Input:[3 1 2 5 3]
//
//        Output:[3, 4]
//
//        A = 3, B = 4

//Given a read only array of n + 1 integers between 1 and n, find one number that repeats in linear time using less than O(n) space and traversing the stream sequentially O(1) times.
//
//        Sample Input:
//
//        [3 4 1 4 1]
//        Sample Output:
//
//        1
//        If there are multiple possible answers ( like in the sample case above ), output any one.
//
//        If there is no duplicate, output -1
import java.util.ArrayList;
//Since the array has 1 to n elements the length is going to be enough
//idea is : iterate and for each elm, find the elm which is at the idx = elm. Make it -ve. You find already visited value
// if you encounter a -ve number.
public class RepeatandMissingNumberArray {

    ArrayList<Integer> repeatAndMission(ArrayList<Integer> arr) {



        return null;


    }

}
