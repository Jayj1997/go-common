/*
 * @Author       : jayj
 * @Date         : 2021-09-09 10:49:30
 * @Description  :
 */
package btree

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Jayj1997/go-common/comparator"
)

/**
特性：
m阶的定义：一个结点能拥有的做大结点数来表示这棵树的阶数
e.g. 如果一个结点最多有n个key，那么这个结点最多就会有n+1个子节点，这棵树就叫做n+1(m=n+1)阶树

1. 每个结点x有如下特性：
	- x.n 标识当前存储在结点x中的关键字个数
	- x.n 的各个关键字本身: x.key1 x.key2 ... 以升序存放， 使得x.key1 <= x.key2 <= ...
	- x.leaf 是一个布尔值，如果x是叶子结点，则为true，如果x为内部节点，则为false
2. 每个'内部节点x'还包含x.n+1个指向它的孩子的指针 x.c1, x.c2...。叶子结点没有孩子结点，所以他的c1属性没有定义
	- key和指针相互间隔，结点两端是指针，所以结点中指针比key多一个
	- 每个指针要么为null，要么指向另一节点
3. 关键字 x.keyI 对存储在各子树中的关键字进行分割：如果ki为任意一个存储在以 x.ci为根的子树的关键字，那么：
k1 <= x.key1 <= k2  x.key2 <= ... <= x.keyX.n <= kx.n+1
难理解可以这么说：

> 如果某个指针在结点node最左边且不为null，则其指向结点的所有key小于(key1)，其中(key1)为node的第一个key的值
> 如果某个指针在结点node最右边且不为null，则其指向结点的所有key大于(keyM)，其中(keyM)为node的最后一个key的值
> 如果某个指针在结点node的左右相邻key分别是keyI和keyI+1且不为null，则其指向结点的所有key小于(keyI+1)且大于(keyI)

4. 每个叶子结点具有相同的深度，即树的高度h

5. 每个节点所包含的关键字个数有上界和下界。用一个被称作B树的最小度数的估计整数t(t>=2)来表示这些界：
除了根节点以外的每个结点必须**至少**有t-1个关键字。因此，除了根节点以外的内部结点**至少**有t个孩子。
（因为上面说了右x.n+1个指向它的孩子的指针）
如果树非空，根节点至少有一个关键字
每个结点**最多**包含2t-1个关键字。因为，一个内部结点至多可以有2t个孩子。当一个结点恰好有2t-1个关键字时，称该结点是满的(full)
*/

type Tree struct {
	Root       *Node                 // root node
	Comparator comparator.Comparator // key comparator
	size       int                   // total number of keys in the tree
	m          int                   // order (maximum number of children)
}

type Node struct {
	Parent   *Node
	Entries  []*Entry // contained keys in node
	Children []*Node  // children nodes
}

type Entry struct {
	Key   interface{}
	Value interface{}
}

// NewWith instantiates a B-tree with order(maximum number of children) and costom key comparator
func NewWith(order int, comparator comparator.Comparator) *Tree {
	if order < 3 {
		panic("Invalid order, should be at least 3")
	}

	return &Tree{m: order, Comparator: comparator}
}

// NewWithIntComparator instantiates a B-tree with the order (maximum number of children) and the IntComparator, i.e. keys are of type int.
func NewWithIntComparator(order int) *Tree {
	return NewWith(order, comparator.IntComparator)
}

// NewWithStringComparator instantiates a B-tree with the order (maximum number of children) and the StringComparator, i.e. keys are of type string.
func NewWithStringComparator(order int) *Tree {
	return NewWith(order, comparator.StringComparator)
}

/** function related */

