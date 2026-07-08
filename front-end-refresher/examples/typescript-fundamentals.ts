// ============================================================
// TYPESCRIPT PLAYGROUND — run with `npm run ts` (uses tsx)
// Type-check only (no execution) with `npm run typecheck`
// Companion code for ../README.md section 2 (TypeScript Deep Dive)
// ============================================================

// ---- 1. Basic types & inference ----
let id: number = 1;
let username: string = "Rio";
let isActive: boolean = true;
let tags: string[] = ["a", "b"];       // array of strings
let coords: [number, number] = [1, 2]; // tuple — fixed length, fixed types per slot

// TS infers types when you don't annotate — prefer letting it infer for locals,
// annotate function boundaries (params/return) explicitly
let inferred = "hello"; // inferred as `string`, NOT the literal type "hello" (that's `let`)
const literal = "hello"; // inferred as the literal type "hello" (const can't change)

console.log({ id, username, isActive, tags, coords, inferred, literal });

// ---- 2. Enums vs literal union types ----
// Numeric enum — generates a real runtime object with reverse mapping
enum StatusEnum { Pending, Active, Closed } // 0, 1, 2 under the hood
console.log(StatusEnum.Active, StatusEnum[1]); // 1, "Active"

// Literal union — the far more common modern preference: no runtime cost,
// works better with plain JS objects/APIs, easier to serialize to JSON
type Status = "pending" | "active" | "closed";
function setStatus(s: Status) { return `status: ${s}`; }
console.log(setStatus("active"));
// setStatus("done"); // compile error — not in the union

// ---- 3. Interfaces vs type aliases ----
// Both describe object shapes; interfaces can be "reopened" (declaration merging),
// type aliases can express unions/primitives/tuples directly. For plain object
// shapes, either works — many teams pick `interface` for public API shapes and
// `type` for unions/aliases.
interface UserInterface {
  id: number;
  name: string;
  age?: number;        // optional
  readonly createdAt: string; // can't be reassigned after creation
}
type UserType = {
  id: number;
  name: string;
};
interface Employee extends UserInterface { role: string; } // interface extension
type Manager = UserType & { reports: number };               // type intersection

const emp: Employee = { id: 1, name: "Dita", createdAt: "2024-01-01", role: "eng" };
console.log(emp);

// ---- 4. Index signatures ----
interface StringDictionary {
  [key: string]: string; // any string key maps to a string value
}
const translations: StringDictionary = { hello: "hola", bye: "adios" };
console.log(translations["hello"]);

// ---- 5. Functions: params, optional/default, overloads ----
function add(a: number, b: number): number { return a + b; }
function greet(name: string = "world"): string { return `Hello, ${name}`; }
function firstEl<T>(arr: T[]): T | undefined { return arr[0]; } // generic — see section 7

// Overloads: same function, different call signatures resolved at compile time
function toArray(value: string): string[];
function toArray(value: number): number[];
function toArray(value: string | number): string[] | number[] {
  return typeof value === "string" ? value.split("") : [value];
}
console.log(add(1, 2), greet(), firstEl([1, 2, 3]), toArray("ab"), toArray(5));

// ---- 6. Union & intersection types ----
type IdType = string | number;                 // "OR" — could be either
type Timestamped = { createdAt: string };
type Named = { name: string };
type NamedAndTimestamped = Timestamped & Named; // "AND" — must satisfy both

const record: NamedAndTimestamped = { name: "Rio", createdAt: "2024-01-01" };
console.log(record);

// ---- 7. Discriminated unions & exhaustiveness ----
// A shared literal field (the "discriminant") lets TS narrow which shape you have.
type Shape =
  | { kind: "circle"; radius: number }
  | { kind: "square"; side: number }
  | { kind: "rectangle"; width: number; height: number };

function area(shape: Shape): number {
  switch (shape.kind) {
    case "circle": return Math.PI * shape.radius ** 2;   // shape narrowed to circle here
    case "square": return shape.side ** 2;                // narrowed to square here
    case "rectangle": return shape.width * shape.height;  // narrowed to rectangle here
    default: {
      const _exhaustive: never = shape; // compile error if a Shape variant is unhandled
      return _exhaustive;
    }
  }
}
console.log(area({ kind: "circle", radius: 2 }));

// ---- 8. Type narrowing (typeof, instanceof, in, custom guards) ----
function describe(value: string | number | Date) {
  if (typeof value === "string") return value.toUpperCase();       // typeof guard
  if (typeof value === "number") return value.toFixed(2);
  if (value instanceof Date) return value.toISOString();           // instanceof guard
}

