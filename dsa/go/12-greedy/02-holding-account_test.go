package greedy

import "testing"

// Payout is one (account, amount) leg of a release instruction.
type Payout struct {
    Account string
    Amount  int
}

// HoldingAccount models a marketplace escrow account. Purchases deposit into it, and release
// instructions pay merchants out once the balance covers the whole instruction.
//
// Approach: a FIFO queue of pending instructions, drained from the front after every deposit or
// release. An instruction that the balance cannot cover blocks everything queued behind it, which
// keeps payouts strictly in submission order.
// time: O(N) amortised across a whole sequence of operations — each instruction is enqueued and
// dequeued once, though a single deposit can unblock every pending group at once
// space: O(N) — the queue holds at most all N unprocessed instructions
type HoldingAccount struct {
    balance int
    pending [][]Payout
}

func ConstructorHoldingAccount() *HoldingAccount {
    return &HoldingAccount{pending: [][]Payout{}}
}

// Deposit adds funds and drains whatever that unblocks.
// time: O(N) worst case, space: O(1)
func (h *HoldingAccount) Deposit(amount int) (int, []Payout) {
    h.balance += amount
    return h.balance, h.processPending()
}

// Release queues a payout instruction and drains whatever is now payable.
// time: O(N) worst case, space: O(1)
func (h *HoldingAccount) Release(pairs []Payout) (int, []Payout) {
    h.pending = append(h.pending, pairs)
    return h.balance, h.processPending()
}

// processPending pays out queued instructions from the front until one is unaffordable.
// time: O(N), space: O(1)
func (h *HoldingAccount) processPending() []Payout {
    released := []Payout{}

    for len(h.pending) > 0 {
        group := h.pending[0]

        total := 0
        for _, p := range group {
            total += p.Amount
        }

        if h.balance < total { break }

        h.pending = h.pending[1:]
        h.balance -= total
        released = group
    }

    return released
}

func equalPayouts(a, b []Payout) bool {
    if len(a) != len(b) { return false }
    for i := range a {
        if a[i] != b[i] { return false }
    }
    return true
}

func TestHoldingAccount(t *testing.T) {
    account := ConstructorHoldingAccount()

    check := func(step string, gotBal int, gotRel []Payout, wantBal int, wantRel []Payout) {
        t.Helper()
        if gotBal != wantBal || !equalPayouts(gotRel, wantRel) {
            t.Fatalf("%s = (%d, %v), want (%d, %v)", step, gotBal, gotRel, wantBal, wantRel)
        }
    }

    bal, rel := account.Deposit(1000)
    check("Deposit(1000)", bal, rel, 1000, []Payout{})

    bal, rel = account.Release([]Payout{{"acctA", 500}, {"acctB", 400}})
    check("Release(A/B)", bal, rel, 100, []Payout{{"acctA", 500}, {"acctB", 400}})

    // blocked: needs 1050, has 100
    bal, rel = account.Release([]Payout{{"acctC", 50}, {"acctA", 1000}})
    check("Release(C/A)", bal, rel, 100, []Payout{})

    // still blocked: needs 1050, has 1000
    bal, rel = account.Deposit(900)
    check("Deposit(900)", bal, rel, 1000, []Payout{})

    // queued behind the blocked instruction even though it is affordable on its own
    bal, rel = account.Release([]Payout{{"acctD", 10}})
    check("Release(D)", bal, rel, 1000, []Payout{})

    // unblocks the first pending group
    bal, rel = account.Deposit(50)
    check("Deposit(50)", bal, rel, 0, []Payout{{"acctC", 50}, {"acctA", 1000}})

    // unblocks the second pending group
    bal, rel = account.Deposit(50)
    check("Deposit(50)", bal, rel, 40, []Payout{{"acctD", 10}})
}
