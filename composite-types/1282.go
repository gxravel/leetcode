// 1282. Group the People Given the Group Size They Belong To
package main

func groupThePeople(groupSizes []int) [][]int {
	var results = make([][]int, 0)
	var iGroup = make(map[int]int)
	for man, size := range groupSizes {
		_, ok := iGroup[size]
		if !ok || (ok && len(results[iGroup[size]]) == size) {
			results = append(results, make([]int, 0, size))
			iGroup[size] = len(results) - 1
		}
		results[iGroup[size]] = append(results[iGroup[size]], man)
	}
	return results
}
