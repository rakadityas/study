// INHERITANCE: simplest example in this folder - Rabbit adds nothing but a
// new method, and reuses Animal's constructor/eat()/sleep()/aging() as-is.
public class Rabbit extends Animal {

    public Rabbit(int age) {
        super(age);
    }

    public void running() {
        System.out.println("this rabbit is running");
    }
}
