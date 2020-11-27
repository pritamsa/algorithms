package arrays;

public class StockProfit {

    int maxProfit(int[] prices) {
        int profit = 0;
        int min = prices[0];
        StringBuffer sb = new StringBuffer();



        for (int i = 0; i < prices.length; i++) {
            profit = Math.max(profit, prices[i] - min);
            min = Math.min(min, prices[i]);
        }
        return profit;
    }
}
