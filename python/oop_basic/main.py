from rabbit import Rabbit
from hawk import Hawk
from fish import Fish

fish = Fish(age = 10, fin = 11)
rabbit = Rabbit(age = 11)
hawk = Hawk(age = 12, wing = 2)

rabbit.running()
fish.swim()

assert rabbit.aging() == 12
