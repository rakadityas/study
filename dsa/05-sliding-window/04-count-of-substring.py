class SolutionOne:
    def countOfSubstring(self, s: str) -> int:
        l, r = 0, 0
        mapHistory = {}
        numSubstring = 0

        if s == "":
            return 0

        while r < len(s):
            if s[r] in mapHistory:
                numSubstring += 1

                while l != r:
                    del mapHistory[s[l]]
                    l += 1
            
            mapHistory[s[r]] = True
            r += 1

        return numSubstring+1
                
if __name__ == "__main__":
    s = SolutionOne()
    assert s.countOfSubstring("abcabcbb") == 4
    assert s.countOfSubstring("bbbbb") == 5
    assert s.countOfSubstring("") == 0