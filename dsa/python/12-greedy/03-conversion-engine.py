# When a customer wants to convert USD to EUR, Airwallex sources liquidity from multiple providers. 
# Each provider offers a rate (EUR received per USD sent) and a maximum volume they can handle. 
# To get the best deal, you must use highest-rate providers first.
# Implement a conversion engine that fills the requested amount greedily across providers.

# providers = [(name, eur_per_usd_rate, max_usd_volume)]
# Return: (total_eur_received, [(provider_name, usd_used, eur_received)])
# If total liquidity is insufficient, fill as much as possible.

# time complexity is O(N log N) dominated by sort, O(N) space

def best_execution(
    amount_usd: int,
    providers: list[tuple[str, float, int]]
) -> tuple[float, list[tuple[str, int, float]]]:
    sorted_providers = sorted(providers, key=lambda x: x[1], reverse=True)

    total = 0.0
    breakdown = []
    remaining = amount_usd

    for i in range(len(sorted_providers)):
        name, rate, max_vol = sorted_providers[i]

        if remaining <= 0:
            break

        usd_used = min(remaining, max_vol)
        if usd_used == 0:
            continue

        eur_received = round(usd_used * rate, 2)
        total += eur_received
        
        breakdown.append((name, usd_used, eur_received))
        remaining -= usd_used

    return total, breakdown

if __name__ == "__main__":
    providers = [
        ("Bank A", 1.08, 10_000),
        ("Bank B", 1.07,  5_000),
        ("Bank C", 1.09,  3_000),
    ]

    # Fill 15,000 USD: Bank C first (best rate, 3k), then Bank A (10k), then Bank B (2k)
    total, breakdown = best_execution(15_000, providers)
    assert breakdown == [("Bank C", 3_000, 3_270.0), ("Bank A", 10_000, 10_800.0), ("Bank B", 2_000, 2_140.0)]
    assert round(total, 2) == 16_210.0

    # Exact fill — uses only one provider
    total, breakdown = best_execution(3_000, providers)
    assert breakdown == [("Bank C", 3_000, 3_270.0)]
    assert round(total, 2) == 3_270.0

    # Insufficient liquidity — fill what's available (3k + 10k + 5k = 18k max, asking 20k)
    total, breakdown = best_execution(20_000, providers)
    assert breakdown == [("Bank C", 3_000, 3_270.0), ("Bank A", 10_000, 10_800.0), ("Bank B", 5_000, 5_350.0)]
    assert round(total, 2) == 19_420.0

    print("All assertions passed.")