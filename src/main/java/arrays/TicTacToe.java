package arrays;

class TicTacToe {

    char[][] board;

    /** Initialize your data structure here. */
    public TicTacToe(int n) {
        board = new char[n][n];
    }

    /** Player {player} makes a move at ({row}, {col}).
     @param row The row of the board.
     @param col The column of the board.
     @param player The player, can be either 1 or 2.
     @return The current winning condition, can be either:
     0: No one wins.
     1: Player 1 wins.
     2: Player 2 wins. */
    public int move(int row, int col, int player) {
        char c = (player == 1) ? 'X' :'O';
        if ((row >= 0 && row <= board.length) && (col >=0 && col <= board[row].length)) {
            board[row][col] = c;
            boolean win = isRowWinning(row, c) || isColWinning(col, c)
                    || isDiagonalWinning (row, col, c);
            if (win) {
                return player;
            }
        }
        return 0;
    }

    private boolean isRowWinning(int row, char c) {
        for (int i = 0; i < board[row].length; i++) {
            if (board[row][i] != c) {
                return false;
            }
        }
        return true;
    }
    private boolean isColWinning(int col, char c) {
        for (int i = 0; i < board.length; i++) {
            if (board[i][col] != c) {
                return false;
            }
        }
        return true;
    }

    private boolean isDiagonalWinning(int row, int col, char c) {
        boolean leftDia = (board[0][0] == c);
        boolean rightDia = (board[0][board.length-1] == c);
        if ((row == col) || (row+col == board.length - 1)) {

            if (leftDia) {
                for (int i = 1; i < board.length; i++) {
                    leftDia = leftDia&&(board[i][i] == c);
                }



                if (leftDia) {
                    return true;
                }
            }

            if (rightDia) {
                for(int i = 1; i < board.length; i++ ) {
                    rightDia = rightDia && (board[i][board.length-1-i] == c);

                }
                if (rightDia) {
                    return true;
                }
            }

        }
        return false;

    }

    public static void main(String[] args) {
        TicTacToe ticTacToe = new TicTacToe(2);
        ticTacToe.move(0, 1, 1);
        ticTacToe.move(1, 1, 2);

        ticTacToe.move(1, 0, 1);
//        ticTacToe.move(1, 1, 2);
//
//        ticTacToe.move(2, 0, 1);
//
//        ticTacToe.move(1, 0, 2);
//        ticTacToe.move(2, 1, 1);
    }
}