package main

func abs(num int) int {
	if num > 0 {
		return num
	}
	return -num
}

func minOperations(boxes string) []int {
	var result = make([]int, len(boxes))
	for i := range boxes {
		for j, box := range boxes {
			if box == '1' {
				result[i] += abs(j - i)
			}
		}
	}
	return result
}
