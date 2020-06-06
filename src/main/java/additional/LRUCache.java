package additional;

import java.util.Comparator;
import java.util.HashMap;
import java.util.PriorityQueue;

class CacheEntry {
    long timeStamp;
    int value;
    int key;

    CacheEntry(int key, int val, long timeStamp) {
        this.timeStamp = timeStamp;
        this.value = val;
        this.key = key;
    }
}
class CacheEntryComp implements Comparator<CacheEntry> {

    @Override
    public int compare(CacheEntry o1, CacheEntry o2) {
        return o1.timeStamp>o2.timeStamp ? 1 : ((o1.timeStamp == o2.timeStamp) ? 0: -1);
    }
}
public class LRUCache {

    public static void main(String[] args) {
        LRUCache cache = new LRUCache(2);

        cache.get(2);

        cache.put(2, 1);
        cache.put(1, 1);
        cache.put(2, 3);
        cache.put(4, 1);
        System.out.println(cache.get(1));
        System.out.println(cache.get(2));
//        cache.put(1, 5);
//        cache.put(1, 2);
//        System.out.println(cache.get(1));
//        System.out.println(cache.get(2));

//        cache.put(2, 2);
//        int i = cache.get(1);       // returns 1
//        System.out.println(i);
//        cache.put(3, 3);    // evicts key 2
//        System.out.println(cache.get(2));       // returns -1 (not found)
//        cache.put(4, 4);    // evicts key 1
//        System.out.println(cache.get(1));       // returns -1 (not found)
//        System.out.println(cache.get(3));       // returns 3
//        System.out.println(cache.get(4));       // returns 4

    }


    PriorityQueue<CacheEntry> minHeap;
    HashMap<Integer, CacheEntry> map;
    int capacity;

    public LRUCache(int capacity) {
        this.capacity = capacity;
        map = new HashMap<>(capacity);
        minHeap = new PriorityQueue<>(new CacheEntryComp());

    }

    public int get(int key) {
        CacheEntry entry = map.get(key);
        if(entry == null) {
            return -1;
        }
        minHeap.remove(entry);
        entry.timeStamp = System.currentTimeMillis();
        minHeap.add(entry);
        return entry.value;

    }

    public void put(int key, int value) {

            long currTime = System.currentTimeMillis();
            CacheEntry newEntry = new CacheEntry(key, value, currTime);
            if (minHeap.size() == capacity) {
                CacheEntry earliestEntry = minHeap.remove();
                map.remove(earliestEntry.key);
            }
            minHeap.add(newEntry);
            map.put(key, newEntry);
        }

}
