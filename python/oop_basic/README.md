# OOP Basics Crash Course

This folder is a tiny working example of the four core OOP pillars, using
`Animal` and its subclasses `Rabbit`, `Fish`, and `Hawk`.

## 1. Inheritance — `Rabbit`, `Fish`, `Hawk` extend `Animal`

Inheritance lets a class reuse behavior from a parent class instead of
rewriting it.

```python
class Rabbit(Animal):   # Rabbit inherits everything Animal has
    pass
```

`Rabbit`, `Fish`, and `Hawk` all inherit `alive`, `_age`, `eat()`, `sleep()`,
and `aging()` from `Animal` for free ([animal.py](animal.py)). Each subclass
then adds its own extra behavior (`Rabbit.running()`, `Fish.swim()`).

## 2. Polymorphism — three different kinds, not just one

"Polymorphism" covers a few distinct mechanisms. Only the first one needs
inheritance — the other two work with zero shared parent class.

### 2a. Subtype polymorphism — `sleep()` behaves differently per class

The same method name produces different behavior depending on the actual
object's class, resolved at **runtime**, and it *requires* a shared
interface/parent (here, `Animal`).

- `Animal.sleep()` prints `"this animal is sleeping"`.
- `Fish.sleep()` overrides it and prints `"this fish is sleeping"` instead.

Both respond to `.sleep()`, but each does it its own way. (`Hawk` and
`Rabbit` don't override `sleep()`, so they fall back to `Animal`'s version.)

### 2b. Ad-hoc polymorphism (overloading) — `feed()` in [animal.py](animal.py)

Same function *name*, different behavior picked by **argument type**, no
inheritance involved. Python doesn't have native overloading (redefining
`def add` just replaces the previous one), so `feed()` uses
`@functools.singledispatch` to get the same effect: `feed(50)` and
`feed("worms")` each dispatch to a different implementation based on type.

### 2c. Parametric polymorphism (generics) — `oldest()` in [animal.py](animal.py)

One implementation that works across many types because the type itself is
a parameter (`TypeVar`), not because of a shared parent class or per-type
branching. `oldest(a, b)` works the same way whether `a`/`b` are `Fish`s,
`Hawk`s, or `Rabbit`s — it only relies on `_age` existing.

## 3. Encapsulation — `_age` is a protected attribute

Encapsulation means bundling data with the methods that operate on it, and
controlling access to that data instead of letting anything modify it
directly.

`Animal._age` is prefixed with a single underscore, which is Python's
convention for "protected — internal to the class hierarchy, please don't
touch this directly from outside." Instead of doing `rabbit._age += 1`
from outside the class, you're expected to go through the `aging()` method:

```python
def aging(self):
    self._age += 1
    return self._age
```

This keeps the *rule* for how age changes in one place (the class), rather
than scattered across every place that happens to touch age.

> Note: Python doesn't enforce this like `private` in Java/C++ — it's a
> convention (single underscore = protected, double underscore = name-mangled
> "private"). Nothing stops you from accessing `_age` directly, but doing so
> breaks the encapsulation contract.

## 4. Abstraction — `Animal` defines a general interface

Abstraction means exposing a simple, general interface (`eat()`, `sleep()`,
`aging()`) while hiding the specific details of *how* each concrete animal
actually does it.

Code that only cares "can this thing eat and sleep" can work with any
`Animal` subclass without knowing whether it's a `Fish` or a `Hawk` — it just
calls `.eat()` / `.sleep()`. The caller is abstracted away from the concrete
implementation details.

> This example doesn't use Python's `abc.ABC` / `@abstractmethod` (which
> would *force* subclasses to implement a method), but `Animal` still plays
> the role of an abstract base conceptually — it defines the shared shape
> that every animal follows.

## Where to look

| Concept        | File                       | Line(s) to look at                          |
|----------------|----------------------------|----------------------------------------------|
| Inheritance    | [rabbit.py](rabbit.py), [fish.py](fish.py), [hawk.py](hawk.py) | `class X(Animal):` |
| Polymorphism (subtype) | [fish.py](fish.py)  | `sleep()` override                          |
| Polymorphism (ad-hoc)  | [animal.py](animal.py) | `@singledispatch feed()`                  |
| Polymorphism (parametric) | [animal.py](animal.py) | `oldest(a, b)`                         |
| Encapsulation  | [animal.py](animal.py)    | `_age`, `aging()`                           |
| Abstraction    | [animal.py](animal.py)    | `eat()`, `sleep()`, `aging()` interface     |
| Everything together | [main.py](main.py)   | creating instances, calling methods         |
