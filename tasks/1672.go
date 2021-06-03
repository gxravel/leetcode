func maximumWealth(accounts [][]int) int {
	var maxWealth int
	for _, banks := range accounts {
		var wealth int
		for _, money := range banks {
			wealth += money
		}
		if maxWealth < wealth {
			maxWealth = wealth
		}
	}
	return maxWealth
}