// Insert inserts key-value pair node into the tree.
// If key already exists, then its value is updated with the new value
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Insert(key, value interface{}) {

	entry := &Entry{Key: key, Value: value}

	if tree.Root == nil {
		tree.Root = &Node{Entries: []*Entry{entry}, Children: []*Node{}}

		tree.size++

		return
	}

	if tree.insert(tree.Root, entry) {
		tree.size++
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree,
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Get(key interface{}) (value interface{}, found bool) {
	node, index, found := tree.searchRecursively(tree.Root, key)
	if found {
		return node.Entries[index].Value, true
	}

	return nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics
func (tree *Tree) Remove(key interface{}) {
	node, index, found := tree.searchRecursively(tree.Root, key)

	if found {
		tree.delete(node, index)
		tree.size--
	}
}

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

func (tree *Tree) Values() []interface{} {

	values := make([]interface{}, tree.size)
	it := tree.Iterator()

	for i := 0; it.Next(); i++ {

		values[i] = it.Value()
	}

	return values
}

// clear removes all nodes from the tree
func (tree *Tree) Clear() {

	tree.Root = nil
	tree.size = 0
}

// Height returns the height of the tree
func (tree *Tree) Height() int {

	return tree.Root.height()
}

// Left returns the left-most(min) node or nil if tree was empty
func (tree *Tree) Left() *Node {

	return tree.left(tree.Root)
}

// LeftKey returns the left-most(min) key or nil if tree was empty
func (tree *Tree) LeftKey() interface{} {

	if left := tree.Left(); left != nil {

		return left.Entries[0].Key
	}

	return nil
}

// LeftValue returns the left-most(min) value or nil if tree was empty
func (tree *Tree) LeftValue() interface{} {

	if left := tree.Left(); left != nil {

		return left.Entries[0].Value
	}

	return nil
}

// Right returns the right-most(max) node or nil if tree was empty
func (tree *Tree) Right() *Node {

	return tree.right(tree.Root)
}

// RightKey returns the right-most(max) key or nil if tree is empty
func (tree *Tree) RightKey() interface{} {

	if right := tree.Right(); right != nil {

		return right.Entries[len(right.Entries)-1].Key
	}

	return nil
}

// RightValue returns the right-most(max) value or nil if tree was empty
func (tree *Tree) RightValue() interface{} {

	if right := tree.Right(); right != nil {

		return right.Entries[len(right.Entries)-1].Value
	}

	return nil
}

// String returns a string representation of container (for debugging purposes)
func (tree *Tree) String() string {

	var buffer bytes.Buffer

	if _, err := buffer.WriteString("BTree\n"); err != nil {
		// ?
	}

	if !tree.Empty() {

		tree.output(&buffer, tree.Root, 0, true)
	}

	return buffer.String()
}

func (entry *Entry) String() string {

	return fmt.Sprintf("%v", entry.Key)
}

func (tree *Tree) output(buffer *bytes.Buffer, node *Node, level int, isTail bool) {

	for e := 0; e < len(node.Entries)+1; e++ {

		if e < len(node.Children) {

			tree.output(buffer, node.Children[e], level+1, true)
		}

		if e < len(node.Entries) {

			if _, err := buffer.WriteString(strings.Repeat("    ", level)); err != nil {
			}

			if _, err := buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n"); err != nil {
			}
		}
	}
}

/** inner function related */

func (node *Node) height() int {

	height := 0

	for ; node != nil; node = node.Children[0] {

		height++

		if len(node.Children) == 0 {
			break
		}
	}

	return height
}

func (tree *Tree) insert(node *Node, entry *Entry) (inserted bool) {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}

	return tree.insertIntoInternal(node, entry)
}

func (tree *Tree) insertIntoLeaf(node *Node, entry *Entry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)

	if found {
		node.Entries[insertPosition] = entry
		return false
	}

	// Insert entry's key in the middle of the node
	node.Entries = append(node.Entries, nil)

	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])

	node.Entries[insertPosition] = entry

	tree.split(node)

	return true
}

func (tree *Tree) insertIntoInternal(node *Node, entry *Entry) (inserted bool) {

	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}

	return tree.insert(node.Children[insertPosition], entry)
}

func (tree *Tree) isLeaf(node *Node) bool {
	return len(node.Children) == 0
}

func (tree *Tree) search(node *Node, key interface{}) (index int, found bool) {

	low, high := 0, len(node.Entries)-1

	var mid int

	// two pointer
	for low <= high {
		mid = (high + low) / 2

		compare := tree.Comparator(key, node.Entries[mid].Key)

		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}

	return low, false
}

// searchRecursively searches recursively down the tree starting at the startNode
func (tree *Tree) searchRecursively(startNode *Node, key interface{}) (node *Node, index int, found bool) {
	if tree.Empty() {
		return nil, -1, false
	}

	node = startNode

	for {
		index, found = tree.search(node, key)

		if found {
			return node, index, true
		}

		if tree.isLeaf(node) {
			return nil, -1, false
		}

		node = node.Children[index]
	}
}

func (tree *Tree) split(node *Node) {

	if !tree.shouldSplit(node) {
		return
	}

	if node == tree.Root {
		tree.splitRoot()
		return
	}

	tree.splitNonRoot(node)
}

func (tree *Tree) splitNonRoot(node *Node) {

	middle := tree.middle()
	parent := node.Parent

	left := &Node{Entries: append([]*Entry(nil), node.Entries[:middle]...), Parent: parent}
	right := &Node{Entries: append([]*Entry(nil), node.Entries[middle+1:]...), Parent: parent}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(node) {
		left.Children = append([]*Node(nil), node.Children[:middle+1]...)
		right.Children = append([]*Node(nil), node.Children[middle+1:]...)

		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)

	// Insert middle key into parent
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.Children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *Tree) shouldSplit(node *Node) bool {
	return len(node.Entries) > tree.maxEntries()
}

func (tree *Tree) minEntries() int {
	return tree.minChildren() - 1
}

func (tree *Tree) maxEntries() int {
	return tree.maxChildren() - 1
}

