// 1302. Deepest Leaves Sum
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func getDeepestLevel(node *TreeNode, level int) int {
	if node == nil {
		return level - 1
	}
	left := getDeepestLevel(node.Left, level+1)
	right := getDeepestLevel(node.Right, level+1)
	if left > right {
		return left
	}
	return right
}

func walk(node *TreeNode, level, maxLevel int) int {
	if node == nil {
		return 0
	}
	if level == maxLevel {
		return node.Val
	}
	left := walk(node.Left, level+1, maxLevel)
	right := walk(node.Right, level+1, maxLevel)
	return left + right
}

func deepestLeavesSum(root *TreeNode) int {
	deepestLevel := getDeepestLevel(root, 0)
	return walk(root, 0, deepestLevel)
}

// -----------------------------------------------

// 1282. Group the People Given the Group Size They Belong To
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

// ----------------------------------------------