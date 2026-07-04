// INHERITANCE: Hawk gets alive/eat()/sleep()/aging() from Animal untouched
// (no override here, so Hawk.sleep() falls back to Animal's version).
public class Hawk extends Animal {
    public int wing;

    public Hawk(int age, int wing) {
        super(age);
        this.wing = wing;
    }
}
