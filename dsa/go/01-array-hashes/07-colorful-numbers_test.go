package array_hashes

import "testing"

// isColorful checks if number is colorful: all products of contiguous subsequences are unique.
func isColorful(number int) bool {
    numStr := []byte((func(n int) string { return fmtInt(n) })(number))
    products := make(map[int]bool)
    for i := 0; i < len(numStr); i++ {
        prod := 1
        for j := i; j < len(numStr); j++ {
            prod *= int(numStr[j] - '0')
            if products[prod] { return false }
            products[prod] = true
        }
    }
    return true
}

// fmtInt converts int to string without importing fmt to keep dependencies minimal.
func fmtInt(n int) string {
    if n == 0 { return "0" }
    neg := n < 0
    if neg { n = -n }
    buf := make([]byte, 0, 20)
    for n > 0 { buf = append(buf, byte('0'+n%10)); n /= 10 }
    // reverse
    for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 { buf[i], buf[j] = buf[j], buf[i] }
    if neg { return "-" + string(buf) }
    return string(buf)
}

func TestIsColorful(t *testing.T) {
    if isColorful(326) != false { t.Fatalf("isColorful(326) should be false") }
    if isColorful(3245) != true { t.Fatalf("isColorful(3245) should be true") }
}

