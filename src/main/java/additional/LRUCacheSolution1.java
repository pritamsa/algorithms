package additional;

import java.util.HashSet;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.Map;

public class LRUCacheSolution1 extends LinkedHashMap<Integer, Integer> {
    private int capacity;

    public LRUCacheSolution1(int capacity) {
        super(capacity, 0.75F, true);
        this.capacity = capacity;
        HashSet<Integer> set = new HashSet<>();
        Iterator<Integer> iter = set.iterator();

        while (iter.hasNext()) {

        }
    }

    public int get(int key) {
        return super.getOrDefault(key, -1);
    }

    public void put(int key, int value) {
        super.put(key, value);
    }

    @Override
    protected boolean removeEldestEntry(Map.Entry<Integer, Integer> eldest) {
        return size() > capacity;
    }
}

