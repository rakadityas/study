// ============================================================
// REACT HOOKS PLAYGROUND (illustrative — not wired to a bundler)
// Companion code for ../README.md section 3 & 4
// Assumes React 18+ function components. No class components used
// on purpose — that's legacy company code you may still encounter,
// see README "Class components (legacy)" section for the mapping.
// ============================================================
import { useState, useEffect, useRef, useMemo, useCallback, useReducer, createContext, useContext, memo, forwardRef, lazy, Suspense } from "react";

// ---- 1. useState — local component state ----
function Counter() {
  const [count, setCount] = useState(0); // [value, setter], initial value only used on first render

  // Functional updater form — always use this when new state depends on old state
  // (batched updates + concurrent rendering can make the "current" count stale otherwise)
  const incrementTwice = () => {
    setCount(prev => prev + 1);
    setCount(prev => prev + 1); // count is now +2, not +1 — proves functional update works
  };

  return (
    <div>
      <p>{count}</p>
      <button onClick={() => setCount(c => c + 1)}>+1</button>
      <button onClick={incrementTwice}>+2</button>
    </div>
  );
}

// ---- 2. useEffect — side effects & lifecycle replacement ----
function UserProfile({ userId }) {
  const [user, setUser] = useState(null);

  useEffect(() => {
    let cancelled = false; // guard against setting state after unmount / stale request

    async function load() {
      const res = await fetch(`/api/users/${userId}`);
      const data = await res.json();
      if (!cancelled) setUser(data);
    }
    load();

    return () => {
      cancelled = true; // cleanup — runs before next effect AND on unmount
    };
  }, [userId]); // dependency array: effect re-runs ONLY when userId changes
  // []              -> run once, like componentDidMount
  // no array at all -> run after EVERY render (rarely what you want)

  if (!user) return <p>Loading...</p>;
  return <p>{user.name}</p>;
}

// ---- 3. useRef — mutable value that survives renders WITHOUT causing one ----
function TextInputWithFocusButton() {
  const inputRef = useRef(null);          // DOM node access
  const renderCount = useRef(0);          // "instance variable" equivalent, no re-render on change
  renderCount.current++;

  return (
    <>
      <input ref={inputRef} />
      <button onClick={() => inputRef.current.focus()}>Focus</button>
      <p>Rendered {renderCount.current} times</p>
    </>
  );
}

// ---- 4. useMemo & useCallback — memoization to avoid wasted work/renders ----
function ExpensiveList({ items, onSelect }) {
  // useMemo: memoize a COMPUTED VALUE, recompute only if items changes
  const sorted = useMemo(() => [...items].sort((a, b) => a.price - b.price), [items]);

  // useCallback: memoize a FUNCTION IDENTITY so children wrapped in React.memo
  // don't re-render just because a new function reference was passed as a prop
  const handleSelect = useCallback((id) => onSelect(id), [onSelect]);

  return (
    <ul>
      {sorted.map(item => (
        <li key={item.id} onClick={() => handleSelect(item.id)}>{item.name}</li>
      ))}
    </ul>
  );
}

// ---- 5. useReducer — state machine style updates (Redux-lite, no library) ----
const initialState = { count: 0, step: 1 };
function reducer(state, action) {
  switch (action.type) {
    case "increment": return { ...state, count: state.count + state.step };
    case "setStep": return { ...state, step: action.payload };
    default: throw new Error(`Unknown action: ${action.type}`);
  }
}
function ReducerCounter() {
  const [state, dispatch] = useReducer(reducer, initialState);
  return (
    <div>
      <p>{state.count}</p>
      <button onClick={() => dispatch({ type: "increment" })}>+{state.step}</button>
      <input type="number" value={state.step}
        onChange={e => dispatch({ type: "setStep", payload: Number(e.target.value) })} />
    </div>
  );
}

// ---- 6. Context — avoid prop drilling ----
const ThemeContext = createContext("light");
function ThemedButton() {
  const theme = useContext(ThemeContext); // subscribes to nearest <ThemeContext.Provider>
  return <button className={theme}>Themed</button>;
}
function App() {
  return (
    <ThemeContext.Provider value="dark">
      <ThemedButton />
    </ThemeContext.Provider>
  );
}

// ---- 7. Custom hook — extract & reuse stateful logic ----
function useDebouncedValue(value, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const id = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(id); // cancel the pending timeout if value changes again
  }, [value, delayMs]);
  return debounced;
}
function SearchBox() {
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query, 400);

  useEffect(() => {
    if (debouncedQuery) console.log("Searching for:", debouncedQuery);
  }, [debouncedQuery]);

  return <input value={query} onChange={e => setQuery(e.target.value)} />;
}

// ---- 8. Controlled vs uncontrolled inputs ----
function ControlledForm() {
  const [value, setValue] = useState(""); // React state is the single source of truth
  return <input value={value} onChange={e => setValue(e.target.value)} />;
}
function UncontrolledForm() {
  const ref = useRef();                    // DOM itself holds the value
  const handleSubmit = () => console.log(ref.current.value);
  return <input ref={ref} defaultValue="init" />;
}

// ---- 9. Conditional rendering & lists ----
function ListExample({ items, isLoading }) {
  if (isLoading) return <Spinner />;
  return (
    <ul>
      {items.length === 0
        ? <li>No items</li>
        : items.map(item => <li key={item.id}>{item.name}</li>) /* key must be stable + unique, never index if list reorders */
      }
    </ul>
  );
}
function Spinner() { return <div>Loading…</div>; }

// ---- 10. React.memo — skip re-render if props are shallow-equal ----
const Row = memo(function Row({ item, onSelect }) {
  console.log("rendering row", item.id); // won't log again if item/onSelect refs are unchanged
  return <li onClick={() => onSelect(item.id)}>{item.name}</li>;
});

// ---- 11. forwardRef — let a parent reach a child's DOM node ----
const FancyInput = forwardRef(function FancyInput(props, ref) {
  return <input ref={ref} {...props} className="fancy" />;
});
function ParentWithRef() {
  const inputRef = useRef(null);
  return (
    <>
      <FancyInput ref={inputRef} placeholder="type here" />
      <button onClick={() => inputRef.current.focus()}>Focus fancy input</button>
    </>
  );
}

// ---- 12. Code splitting — React.lazy + Suspense ----
const SettingsPanel = lazy(() => import("./SettingsPanel")); // hypothetical module
function SettingsPage() {
  return (
    <Suspense fallback={<Spinner />}>
      <SettingsPanel />
    </Suspense>
  );
}

export { Counter, UserProfile, TextInputWithFocusButton, ExpensiveList, ReducerCounter, App, SearchBox, ControlledForm, UncontrolledForm, ListExample, Row, FancyInput, ParentWithRef, SettingsPage };
