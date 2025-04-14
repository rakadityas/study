from animal import Animal

class Hawk(Animal):
    wing = None

    def __init__(self, age, wing):
        self._age = age
        self._wing = wing