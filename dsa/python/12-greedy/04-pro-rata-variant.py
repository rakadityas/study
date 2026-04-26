# their_share / total_owed * available_balance

def release_proportional(
    balance: int,
    pairs: list[tuple[str, int]]
) -> tuple[int, list[tuple[str, float]]]:
    total_owed = 0
    for i in range(len(pairs)):
        total_owed += pairs[i][1]

    if total_owed == 0:
        return (balance, [])

    to_distribute = min(balance, total_owed)

    result = []

    for i in range(len(pairs)):
        pairs_acc, pairs_amt = pairs[i]
        prorated_amt = (pairs_amt / total_owed) * to_distribute
        result.append((pairs_acc, prorated_amt))

    return (balance - to_distribute, result)

assert release_proportional(1000, [("A", 600), ("B", 400)]) == (0, [("A", 600), ("B", 400)])
assert release_proportional(500,  [("A", 600), ("B", 400)]) == (0, [("A", 300.0), ("B", 200.0)])
assert release_proportional(100,  [("A", 600), ("B", 400)]) == (0, [("A", 60.0), ("B", 40.0)])
assert release_proportional(1500, [("A", 600), ("B", 400)]) == (500, [("A", 600), ("B", 400)])