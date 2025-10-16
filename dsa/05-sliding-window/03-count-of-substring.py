class SolutionOne:
    def countOfSubstring(self, s: str) -> int:
        l = 0
        mapHistory = {}
        numSubstring = 0

        if s == "":
            return 0

        for i in range(len(s)):
            if s[i] in mapHistory:
                numSubstring += 1

                while l != i:
                    del mapHistory[s[l]]
                    l += 1
            
            mapHistory[s[i]] = True

        return numSubstring+1
                
if __name__ == "__main__":
    s = SolutionOne()
    assert s.countOfSubstring("abcabcbb") == 4
    assert s.countOfSubstring("bbbbb") == 5
    assert s.countOfSubstring("") == 0