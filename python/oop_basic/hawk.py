from animal import Animal

# INHERITANCE: Hawk gets alive/eat()/sleep()/aging() from Animal untouched
# (no override here, so Hawk.sleep() falls back to Animal's version).
class Hawk(Animal):
    wing = None

    def __init__(self, age, wing):
        self._age = age
        self._wing = wing