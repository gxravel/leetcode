func shuffle(nums []int, n int) []int {
	result := make([]int, n*2)
	for i := 0; i < n; i++ {
		var t = 2 * i
		result[t] = nums[i]
		result[t+1] = nums[n+i]
	}
	return result
}