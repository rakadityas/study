# time complexity: O(N) amortised per sequence of operations — each pair is enqueued once and dequeued once, so the total work across all deposit/release calls is O(N) where N = total pairs submitted. A single deposit can cost O(N) in the worst case (unblocking every pending group at once).
# space complexity: O(N) — the pending deque holds at most all N unprocessed pairs at any point.

# You are building a funds management system for an online marketplace platform.
# Whenever a customer makes a purchase, the payment is first routed into a centralized "Holding Account" managed by the platform. Once the customer confirms that the goods or services have been received, the system will initiate a funds distribution, releasing the money  holding account to one or more merchant bank accounts.
# Your task is to implement a transaction processing system that manages this holding account and processes instructions provided by the platform.

from collections import deque

class HoldingAccount:
    def __init__(self):
        self.balance = 0
        # Each entry is a list of (account, amount) pairs representing one release instruction
        self.pending: deque[list[tuple[str, int]]] = deque()

    def deposit(self, amount: int) -> tuple[int, list[tuple[str, int]]]:
        self.balance += amount
        released = self._process_pending()
        return self.balance, released

    def release(self, pairs: list[tuple[str, int]]) -> tuple[int, list[tuple[str, int]]]:
        self.pending.append(pairs)
        released = self._process_pending()
        return self.balance, released

    def _process_pending(self) -> list[tuple[str, int]]:
        released = []
        while self.pending:
            group = self.pending[0]

            total = 0 
            for i in range(len(group)):
                total += group[i][1]

            if self.balance < total:
                break

            self.pending.popleft()
            for i in range(len(group)):
                account, amount = group[i]
                self.balance -= amount
                released.append((account, amount))

        return released

if __name__ == "__main__":
    account = HoldingAccount()

    assert account.deposit(1000) == (1000, [])
    assert account.release([("acctA", 500), ("acctB", 400)]) == (100, [("acctA", 500), ("acctB", 400)])
    assert account.release([("acctC", 50), ("acctA", 1000)]) == (100, [])   # blocked: needs 1050, has 100
    assert account.deposit(900) == (1000, [])                                # still blocked: needs 1050, has 1000
    assert account.release([("acctD", 10)]) == (1000, [])                   # queued behind the blocked release
    assert account.deposit(50) == (0, [("acctC", 50), ("acctA", 1000)])     # unblocks first pending group
    assert account.deposit(50) == (40, [("acctD", 10)])                     # unblocks second pending group

    print("All assertions passed.")
