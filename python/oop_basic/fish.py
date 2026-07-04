from animal import Animal

# INHERITANCE: Fish gets alive/_age/eat()/aging() from Animal for free.
class Fish(Animal):
    fin = None

    def __init__(self, age, fin):
        self._age = age
        self.fin = fin


    # POLYMORPHISM: overrides Animal.sleep() with fish-specific behavior.
    # Calling code can still just do fish.sleep() without caring it's a Fish.
    def sleep(self):
        print("this fish is sleeping")

    # Fish-only behavior, not part of the shared Animal interface.
    def swim(self):
        print("this fish is swimming")