# Front-End Refresher — TS + React + JS (company interview prep)

Personal refresher notes. Priority order reflects how this stack actually gets
used day to day: **TypeScript > React > plain JavaScript**. Goal: cover basic
→ medium front-end engineering fundamentals well enough to be interview-ready.
Runnable companion code (`npm install` once, then):

- [examples/typescript-fundamentals.ts](examples/typescript-fundamentals.ts) — `npm run ts` (runs), `npm run typecheck` (type-checks all examples)
- [examples/react-typescript-patterns.tsx](examples/react-typescript-patterns.tsx) — illustrative, type-checked but needs a bundler to actually render
- [examples/js-syntax.js](examples/js-syntax.js) — `npm run js`
- [examples/react-hooks.jsx](examples/react-hooks.jsx) — illustrative, needs a bundler (Vite/CRA) to actually run

---

## Table of Contents

1. [JavaScript Syntax & Variations](#1-javascript-syntax--variations)
2. [TypeScript Deep Dive](#2-typescript-deep-dive)
3. [How Functions Work in JS (and what's different from backend)](#3-how-functions-work-in-js-and-whats-different-from-backend)
4. [React Fundamentals](#4-react-fundamentals)
5. [React Hooks Deep Dive](#5-react-hooks-deep-dive)
6. [React + TypeScript Patterns](#6-react--typescript-patterns)
7. [State Management & Data Flow](#7-state-management--data-flow)
8. [Rendering, the DOM, and Performance](#8-rendering-the-dom-and-performance)
9. [Front-End vs Backend — Key Mental Model Shifts](#9-front-end-vs-backend--key-mental-model-shifts)
10. [CSS / Layout Essentials](#10-css--layout-essentials)
11. [Browser & Web Platform Basics](#11-browser--web-platform-basics)
12. [Accessibility (a11y) Essentials](#12-accessibility-a11y-essentials)
13. [JavaScript Coding Interview Patterns](#13-javascript-coding-interview-patterns)
14. [Interview Q&A — Basic to Medium](#14-interview-qa--basic-to-medium)
15. [Quick Reference Tables](#15-quick-reference-tables)

---

## 1. JavaScript Syntax & Variations

### Variable declarations

> **ELI5:** `var` is an old box that leaks out of the room (function) it was made in. `let` is a labeled box that only exists in the `{ }` room it was made in, and you can put a different thing in it later. `const` is the same room-scoped box, but taped shut — you can't swap out what's in the box, though if the box holds another box (an object), you can still rummage around inside that inner box.

| Keyword | Scope | Reassignable | Hoisting behavior |
|---|---|---|---|
| `var` | function | yes | hoisted + initialized as `undefined` |
| `let` | block | yes | hoisted, but in "temporal dead zone" (TDZ) until the line runs |
| `const` | block | no (binding only, object contents still mutable) | same as `let` |

Rule of thumb: **default to `const`**, use `let` only when you must reassign, avoid `var` entirely (legacy code will still have it — recognize it, don't write it).

### Function variations

> **ELI5:** All four lines below do the same job — add two numbers — they're just different ways to write "a machine that takes inputs and gives an output." A `function` declaration is a machine you can use before you even build it (JS builds it first). A `const x = function(){}` or arrow version only exists once JS reaches that line.

```js
function add(a, b) { return a + b; }      // declaration — hoisted, callable before definition
const sub = function (a, b) { return a - b; }; // expression — not hoisted
const mul = (a, b) => a * b;              // arrow — implicit return, lexical `this`
const sq  = x => x * x;                    // single param, no parens
const obj = () => ({ x: 1 });              // must wrap object literal in ()
```

Params:
```js
function greet(name = "world") {}          // default param
function sum(...nums) {}                   // rest param -> array
function print({ name, age = 18 }) {}      // destructured param
```

### Destructuring & spread
```js
const [a, b, ...rest] = [1, 2, 3, 4];
const { x, ...others } = { x: 1, y: 2 };
const merged = { ...obj1, ...obj2 };       // shallow merge, later keys win
const clone = [...arr];                     // shallow clone
```

### Equality

> **ELI5:** `===` asks "are you the exact same type AND value?" `==` first tries to convert one side to match the other, which causes weird surprises (a string "0" and number 0 suddenly look "equal"). `??` means "use this value, but if it's completely empty (null/undefined), fall back to a default." `?.` means "try to go inside this object, but if it doesn't exist, just give me nothing instead of crashing."

- `==` coerces types (`0 == "0"` → true) — **avoid**.
- `===` strict, no coercion — **always prefer**.
- `??` (nullish coalescing) only falls back on `null`/`undefined`, unlike `||` which falls back on any falsy value (`0`, `""`, `NaN`).
- `?.` (optional chaining) short-circuits to `undefined` instead of throwing on `null`/`undefined` access.

### Template literals
```js
const name = "Rio";
`Hello, ${name}!`;         // interpolation
`multi
line`;                     // real newlines preserved
```

### Array methods (interview staples)
`map`, `filter`, `reduce`, `find`, `some`, `every`, `sort` (mutates! clone with `[...arr]` first), `flat`, `flatMap`.

### Modules (ESM vs CommonJS)

> **ELI5:** A module is just a file that can export things it wants to share, and import things from other files. `import`/`export` (ESM, what the browser and modern bundlers use) is checked at build time and is more static; `require`/`module.exports` (CommonJS, classic Node) is resolved at runtime and more dynamic.

```js
// ESM (Vite/webpack/browsers with type="module", modern Node with "type": "module")
export function add(a, b) { return a + b; }
export default function App() {}
import App, { add } from "./App.js";

// CommonJS (classic Node/Jest-without-babel)
function add(a, b) { return a + b; }
module.exports = { add };
const { add } = require("./math");
```
Key differences: ESM imports are hoisted and statically analyzable (enables tree-shaking — unused exports get dropped from the bundle); CommonJS `require()` can be called conditionally/anywhere since it's just a function call.

### Classes (sugar over prototypes)

> **ELI5:** A `class` is a blueprint for making objects that all share the same behaviors. `extends` means "start from this other blueprint, then add/override some parts." Under the hood JS doesn't really have classes — it's all objects pointing to other objects (`prototype`) for shared methods — `class` is just a nicer way to write that.

```js
class Animal {
  #secret = "hidden";       // private field (ES2022)
  static count = 0;         // shared across all instances
  constructor(name) { this.name = name; }
  speak() { return `${this.name} makes a sound`; }
  get label() { return `Animal: ${this.name}`; }
}
class Dog extends Animal {
  speak() { return `${super.speak()}, a bark`; }
}
```

---

## 2. TypeScript Deep Dive

Runnable/type-checked companion: [examples/typescript-fundamentals.ts](examples/typescript-fundamentals.ts) — `npm run ts` to execute, `npm run typecheck` to see it type-check clean.

> **ELI5:** TypeScript is JavaScript with labels on everything ("this is a number," "this is a User"). It only exists at build time — the labels are checked and then completely erased before the code runs, so at runtime it's just plain JS. The whole point is catching mistakes (like calling `.toUpperCase()` on a number) while you're typing, instead of finding out when a user hits it in production.

### Basic types & inference
```ts
let id: number = 1;
let tags: string[] = ["a", "b"];
let pair: [string, number] = ["x", 1];        // tuple — fixed length, fixed type per slot
let inferred = "hello";                        // TS infers `string` — don't over-annotate locals
```
Rule of thumb: let TS **infer** local variable types; **annotate** function parameters and return types (the boundaries), since inference can't see into a function you haven't written yet.

### Enums vs literal union types

> **ELI5:** An `enum` is TypeScript's own named-constants feature — but it actually generates a real JS object at runtime. A literal union (`"pending" | "active"`) is just a type-level list of allowed strings, with zero runtime footprint. Most modern TS codebases prefer literal unions — they're lighter, and any plain JSON/API value can satisfy them directly.

```ts
enum StatusEnum { Pending, Active, Closed }     // generates a real runtime object + reverse mapping
type Status = "pending" | "active" | "closed";  // no runtime cost, just a compile-time check
function setStatus(s: Status) { /* ... */ }
setStatus("active");   // OK
// setStatus("done");  // compile error — "done" isn't in the union
```

### Interfaces vs type aliases

> **ELI5:** Both describe "the shape of an object." `interface` is like a form that can be extended later or reopened elsewhere in the codebase (declaration merging) — good for public/library shapes. `type` is more like a nickname for any type, including unions and primitives, which `interface` can't express. For a plain object shape, either is fine; pick one convention and stay consistent.

```ts
interface User { id: number; name: string; age?: number; readonly createdAt: string; }
type UserType = { id: number; name: string };

interface Employee extends User { role: string; }   // interface extension
type Manager = UserType & { reports: number };       // type intersection ("AND")
type Id = string | number;                            // union ("OR") — only `type` can do this directly
```
- `?` marks a property **optional**.
- `readonly` blocks reassignment after the object is created (compile-time only — doesn't freeze it at runtime like `Object.freeze`).

### Functions: params, overloads, generics preview
```ts
function add(a: number, b: number): number { return a + b; }

// Overloads — same function name, different call signatures resolved at compile time
function toArray(value: string): string[];
function toArray(value: number): number[];
function toArray(value: string | number): string[] | number[] {
  return typeof value === "string" ? value.split("") : [value];
}
```

### Discriminated unions & exhaustiveness checking

> **ELI5:** Give each variant of a union a shared "tag" field (like `kind: "circle"`), and TypeScript will automatically figure out exactly which shape you're holding inside an `if`/`switch` — no manual casting needed. The `never` trick at the end of a `switch` is a tripwire: if someone adds a new variant later and forgets to handle it, the compiler errors instead of silently doing nothing at runtime.

```ts
type Shape =
  | { kind: "circle"; radius: number }
  | { kind: "square"; side: number };

function area(shape: Shape): number {
  switch (shape.kind) {
    case "circle": return Math.PI * shape.radius ** 2; // narrowed to circle here
    case "square": return shape.side ** 2;               // narrowed to square here
    default: { const _exhaustive: never = shape; return _exhaustive; } // catches unhandled variants
  }
}
```

### Type narrowing
```ts
function describe(value: string | number | Date) {
  if (typeof value === "string") return value.toUpperCase(); // typeof guard
  if (value instanceof Date) return value.toISOString();      // instanceof guard
}
interface Cat { meow(): void }
interface Bird { fly(): void }
function isCat(pet: Cat | Bird): pet is Cat { return "meow" in pet; } // custom type guard ("in" + `is`)
```

### Generics

> **ELI5:** A generic is a type-level variable — "I don't know the exact type yet, call it `T`, and I'll figure it out from whatever you actually pass in." It lets one function/interface work correctly for many types without losing type safety (unlike `any`, which just gives up on checking).

```ts
function identity<T>(value: T): T { return value; }              // T inferred from the argument
interface ApiResponse<T> { data: T; error?: string; }
function logLength<T extends { length: number }>(v: T) { return v.length; } // constrained generic
interface Box<T = string> { value: T; }                            // default generic param
```

### Utility types (built into TS, no import needed)
```ts
interface Article { id: number; title: string; body: string; published: boolean; }

type PartialArticle = Partial<Article>;         // every field optional — e.g. a PATCH payload
type Preview = Pick<Article, "id" | "title">;    // only these keys
type NoBody = Omit<Article, "body">;             // all except this key
type Frozen = Readonly<Article>;                 // every field readonly
type ArticleMap = Record<number, Article>;       // { [id: number]: Article }
type Excluded = Exclude<"a" | "b" | "c", "c">;   // "a" | "b"
type Only = Extract<string | number, string>;    // string
type Fetched = Awaited<ReturnType<typeof fetch>>; // unwraps a function's Promise<T> return down to T
```

### `keyof`, `typeof`, mapped & conditional types

> **ELI5:** `keyof` asks an object type "what are your property names?" and gives you a union of them. A mapped type says "build a new type by looping over every key of an existing one." A conditional type is an `if` statement for types — "if T fits this shape, resolve to X, otherwise Y."

```ts
const config = { retries: 3, timeoutMs: 5000 };
type ConfigKey = keyof typeof config;                  // "retries" | "timeoutMs"

type Optional<T> = { [K in keyof T]?: T[K] };           // mapped type — same shape, all optional
type IsString<T> = T extends string ? "yes" : "no";     // conditional type
type UnwrapArray<T> = T extends (infer U)[] ? U : T;    // `infer` pulls out a nested type — array's element type
```

### `any` vs `unknown` vs `never` vs `void`
- **`any`** — disables type checking entirely for that value. Compiles, but can blow up at runtime with zero warning. Avoid; it's an escape hatch, not a default.
- **`unknown`** — the type-safe version of `any`: you must **narrow** it (via `typeof`/`instanceof`/a guard) before doing anything with it. Prefer this over `any` for genuinely-unknown values (e.g. `catch (e: unknown)`, JSON parsing).
- **`never`** — a function that never returns (always throws or loops forever), or an impossible/unreachable type (useful for exhaustiveness checks above).
- **`void`** — a function's return value that should be ignored (usually just side effects, like `console.log`).

### Classes: access modifiers & abstract
```ts
abstract class Animal {
  protected readonly name: string;   // visible in this class + subclasses only, can't reassign
  private secret = "hidden";         // visible only inside Animal itself
  constructor(name: string) { this.name = name; }
  abstract speak(): string;           // subclasses MUST implement this
}
class Dog extends Animal { speak() { return "Woof"; } }
// new Animal("x"); // compile error — can't instantiate an abstract class directly
```

### Assertions & `satisfies`
```ts
const input = document.querySelector("input") as HTMLInputElement | null; // `as` — "trust me" type assertion
const val = input!.value; // `!` — non-null assertion, only when you're SURE it isn't null (otherwise it'll throw at runtime)

// `satisfies` checks a value matches a type WITHOUT widening it (unlike `: Type`)
const palette = { primary: "#000", secondary: "#fff" } satisfies Record<string, string>;
// palette.primary keeps the literal type "#000" here — a `: Record<string,string>` annotation
// would have widened it to plain `string`, losing that precision
```

### Type-only imports
```ts
import type { User } from "./types"; // fully erased at compile time — zero runtime cost, avoids circular-import issues
```

---

## 3. How Functions Work in JS (and what's different from backend)

### `this` binding — the biggest JS-specific trap

> **ELI5:** `this` is like the word "I" — its meaning depends on who's talking, not who wrote the sentence. If `obj.regular()` is called, "I" = `obj`. Arrow functions don't have their own "I" — they just borrow whatever "I" means in the place they were written, permanently.

Unlike most backend languages where `this`/`self` is bound to the instance at
definition time, in JS `this` is determined by **how a function is called**,
not where it's defined — except for arrow functions.

```js
const obj = {
  count: 0,
  regular() { this.count++; },   // `this` = obj (whoever calls obj.regular())
  arrow: () => { this.count++; } // `this` = enclosing lexical scope, NOT obj
};
```

- `fn.call(thisArg, a, b)` — invoke with explicit `this`, args individually.
- `fn.apply(thisArg, [a, b])` — same but args as array.
- `fn.bind(thisArg)` — returns a new function permanently bound to `thisArg`.

### Closures

> **ELI5:** Imagine a function is handed a backpack of variables when it's created. Even after it leaves the room (the outer function returns), it keeps carrying that backpack around and can still use what's inside — nobody else can peek into that specific backpack.

A function "remembers" the variables from its enclosing scope even after that
scope has returned. This is how private state and factories work in JS
without classes:

```js
function makeCounter() {
  let count = 0;                  // private to this closure
  return { inc: () => ++count };
}
```

### Hoisting
Function **declarations** are hoisted with their full body (usable before
their line). Function **expressions** and arrow functions are only hoisted
as `undefined`/TDZ bindings — calling them before assignment throws.

### Single-threaded, event-loop concurrency

> **ELI5:** Picture one waiter (the JS engine) in a restaurant with one line of customers (the call stack). The waiter never stands around waiting for the kitchen — they take an order, hand it to the kitchen (the browser/Node), and immediately go help the next customer. When the kitchen finishes a dish, it's added to a to-do list the waiter checks between customers. That's the event loop: the waiter isn't literally multitasking, they're just never idle waiting on slow things.

There is no multithreading in normal JS (no locks/mutexes needed for JS code
itself). Concurrency is achieved via the **event loop**:

1. Call stack runs synchronous code to completion.
2. **Microtask queue** drains fully (Promises, `queueMicrotask`).
3. One task from the **macrotask queue** runs (`setTimeout`, I/O, UI events).
4. Repeat.

```js
console.log("A");
setTimeout(() => console.log("C"), 0);
Promise.resolve().then(() => console.log("B"));
console.log("D");
// A, D, B, C
```

This is the #1 conceptual difference from backend request-per-thread models
(e.g. a typical Java servlet thread) — JS never blocks the main thread
waiting on I/O; everything I/O-bound is callback/Promise-based.

### Async/await

> **ELI5:** `await` means "pause here and let other stuff happen until this is ready, then continue right where you left off" — it just makes async code *read* like normal top-to-bottom code, even though it's still non-blocking underneath.

Syntax sugar over Promises. `await` pauses the *async function* (not the whole
program) until the Promise settles.
```js
async function run() {
  try {
    const res = await fetch("/api");
    const data = await res.json();
  } catch (err) { /* rejections become throws */ }
}
```

### Promise combinators

> **ELI5:** These answer "I fired off several async tasks — how do I know when I'm done, and do I care if some fail?" `all` = "wait for everyone, but bail the instant anyone fails." `allSettled` = "wait for everyone, and just tell me who succeeded and who failed." `race`/`any` = "just give me the first one to finish."

```js
await Promise.all([p1, p2, p3]);        // rejects as soon as ANY promise rejects
await Promise.allSettled([p1, p2, p3]); // never rejects; returns [{status, value|reason}, ...] per promise
await Promise.race([p1, p2, p3]);       // settles as soon as the FIRST promise settles (resolve or reject)
await Promise.any([p1, p2, p3]);        // resolves as soon as the FIRST promise resolves; rejects only if ALL reject
```
`Promise.all` is the one to reach for when fetching several independent resources in parallel and you need all of them to render (e.g. `Promise.all([fetchUser(), fetchPosts()])`). Use `allSettled` when partial failure is acceptable (e.g. a dashboard with independent widgets).

---

## 4. React Fundamentals

### What React actually is

> **ELI5:** Instead of manually telling the browser "now update this text, now hide that button," you just describe what the screen *should* look like for the current data, and React figures out what actually needs to change on screen — like giving an artist a finished painting instead of a list of brush strokes.

A library for building UI as a function of state: `UI = f(state)`. You
describe *what* the UI should look like for a given state; React figures out
the minimal DOM mutations (via the **Virtual DOM diff / reconciliation**).

### JSX

> **ELI5:** JSX lets you write what looks like HTML directly inside your JS, but it's really just a shorthand — a compiler turns those tags into plain JS function calls before the browser ever sees them.

JSX is not HTML — it's syntactic sugar for `React.createElement(type, props, children)`.
```jsx
<div className="box">{count}</div>
// compiles to:
React.createElement("div", { className: "box" }, count)
```
Key JSX rules:
- `className` not `class`, `htmlFor` not `for` (reserved JS words).
- Any JS expression inside `{ }`; statements (if/for) are NOT allowed inline — use ternaries, `&&`, or extract logic above the return.
- Must return a **single root** (or use a Fragment `<>...</>`).
- Every list item needs a stable, unique `key` prop (not array index if the list can reorder/filter).

### Components: function vs class (legacy)
Modern React (company's current codebases very likely too) uses **function
components + Hooks**. You should still recognize class components since
older code exists:

| Concept | Function component | Class component (legacy) |
|---|---|---|
| State | `useState` | `this.state` / `this.setState` |
| Mount effect | `useEffect(() => {...}, [])` | `componentDidMount` |
| Update effect | `useEffect(() => {...}, [deps])` | `componentDidUpdate` |
| Cleanup | return fn from `useEffect` | `componentWillUnmount` |
| Ref | `useRef` | `React.createRef()` / callback refs |

### Props vs State
- **Props**: read-only, passed from parent → child. Component must never mutate its own props.
- **State**: local, mutable (via setter), owned by the component; triggers re-render on change.

### One-way data flow

> **ELI5:** Think of it like a company org chart: instructions/data flow down from manager to employee (props), and if an employee needs to tell the manager something, they can't edit the manager's notes directly — they call a phone number the manager gave them (a callback prop) and let the manager decide what to do with that information.

Data flows **down** via props. To send data "up," a parent passes a callback
function down as a prop, and the child calls it — this is the standard
pattern (no two-way binding like Angular/Vue by default).

```jsx
function Parent() {
  const [value, setValue] = useState("");
  return <Child value={value} onChange={setValue} />;
}
function Child({ value, onChange }) {
  return <input value={value} onChange={e => onChange(e.target.value)} />;
}
```

---

## 5. React Hooks Deep Dive

See runnable-ish examples in [examples/react-hooks.jsx](examples/react-hooks.jsx).

### `useState`

> **ELI5:** A sticky note attached to your component that React remembers between renders. `setCount` doesn't just change the note — it tells React "redraw this component with the new note value."

```jsx
const [count, setCount] = useState(0);
setCount(c => c + 1);   // functional update — always safe when new state depends on old
```
Rules of Hooks: only call at the top level of a component (never in loops/conditions/nested functions), only from React functions.

### `useEffect`

> **ELI5:** "After you finish drawing the screen, go do this side task" (fetch data, set a timer, subscribe to something). The dependency array is the answer to "when should I redo this task?" — an empty list means "just once, right after the first drawing."

Runs *after* the browser paints. Dependency array controls when it re-runs:
- `useEffect(fn)` — no array → every render.
- `useEffect(fn, [])` — empty array → once, like mount.
- `useEffect(fn, [dep])` — re-runs when `dep` changes (reference equality!).
- Return a cleanup function — runs before the next effect and on unmount. Critical for subscriptions, timers, aborting fetches.

### `useRef`

> **ELI5:** A sticky note React does NOT watch. You can scribble on it and read it anytime, but changing it never makes React redraw anything — useful for things like "grab this actual input element" or "remember a value quietly in the background."

Mutable `.current` box that persists across renders **without** causing a re-render when changed. Two main uses: (1) direct DOM node access, (2) "instance variable" that shouldn't trigger a render.

### `useMemo` / `useCallback`

> **ELI5:** Both mean "don't redo this work unless the ingredients actually changed." `useMemo` caches a *result* (like a sorted list), `useCallback` caches the *function itself* so a child component doesn't think "oh, a brand new function, I must have new props" every single render.

Both are performance escape hatches, not correctness tools:
- `useMemo(fn, deps)` memoizes a **computed value**.
- `useCallback(fn, deps)` memoizes a **function reference** — matters when passing callbacks to `React.memo`-wrapped children, otherwise a new function each render defeats the memoization.

### `useReducer`

> **ELI5:** Instead of directly editing state yourself, you send a labeled "request" (an action, like `{ type: "increment" }`) to a single function (the reducer) that decides exactly how state should change — handy once state updates get complicated enough that scattered `setState` calls get hard to follow.

Prefer over `useState` when state transitions are complex/interrelated (state machine style), analogous to Redux reducers but local.

### `useContext`

> **ELI5:** Instead of passing a value down through every single component in between (prop drilling — like relaying a message person-to-person down a long line), `useContext` lets any descendant component tune in directly to a shared broadcast, no matter how deep it is.

Reads the nearest matching `<Context.Provider value={...}>` above it in the tree — solves prop drilling. Overusing Context for frequently-changing state causes broad re-renders (every consumer re-renders on any value change) — for that, prefer a state library or split contexts.

### Custom hooks
Any function starting with `use` that calls other hooks — the mechanism for extracting/reusing stateful logic (e.g. `useDebouncedValue`, `useFetch`).

### `useLayoutEffect` vs `useEffect`

> **ELI5:** `useEffect` says "do this after the screen is already painted, no rush." `useLayoutEffect` says "do this BEFORE the browser paints, even if it delays showing anything" — use it only when you must measure/mutate the DOM synchronously (e.g. reading an element's size and repositioning it) to avoid a visible flicker.

`useLayoutEffect` runs synchronously after DOM mutations but before the browser paints; it blocks visual updates, so overusing it can hurt performance. Default to `useEffect`; reach for `useLayoutEffect` only for layout measurement/flicker fixes.

---

## 6. React + TypeScript Patterns

Type-checked companion: [examples/react-typescript-patterns.tsx](examples/react-typescript-patterns.tsx) — `npm run typecheck`.

> **ELI5:** This is where the two priorities meet — the React patterns from section 5, but every prop, piece of state, and event handler gets a type so the compiler catches mismatches (passing a `number` where a component expects a `string`, forgetting a required prop, etc.) before you ever open the browser.

### Typing props
```tsx
interface ButtonProps {
  label: string;
  onClick?: () => void;                    // optional callback
  variant?: "primary" | "secondary";        // literal union beats a loose `string`
  children?: React.ReactNode;               // covers strings, elements, arrays, fragments, null
}
function Button({ label, onClick, variant = "primary", children }: ButtonProps) {
  return <button className={variant} onClick={onClick}>{label}{children}</button>;
}
```

### Typing `useState`, `useRef`, `useReducer`

> **ELI5:** `useState(0)` already knows it holds a number from the initial value — you only need to spell out `useState<T>()` explicitly when the initial value doesn't tell the whole story (like starting at `null` but later holding a `User`).

```tsx
const [count, setCount] = useState<number>(0);                    // usually inferred, shown explicit here
const [user, setUser] = useState<{ name: string } | null>(null);   // union needed — null doesn't imply the full type
const inputRef = useRef<HTMLInputElement>(null);                   // DOM ref — TS now knows .current is an <input> or null

interface CounterState { count: number; step: number; }
type CounterAction = { type: "increment" } | { type: "setStep"; payload: number }; // discriminated union
function reducer(state: CounterState, action: CounterAction): CounterState {
  switch (action.type) {
    case "increment": return { ...state, count: state.count + state.step };
    case "setStep": return { ...state, step: action.payload }; // `payload` only exists on this branch — TS enforces it
  }
}
```

### Typing event handlers

> **ELI5:** The browser's event objects are generic (`Event`) — React's typed versions tell TypeScript exactly which element and event kind you're dealing with, so `e.target.value` type-checks instead of being `any`.

```tsx
function handleChange(e: React.ChangeEvent<HTMLInputElement>) { console.log(e.target.value); }
function handleClick(e: React.MouseEvent<HTMLButtonElement>) { e.preventDefault(); }
```

### Generic components

> **ELI5:** Same idea as a generic function, applied to a component — a `<List>` that works for a list of users, products, or anything else, while still knowing the exact shape of each item so `renderItem`/`keyExtractor` are fully type-checked.

```tsx
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
  keyExtractor: (item: T) => string | number;
}
function List<T>({ items, renderItem, keyExtractor }: ListProps<T>) {
  return <ul>{items.map(item => <li key={keyExtractor(item)}>{renderItem(item)}</li>)}</ul>;
}
// <List items={users} keyExtractor={u => u.id} renderItem={u => u.name} /> — T inferred from `users`
```

### `forwardRef` with typed props + ref
```tsx
interface FancyInputProps { placeholder?: string; }
const FancyInput = React.forwardRef<HTMLInputElement, FancyInputProps>(function FancyInput(props, ref) {
  return <input ref={ref} placeholder={props.placeholder} />;
});
```

### Typing a custom hook's return value
```tsx
function useToggle(initial = false): [boolean, () => void] { // tuple return, order matters (like useState)
  const [value, setValue] = useState(initial);
  return [value, () => setValue(v => !v)];
}
```

---

## 7. State Management & Data Flow

> **ELI5:** "State management" is just: where does the app remember things, and who's allowed to know about it? Keep it as local/private as possible (a sticky note only one component sees) and only make it more "public" (lifted up, Context, or a global library) once more than one part of the app genuinely needs to see the same note.

- **Local state** (`useState`/`useReducer`) — component-owned, simplest, default choice.
- **Lifted state** — move state up to the nearest common ancestor when siblings need to share it.
- **Context** — for state, several components deep, that changes rarely (theme, auth user, locale).
- **External libraries** (Redux Toolkit, Zustand, Recoil, Jotai) — for large apps with complex cross-cutting state, dev tools/time-travel debugging, or state that outlives component trees.
- **Server state** (React Query / SWR / RTK Query) — a distinct category from client state: caching, revalidation, dedupe of async server data. Don't manage this manually with `useState` + `useEffect` in nontrivial apps.

Controlled vs uncontrolled inputs:
```jsx
// controlled — React state is the single source of truth
<input value={value} onChange={e => setValue(e.target.value)} />
// uncontrolled — DOM holds the value, read via ref when needed
<input ref={ref} defaultValue="init" />
```

---

## 8. Rendering, the DOM, and Performance

### Virtual DOM & reconciliation

> **ELI5:** Rather than repainting the whole picture every time one pixel changes, React sketches the new picture on scratch paper first, compares it to the last sketch, and only touches up the real canvas (the actual browser DOM) where something actually differs.

On every state change, React builds a new virtual tree, diffs it against the
previous one, and applies the minimal set of real DOM mutations. Elements are
matched by **type + key** — this is why stable `key`s matter for list items
(wrong keys cause React to unmount/remount or misapply state to the wrong
node).

### Why components re-render
A component re-renders when: its own state changes, its parent re-renders
(and it's not memoized), or a context it consumes changes.

### `React.memo`

> **ELI5:** "Don't bother re-drawing this component if nothing it actually looks at has changed." It's a shortcut React takes only if you ask for it — by default React re-renders children whenever the parent re-renders, memoized or not, unless props are unchanged.

Wraps a function component to skip re-rendering if its props are shallow-equal
to the previous render. Only useful when the component is expensive AND
receives stable props (pair with `useMemo`/`useCallback` for object/function props).

```jsx
const Row = React.memo(function Row({ item, onSelect }) {
  return <li onClick={() => onSelect(item.id)}>{item.name}</li>;
});
```

### `forwardRef`

> **ELI5:** Normally a parent can't reach "inside" a child component to grab its DOM node — `forwardRef` is the child explicitly saying "here, parent, you can hold onto this specific piece of me."

```jsx
const FancyInput = React.forwardRef(function FancyInput(props, ref) {
  return <input ref={ref} {...props} className="fancy" />;
});
// parent: const inputRef = useRef(); <FancyInput ref={inputRef} />
```

### Error boundaries

> **ELI5:** A safety net around part of your UI — if anything inside throws while rendering, the boundary catches it and shows a fallback instead of the whole app going blank.

```jsx
class ErrorBoundary extends React.Component {
  state = { hasError: false };
  static getDerivedStateFromError() { return { hasError: true }; }
  componentDidCatch(error, info) { console.error(error, info); }
  render() {
    return this.state.hasError ? <p>Something went wrong.</p> : this.props.children;
  }
}
// usage: <ErrorBoundary><Widget /></ErrorBoundary>
```
Must be a class component (no Hook equivalent exists yet); only catches errors during **render**, not in event handlers or async code (catch those with a plain try/catch).

### Code splitting: `React.lazy` + `Suspense`

> **ELI5:** Instead of shipping the whole app as one giant file, `lazy` says "only download this component's code when it's actually needed," and `Suspense` says "while it's downloading, show this loading UI instead."

```jsx
const Settings = React.lazy(() => import("./Settings"));
<Suspense fallback={<Spinner />}>
  <Settings />
</Suspense>
```

### Common performance mistakes
- Creating new objects/arrays/functions inline as props every render, defeating `React.memo`.
- Using array index as `key` in a reorderable list.
- Putting all state at the top of a large tree instead of colocating it near where it's used (causes wide re-render blast radius).
- Missing `useEffect` dependencies (stale closures) or over-including them (extra re-runs).

### Core Web Vitals (how "fast" gets measured)
- **LCP** (Largest Contentful Paint) — time until the biggest visible element renders; measures perceived load speed.
- **INP** (Interaction to Next Paint, replaced FID) — responsiveness to user interaction (clicks/taps/keys).
- **CLS** (Cumulative Layout Shift) — how much visible content unexpectedly jumps around while loading.
Good talking point in interviews: tie these back to concrete causes (large unoptimized images → poor LCP, heavy main-thread JS → poor INP, images/ads without reserved dimensions → poor CLS).

### Keys, list, conditional rendering
```jsx
{items.map(item => <Row key={item.id} {...item} />)}
{isLoading ? <Spinner /> : <Content />}
{error && <ErrorBanner message={error} />}
```

---

## 9. Front-End vs Backend — Key Mental Model Shifts

> **ELI5:** A backend handles one request, answers it, and forgets about it — like a customer service rep who takes one call at a time. A front-end is more like running a whole live TV show that never really "ends" while the tab is open — it has to keep reacting to the user, remember things across the whole session, and never freeze up no matter what's happening in the background.

| Aspect | Backend (typical) | Front-end (JS/React) |
|---|---|---|
| Execution model | multi-threaded / per-request process | single-threaded event loop, one thread does everything including rendering |
| State lifetime | request-scoped, often stateless | long-lived in-memory UI state across a session |
| I/O | blocking or thread-pooled | always async (Promises/callbacks), never blocks the main thread |
| "Rendering" | server returns a response body once | UI re-renders continuously as state changes, driven by a diffing algorithm |
| Source of truth | database / persistent store | often a mix of local component state + server cache + URL — many truths to reconcile |
| Errors | exceptions/HTTP status codes, logged server-side | must be caught and rendered to the user gracefully (error boundaries), silent failures = blank UI |
| Security surface | server fully trusted | all client code, tokens, and requests are visible/tamperable by the end user — never trust client-side validation alone |
| Performance target | throughput / latency per request | perceived responsiveness — time-to-interactive, avoiding layout thrash, minimizing re-renders |
| Deployment unit | server process/binary | bundle shipped to an untrusted, wildly heterogeneous runtime (browsers/devices) |

---

## 10. CSS / Layout Essentials

### Box model

> **ELI5:** Every element is like a picture frame: the `content` is the photo, `padding` is the mat around the photo, `border` is the actual frame, and `margin` is the gap between this frame and the next one on the wall.

`content` → `padding` → `border` → `margin`. `box-sizing: border-box` makes `width`/`height` include padding+border (nearly always what you want; many resets set this globally).

### Flexbox (1-D layout)

> **ELI5:** Flexbox is for lining things up in a single row or column and deciding how leftover space gets shared out — like arranging books on one shelf.

```css
.row { display: flex; justify-content: space-between; align-items: center; gap: 8px; }
.item { flex: 1; }              /* grow/shrink/basis shorthand */
```
- `justify-content` = main axis, `align-items` = cross axis.
- `flex-direction: column` swaps which axis is "main."

### Grid (2-D layout)

> **ELI5:** Grid is for laying things out in rows AND columns at once — like a full bookshelf with multiple shelves and columns, not just one row of books.

```css
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
```

### Positioning
- `static` (default) → `relative` (offsets from normal position, becomes a positioning context for children) → `absolute` (removed from flow, positioned relative to nearest positioned ancestor) → `fixed` (relative to viewport) → `sticky` (hybrid, sticks within its scroll container).

### Specificity (low → high)
`element` < `.class` < `#id` < inline `style=""` < `!important` (avoid).

### Responsive
```css
@media (max-width: 768px) { .row { flex-direction: column; } }
```
Mobile-first: write base styles for small screens, add `min-width` media queries upward.

---

## 11. Browser & Web Platform Basics

> **ELI5:** The DOM is the browser's live model of the page (what React quietly edits for you). CORS is a bouncer at a website's door checking "is this other website on the guest list to read my stuff?" localStorage/cookies are just the browser's own little filing cabinets for remembering things between visits.

- **DOM**: tree representation of the page; `document.querySelector`, `addEventListener` are the imperative APIs React abstracts over.
- **Event loop / rendering pipeline**: JS runs → style/layout recalculated → paint → composite. Long JS tasks block rendering (jank).
- **HTTP basics**: verbs (GET/POST/PUT/PATCH/DELETE), status codes (2xx success, 3xx redirect, 4xx client error, 5xx server error), CORS (server must explicitly allow cross-origin requests via headers).
- **Storage**: `localStorage` (persists, ~5-10MB, synchronous, string-only), `sessionStorage` (tab-scoped), cookies (sent with every matching request, can be `httpOnly`/secure, used for auth tokens).
- **Same-origin policy & CORS**: browsers block cross-origin reads by default; the server opts in via `Access-Control-Allow-Origin` etc.

### Event bubbling, capturing & delegation

> **ELI5:** Click a button inside a `div` inside a `body`, and the click doesn't just happen on the button — it "bubbles" up through every ancestor, like a shout traveling from a bathroom, up the stairs, out to the whole house. Capturing is the same trip in reverse (top-down, on the way in, before bubbling). Delegation is the trick of listening on a big ancestor instead of every single small child, then figuring out which child was actually clicked.

```js
document.getElementById("list").addEventListener("click", (e) => {
  if (e.target.matches("li")) console.log("clicked item:", e.target.textContent);
});
// one listener handles clicks on ALL <li>, even ones added later — no need to
// attach/re-attach a listener per item
e.stopPropagation(); // stop the bubble from going further up
e.preventDefault();  // stop the browser's default action (e.g. link navigation, form submit)
```
Delegation is why React's synthetic event system attaches one listener at the root and dispatches from there instead of one native listener per element — cheaper for large lists.

---

## 12. Accessibility (a11y) Essentials

> **ELI5:** Accessibility means the app also works for people who can't use a mouse, can't see the screen well, or use a screen reader. Practically: use real HTML elements for what they're for (a `<button>` not a `<div onClick>`), because browsers and screen readers already know how buttons should behave.

- **Semantic HTML first**: `<button>`, `<nav>`, `<header>`, `<label>` carry built-in keyboard/screen-reader behavior a `<div>` doesn't. Prefer them over generic elements + ARIA whenever possible ("no ARIA is better than bad ARIA").
- **Keyboard navigation**: every interactive element must be reachable/operable via `Tab`/`Enter`/`Space` alone — a `<div onClick>` is invisible to keyboard users unless you add `tabIndex={0}` and a key handler, which is more work than just using a `<button>`.
- **ARIA attributes**: `aria-label` (accessible name when there's no visible text), `aria-hidden="true"` (hide decorative elements from screen readers), `role` (override/clarify an element's semantic role — use sparingly).
- **Labels & forms**: every `<input>` needs an associated `<label htmlFor="id">` (or `aria-label`) — screen readers announce the label when the input is focused.
- **Focus management**: after a modal opens, focus should move into it; after it closes, focus should return to the triggering element — otherwise keyboard/screen-reader users get lost.
- **Color contrast**: text needs sufficient contrast against its background (WCAG AA ≈ 4.5:1 for normal text) — don't rely on color alone to convey meaning (e.g. red/green only for errors).

---

## 13. JavaScript Coding Interview Patterns

Classic "write this from scratch" prompts. Runnable versions with test calls are in [examples/js-syntax.js](examples/js-syntax.js) (see the "Coding interview patterns" section at the bottom).

### Debounce

> **ELI5:** "Wait until the user stops typing for a bit before doing the expensive thing" — every new keystroke resets the clock.

```js
function debounce(fn, delayMs) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delayMs);
  };
}
```

### Throttle

> **ELI5:** "Do the thing at most once every X milliseconds, no matter how often it's triggered" — like a bouncer only letting one person in per second regardless of how many are pushing at the door.

```js
function throttle(fn, limitMs) {
  let inCooldown = false;
  return (...args) => {
    if (inCooldown) return;
    fn(...args);
    inCooldown = true;
    setTimeout(() => (inCooldown = false), limitMs);
  };
}
```

### Deep clone

> **ELI5:** A shallow copy (`{ ...obj }`) copies the outer box but still shares inner boxes (nested objects) with the original — mutate a nested object and both "copies" see it. A deep clone recursively copies everything, all the way down, so they're fully independent.

```js
function deepClone(value) {
  if (value === null || typeof value !== "object") return value;
  if (Array.isArray(value)) return value.map(deepClone);
  return Object.fromEntries(Object.entries(value).map(([k, v]) => [k, deepClone(v)]));
}
// modern shortcut (structured data only, no functions/DOM nodes): structuredClone(value)
```

### Curry

> **ELI5:** Turns a function that takes several arguments into a chain of functions that each take one argument — `add(1)(2)(3)` instead of `add(1, 2, 3)` — letting you "partially fill in" arguments now and supply the rest later.

```js
function curry(fn) {
  return function curried(...args) {
    if (args.length >= fn.length) return fn(...args);
    return (...more) => curried(...args, ...more);
  };
}
const add3 = curry((a, b, c) => a + b + c);
add3(1)(2)(3); // 6
add3(1, 2)(3); // 6
```

### Flatten a nested array

```js
function flatten(arr) {
  return arr.reduce((flat, item) =>
    flat.concat(Array.isArray(item) ? flatten(item) : item), []);
}
flatten([1, [2, [3, 4], 5]]); // [1, 2, 3, 4, 5]
// built-in shortcut: arr.flat(Infinity)
```

### Memoize

> **ELI5:** Remember the answer to a question you've already answered before, keyed by the exact inputs asked — so asking the same question again is instant instead of redone from scratch.

```js
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
```

---

## 14. Interview Q&A — Basic to Medium

**Q: What's the difference between `let`, `const`, and `var`?**
Scope (function vs block) and hoisting/TDZ behavior — see [table above](#variable-declarations). `const` prevents reassignment, not mutation.

**Q: Explain closures with an example.**
A function retains access to variables from its defining scope even after that scope exits. Classic example: a counter factory (`makeCounter`) where each call produces an independently-scoped `count`.

**Q: What is the event loop / why doesn't JS block on network calls?**
JS is single-threaded; async I/O is delegated to the browser/Node runtime, and completion callbacks are queued (micro/macrotask) to run when the call stack is empty. This is why `fetch` doesn't freeze the UI.

**Q: Why does `useEffect` need a dependency array, and what happens if you get it wrong?**
It controls when the effect re-runs. Missing a dependency creates a **stale closure** bug (effect uses an old value). Over-including causes redundant re-runs. ESLint's `react-hooks/exhaustive-deps` rule catches most mistakes.

**Q: What's the Virtual DOM and why does React use it?**
An in-memory representation of the UI tree. React diffs the new tree against the old one and applies only the minimal real DOM changes, which is (usually) cheaper than direct DOM manipulation for frequent updates, and lets React express UI declaratively.

**Q: Why do list items need a `key`, and why not use the array index?**
Keys let React match elements between renders. Index keys break when the list is reordered/filtered/inserted-in-the-middle — React can misattribute state (e.g. input focus/value) to the wrong row.

**Q: Controlled vs uncontrolled components?**
Controlled: React state is the single source of truth (`value` + `onChange`). Uncontrolled: the DOM itself holds the value, read imperatively via `ref`. Controlled is preferred for anything needing validation/derived UI; uncontrolled is fine for simple forms/file inputs.

**Q: When would you reach for `useMemo`/`useCallback`, and when is it a waste?**
Only when there's a measurable cost: an expensive computation, or preventing a memoized child (`React.memo`) from re-rendering due to a new prop reference. Wrapping everything by default adds complexity and its own (small) overhead for no benefit.

**Q: What causes unnecessary re-renders, and how do you find them?**
Parent re-renders cascading to non-memoized children, new inline object/array/function props each render, or broad Context consumers. Diagnose with React DevTools Profiler.

**Q: Explain `this` in JS — why do arrow functions behave differently?**
Regular function `this` is dynamic (determined by call-site: `obj.method()` binds `this` to `obj`). Arrow functions don't have their own `this` — they close over the `this` of their enclosing scope at definition time. This is why arrow functions are preferred for callbacks inside class/object methods.

**Q: What's the difference between `==` and `===`? Between `null` and `undefined`?**
`==` coerces types before comparing (footguns like `[] == false`); `===` does not. `undefined` = a variable declared but not assigned (or a missing object property); `null` = an explicit "no value" assigned intentionally.

**Q: How does one-way data binding work in React, and how do children communicate with parents?**
Props flow down only. Children "communicate up" by invoking a callback function passed down as a prop (the parent owns and updates the actual state).

**Q: What's a React error boundary, and why can't you catch render errors with try/catch?**
A class component implementing `static getDerivedStateFromError`/`componentDidCatch` that catches errors thrown during rendering in its child tree and renders a fallback UI instead of crashing the whole app. try/catch doesn't work across render because React's render phase runs your function body, throws propagate up through React's own call stack, not a `try` block you control at the call site.

**Q: What's CORS and why do you sometimes see it fail only in the browser?**
A browser security mechanism restricting cross-origin requests unless the server explicitly allows it via response headers. It's enforced by the browser only — server-to-server or Postman/curl calls aren't subject to it, which is why "it works in Postman but not in the browser" is a classic CORS symptom.

**Q: What's the difference between event bubbling and event capturing, and what is event delegation?**
Bubbling: the event fires on the target, then travels up through its ancestors. Capturing: the reverse, top-down, before bubbling (opt in with `addEventListener(type, fn, { capture: true })`). Delegation exploits bubbling — attach one listener on a common ancestor instead of one per child, and check `e.target` to find which child triggered it. Cheaper for long/dynamic lists.

**Q: `Promise.all` vs `Promise.allSettled` vs `Promise.race`?**
`all` resolves when every promise resolves, rejects immediately if any one rejects. `allSettled` always resolves once every promise has settled, giving you each outcome (fulfilled or rejected) individually — used when partial failure is acceptable. `race`/`any` settle as soon as the first promise settles/resolves respectively.

**Q: Shallow copy vs deep copy — why does it matter in React?**
A shallow copy (`{ ...state }`) duplicates only the top level; nested objects/arrays are still shared references with the original. Mutating a nested value after a shallow copy silently mutates the original too — a common source of "React isn't re-rendering" bugs, since React's `useState`/`React.memo` compare by reference and won't see a change if the nested object's reference stayed the same.

**Q: How would you make a custom `<div>` clickable element accessible?**
Prefer a real `<button>` — it's focusable, keyboard-activatable (`Enter`/`Space`), and announced correctly by screen readers for free. If a `<div>` is unavoidable, you'd need `role="button"`, `tabIndex={0}`, and manual `onKeyDown` handling for `Enter`/`Space` — which is exactly why using the real semantic element is almost always less work and more correct.

**Q: `interface` vs `type` — when do you pick one over the other?**
Both describe object shapes and are largely interchangeable for that case. `interface` supports declaration merging (reopening the same name to add members) and reads slightly better in public library APIs; `type` is required for unions, intersections, tuples, and mapped/conditional types, which `interface` can't express. Common convention: `interface` for object/component-prop shapes, `type` for everything else.

**Q: `any` vs `unknown` — why prefer `unknown`?**
Both accept anything, but `any` disables type checking on everything you do with that value afterward (so mistakes only surface at runtime). `unknown` forces you to narrow the type (via `typeof`, `instanceof`, or a custom guard) before you're allowed to use it, so the compiler still protects you. Use `unknown` for genuinely-unknown input (parsed JSON, `catch (e: unknown)`); reach for `any` only as a last resort, not a default.

**Q: What is a discriminated union, and why use `never` in the default branch of a switch?**
A union of object types that share a common literal "tag" field (e.g. `kind: "circle" | "square"`), which lets TypeScript narrow to the exact variant inside an `if`/`switch` without manual casts. Assigning the switch's default-case value to a variable typed `never` creates an "exhaustiveness check" — if a new variant is added to the union later and a case is missed, the assignment fails to compile instead of silently doing nothing at runtime.

**Q: Does TypeScript run in the browser, or does it get compiled away?**
TypeScript has no runtime of its own — types are checked at compile time and then fully erased. The output is plain JavaScript; by the time it reaches the browser (or Node), there's no TypeScript left, no performance cost, and no runtime type checking (which is why `unknown`/user-defined type guards matter — they're what actually check something at runtime, not just the type annotations).

**Q: How would you type a `useReducer` action so an invalid action shape is a compile error?**
Model the actions as a discriminated union (`type Action = { type: "increment" } | { type: "setStep"; payload: number }`) and type the reducer's second parameter as that union. Dispatching `{ type: "setStep" }` without `payload`, or `{ type: "bogus" }`, then fails to compile — the reducer's `switch` also gets exhaustiveness checking for free.

---

## 15. Quick Reference Tables

**TypeScript utility types**

| Utility | What it does |
|---|---|
| `Partial<T>` | every field optional |
| `Required<T>` | every field required (opposite of `Partial`) |
| `Readonly<T>` | every field readonly |
| `Pick<T, K>` | keep only keys `K` |
| `Omit<T, K>` | drop keys `K` |
| `Record<K, V>` | object type with keys `K`, all values `V` |
| `Exclude<T, U>` | remove members of union `T` assignable to `U` |
| `Extract<T, U>` | keep only members of union `T` assignable to `U` |
| `ReturnType<F>` | the return type of function type `F` |
| `Parameters<F>` | tuple of a function type `F`'s parameter types |
| `Awaited<T>` | unwraps a `Promise<T>` down to `T` |

**Array method returns**

| Method | Returns | Mutates original? |
|---|---|---|
| `map` | new array | no |
| `filter` | new array | no |
| `reduce` | any accumulated value | no |
| `forEach` | `undefined` | no |
| `sort` | same array, sorted | **yes** |
| `splice` | removed elements | **yes** |
| `slice` | new array (shallow) | no |
| `push`/`pop` | new length / removed item | **yes** |

**Hook cheat sheet**

| Hook | Purpose | Re-renders on change? |
|---|---|---|
| `useState` | local state | yes |
| `useReducer` | complex state transitions | yes |
| `useEffect` | side effects / lifecycle | no (runs after render) |
| `useRef` | mutable value / DOM ref | no |
| `useMemo` | memoize computed value | only if deps change |
| `useCallback` | memoize function identity | only if deps change |
| `useContext` | read nearest Provider value | yes, when value changes |

**HTTP status codes**

| Range | Meaning |
|---|---|
| 2xx | success (200 OK, 201 Created, 204 No Content) |
| 3xx | redirect (301 Moved, 304 Not Modified) |
| 4xx | client error (400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 429 Too Many Requests) |
| 5xx | server error (500 Internal Server Error, 503 Service Unavailable) |
