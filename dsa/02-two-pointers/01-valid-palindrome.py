
class Solution:
    def isPalindrome(self, s: str) -> bool:
        
        l = 0
        r = len(s)-1

        s = s.lower()

        while l < r:
            if s[l].isalnum() == False:
                l += 1
                continue

            if s[r].isalnum() == False:
                r -= 1
                continue

            if s[l] != s[r]:
                return False
            
            l += 1
            r -= 1
        
        return True

if __name__ == "__main__":
    solution = Solution()

    assert solution.isPalindrome("A man, a plan, a canal: Panama") == True
    assert solution.isPalindrome("race a car") == False
    assert solution.isPalindrome(" ") == True




