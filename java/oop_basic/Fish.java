// INHERITANCE: Fish gets alive/age/eat()/aging() from Animal for free.
public class Fish extends Animal {
    public int fin;

    public Fish(int age, int fin) {
        super(age);
        this.fin = fin;
    }

    // POLYMORPHISM (subtype): overrides Animal.sleep() with fish-specific
    // behavior. Calling code can still just do fish.sleep() without caring
    // it's a Fish.
    @Override
    public void sleep() {
        System.out.println("this fish is sleeping");
    }

    // Fish-only behavior, not part of the shared Animal interface.
    public void swim() {
        System.out.println("this fish is swimming");
    }
}
