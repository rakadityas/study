// ============================================================
// JS SYNTAX PLAYGROUND — run with `node js-syntax.js`
// Companion code for ../README.md section 1 & 2
// ============================================================

// ---- 1. Variable declarations ----
var a = 1;      // function-scoped, hoisted + initialized as undefined, redeclarable
let b = 2;      // block-scoped, hoisted but in "temporal dead zone" until assigned
const c = 3;    // block-scoped, cannot be reassigned (object CONTENTS can still mutate)

const obj = { x: 1 };
obj.x = 2;      // OK — const only locks the binding, not the object
// obj = {};    // TypeError

console.log({ a, b, c, obj });

// ---- 2. Function variations ----

// 2a. Function declaration — hoisted entirely, can be called before defined
function add(x, y) { return x + y; }

// 2b. Function expression — not hoisted (the variable is, the value isn't)
const subtract = function (x, y) { return x - y; };

// 2c. Arrow function — lexical `this`, no own `arguments`, cannot be a constructor
const multiply = (x, y) => x * y;
const square = x => x * x;             // single param, no parens needed
const noop = () => {};                  // returns undefined
const makeObj = () => ({ x: 1 });       // must wrap object literal in parens

// 2d. Default params, rest params, destructured params
function greet(name = "world") { return `Hello, ${name}`; }
function sumAll(...nums) { return nums.reduce((acc, n) => acc + n, 0); }
function printUser({ name, age = 18 }) { return `${name} (${age})`; }

console.log(add(2, 3), subtract(5, 2), multiply(4, 3), square(5));
console.log(greet(), sumAll(1, 2, 3), printUser({ name: "Rio" }));

// ---- 3. `this` binding — the #1 JS interview gotcha ----
const counter = {
  count: 0,
  incRegular: function () {
    // `this` = the object that called it (counter)
    this.count++;
  },
  incArrow: () => {
    // `this` = enclosing lexical scope (NOT counter) — arrow fns don't bind `this`
    this.count++; // buggy on purpose, demonstrates the trap
  },
};
counter.incRegular();
console.log("counter.count after regular fn:", counter.count); // 1

// call / apply / bind — explicitly setting `this`
function whoAmI() { return this.name; }
const person = { name: "Dita" };
console.log(whoAmI.call(person));           // call: args individually
console.log(whoAmI.apply(person));          // apply: args as array
const bound = whoAmI.bind(person);
console.log(bound());                       // bind: returns new fn with `this` fixed

// ---- 4. Closures ----
function makeCounter() {
  let count = 0;              // captured in closure, private to each makeCounter() call
  return {
    increment: () => ++count,
    reset: () => (count = 0),
  };
}
const counterA = makeCounter();
const counterB = makeCounter();
counterA.increment();
counterA.increment();
console.log(counterA.increment(), counterB.increment()); // 3, 1 — independent closures

// ---- 5. Destructuring & spread ----
const [first, second, ...rest] = [1, 2, 3, 4, 5];
const { x: renamedX, ...others } = { x: 1, y: 2, z: 3 };
const merged = { ...{ a: 1 }, ...{ b: 2 } };       // shallow merge
const clonedArr = [...[1, 2, 3]];                   // shallow clone

console.log({ first, second, rest, renamedX, others, merged, clonedArr });

// ---- 6. Equality & truthiness ----
console.log(0 == "0", 0 === "0");     // true, false — always prefer ===
console.log(null == undefined, null === undefined); // true, false
console.log([] == false);              // true (coercion trap — avoid == entirely)
console.log(Boolean(""), Boolean(0), Boolean(NaN), Boolean(null), Boolean(undefined)); // all false

// ---- 7. Optional chaining & nullish coalescing (ES2020) ----
const user = { profile: { name: "Rio" } };
console.log(user?.profile?.name);       // "Rio" — safe access
console.log(user?.settings?.theme);     // undefined, no throw
console.log(user.age ?? 18);            // 18 — only falls back on null/undefined
console.log(0 ?? 18);                   // 0 — unlike `||`, 0 is NOT overridden

// ---- 8. Array methods every interview expects ----
const nums = [1, 2, 3, 4, 5];
console.log(nums.map(n => n * 2));                 // [2,4,6,8,10]
console.log(nums.filter(n => n % 2 === 0));         // [2,4]
console.log(nums.reduce((acc, n) => acc + n, 0));   // 15
console.log(nums.find(n => n > 3));                 // 4
console.log(nums.some(n => n > 4), nums.every(n => n > 0)); // true, true
console.log([...nums].sort((a, b) => b - a));       // [5,4,3,2,1] — copy first! sort() mutates

