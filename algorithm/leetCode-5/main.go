package main

import "log"

func checkRevetStr(val string) string {
	bstr := []byte(val)

	var n, m int // 回文中心位置
	// 以n为中心，向两边扩散
	var i int
	var maxlen int = 0
	var start, end int
	for n = 0; n < len(val); n++ {
		for i = 1; i < len(val)-n; i++ {
			if bstr[n] != bstr[n+i] {
				break
			}
		}
		m = n + i - 1
		for i = 1; i < len(val)-n && n >= i; i++ {
			if (m+i >= len(val)) || (bstr[n-i] != bstr[m+i]) {
				break
			}
		}

		if m+i-1-(n-i+1) > maxlen {
			maxlen = m + i - 1 - (n - i + 1)
			start = n - i + 1
			end = m + i - 1
		}
	}
	//log.Println("Get revert string", string(bstr[start:end+1]))
	if maxlen > 0 {
		return string(bstr[start : end+1])
	}
	return ""
}

const (
	str1 = "facbcadef"
	str2 = "fcbcc"
	str3 = "bacc"
	str4 = "ddfccaaaaabbbbbccccc"
)

func main() {
	log.Println(checkRevetStr(str1))
	log.Println(checkRevetStr(str2))
	log.Println(checkRevetStr(str3))
	log.Println(checkRevetStr(str4))

}
