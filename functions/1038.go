/*
  1038. Binary Search Tree to Greater Sum Tree

  Given the root of a Binary Search Tree (BST), convert it to a Greater Tree
  such that every key of the original BST is changed to the original key
  plus sum of all keys greater than the original key in BST.

*/
package main

func bstToGst(root *TreeNode) *TreeNode {
	if root.Right != nil {
		bstToGst(root.Right)
		if root.Right.Left != nil {
			root.Val += root.Right.Left.Val
		} else {
			root.Val += root.Right.Val
		}
	}
	if root.Left != nil {
		bstToGst(root.Left)
		root.Left.Val += root.Val
	}

	return root
}