// ---- 9. Synchronous vs asynchronous ----
console.log("A");
setTimeout(() => console.log("C - macrotask (setTimeout)"), 0);
Promise.resolve().then(() => console.log("B - microtask (Promise)"));
console.log("D");
// Output order: A, D, B, C
// Why: call stack finishes (A, D) -> microtask queue drains (B) -> macrotask queue runs (C)

// async/await is sugar over Promises
async function fetchUserFake() {
  await new Promise(res => setTimeout(res, 10));
  return { id: 1, name: "Rio" };
}
async function run() {
  try {
    const u = await fetchUserFake();
    console.log("fetched:", u);
  } catch (err) {
    console.error("failed:", err);
  }
}
run();

// ---- 10. Prototypes & classes (syntax sugar over prototypes) ----
class Animal {
  #privateSecret = "hidden"; // true private field (ES2022)
  static count = 0;          // static property, shared across instances

  constructor(name) {
    this.name = name;
    Animal.count++;
  }
  speak() { return `${this.name} makes a sound`; }        // on Animal.prototype
  get description() { return `Animal: ${this.name}`; }      // getter
}
class Dog extends Animal {
  speak() { return `${super.speak()}, specifically a bark`; } // super = parent method
}
const dog = new Dog("Rex");
console.log(dog.speak(), dog.description, Animal.count);
console.log(dog instanceof Animal, dog instanceof Dog);

// ---- 11. Promise combinators ----
const ok = (val, ms) => new Promise(res => setTimeout(() => res(val), ms));
const fail = (err, ms) => new Promise((_, rej) => setTimeout(() => rej(err), ms));

Promise.all([ok(1, 10), ok(2, 5)]).then(r => console.log("Promise.all:", r));
Promise.allSettled([ok(1, 10), fail("boom", 5)]).then(r => console.log("Promise.allSettled:", r));
Promise.race([ok("slow", 20), ok("fast", 5)]).then(r => console.log("Promise.race:", r));

// ============================================================
// Coding interview patterns — see README.md section 12
// ============================================================

// ---- Debounce ----
function debounce(fn, delayMs) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delayMs);
  };
}
const logDebounced = debounce((msg) => console.log("debounced:", msg), 30);
logDebounced("a"); logDebounced("b"); logDebounced("c"); // only "c" logs, ~30ms later

// ---- Throttle ----
function throttle(fn, limitMs) {
  let inCooldown = false;
  return (...args) => {
    if (inCooldown) return;
    fn(...args);
    inCooldown = true;
    setTimeout(() => (inCooldown = false), limitMs);
  };
}
const logThrottled = throttle((msg) => console.log("throttled:", msg), 30);
logThrottled("x"); logThrottled("y"); logThrottled("z"); // only "x" logs immediately

// ---- Deep clone ----
function deepClone(value) {
  if (value === null || typeof value !== "object") return value;
  if (Array.isArray(value)) return value.map(deepClone);
  return Object.fromEntries(Object.entries(value).map(([k, v]) => [k, deepClone(v)]));
}
const original = { a: 1, nested: { b: 2 } };
const shallow = { ...original };
const deep = deepClone(original);
shallow.nested.b = 999; // mutates original.nested too — shallow copies share nested refs
console.log("shallow copy leaked mutation:", original.nested.b === 999);
console.log("deep clone stayed independent:", deep.nested.b === 2);

// ---- Curry ----
function curry(fn) {
  return function curried(...args) {
    if (args.length >= fn.length) return fn(...args);
    return (...more) => curried(...args, ...more);
  };
}
const add3 = curry((a, b, c) => a + b + c);
console.log(add3(1)(2)(3), add3(1, 2)(3), add3(1, 2, 3)); // 6 6 6

// ---- Flatten ----
function flatten(arr) {
  return arr.reduce((flat, item) =>
    flat.concat(Array.isArray(item) ? flatten(item) : item), []);
}
console.log(flatten([1, [2, [3, 4], 5]])); // [1,2,3,4,5]

// ---- Memoize ----
function memoize(fn) {
  const cache = new Map();
  return (...args) => {
    const key = JSON.stringify(args);
    if (cache.has(key)) return cache.get(key);
    const result = fn(...args);
    cache.set(key, result);
    return result;
  };
}
let calls = 0;
const slowSquare = memoize((n) => { calls++; return n * n; });
slowSquare(4); slowSquare(4); slowSquare(4);
console.log("memoized calls (expect 1):", calls);

module.exports = {}; // keep node quiet if required elsewhere
