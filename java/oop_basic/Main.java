public class Main {
    public static void main(String[] args) {
        Fish fish = new Fish(10, 11);
        Rabbit rabbit = new Rabbit(11);
        Hawk hawk = new Hawk(12, 2);

        // Each subclass's own extra method - not part of the shared Animal interface.
        rabbit.running();
        fish.swim();

        // ENCAPSULATION in action: age is bumped through aging(), not by
        // touching rabbit.age directly. 11 (constructor) + 1 (aging()) == 12.
        assert rabbit.aging() == 12;

        // POLYMORPHISM (ad-hoc): same feed() call, dispatched by argument type
        // at COMPILE time (the compiler picks the overload to call).
        Animal.feed(50);       // -> "feeding 50 grams of food"
        Animal.feed("worms");  // -> "feeding some worms"

        // POLYMORPHISM (parametric): same oldest(), works for any Animal subtype.
        assert Animal.oldest(fish, hawk) == hawk;

        System.out.println("all assertions passed");
    }
}
