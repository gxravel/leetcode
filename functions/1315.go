/*
1315. Sum of Nodes with Even-Valued Grandparent

Given a binary tree, return the ch of values of nodes with even-valued grandparent.
(A grandparent of a node is the parent of its parent, if it exists.)

If there are no nodes with an even-valued grandparent, return 0.
*/
package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func sumEvenGrandparent(root *TreeNode) int {
	type Even struct {
		prev, curr bool
	}

	var walk func(*TreeNode, Even, chan int)

	t := root

	walk = func(root *TreeNode, even Even, ch chan int) {
		if root == nil {
			return
		}
		if even.prev {
			ch <- root.Val
		}
		even.prev = even.curr
		even.curr = root.Val%2 == 0
		walk(root.Left, even, ch)
		walk(root.Right, even, ch)
		if t == root {
			close(ch)
		}
	}

	ch := make(chan int)

	go walk(root, Even{}, ch)

	var result int

	for val := range ch {
		result += val
	}

	return result
}

func main() {
	var root = &TreeNode{6, &TreeNode{7, &TreeNode{2, &TreeNode{9, nil, nil}, nil}, &TreeNode{7, &TreeNode{1, nil, nil}, &TreeNode{4, nil, nil}}}, &TreeNode{8, &TreeNode{1, nil, nil}, &TreeNode{3, nil, &TreeNode{5, nil, nil}}}}
	fmt.Println(sumEvenGrandparent(root))
}
