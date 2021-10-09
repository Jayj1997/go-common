/*
 * @Author       : jayj
 * @Date         : 2021-09-02 16:49:29
 * @Description  :
 */
package rbt

import (
	"fmt"

	"github.com/Jayj1997/go-common/comparator"
)

// Time Complexity: O(lgn)

// 红黑树的特性
// 1. 每个结点或者是黑色 或者是红色
// 2. 根结点是黑色
// 3. 每个叶子结点(NIL)是黑色。注意：这里叶子结点，是指为空的叶子结点
// 4. 如果一个结点是红色的，则它的子结点必须是黑色的
// 5. 从一个结点到该结点的子孙结点的所有路径上包含相同数目的黑结点
// *. 红黑树同时是一个二叉查找树，左右旋转之后仍然是二叉查找树

// 定理：一颗含有n个结点的红黑树的高度至多为2log(n+1)

type color bool

const (
	black, red color = true, false
)

// tree holds elements of the red-black tree
type Tree struct {
	Root       *Node
	size       int
	Comparator comparator.Comparator
}

type Node struct {
	Key    interface{}
	Value  interface{}
	color  color
	Left   *Node
	Right  *Node
	Parent *Node
}

// NewWith instantiates a red-black tree with the IntComparator,
// i.e. keys are of type string.
func NewWith(comparator comparator.Comparator) *Tree {
	return &Tree{Comparator: comparator}
}

// NewWithIntComparator instantiates a red-black tree with IntComparator,
// i.e. keys are of type int.
func NewWithIntComparator() *Tree {
	return &Tree{Comparator: comparator.IntComparator}
}

// NewWithStringComparator instantiates a red-black tree with the StringComparator,
// i.e. keys are of type string.
func NewWithStringComparator() *Tree {
	return &Tree{Comparator: comparator.StringComparator}
}

/** function related */

// Insert inserts node into the tree
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Insert(key, value interface{}) {

	var insertedNode *Node

	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		tree.Comparator(key, key)
		tree.Root = &Node{Key: key, Value: value, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true

		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				// overwrite
				node.Key = key
				node.Value = value
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &Node{Key: key, Value: value, color: red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &Node{Key: key, Value: value, color: red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}

	tree.insertCase1(insertedNode)
	tree.size++
}

// Get searchs the node in the tree by key and returns its value or nil if key is not found in tree,
// Second return parameter is true if key was found, otherwise false
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Get(key interface{}) (value interface{}, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Value, true
	}

	return nil, false
}

// Remove remove the node from the tree by key
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key interface{}) {
	var child *Node
	node := tree.lookup(key)
	if node == nil {
		return
	}

	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Value = pred.Value
		node = pred // ???
	}

	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}

		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}

		tree.replaceNode(node, child)

		if node.Parent == nil && child != nil {
			child.color = black
		}
	}

	tree.size--
}

// Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree
func (tree *Tree) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *Tree) Keys() []interface{} {
	keys := make([]interface{}, tree.size)

	it := tree.Iterator()

	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}

	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree) Values() []interface{} {
	values := make([]interface{}, tree.size)

	it := tree.Iterator()

	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}

	return values
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or all nodes in the tree are
// larger than the given node
//
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Floor(key interface{}) (floor *Node, found bool) {
	found = false

	node := tree.Root

	for node != nil {
		compare := tree.Comparator(key, node.Key)

		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor, found = node, true
			node = node.Right
		}
	}

	if found {
		return floor, true
	}

	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if node ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node
// A ceiling node may not be found, either because the tree is empty, or all nodes in the tree
// are smaller than the given node
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Ceiling(key interface{}) (ceiling *Node, found bool) {

	found = false

	node := tree.Root

	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}

	if found {
		return ceiling, true
	}

	return nil, false
}

// Clear removes all nodes from the tree
func (tree *Tree) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String() string {
	str := "REDBLACKTREE\n"

	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}

	return str
}

func (node *Node) String() string {
	return fmt.Sprintf("%v", node.Key)
}

/** inner function related */

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "|   "
		} else {
			newPrefix += "    "
		}

		output(node.Right, newPrefix, false, str)
	}

	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}

	*str += node.String() + "\n"

	if node.Left != nil {
		newPrefix := prefix

		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "|   "
		}

		output(node.Left, newPrefix, true, str)
	}
}

func (tree *Tree) lookup(key interface{}) *Node {

	node := tree.Root

	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

// get maximum Node
func (node *Node) maximumNode() *Node {
	if node == nil {
		return nil
	}

	for node.Right != nil {
		node = node.Right
	}

	return node
}

// those cases use to balance tree

func (tree *Tree) insertCase1(node *Node) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *Tree) insertCase2(node *Node) {
	if nodeColor(node.Parent) == black {
		return
	}

	tree.insertCase3(node)
}

func (tree *Tree) insertCase3(node *Node) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *Tree) insertCase4(node *Node) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.leftRotate(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rightRotate(node.Parent)
		node = node.Right
	}

	tree.insertCase5(node)
}

func (tree *Tree) insertCase5(node *Node) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rightRotate(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.leftRotate(grandparent)
	}
}

func (tree *Tree) deleteCase1(node *Node) {
	if node.Parent == nil {
		return
	}

	tree.deleteCase2(node)
}

func (tree *Tree) deleteCase2(node *Node) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black

		if node == node.Parent.Left {
			tree.leftRotate(node.Parent)
		} else {
			tree.rightRotate(node.Parent)
		}
	}

	tree.deleteCase3(node)
}

func (tree *Tree) deleteCase3(node *Node) {

	sibling := node.sibling()

	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *Tree) deleteCase4(node *Node) {
	sibling := node.sibling()

	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *Tree) deleteCase5(node *Node) {
	sibling := node.sibling()

	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		tree.rightRotate(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		tree.leftRotate(sibling)
	}

	tree.deleteCase6(node)
}

func (tree *Tree) deleteCase6(node *Node) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		tree.leftRotate(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		tree.rightRotate(node.Parent)
	}
}

/** direction related */

// Left returns the left-most (min) node or nil it tree is empty
func (tree *Tree) Left() *Node {
	var parent *Node

	current := tree.Root

	for current != nil {
		parent = current
		current = current.Left
	}

	return parent
}

// Right returns the right-most (max) node or nil if tree was empty
func (tree *Tree) Right() *Node {
	var parent *Node

	current := tree.Root

	for current != nil {
		parent = current
		current = current.Right
	}

	return parent
}

func (tree *Tree) leftRotate(node *Node) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *Tree) rightRotate(node *Node) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}

	left.Right = node
	node.Parent = left
}

func (tree *Tree) replaceNode(old *Node, new *Node) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}

	if new != nil {
		new.Parent = old.Parent
	}
}

/** relationship related */
func (node *Node) grandparent() *Node {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}

	return nil
}

func (node *Node) uncle() *Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}

	return node.Parent.sibling()
}

func (node *Node) sibling() *Node {
	if node == nil || node.Parent == nil {
		return nil
	}

	if node == node.Parent.Left {
		return node.Parent.Right
	}

	return node.Parent.Left
}

/** color related */

func nodeColor(node *Node) color {
	if node == nil {
		return black
	}

	return node.color
}
