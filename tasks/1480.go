func runningSum(nums []int) []int {
	result := make([]int, 0, len(nums))
	var acc int
	for _, val := range nums {
		acc += val
		result = append(result, acc)
	}
	return result
}