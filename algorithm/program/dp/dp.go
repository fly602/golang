package dp

import "fmt"

/*
动态分布算法题：
有100层楼梯，一次可以上1步或者2步，求共有多少种走法
每一次有2中走法，f(0)=1,f(1)=1
f(n)= f(n-1) + f(n-2)
f(n-1)= f(n-2)+f(n-3)
*/

// 使用递归
func step(n int) int {
	fmt.Println("get n=", n)
	if n == 0 || n == 1 {
		return 1
	}
	return step(n-1) + step(n-2)
}

/*
使用动态分布算法
将已经计算过的结果存储到内存中，
下次用到直接从内存中读取，
不用重复计算，提高效率
*/
var arr []int

func setStep(n int) {
	if arr[n] == 0 {
		arr[n] = getStep(n)
	}
}
func getStep(n int) int {
	if n == 0 || n == 1 {
		arr[n] = 1
		return arr[n]
	}
	setStep(n - 1)
	setStep(n - 2)
	return arr[n-1] + arr[n-2]
}

func dp(n int) int {
	arr = make([]int, n)
	return getStep(n)
}
