// Copyright 2015, Hu Keping. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbtree

type Iterator func(i Item) bool

// Ascend will call iterator once for each element greater or equal than pivot
// in ascending order. It will stop whenever the iterator returns false.
func (t *Rbtree) Ascend(pivot Item, iterator Iterator) {
	t.ascend(t.root, pivot, iterator)
}

func (t *Rbtree) ascend(x *Node, pivot Item, iterator Iterator) bool {
	if x == t.NIL {
		return true
	}

	if !less(x.Item, pivot) {
		if !t.ascend(x.Left, pivot, iterator) {
			return false
		}
		if !iterator(x.Item) {
			return false
		}
	}

	return t.ascend(x.Right, pivot, iterator)
}

// Descend will call iterator once for each element less or equal than pivot
// in descending order. It will stop whenever the iterator returns false.
func (t *Rbtree) Descend(pivot Item, iterator Iterator) {
	t.descend(t.root, pivot, iterator)
}

func (t *Rbtree) descend(x *Node, pivot Item, iterator Iterator) bool {
	if x == t.NIL {
		return true
	}

	if !less(pivot, x.Item) {
		if !t.descend(x.Right, pivot, iterator) {
			return false
		}
		if !iterator(x.Item) {
			return false
		}
	}

	return t.descend(x.Left, pivot, iterator)
}

// AscendRange will call iterator once for elements greater or equal than @ge
// and less than @lt, which means the range would be [ge, lt).
// It will stop whenever the iterator returns false.
func (t *Rbtree) AscendRange(ge, lt Item, iterator Iterator) {
	t.ascendRange(t.root, ge, lt, iterator)
}

func (t *Rbtree) ascendRange(x *Node, inf, sup Item, iterator Iterator) bool {
	if x == t.NIL {
		return true
	}

	if !less(x.Item, sup) {
		return t.ascendRange(x.Left, inf, sup, iterator)
	}
	if less(x.Item, inf) {
		return t.ascendRange(x.Right, inf, sup, iterator)
	}

	if !t.ascendRange(x.Left, inf, sup, iterator) {
		return false
	}
	if !iterator(x.Item) {
		return false
	}
	return t.ascendRange(x.Right, inf, sup, iterator)
}

func (t *Rbtree) Next(x *Node) *Node {
	return t.successor(x)
}

func (t *Rbtree) Prev(x *Node) *Node {
	return t.desuccessor(x)

}

func (t *Rbtree) First() *Node {
	return t.min(t.root)
}
