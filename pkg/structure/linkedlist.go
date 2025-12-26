package structure

import (
  "playground/internal/log"
)

func LinkedList() {
  var list List[int]

  list.Append(10)
  list.Prepend(20)

  list.print()
}

type Node[T any] struct {
  Value T
  next  *Node[T]
}

type List[T any] struct {
  head, tail *Node[T]
  len        int
}

func (this *List[T]) Grow() {
  this.len++
}

func (this *List[T]) Prepend(value T) {
  var node Node[T]
  node.Value = value
  node.next = this.head // Nil if empty, current head otherwise.

  this.head = &node
  // Align head and tail if empty.
  if this.tail == nil {
    this.tail = &node
  }

  log.Debug("prepended: %v", value)
  this.Grow()
}

func (this *List[T]) Append(value T) {
  var node Node[T]
  node.Value = value
  node.next = nil // NOTE: Explicit nil for new tail.

  if this.tail != nil {
    // Give current tail next node.
    this.tail.next = &node
  } else {
    // If tail is nil, head must be too so align them.
    this.head = &node
  }

  log.Debug("appended: %v", value)
  this.tail = &node
  this.Grow()
}

func (this *List[T]) print() {
  for node := this.head; node != nil; node = node.next {
    log.Info("value: %v", node.Value)
  }

  log.Info("len: %d", this.len)
}
