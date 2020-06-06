package executorservice;

import java.util.concurrent.*;

public class ExecutorServiceDemo {

    public void execudemo() throws InterruptedException, ExecutionException, TimeoutException {
        ExecutorService exSvc = Executors.newFixedThreadPool(10);
        Future<Double> res = exSvc.submit(new ServiceHandler());

        Double val = res.get(4, TimeUnit.MILLISECONDS);

    }

    public void execudemo1() throws InterruptedException, ExecutionException, TimeoutException {
        ExecutorService exSvc = Executors.newFixedThreadPool(10);
        Future<Double> res = exSvc.submit(new ServiceHandler());

        Double val = res.get(4, TimeUnit.MILLISECONDS);

    }
}
