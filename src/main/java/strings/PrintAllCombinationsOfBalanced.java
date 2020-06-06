package strings;

//Write a function to generate all possible n pairs of balanced parentheses.
public class PrintAllCombinationsOfBalanced {
    static void _printParenthesis(char str[], int pos, int n, int open, int close)
    {
        if(close == n)
        {
            // print the possible combinations
            for(int i=0;i<str.length;i++)
                System.out.print(str[i]);
            System.out.println();
            return;
        }
        else
        {
            if(open > close) {
                str[pos] = '}';
                _printParenthesis(str, pos+1, n, open, close+1);
            }
            if(open < n) {
                str[pos] = '{';
                _printParenthesis(str, pos+1, n, open+1, close);
            }
        }
    }

    // Wrapper over _printParenthesis()
    static void printParenthesis(char str[], int n)
    {
        if(n > 0)
            printAll(str, 0, 0, 0, n);
           // _printParenthesis(str, 0, n, 0, 0);
        return;
    }

    public static void printAll(char[] arr, int pos, int open, int close, int n) {
        if (close == n) {
            // print the possible combinations
            for(int i=0;i<arr.length;i++)
                System.out.print(arr[i]);
            System.out.println();
            return;
        } else {
            if(open > close) {
                arr[pos] = '}';
                printAll(arr, pos + 1, open, close + 1, n);

            }

            if (open < n) {
                arr[pos] = '{';
                printAll(arr, pos+1, open+1, close, n);
            }

        }
    }

    // driver program
    public static void main (String[] args)
    {
        int n = 4;
        char[] str = new char[2 * n];
        //printAll(str, 0, 0, 0, n);
        printParenthesis(str, n);
    }
}
