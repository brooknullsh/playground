package structure

import (
  "cmp"
  "playground/internal/log"
)

func BinarySearchTree() {
  var tree BSTree[int]

  tree.Insert(10)
  tree.Insert(2)
  tree.Insert(30)
  tree.Insert(14)
  tree.Insert(14)
  tree.Insert(1)
  tree.print()
}

type BSTNode[T cmp.Ordered] struct {
  Value       T
  left, right *BSTNode[T]
}

type BSTree[T cmp.Ordered] struct {
  root *BSTNode[T]
}

func (this *BSTree[T]) Insert(value T) {
  var node BSTNode[T]
  node.Value = value

  if this.root == nil {
    log.Debug("root: %v", value)
    this.root = &node
    return
  }

  curr := this.root
  for {
    // Avoid duplicates, otherwise duplicates will always be inserted at the
    // bottom of the "greater" path.
    if value == curr.Value {
      log.Debug("duplicate: %v", value)
      return
    }

    if value < curr.Value {
      // Bottom of the "lesser" path.
      if curr.left == nil {
        curr.left = &node
        log.Debug("left insert: %v child of: %v", value, curr.Value)
        return
      }

      // Continue down the "lesser" path.
      curr = curr.left
      continue
    }

    // Bottom of the "greater" path.
    if curr.right == nil {
      log.Debug("right insert: %v child of: %v", value, curr.Value)
      curr.right = &node
      return
    }

    // Continue down the "greater" path.
    curr = curr.right
  }
}

func (this *BSTree[T]) print() {
  stack := make([]*BSTNode[T], 0)

  // 1. Populate stack from the "lesser" path from root to bottom.
  // 2. Returning from bottom, print & pop node from stack.
  // 3. If node has "greater" children, repeat above.
  curr := this.root
  for curr != nil || len(stack) > 0 {
    // Walk the "lesser" path from the node.
    for curr != nil {
      stack, curr = append(stack, curr), curr.left
    }

    idx := len(stack) - 1                 // Most recent node.
    curr, stack = stack[idx], stack[:idx] // Pop the stack.

    log.Info("%v", curr.Value)
    curr = curr.right // Walk the "greater" path from the node.
  }
}
