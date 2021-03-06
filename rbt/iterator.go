/*
 * @Author       : jayj
 * @Date         : 2021-09-06 15:31:45
 * @Description  :
 */
package rbt

type Iterator struct {
	tree     *Tree
	node     *Node
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *Tree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

// IteratorAt returns a stateful iterator whose elements are key/value pairs that is initialized at a particular node.
func (tree *Tree) IteratorAt(node *Node) Iterator {
	return Iterator{tree: tree, node: node, position: between}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// if next() returns true, the next element's key and value can be retrieved by Key() and Value()
// if next() was called for the first time, then it will point the iterator to the first element if it exists
// Modifies the state of the iterator
func (iterator *Iterator) Next() bool {

	if iterator.position == end {
		goto end
	}

	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}

		iterator.node = left
		goto between
	}

	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right

		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}

		goto between
	}

	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if iterator.tree.Comparator(node.Key, iterator.node.Key) <= 0 {
				goto between
			}
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false
between:
	iterator.position = between
	return true
}

// Previous moves the iterator to the previous element and returns true if there was a previous element in the container.
// if Previous() returns true, the previous element's key and value can be retrieved by Key() and Value()
func (iterator *Iterator) Previous() bool {

	if iterator.position == begin {
		goto begin
	}

	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}

		iterator.node = right
		goto between
	}

	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left

		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}

		goto between
	}

	if iterator.node.Parent != nil {
		node := iterator.node

		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent

			if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
				goto between
			}
		}
	}

begin:
	iterator.node = nil
	iterator.position = begin
	return false
between:
	iterator.position = between
	return true
}

// Value returns the current element's value.
// Dose not modify the state of the iterator
func (iterator *Iterator) Value() interface{} {
	return iterator.node.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator *Iterator) Key() interface{} {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any
func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Previous() to fetch the last element if any
func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, the first element's key and value can be retrieved by Key() and Value()
// Modifies the state of the iterator
func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value()
// Modifies the state of the iterator
func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Previous()
}
