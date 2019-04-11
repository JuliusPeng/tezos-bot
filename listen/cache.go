package listen

import (
	"container/list"
)

type cache struct {
	queue *list.List
}

func newCache() *cache {
	return &cache{
		queue: list.New(),
	}
}

func (c *cache) Has(item string) bool {
	current := c.queue.Front()
	for current != nil {
		v, _ := current.Value.(string)
		if item == v {
			return true
		}
		current = current.Next()
	}
	return false
}

func (c *cache) Add(item string) {
	if !c.Has(item) {
		if c.queue.Len() > 3 {
			c.queue.Remove(c.queue.Front())
		}
		c.queue.PushBack(item)
	}
}
