// 1302. Deepest Leaves Sum
package main

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
