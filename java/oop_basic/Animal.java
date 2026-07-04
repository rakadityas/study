// ABSTRACTION: Animal defines the general interface (eat/sleep/aging) that
// every subclass shares, without callers needing to know the concrete type.
public class Animal {
    public boolean alive = true;
    // ENCAPSULATION: protected - internal state that should be changed via
    // a method (aging()) rather than set directly from outside the class.
    protected int age = 0;

    public Animal(int age) {
        this.age = age;
    }

    public void eat() {
        System.out.println("this animal is eating");
    }

    // POLYMORPHISM (subtype): subclasses (e.g. Fish) can override this to
    // behave differently while still responding to the same .sleep() call,
    // resolved at runtime based on the object's actual class.
    public void sleep() {
        System.out.println("this animal is sleeping");
    }

    // ENCAPSULATION: the only sanctioned way to mutate age.
    public int aging() {
        this.age += 1;
        return this.age;
    }

    public int getAge() {
        return this.age;
    }

    // POLYMORPHISM (ad-hoc / overloading): same name "feed", different
    // behavior picked by the argument's TYPE, resolved at COMPILE time.
    // No inheritance involved at all - contrast with sleep() above.
    public static void feed(int grams) {
        System.out.println("feeding " + grams + " grams of food");
    }

    public static void feed(String food) {
        System.out.println("feeding some " + food);
    }

    // POLYMORPHISM (parametric / generics): one implementation that works
    // across any Animal subtype because the type is a parameter (<T>), not
    // because of per-type branching or a shared method override.
    public static <T extends Animal> T oldest(T a, T b) {
        return a.getAge() > b.getAge() ? a : b;
    }
}
