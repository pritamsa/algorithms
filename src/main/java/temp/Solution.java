package temp;


import java.io.*;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

class Comp implements Comparator<Long> {
    //@Override
    public int compare(Long o1, Long o2) {

        return o2.compareTo(o1) > 1 ? 1 : o2.compareTo(o1) < 0 ? -1 : 0;
    }
}
class Result {

    /*
     * Complete the 'getMaxUnits' function below.
     *
     * The function is expected to return a LONG_INTEGER.
     * The function accepts following parameters:
     *  1. LONG_INTEGER_ARRAY boxes
     *  2. LONG_INTEGER_ARRAY unitsPerBox
     *  3. LONG_INTEGER truckSize
     */

    public static long getMaxUnits(List<Long> boxes, List<Long> unitsPerBox, long truckSize) {
        // Write your code here

        if (boxes == null || unitsPerBox == null || truckSize <=0 ) {
            return -1;
        }

        TreeMap<Long, Long> unitsSorted = new TreeMap<>(Comparator.reverseOrder());
        PriorityQueue<Long> sortedUnits = new PriorityQueue<>(new Comp());

        for (Long unit: unitsPerBox) {
            unitsSorted.put(unit, getBoxesOfUnits(unitsPerBox, boxes, unit));

        }
        //Prepare a max heap to get max units every time
        for (Long unit: unitsPerBox) {
            sortedUnits.add(unit);
        }

        long units = 0;
        while (truckSize > 0) {
            Long maxUnits = unitsSorted.firstKey();//.peek();
            Long availBoxes = 0L;

            if (maxUnits != null && maxUnits > 0) {
                availBoxes = unitsSorted.get(maxUnits);


                    if (truckSize > availBoxes) {
                        truckSize -= availBoxes;
                        //boxes.add(j, 0L);
                        if (availBoxes == 0) {
                            sortedUnits.remove();
                        }
                        units += availBoxes*maxUnits;
                        unitsSorted.put(maxUnits, availBoxes);
                    } else {
                        truckSize = 0;
                        availBoxes -= truckSize;
                        units += availBoxes*maxUnits;
                        unitsSorted.put(maxUnits, availBoxes);

                    }
                //}
            }

        }
        return units;

    }

    private static long getBoxesOfUnits(List<Long> unitsPerBoxes, List<Long> boxes, Long units) {
        long boxesOfUnits = 0L;
        int j = 0;
        for (int i = 0; i < unitsPerBoxes.size(); i++) {
            if (unitsPerBoxes.get(i) != null && unitsPerBoxes.get(i).equals(units)
                    && boxes.get(i) != null && boxes.get(i) > 0L) {
                boxesOfUnits += boxes.get(i);

                j++;
            }

        }
        return boxesOfUnits;
    }

    private static ArrayList<Integer> getIndicesOfUnits(List<Long> unitsPerBoxes, List<Long> boxes, Long units) {

        ArrayList<Integer> indOfUnits = new ArrayList<>();
        int j = 0;
        for (int i = 0; i < unitsPerBoxes.size(); i++) {
            if (unitsPerBoxes.get(i) != null && unitsPerBoxes.get(i).equals(units)
                    && boxes.get(i) != null && boxes.get(i) > 0L) {
                indOfUnits.add(i);
                j++;
            }

        }
        return indOfUnits;
    }

}

public class Solution {
    public static void main(String[] args) throws IOException {
        BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(System.in));
        //BufferedWriter bufferedWriter = new BufferedWriter(new FileWriter(System.getenv("OUTPUT_PATH")));

        int boxesCount = Integer.parseInt(bufferedReader.readLine().trim());

        List<Long> boxes = IntStream.range(0, boxesCount).mapToObj(i -> {
            try {
                return bufferedReader.readLine().replaceAll("\\s+$", "");
            } catch (IOException ex) {
                throw new RuntimeException(ex);
            }
        })
                .map(String::trim)
                .map(Long::parseLong)
                .collect(Collectors.toList());

        int unitsPerBoxCount = Integer.parseInt(bufferedReader.readLine().trim());

        List<Long> unitsPerBox = IntStream.range(0, unitsPerBoxCount).mapToObj(i -> {
            try {
                return bufferedReader.readLine().replaceAll("\\s+$", "");
            } catch (IOException ex) {
                throw new RuntimeException(ex);
            }
        })
                .map(String::trim)
                .map(Long::parseLong)
                .collect(Collectors.toList());

        long truckSize = Long.parseLong(bufferedReader.readLine().trim());

        long result = Result.getMaxUnits(boxes, unitsPerBox, truckSize);

        //bufferedWriter.write(String.valueOf(result));
        //bufferedWriter.newLine();

        bufferedReader.close();
        //bufferedWriter.close();
    }
}
