package structure

import "playground/internal/log"

func LeastRecentlyUsed() {
  var cache LRU[int, int]
  cache.capacity = 5
  cache.entries = make(map[int]*Entry[int, int])

  cache.Put(0, 0)
  cache.Put(3, 9)
  cache.Put(12, 4)

  cache.print()
}

type Entry[K comparable, V any] struct {
  key        K
  value      V
  prev, next *Entry[K, V]
}

type LRU[K comparable, V any] struct {
  capacity   int
  entries    map[K]*Entry[K, V]
  head, tail *Entry[K, V]
}

func (this *LRU[K, V]) Get(key K) (V, bool) {
  if entry, found := this.entries[key]; found {
    this.updateMRU(entry) // Move entry to the head of the cache.
    return entry.value, true
  }

  var zero V
  return zero, false
}

func (this *LRU[K, V]) Put(key K, value V) {
  if entry, found := this.entries[key]; found {
    entry.value = value
    this.updateMRU(entry) // Assure entry is the head of the cache.
    return
  }

  var fresh Entry[K, V]
  fresh.key, fresh.value = key, value

  this.entries[key] = &fresh
  this.newMRU(&fresh)

  // Check the cache size once new head is added.
  if len(this.entries) > this.capacity {
    this.evict()
  }
}

func (this *LRU[K, V]) newMRU(entry *Entry[K, V]) {
  // Set entry before head and update head to point back if needed.
  entry.next = this.head
  if this.head != nil {
    this.head.prev = entry
  }

  // Set entry as head and set tail for first entry.
  this.head = entry
  if this.tail == nil {
    this.tail = entry
  }
}

func (this *LRU[K, V]) updateMRU(entry *Entry[K, V]) {
  if entry == this.head {
    return
  }

  this.remove(entry)
  this.newMRU(entry)
}

func (this *LRU[K, V]) remove(entry *Entry[K, V]) {
  // Point previous entry to next entry, else assume we're head.
  if entry.prev != nil {
    entry.prev.next = entry.next
  } else {
    this.head = entry.next
  }

  // Point next entry to previous entry, else assume we're tail.
  if entry.next != nil {
    entry.next.prev = entry.prev
  } else {
    this.tail = entry.prev
  }

  // Nil out entry for GC and safety.
  entry.prev, entry.next = nil, nil
}

func (this *LRU[K, V]) evict() {
  if this.tail == nil {
    return
  }

  delete(this.entries, this.tail.key)
  this.remove(this.tail)
}

func (this *LRU[K, V]) print() {
  for entry := this.head; entry != nil; entry = entry.next {
    log.Info("%v:%v", entry.key, entry.value)
  }

  log.Info("capacity: %d/%d", len(this.entries), this.capacity)
}
