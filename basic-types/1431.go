// 1431. Kids With the Greatest Number of Candies
package main

func kidsWithCandies(candies []int, extraCandies int) []bool {
	var maxCandies int
	var result = make([]bool, len(candies))
	for _, candy := range candies {
		if maxCandies < candy {
			maxCandies = candy
		}
	}
	for i, candy := range candies {
		if candy+extraCandies >= maxCandies {
			result[i] = true
		}
	}
	return result
}
