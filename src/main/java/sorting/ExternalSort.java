package sorting;

import java.io.*;
import java.util.*;

public class ExternalSort {

    public static void main(String[] args) throws IOException {
        File file = new File("C:\\Users\\preetam\\strs.txt");

        ExternalSort e = new ExternalSort();
        e.makeSmallBlocks(file);

    }

    public void makeSmallBlocks(File longFile) throws IOException {
        List<File> files = new LinkedList<>();
        long blockSize = getBlockSize(longFile);

        FileReader fd = new FileReader(longFile);
        BufferedReader bfd = new BufferedReader(fd);
        String line = "";
        List<String> lst = new LinkedList<>();
        long currBlockSize = 0;
        while (line != null) {

            while (currBlockSize < blockSize && line != null) {
                line = bfd.readLine();
                if (line != null) {
                    currBlockSize += line.length();
                    lst.add(line);
                }
            }
            if (lst != null && lst.size() > 0) {
                File fl = sortAndSave(lst);
                files.add(fl);
                currBlockSize = 0;
                lst.clear();
            }
        }
        File opFile = new File("C:\\\\Users\\\\preetam\\\\sorted.txt");
        mergeSortedFiles(files, opFile);

    }

    public File sortAndSave(List<String> lst) throws IOException {
        Collections.sort(lst);
        File tempFile = File.createTempFile("C:\\Users\\preetam\\tmp" + System.currentTimeMillis(), ".txt");
        tempFile.deleteOnExit();
        BufferedWriter fbw = new BufferedWriter(new FileWriter(tempFile));


        try {
            for (int i = 0; i < lst.size(); i++) {
                fbw.write(lst.get(i));
                fbw.newLine();
            }
        } finally {
            fbw.close();
        }
        return tempFile;
    }

    public long getBlockSize(File file) {
        int MAX_OPEN_FILES = 1024;

        long fileSize = file.length();
        long blockSize = fileSize/MAX_OPEN_FILES;
        if (Runtime.getRuntime().freeMemory()/2 < blockSize) {
            return Runtime.getRuntime().freeMemory()/2;
        } else {
            return 30;
        }
    }


    public void mergeSortedFiles(List<File> files, File opFile) throws IOException {

        if (files != null && files.size() > 0) {
            BufferedWriter fbw = new BufferedWriter(new FileWriter(opFile));

            PriorityQueue<FileBuffer> pq = new PriorityQueue<>();

            for (int i = 0; i < files.size() ; i++) {
                FileBuffer fbr = new FileBuffer(files.get(i));
                pq.add(fbr);
            }

            while(!pq.isEmpty()) {
                FileBuffer smallestFbr = pq.poll();
                if (smallestFbr.empty) {
                    continue;
                } else {
                    Integer val = smallestFbr.pop();
                    fbw.write(val);
                }
            }
            fbw.close();
        }


    }

}

class FileBufferComp implements Comparator<FileBuffer> {



    @Override
    public int compare(FileBuffer o1, FileBuffer o2) {

        try {
            if (o1.peek().compareTo(o2.peek()) > 0) {
                return 1;
            } else if (o1.peek().compareTo(o2.peek()) < 0) {
                return -1;
            } else {
                return 0;
            }

        } catch (IOException e) {
            e.printStackTrace();
        }
        return 0;
    }
}

//Create a fileBuffer first
class FileBuffer {

    private static final Integer MAX_CAPACITY = 2048;
    Integer topLine;
    int currentPointer;
    File f;
    boolean empty;


    FileBuffer(File file) throws IOException {
        this.f = file;
        currentPointer = 0;
        topLine = new Integer("");
        reLoad();
    }

    public void reLoad() throws IOException {
        BufferedReader br = new BufferedReader(new FileReader(f), MAX_CAPACITY);
        String strLine = br.readLine();
        if (strLine == null) {
            empty = true;
        } else {
            topLine = Integer.parseInt(strLine);
            empty = false;
        }
    }

    public Integer peek() throws IOException {

        if (empty) {
            return null;
        }
        return topLine;
    }

    public Integer pop() throws IOException {
        Integer tmp = peek();
        reLoad();
        return tmp;

    }


}