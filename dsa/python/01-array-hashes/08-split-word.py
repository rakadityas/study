from typing import List

class Solution:
    # Brute force: for every column index, compare every word's char against
    # the first word's char at that column. O(n*m), n = #words, m = shortest word length.
    def GetSameLetter(self, words: List[str]) -> str:
        if not words:
            return ""

        res = []
        smallestLen = len(words[0])
        for i in range(len(words)):
            smallestLen = min(len(words[i]), smallestLen)

        for i in range(0, smallestLen, 1):
            prevWord = ""
            for j in range(len(words)):
                if prevWord == "":
                    prevWord = words[j][i]
                    continue

                if prevWord != words[j][i]:
                    return "".join(res)


            res.append(words[0][i])

        return "".join(res)

    # The common prefix of the whole list equals the common prefix of just
    # its lexicographic min and max (any divergent word would sit strictly
    # between them, so it can't shrink the prefix further). min()/max() find
    # those two extremes in a single O(n) pass each, no full sort needed
    # (each comparison up to O(m)), then O(m) to compare the two ends.
    def GetSameLetterOptimal(self, words: List[str]) -> str:
        if not words:
            return ""

        first, last = min(words), max(words)
        smallestLen = min(len(first), len(last))

        res = []
        for i in range(smallestLen):
            if first[i] != last[i]:
                break
            res.append(first[i])

        return "".join(res)

    # True O(n) where n = total number of characters across all words.
    # zip(*words) walks column-by-column without an inner index loop; each
    # character is visited at most once, and we bail out on the first
    # mismatched column. No sorting overhead like GetSameLetterOptimal.
    def GetSameLetterLinear(self, words: List[str]) -> str:
        if not words:
            return ""

        res = []
        for chars in zip(*words):
            if any(c != chars[0] for c in chars):
                break
            res.append(chars[0])

        return "".join(res)

if __name__ == "__main__":
    sol = Solution()
    print(sol.GetSameLetter(["APPLE", "APP", "APPLAUD"]), "hello")
    assert sol.GetSameLetter(["APPLE", "APP", "APPLAUD"]) == "APP"
    assert sol.GetSameLetter(["BANANA", "BAN", "BLUE"]) == "B"
    assert sol.GetSameLetter(["BANANA", "APPLE", "CAR"]) == ""

    assert sol.GetSameLetterOptimal(["APPLE", "APP", "APPLAUD"]) == "APP"
    assert sol.GetSameLetterOptimal(["BANANA", "BAN", "BLUE"]) == "B"
    assert sol.GetSameLetterOptimal(["BANANA", "APPLE", "CAR"]) == ""

    assert sol.GetSameLetterLinear(["APPLE", "APP", "APPLAUD"]) == "APP"
    assert sol.GetSameLetterLinear(["BANANA", "BAN", "BLUE"]) == "B"
    assert sol.GetSameLetterLinear(["BANANA", "APPLE", "CAR"]) == ""


