package additional;

public class Solution {

//  Given a sorted array of doubles, return the squares of these numbers that are still sorted.
//[-9, 2] -> [81, 4]
//    [4, 81]
//
//
//  All -ves are in the begining
//  create a temp arrays
//  copy all -ves to temp1 and all +ve to temp2
//  square both arrays
//  Merge the two arrays so that the result is sorted

//[-8, -5, -4, -3, 0, 1, 12, 15 , 17]
//
//    [-8, -5, -4, -3]
  public double[] sortSquares(double[] arr) {

    if (arr == null || arr.length == 0) {
      return arr;
    }
    int i = 0;
    int idx = -1;

//find 1st +ve : idx = 4; arr of at least 1 element
    while (i < arr.length) {
      if (arr[i] >= 0 ) {
        idx = i;

        break;
      }
      i++;
    }

//Simple case
    if (idx <= 0) {
      if (idx < 0) {
        for (int j = 0; j < arr.length; j++) {
          arr[j] = Math.pow(arr[j],2);
        }

      }
      for (int j = 0; j < arr.length; j++) {
        arr[j] = Math.pow(arr[j],2);
      }
    } else {

      arr = merge(arr,idx);
    }
    return arr;

  }

  //Merge while squaring the values. [-8, -5, -4, -3, 0, 1, 12, 15 , 17]
  //idx = 4 //O(n) time //O(n) space [64, 25, 16, 9, 0, 1, 144, 225, 289]
//left = [9, 16, 25, 64]
  public double[] merge(double[] arr, int idx) {

    double[] left = new double[idx];
    double[] right = new double[arr.length - idx];

//Square left half
    for (int i = 0; i < idx; i++) {
      left[i] = Math.pow(arr[idx-1-i],2);
    }

//Square right half
    int j = 0;
    for (int i = idx; i < arr.length; i++) {
      right[j++] = Math.pow(arr[i],2);
    }

    j = 0;
    int i = 0;
    int k = 0;

//Use the same space without allocating a new array
    while(i < left.length && j < right.length) {
      if (left[i] <= right[j]) {
        arr[k++] = left[i++];

      } else {
        arr[k++] = right[j++];
      }
    }

    while(i < left.length) {
      arr[k++] = left[i++];
    }

    while(j < right.length) {
      arr[k++] = right[j++];
    }


    return arr;
  }




}
