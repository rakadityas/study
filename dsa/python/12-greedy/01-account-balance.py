# Say we have a number of bank accounts with name and balance. ie. ( name, balance )
# At the end of each day we want to transfer money between accounts to ensure that each account
# has at least 100 dollars in it.
# Write a program to determine the transfers required to satisfy the requirement. Each transfer 
# should be presented as ( from_account_name, to_account_name, amount ).
# For example:
#   bank_accounts = [("a", 80), ("b", 180)]
#   transfers = determine_transfers(bank_accounts)
#   # expected result
#   # transfers = [("b", "a", 20)]
# Ie. you are expected to implement the function below, appropritely adjusted for your solution:
#   List<Transfer> determine_transfers(List<BankAccount> bank_accounts)

from typing import List

class Solution:
    def determine_transfers(self, accounts: List[List]):
        sources = []
        targets = []

        for i in range(len(accounts)):
            if accounts[i][1] > 100:
                sources.append(accounts[i])
            else:
                targets.append(accounts[i])
        
        res = []
        sourceIdx = 0
        for i in range(len(targets)):
            gap = 100-targets[i][1]

            while gap > 0 and sourceIdx < len(sources):
                sourceAcc, sourceBal = sources[sourceIdx]                
                surplus = sourceBal-100
                
                trfAmt = min(surplus, gap)    
                gap -= trfAmt
                
                res.append([sourceAcc, targets[i][0], trfAmt])

                sources[sourceIdx][1] -= trfAmt
                if sources[sourceIdx][1] == 100:
                    sourceIdx += 1

        return res

if __name__ == "__main__":
    mySol = Solution()

    # two accounts available, and only deduct once
    assert mySol.determine_transfers([["a", 80], ["b", 120]]) == [["b", "a", 20]]

    # more than one accounts deductable
    assert mySol.determine_transfers([["a", 140], ["b", 80],  ["c", 80]]) == [["a", "b", 20], ["a", "c", 20]]
    
    # multiple balance deduction
    assert mySol.determine_transfers([["a", 110], ["b", 95],  ["c", 90], ["d", 105]]) == [['a', 'b', 5], ['a', 'c', 5], ['d', 'c', 5]]

    assert mySol.determine_transfers([["a", 200], ["b", 95],  ["c", 90], ["d", 50]]) == [['a', 'b', 5], ['a', 'c', 10], ['a', 'd', 50]]

    # no source transfer balance
    assert mySol.determine_transfers([["a", 50], ["b", 95],  ["c", 90], ["d", 50]]) == []
