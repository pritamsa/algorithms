import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.*;

class PrintTask implements Callable<String> {

    int val;
    PrintTask(int val) {
        this.val = val;
    }


    @Override
    public String call() throws Exception {
        return "PrintTask is running" + val;
    }
}

public class ExecutorServiceDemo {

    public void executionDemo() throws InterruptedException, ExecutionException {
        ExecutorService executorService = Executors.newFixedThreadPool(2);
        Collection<PrintTask> runtasks = new ArrayList<>();

        PrintTask printTask = new PrintTask(99);
        PrintTask printTask1 = new PrintTask(11);

        runtasks.add(printTask);
        runtasks.add(printTask1);

        List<Future<String>> futures = executorService.invokeAll(runtasks);

        for (Future<String> future: futures) {
            System.out.println(future.get());
        }

    }



}
