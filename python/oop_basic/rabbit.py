# INHERITANCE: simplest example in this folder — Rabbit adds nothing but a
# new method, and reuses Animal's __init__/eat()/sleep()/aging() as-is.
from animal import Animal

class Rabbit(Animal):
    pass

    def running(self):
        print("this rabbit is running")