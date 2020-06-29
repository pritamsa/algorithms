package additional;

import java.util.Arrays;
import java.util.PriorityQueue;

//There are N workers.  The i-th worker has a quality[i] and a minimum wage expectation wage[i].
//
//        Now we want to hire exactly K workers to form a paid group.  When hiring a group of K workers, we must pay them according to the following rules:
//
//        Every worker in the paid group should be paid in the ratio of their quality compared to other workers in the paid group.
//        Every worker in the paid group must be paid at least their minimum wage expectation.
//        Return the least amount of money needed to form a paid group satisfying the above conditions.
//Idea is to sort wage to quality ratios. You can achieve this by Arrays.sort with a comparator that
// compares wage/quality ratio
//Then iterate through wage to quality ratios to take the min first, calcuate wage(s) for each worker by multiplying
//quality of that worker by wage/uality min val. If during this loop, any of the worker ers less than expected, continue to
//next wage to quality ratio.
//For each wage to quality, calculate and sort prices. Give addition of min K prices.
//
class Worker implements Comparable<Worker> {
    public int quality, wage;
    public Worker(int q, int w) {
        quality = q;
        wage = w;
    }

    public double ratio() {
        return (double) wage / quality;
    }

    public int compareTo(Worker other) {
        return Double.compare(ratio(), other.ratio());
    }
}
public class MinCostOfHire {

    public static void main(String[] args) {
        int[] quality = {10,20,5};
        int[] wage = {70,50,30};
                int K = 2;
        mincostToHireWorkers(quality, wage, K);
    }


    public static double mincostToHireWorkers(int[] quality, int[] wage, int K) {
        int N = quality.length;
        double ans = 1e9;

        for (int captain = 0; captain < N; ++captain) {
            // Must pay at least wage[captain] / quality[captain] per qual
            double factor = (double) wage[captain] / quality[captain];
            double prices[] = new double[N];
            int t = 0;
            for (int worker = 0; worker < N; ++worker) {
                double price = factor * quality[worker];
                if (price < wage[worker]) continue;
                prices[t++] = price;
            }

            if (t < K) continue;
            Arrays.sort(prices, 0, t);
            double cand = 0;
            for (int i = 0; i < K; ++i)
                cand += prices[i];
            ans = Math.min(ans, cand);
        }

        return ans;
    }
}
