# OOP Basics Crash Course (Java)

Java port of [`python/oop_basic`](../../python/oop_basic) — same `Animal` /
`Rabbit` / `Fish` / `Hawk` example, showing the four core OOP pillars plus
all three kinds of polymorphism.

## 1. Inheritance — `Rabbit`, `Fish`, `Hawk` extend `Animal`

```java
public class Rabbit extends Animal {   // Rabbit inherits everything Animal has
    ...
}
```

`Rabbit`, `Fish`, and `Hawk` all inherit `alive`, `age`, `eat()`, `sleep()`,
and `aging()` from `Animal` ([Animal.java](Animal.java)). Each subclass then
adds its own extra behavior (`Rabbit.running()`, `Fish.swim()`).

## 2. Polymorphism — three different kinds, not just one

Only the first kind needs inheritance — the other two work with zero shared
parent class.

### 2a. Subtype polymorphism — `sleep()` behaves differently per class

Resolved at **runtime**, based on the object's actual class, and it
*requires* a shared parent/interface (here, `Animal`).

- `Animal.sleep()` prints `"this animal is sleeping"`.
- `Fish.sleep()` (marked `@Override`) prints `"this fish is sleeping"` instead.

(`Hawk` and `Rabbit` don't override `sleep()`, so they fall back to
`Animal`'s version.)

### 2b. Ad-hoc polymorphism (overloading) — `feed()` in [Animal.java](Animal.java)

Same method *name*, different behavior picked by **argument type**, resolved
at **compile time**, no inheritance involved. Java supports this natively —
`feed(int grams)` and `feed(String food)` coexist as separate overloads, and
the compiler picks which one to call based on the argument type:

```java
Animal.feed(50);        // -> "feeding 50 grams of food"
Animal.feed("worms");   // -> "feeding some worms"
```

### 2c. Parametric polymorphism (generics) — `oldest()` in [Animal.java](Animal.java)

One implementation that works across many types because the type itself is
a parameter (`<T extends Animal>`), not because of a shared parent method
override or per-type branching:

```java
public static <T extends Animal> T oldest(T a, T b) {
    return a.getAge() > b.getAge() ? a : b;
}
```

`oldest(fish, hawk)` works the same way regardless of which `Animal`
subtype `fish`/`hawk` actually are.

## 3. Encapsulation — `age` is a protected field

`Animal.age` is `protected`, so it can only be touched from within `Animal`
itself or its subclasses — not from arbitrary outside code. Instead of doing
`rabbit.age += 1` from `Main`, you're expected to go through `aging()`:

```java
public int aging() {
    this.age += 1;
    return this.age;
}
```

This keeps the *rule* for how age changes in one place, rather than
scattered across every caller. `getAge()` exposes a read-only view for
code (like `oldest()`) that needs to compare ages without being able to
mutate them.

## 4. Abstraction — `Animal` defines a general interface

Code that only cares "can this thing eat and sleep" can work with any
`Animal` subclass without knowing whether it's a `Fish` or a `Hawk` — it
just calls `.eat()` / `.sleep()`. The caller is abstracted away from the
concrete implementation details.

> This example doesn't use Java's `abstract class` / `abstract` methods
> (which would *force* subclasses to implement a method), but `Animal`
> still plays the role of an abstract base conceptually — it defines the
> shared shape that every animal follows.

## Where to look

| Concept                    | File                                    | What to look at                     |
|----------------------------|------------------------------------------|--------------------------------------|
| Inheritance                | [Rabbit.java](Rabbit.java), [Fish.java](Fish.java), [Hawk.java](Hawk.java) | `extends Animal` |
| Polymorphism (subtype)     | [Fish.java](Fish.java)                  | `@Override sleep()`                  |
| Polymorphism (ad-hoc)      | [Animal.java](Animal.java)              | overloaded `feed(int)` / `feed(String)` |
| Polymorphism (parametric)  | [Animal.java](Animal.java)              | `<T extends Animal> oldest(T, T)`    |
| Encapsulation              | [Animal.java](Animal.java)              | `protected age`, `aging()`, `getAge()` |
| Abstraction                | [Animal.java](Animal.java)              | `eat()`, `sleep()`, `aging()` interface |
| Everything together        | [Main.java](Main.java)                  | creating instances, calling methods  |

## Running it

```bash
cd java/oop_basic
javac *.java
java -ea Main   # -ea enables the assert statements
```
