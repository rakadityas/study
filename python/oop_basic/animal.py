from functools import singledispatch
from typing import TypeVar

# ABSTRACTION: Animal defines the general interface (eat/sleep/aging) that
# every subclass shares, without callers needing to know the concrete type.
class Animal:
    alive = True
    # ENCAPSULATION: leading underscore = "protected", internal state that
    # should be changed via a method (aging()) rather than set directly.
    _age = 0

    def __init__(self, age):
        self._age = age

    def eat(self):
        print("this animal is eating")

    # POLYMORPHISM: subclasses (e.g. Fish) can override this to behave
    # differently while still responding to the same .sleep() call.
    def sleep(self):
        print("this animal is sleeping")

    # ENCAPSULATION: the only sanctioned way to mutate _age.
    def aging(self):
        self._age += 1
        return self._age


# POLYMORPHISM (ad-hoc / overloading): same name "feed", behavior picked by
# the food's TYPE, resolved at call time. No inheritance involved at all —
# contrast with sleep() above, which needs the Animal/Fish hierarchy.
@singledispatch
def feed(food):
    raise TypeError(f"don't know how to feed: {type(food)}")

@feed.register
def _(food: int):
    print(f"feeding {food} grams of food")

@feed.register
def _(food: str):
    print(f"feeding some {food}")


# POLYMORPHISM (parametric / generics): one implementation that works across
# any Animal subtype because the type is a parameter (T), not because of
# per-type branching or a shared method override.
T = TypeVar("T", bound=Animal)

def oldest(a: T, b: T) -> T:
    return a if a._age > b._age else b
