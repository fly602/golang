package luoxuansqu

import (
	"fmt"
	"log"
)

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func NewSqu(imax, jmax int) [][]int {
	var a = make([][]int, imax)
	for i := 0; i < imax; i++ {
		a[i] = make([]int, jmax)
		for j := 0; j < jmax; j++ {
			a[i][j] = i*jmax + j
		}
	}
	log.Println("New squ=", a)
	return a
}

func PrintSqu(squ [][]int) {
	jmax := len(squ[0]) - 1 //多少列 -1
	imax := len(squ) - 1    //多少行 -1
	log.Println("imax,jmax=", imax, jmax)
	var i, j int
	// 第n圈
	for n := 0; n <= min(imax, jmax); n++ {
		// 打印第一行
		fmt.Printf("第%v圈第一行：", n)
		for j = n; j < jmax-n; j++ {
			fmt.Printf("%v ", squ[n][j])
		}
		fmt.Printf("\n第%v圈最后一列：", n)
		// 打印最后一列
		for i = n; i < imax-n; i++ {
			fmt.Printf("%v ", squ[i][jmax-n])
		}
		fmt.Printf("\n第%v圈最后一行：", n)
		// 打印最后一行
		for j = jmax - n; j > n; j-- {
			fmt.Printf("%v ", squ[imax-n][j])
		}
		fmt.Printf("\n第%v圈第一列：", n)
		// 打印第一列
		for i = imax - n; i > n; i-- {
			fmt.Printf("%v ", squ[i][n])
		}
		fmt.Println("")
	}
}
