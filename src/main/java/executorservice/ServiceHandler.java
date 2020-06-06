package executorservice;

import java.util.concurrent.Callable;

public class ServiceHandler implements Callable {
    @Override
    public Double call() throws Exception {
        return Math.random();
    }
}
