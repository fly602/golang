package main

type arrInfo struct {
	arr     []int
	arrlen  int
	amap    map[int][]int
	stepmap map[int]int
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

func (ai *arrInfo) turnRight(n int) int {
	return ai.step(n+1) + 1
}

func (ai *arrInfo) turnJump(n, ori int) int {
	minJump := ai.arrlen
	for _, v := range ai.amap[ai.arr[n]] {
		// 当-1 操作时，跳跃值比原点小，不是最优解
		if v > n {
			if v <= ori {
				return minJump
			}
			jump := ai.step(v)
			if jump < minJump {
				minJump = jump
			}
		}

	}

	return minJump + 1
}

func (ai *arrInfo) step(n int) int {
	if n < ai.arrlen-2 && ai.stepmap[n] == 0 {

		// 初始化，step n 不应该比数组长度长
		// 计算+1 后需要的最小步数
		right := ai.turnRight(n)
		// 计算跳跃后需要的最小步数
		jump := ai.turnJump(n, 0)

		// 计算+1和跳跃的最小值

		var tmpmin int
		if jump > right {
			tmpmin = right
		} else {
			tmpmin = jump
		}

		// 计算-1 后需要的最小步数
		for i := 1; i < tmpmin && n-i > 0; i++ {
			tmpjump := ai.turnJump(n-i, n)
			tmpjump += i
			if tmpjump < tmpmin {
				tmpmin = tmpjump
			}
		}

		ai.stepmap[n] = tmpmin
	}
	return ai.stepmap[n]
}

func count(arr []int) int {
	ai := &arrInfo{}
	ai.arr = arr
	ai.arrlen = len(arr)
	ai.amap = ai.arrToMap()
	ai.stepmap = make(map[int]int)
	// 下角标为n-1
	ai.stepmap[ai.arrlen-1] = 0
	if ai.arrlen > 1 {
		ai.stepmap[ai.arrlen-2] = 1

	}

	return ai.step(0)
}
