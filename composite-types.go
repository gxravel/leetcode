package main

import "fmt"

// 1302. Deepest Leaves Sum
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func getMaxLevel(node *TreeNode, level int) int {
	if node == nil {
		return level - 1
	}
	left := getMaxLevel(node.Left, level+1)
	right := getMaxLevel(node.Right, level+1)
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
		fmt.Println(node.Val)
		return node.Val
	}
	left := walk(node.Left, level+1, maxLevel)
	right := walk(node.Right, level+1, maxLevel)
	return left + right
}

func deepestLeavesSum(root *TreeNode) int {
	deepestLevel := getMaxLevel(root, 0)
	return walk(root, 0, deepestLevel)
}

func main() {
	root := &TreeNode{1, &TreeNode{2, &TreeNode{4, &TreeNode{7, nil, nil}, nil}, &TreeNode{5, nil, nil}}, &TreeNode{3, nil, &TreeNode{6, nil, &TreeNode{8, nil, nil}}}}
	fmt.Println(deepestLeavesSum(root))
}
