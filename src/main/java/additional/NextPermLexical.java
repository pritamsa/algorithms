package additional;

//In an array to find next lexically ordered permutation
//find largest i such that a[i-1] < a[i]
//find largest j >= i such that a[j] > a[i-1]
//swap a[i-1] and a[j]
//reverse order of elements froma[i]
public class NextPermLexical {

    public void getNextPermLexical(int[] num) {

        int len = num.length - 1;

        int firstIdx = 0;
        int secondIdx = 0;
        for (int i = len; i >= 0;i--) {
            if( i > 0 && num[i-1] < num[i]) {
                firstIdx = i-1;
                secondIdx = i;
                break;
            }
        }
        int id = secondIdx;
        if (firstIdx != -1) {
            for (int i = num.length - 1; i >= id; i--) {
                if(num[i] > num[firstIdx]) {
                    secondIdx = i;
                    break;
                }
            }
        }

        int temp = num[firstIdx];
        num[firstIdx] = num[secondIdx];
        num[secondIdx] = temp;

        //reverse the order after firstIdx+1
        reversePart(num,firstIdx+1, num.length-1 );


    }

    private void reversePart(int[] arr, int st, int en) {
        while(st<en) {
            int temp = arr[st];
            arr[st] = arr[en];
            arr[en] = temp;
            st++;
            en--;
        }
    }
}
