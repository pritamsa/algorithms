package additional;



    public class DoubleTon {
        private static final DoubleTon[] doubleTon = {new DoubleTon(),new DoubleTon()};
        private static int index =0;

        private DoubleTon(){
        }

        public static DoubleTon getInstance(){
            synchronized (doubleTon) {
                return doubleTon[index++%2];
            }
        }

    }

