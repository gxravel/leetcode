/*
  1038. Binary Search Tree to Greater Sum Tree

  Given the root of a Binary Search Tree (BST), convert it to a Greater Tree
  such that every key of the original BST is changed to the original key
  plus sum of all keys greater than the original key in BST.

*/
package main

func bstToGst(root *TreeNode) *TreeNode {
	var walk func(*TreeNode, int, bool)
	walk = func(root *TreeNode, upperRight int, wasLeft bool) {
		if root.Right != nil {
			walk(root.Right, upperRight, wasLeft)
			var r *TreeNode
			for r = root.Right; r.Left != nil; r = r.Left {
			}
			root.Val += r.Val
		} else if wasLeft {
			root.Val += upperRight
		}

		if root.Left != nil {
			walk(root.Left, root.Val, true)
		}
	}

	walk(root, 0, false)

	return root
}
