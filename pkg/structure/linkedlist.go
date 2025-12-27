package structure

import "playground/internal/log"

func LinkedList() {
  var list List[int]

  list.Append(10)
  list.Prepend(20)
  list.print()
}

type LLNode[T any] struct {
  Value T
  next  *LLNode[T]
}

type List[T any] struct {
  head, tail *LLNode[T]
  len        int
}

func (this *List[T]) Prepend(value T) {
  var node LLNode[T]
  node.Value = value
  node.next = this.head // Nil if empty, current head otherwise.

  // Align head and tail if empty.
  if this.tail == nil {
    this.tail = &node
  }

  this.head = &node
  this.grow()

  log.Debug("prepended: %v", value)
}

func (this *List[T]) Append(value T) {
  var node LLNode[T]
  node.Value = value
  node.next = nil // NOTE: Explicit nil for new tail.

  if this.tail != nil {
    // Give current tail next node.
    this.tail.next = &node
  } else {
    // If tail is nil, head must be too so align them.
    this.head = &node
  }

  this.tail = &node
  this.grow()

  log.Debug("appended: %v", value)
}

func (this *List[T]) grow() {
  this.len++
}

func (this *List[T]) print() {
  for node := this.head; node != nil; node = node.next {
    log.Info("value: %v", node.Value)
  }

  log.Info("len: %d", this.len)
}