interface Cat { meow(): void }
interface Bird { fly(): void }
function isCat(pet: Cat | Bird): pet is Cat { return "meow" in pet; } // custom type guard ("in" + user-defined predicate)
function handlePet(pet: Cat | Bird) {
  if (isCat(pet)) pet.meow(); else pet.fly();
}
handlePet({ meow: () => console.log("meow!") });
console.log(describe("hello"), describe(42), describe(new Date(0)));

// ---- 9. Generics: functions, interfaces, constraints, defaults ----
function identity<T>(value: T): T { return value; }

interface ApiResponse<T> { data: T; error?: string; }
const res: ApiResponse<UserType> = { data: { id: 1, name: "Rio" } };

// Constrained generic — T must have a `.length` property
function logLength<T extends { length: number }>(value: T): number { return value.length; }
console.log(logLength("hello"), logLength([1, 2, 3]));

// Default generic parameter
interface Box<T = string> { value: T; }
const box: Box = { value: "default is string" };

console.log(identity(42), res, box);

// ---- 10. Utility types ----
interface Article { id: number; title: string; body: string; published: boolean; }

type PartialArticle = Partial<Article>;        // every field optional (e.g. PATCH payload)
type ArticlePreview = Pick<Article, "id" | "title">; // only these keys
type ArticleNoBody = Omit<Article, "body">;    // all except this key
type ReadonlyArticle = Readonly<Article>;      // every field readonly
type ArticleMap = Record<number, Article>;     // { [id: number]: Article }
type ArticleStatus = Exclude<"draft" | "published" | "archived", "archived">; // "draft" | "published"
type OnlyStrings = Extract<string | number | boolean, string>; // string

function fetchArticle(): Promise<Article> { return Promise.resolve({} as Article); }
type FetchedArticle = Awaited<ReturnType<typeof fetchArticle>>; // Article (unwraps Promise + gets return type)

const preview: ArticlePreview = { id: 1, title: "Hello" };
console.log(preview);

// ---- 11. keyof, typeof, indexed access ----
const config = { retries: 3, timeoutMs: 5000 };
type ConfigKey = keyof typeof config;          // "retries" | "timeoutMs"
function getConfigValue<K extends ConfigKey>(key: K): typeof config[K] {
  return config[key];
}
console.log(getConfigValue("retries"));

// ---- 12. Mapped & conditional types ----
type Optional<T> = { [K in keyof T]?: T[K] };              // mapped type — same shape, all optional
type NonFunctionKeys<T> = { [K in keyof T]: T[K] extends Function ? never : K }[keyof T];

type IsString<T> = T extends string ? "yes" : "no";        // conditional type
type A = IsString<"hi">;  // "yes"
type B = IsString<42>;    // "no"

type UnwrapArray<T> = T extends (infer U)[] ? U : T;        // `infer` extracts a nested type
type Item = UnwrapArray<string[]>; // string

// ---- 13. any vs unknown vs never vs void ----
let anything: any = 5;
try {
  anything.foo.bar.baz; // compiles fine (any disables checking entirely)! Only blows up at RUNTIME — avoid `any`.
} catch (e) {
  console.log("any let a bad access through the type checker:", (e as Error).message);
}

let notSure: unknown = 5;
if (typeof notSure === "number") console.log(notSure.toFixed(1)); // must narrow before use — safer than any

function fail(msg: string): never { throw new Error(msg); }        // never returns (always throws/loops)
function logMsg(msg: string): void { console.log(msg); }            // returns nothing meaningful

// ---- 14. Classes: access modifiers, readonly, abstract ----
abstract class Animal {
  protected readonly name: string;               // accessible in this class + subclasses only, can't reassign
  private secret = "hidden";                     // only accessible inside Animal itself
  constructor(name: string) { this.name = name; }
  abstract speak(): string;                      // subclasses MUST implement this
  describe(): string { return `${this.name}: ${this.speak()}`; }
}
class Dog extends Animal {
  speak(): string { return "Woof"; }
}
console.log(new Dog("Rex").describe());
// new Animal("x"); // compile error — can't instantiate an abstract class

// ---- 15. Assertions & the `satisfies` operator ----
// Type assertion (`as`) — "trust me, treat this as X" (browser-only API, shown for reference):
// const input = document.querySelector("input") as HTMLInputElement | null;
// const forced = input!.value; // non-null assertion (!) — only when you're SURE it's not null
const maybeString = Math.random() > 2 ? undefined : "value"; // always "value" here, just for the demo
const asserted = maybeString as string; // "trust me" — bypasses the undefined check

// `satisfies` checks a value matches a type WITHOUT widening/losing the literal type
const palette = { primary: "#000", secondary: "#fff" } satisfies Record<string, string>;
// palette.primary is still typed as the literal "#000", unlike `: Record<string,string>` which would widen it to `string`

console.log(asserted, palette);

// ---- 16. Type-only imports (erased at compile time, zero runtime cost) ----
// import type { SomeType } from "./types"; // stripped entirely from the emitted JS
