package lru

import "container/list"

type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

type OnEliminated func(key string, value Value)

type Cache struct {
	capacity         int64
	length           int64
	hashmap          map[string]*list.Element
	doublyLinkedList *list.List
	callback         OnEliminated
}

func New(maxBytes int64, callback OnEliminated) *Cache {
	return &Cache{
		capacity:         maxBytes,
		hashmap:          make(map[string]*list.Element),
		doublyLinkedList: list.New(),
		callback:         callback,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.hashmap[key]; ok {
		c.doublyLinkedList.MoveToFront(elem)
		entry := elem.Value.(*entry)
		return entry.value, true
	}
	return
}

func (c *Cache) Add(key string, value Value) {
	kvSize := int64(len(key)) + int64(value.Len())
	for c.capacity != 0 && c.length+kvSize > c.capacity {
		c.Remove()
	}
	if elem, ok := c.hashmap[key]; ok {
		c.doublyLinkedList.MoveToFront(elem)
		oldEntry := elem.Value.(*entry)
		c.length += int64(value.Len()) - int64(oldEntry.value.Len())
		oldEntry.value = value
	} else {
		elem := c.doublyLinkedList.PushFront(&entry{key: key, value: value})
		c.hashmap[key] = elem
		c.length += kvSize
	}
}

func (c *Cache) Remove() {
	tailElem := c.doublyLinkedList.Back()
	if tailElem != nil {
		entry := tailElem.Value.(*entry)
		k, v := entry.key, entry.value
		delete(c.hashmap, k)
		c.doublyLinkedList.Remove(tailElem)
		c.length -= int64(len(k)) + int64(v.Len())
		if c.callback != nil {
			c.callback(k, v)
		}
	}
}