func (tree *Tree) minChildren() int {
	return (tree.m + 1) / 2
}

func (tree *Tree) maxChildren() int {
	return tree.m
}

func (tree *Tree) middle() int {
	return (tree.m - 1) / 2
}

func (tree *Tree) splitRoot() {

	middle := tree.middle()

	left := &Node{Entries: append([]*Entry(nil), tree.Root.Entries[:middle]...)}
	right := &Node{Entries: append([]*Entry(nil), tree.Root.Entries[middle+1:]...)}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(tree.Root) {
		left.Children = append([]*Node(nil), tree.Root.Children[:middle+1]...)
		right.Children = append([]*Node(nil), tree.Root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	// Root is a node with one entry and two children (left & right)
	newRoot := &Node{
		Entries:  []*Entry{tree.Root.Entries[middle]},
		Children: []*Node{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	tree.Root = newRoot
}

func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (tree *Tree) left(node *Node) *Node {

	if tree.Empty() {
		return nil
	}

	current := node

	for {
		if tree.isLeaf(current) {
			return current
		}

		current = current.Children[0]
	}
}

func (tree *Tree) right(node *Node) *Node {

	if tree.Empty() {
		return nil
	}

	current := node

	for {
		if tree.isLeaf(current) {
			return current
		}

		current = current.Children[len(current.Children)-1]
	}
}

// leftSibling returns the node's left sibling and child index (in parent) if it exists, otherwise (nil, -1)
// key is any of keys in node (could even be deleted)
func (tree *Tree) leftSibling(node *Node, key interface{}) (*Node, int) {

	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index--

		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}

	return nil, -1
}

// rightSibling returns the node's right sibling and child index (in parent) if it exists, otherwise (nil, -1)
// key is any of keys in node (could even be deleted)
func (tree *Tree) rightSibling(node *Node, key interface{}) (*Node, int) {

	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index++

		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}

	return nil, -1
}

// delete deletes an entry in node at entries' index
// ref: http://en.wikipedia.org/wiki/B-tree#Deletion
func (tree *Tree) delete(node *Node, index int) {

	// deleting from a leaf node
	if tree.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		tree.deleteEntry(node, index)
		tree.rebalance(node, deletedKey)

		if len(tree.Root.Entries) == 0 {
			tree.Root = nil
		}

		return
	}

	// deleting from on internal node
	leftLargestNode := tree.right(node.Children[index]) // largest node in the left sub-tree(assumed to exists)
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	tree.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	tree.rebalance(leftLargestNode, deletedKey)
}

// rebalance rebalances the tree after deletion if necessary and returns true, otherwise false.
// Note the we first delete the entry and then call rebalance, thus the passed deleted key as reference
func (tree *Tree) rebalance(node *Node, deletedKey interface{}) {

	// check if rebalancing is needed
	if node == nil || len(node.Entries) >= tree.minEntries() {
		return
	}

	// try to borrow from left sibling
	leftSibling, leftSiblingIndex := tree.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > tree.minEntries() {

		// rotate right
		// prepend parent's separator entry to node's entries
		node.Entries = append([]*Entry{node.Parent.Entries[leftSiblingIndex]}, node.Entries...)
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		tree.deleteEntry(leftSibling, len(leftSibling.Entries)-1)

		if !tree.isLeaf(leftSibling) {

			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node{leftSiblingRightMostChild}, node.Children...)
			tree.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}

		return
	}

	// try to borrow from right sibling
	rightSibling, rightSiblingIndex := tree.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > tree.minEntries() {

		// rorate left
		// append parent's separator entry to node's entries
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		tree.deleteEntry(rightSibling, 0)

		if !tree.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			tree.deleteChild(rightSibling, 0)
		}

		return
	}

	// merge with siblings
	if rightSibling != nil {

		// merge with right sibling
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		tree.deleteEntry(node.Parent, rightSiblingIndex-1)
		tree.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		tree.deleteChild(node.Parent, rightSiblingIndex)

	} else if leftSibling != nil {

		// merge with left sibling
		entries := append([]*Entry(nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		tree.deleteEntry(node.Parent, leftSiblingIndex)
		tree.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		tree.deleteChild(node.Parent, leftSiblingIndex)
	}

	// make the merged node the root if its parent was the root and the root is empty
	if node.Parent == tree.Root && len(tree.Root.Entries) == 0 {

		tree.Root = node
		node.Parent = nil

		return
	}

	// parent might underflow, so try to rebalance if necessary
	tree.rebalance(node.Parent, deletedKey)

}

func (tree *Tree) prependChildren(fromNode *Node, toNode *Node) {

	children := append([]*Node(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) appendChildren(fromNode *Node, toNode *Node) {

	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) deleteEntry(node *Node, index int) {

	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (tree *Tree) deleteChild(node *Node, index int) {

	if index >= len(node.Children) {
		return
	}

	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}
