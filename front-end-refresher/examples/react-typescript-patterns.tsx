// ============================================================
// REACT + TYPESCRIPT PATTERNS (illustrative — type-checked via
// `npm run typecheck`, not meant to run standalone; needs a bundler)
// Companion code for ../README.md section 6 (React + TypeScript Patterns)
// ============================================================
import {
  useState, useRef, useReducer, forwardRef, type ReactNode, type ChangeEvent, type MouseEvent,
} from "react";

// ---- 1. Typing props ----
interface ButtonProps {
  label: string;                 // required
  onClick?: () => void;          // optional callback prop
  variant?: "primary" | "secondary"; // literal union instead of a loose `string`
  children?: ReactNode;          // covers strings, elements, arrays, fragments, null...
}
function Button({ label, onClick, variant = "primary", children }: ButtonProps) {
  return (
    <button className={variant} onClick={onClick}>
      {label}
      {children}
    </button>
  );
}

// ---- 2. Typing useState ----
function Counter() {
  const [count, setCount] = useState<number>(0);       // explicit — usually inferred fine from useState(0)
  const [user, setUser] = useState<{ name: string } | null>(null); // union needed when initial value is null
  return (
    <div>
      <p>{count}</p>
      <button onClick={() => setCount(c => c + 1)}>+1</button>
      {user && <p>{user.name}</p>}
    </div>
  );
}

// ---- 3. Typing useRef ----
function TextInput() {
  const inputRef = useRef<HTMLInputElement>(null);      // DOM ref — starts null, TS knows the element type
  const renderCount = useRef<number>(0);                // mutable value ref, not tied to the DOM
  renderCount.current++;
  return <input ref={inputRef} onFocus={() => inputRef.current?.select()} />;
}

// ---- 4. Typing event handlers ----
function SearchBox() {
  const [query, setQuery] = useState("");
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => setQuery(e.target.value);
  const handleSubmitClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    console.log("searching for", query);
  };
  return (
    <>
      <input value={query} onChange={handleChange} />
      <button onClick={handleSubmitClick}>Search</button>
    </>
  );
}

// ---- 5. Discriminated union action types with useReducer ----
interface CounterState { count: number; step: number; }
type CounterAction =
  | { type: "increment" }
  | { type: "setStep"; payload: number };

function reducer(state: CounterState, action: CounterAction): CounterState {
  switch (action.type) {
    case "increment": return { ...state, count: state.count + state.step };
    case "setStep": return { ...state, step: action.payload }; // payload only exists on this branch — TS enforces it
  }
}
function ReducerCounter() {
  const [state, dispatch] = useReducer(reducer, { count: 0, step: 1 });
  return <button onClick={() => dispatch({ type: "increment" })}>{state.count}</button>;
}

// ---- 6. Generic components ----
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => ReactNode;
  keyExtractor: (item: T) => string | number;
}
function List<T>({ items, renderItem, keyExtractor }: ListProps<T>) {
  return <ul>{items.map(item => <li key={keyExtractor(item)}>{renderItem(item)}</li>)}</ul>;
}
// usage: <List items={users} keyExtractor={u => u.id} renderItem={u => u.name} />
// T is inferred from whatever `items` array you pass in — no need to specify <User> manually

// ---- 7. forwardRef with typed props + ref ----
interface FancyInputProps { placeholder?: string; }
const FancyInput = forwardRef<HTMLInputElement, FancyInputProps>(function FancyInput(props, ref) {
  return <input ref={ref} placeholder={props.placeholder} className="fancy" />;
});

// ---- 8. Typing a custom hook's return value ----
function useToggle(initial = false): [boolean, () => void] {
  const [value, setValue] = useState(initial);
  const toggle = () => setValue(v => !v);
  return [value, toggle]; // tuple return — order matters, like useState's [value, setter]
}
function Feature() {
  const [enabled, toggleEnabled] = useToggle();
  return <button onClick={toggleEnabled}>{enabled ? "On" : "Off"}</button>;
}

export { Button, Counter, TextInput, SearchBox, ReducerCounter, List, FancyInput, useToggle, Feature };
