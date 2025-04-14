from animal import Animal

class Fish(Animal):
    fin = None

    def __init__(self, age, fin):
        self._age = age
        self.fin = fin
        

    def sleep(self):
        print("this fish is sleeping")

    def swim(self):
        print("this fish is swimming")