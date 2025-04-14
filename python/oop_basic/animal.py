class Animal:
    alive = True
    _age = 0

    def __init__(self, age):
        self._age = age

    def eat(self):
        print("this animal is eating")

    def sleep(self):
        print("this animal is sleeping")

    def aging(self):
        self._age += 1
        return self._age
    