package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 给你一个整数数组arr，你一开始再数组的第一个元素（下表为0），
// 每一步你可以从下表i调到下标i+1、i-1或者j。
// j 需满足：arr[i] == arr[j] 且 i != j
// 请你返回达到数组最后一个元素的下表处所需的操作次数。

type arrInfo struct {
	arr      []int
	arrlen   int
	amap     map[int][]int
	stepmap  map[int]int
	stepnext map[int]int
}

func (ai *arrInfo) arrToMap() map[int][]int {
	aMap := make(map[int][]int)
	for i, v := range ai.arr {
		if aMap[v] == nil {
			aMap[v] = make([]int, 0, 0)
		}
		aMap[v] = append(aMap[v], i)
	}
	return aMap
}

// 计算第n步所需步长，vector 表示到达第n步动作
// 当vector 是 -1时，下一步无需+1
// 当vector 是 +1时，下一步无需-1
// 当vector 是 jump 时，下一步可以+1 -1 jump

func (ai *arrInfo) turnRight(n int) int {
	return ai.step(n+1) + 1
}

// n :当前跳跃点
// ori : 跳跃原点
func (ai *arrInfo) turnJump(n, ori int) (int, int) {
	minJump := ai.arrlen
	jumpNext := -1
	for _, v := range ai.amap[ai.arr[n]] {
		// 当-1 操作时，跳跃值比原点小，不是最优解
		if v > n {
			if v <= ori {
				return minJump, jumpNext
			}
			jump := ai.step(v)
			if jump < minJump {
				minJump = jump
				jumpNext = v
			}
		}

	}

	return minJump + 1, jumpNext
}

func (ai *arrInfo) step(n int) int {
	if n < ai.arrlen-2 && ai.stepmap[n] == 0 {

		// 初始化，step n 不应该比数组长度长
		// 计算+1 后需要的最小步数
		right := ai.turnRight(n)
		// 计算跳跃后需要的最小步数
		jump, jumpNext := ai.turnJump(n, 0)

		// 计算+1和跳跃的最小值

		var tmpmin, minNext int
		if jumpNext == -1 || jump > right {
			minNext = n + 1
			tmpmin = right
		} else {
			tmpmin = jump
			minNext = jumpNext
		}

		// 计算-1 后需要的最小步数
		for i := 1; i < tmpmin && n-i > 0; i++ {
			tmpjump, jumpNext2 := ai.turnJump(n-i, n)
			tmpjump += i
			if tmpjump < tmpmin {
				tmpmin = tmpjump
				minNext = jumpNext2
			}
		}

		ai.stepmap[n] = tmpmin
		ai.stepnext[n] = minNext

	}
	return ai.stepmap[n]
}

func count(arr []int) int {
	ai := &arrInfo{}
	ai.arr = arr
	ai.arrlen = len(arr)
	ai.amap = ai.arrToMap()
	ai.stepmap = make(map[int]int)
	ai.stepnext = make(map[int]int)
	// 下角标为n-1
	ai.stepmap[ai.arrlen-1] = 0
	if ai.arrlen > 1 {
		ai.stepmap[ai.arrlen-2] = 1
		ai.stepnext[ai.arrlen-2] = ai.arrlen - 1

	}

	res := ai.step(0)
	log.Println("route =", ai.stepnext)
	fmt.Printf("Route :arr[0]:%v ", ai.arr[0])
	n := 0
	for {
		next := ai.stepnext[n]
		fmt.Printf("arr[%v]:%v ", next, ai.arr[next])
		n = next
		if next == ai.arrlen-1 {
			break
		}
	}
	fmt.Println()
	return res
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	n := 10000
	var arr []int = make([]int, n)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
	}
	log.Println(arr)
	log.Println("result =", count(arr))
}
