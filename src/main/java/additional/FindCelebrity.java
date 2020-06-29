package additional;

//At a party there is one celebrity. Your job is to find the celeb or report of there is no celeb.
//A celebrity is a person who knows no one in the part but everyone at the party knows the celebrity.
public class FindCelebrity {

    public static void main(String[] args) {
        int[][] mat = new int[3][3];
        mat[0][0] = 1;
        mat[0][1] = 0;
        mat[0][2] = 1;

        mat[1][0] = 1;
        mat[1][1] = 1;
        mat[1][2] = 0;

        mat[2][0] = 0;
        mat[2][1] = 1;
        mat[2][2] = 1;

    }
    public int findCelebrity(int n, int[][] mat) {

        if (n <= 0) {
            return -1;
        } else if (n == 1) {
            return n;
        }

        //Assume celeb as 0
        int celeb = 0;
        //As soon as assumed celeb "knows" someone, that is out of celeb status. We now assume i as next celeb
        for (int i = 1; i < n; i++) {
            if (knows(celeb, i, mat)) {
                celeb = i;
            }
        }

        //Check if this assumed celeb is really a celeb by checking the celeb does not know anyone
        for (int i = 0; i < n; i++) {
            if(i != celeb && knows(celeb, i, mat)) {
                return -1;
            }
        }

        //Complete one more check by confirming everyone else "knows" the celeb.
        for (int i = 0; i < n; i++) {
            if(i != celeb && !knows(i, celeb, mat)) {
                return -1;
            }
        }
        return celeb;

    }

    private boolean knows(int a , int b, int[][] mat) {
        return mat[a][b] == 1;
    }

}
