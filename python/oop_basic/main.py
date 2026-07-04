from animal import feed, oldest
from rabbit import Rabbit
from hawk import Hawk
from fish import Fish

fish = Fish(age = 10, fin = 11)
rabbit = Rabbit(age = 11)
hawk = Hawk(age = 12, wing = 2)

# Each subclass's own extra method — not part of the shared Animal interface.
rabbit.running()
fish.swim()

# ENCAPSULATION in action: age is bumped through aging(), not by touching
# rabbit._age directly. 11 (constructor) + 1 (aging()) == 12.
assert rabbit.aging() == 12

# POLYMORPHISM (ad-hoc): same feed() call, dispatched by argument type.
feed(50)          # -> "feeding 50 grams of food"
feed("worms")     # -> "feeding some worms"

# POLYMORPHISM (parametric): same oldest(), works for any Animal subtype.
assert oldest(fish, hawk) is hawk
