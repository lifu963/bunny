package bunnyDistCache

import (
	"fmt"
	"log"
	"sync"
)

type Retriever interface {
	retrieve(string) ([]byte, error)
}

type RetrieverFunc func(key string) ([]byte, error)

func (f RetrieverFunc) retrieve(key string) ([]byte, error) {
	return f(key)
}

type group struct {
	name      string
	cache     *cache
	retriever Retriever
}

func (g *group) get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("[group %s]: key is required!\n", g.name)
	}

	if value, ok := g.cache.get(key); ok {
		log.Printf("[group %s]: cache hit!\n", g.name)
		return value, nil
	}

	bytes, err := g.retriever.retrieve(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{B: cloneBytes(bytes)}
	g.cache.add(key, value)
	return value, nil
}

type BunnyCache struct {
	mu     sync.RWMutex
	groups map[string]*group
}

func NewBunnyCache() *BunnyCache {
	return &BunnyCache{groups: make(map[string]*group)}
}

func (b *BunnyCache) NewGroup(name string, maxBytes int64, retriever Retriever) error {
	if retriever == nil {
		return fmt.Errorf("[BunnyCache]: group retriever must be existed!\n")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	g := &group{
		name:      name,
		cache:     newCache(maxBytes),
		retriever: retriever,
	}
	b.groups[name] = g
	return nil
}

func (b *BunnyCache) get(groupName, key string) (ByteView, error) {
	g, err := b.getGroup(groupName)
	if err != nil {
		return ByteView{}, err
	}
	return g.get(key)
}

func (b *BunnyCache) getGroup(groupName string) (*group, error) {
	if groupName == "" {
		return nil, fmt.Errorf("[BunnyCache]: groupName is required!\n")
	}
	b.mu.RLock()
	defer b.mu.RUnlock()

	if g, ok := b.groups[groupName]; ok {
		return g, nil
	}

	return nil, fmt.Errorf("[BunnyCache]: %s is not exist!\n", groupName)
}

type Args struct {
	Group string
	Key   string
}

func (d *BunnyCache) Get(args Args, reply *ByteView) error {
	view, err := d.get(args.Group, args.Key)
	if err != nil {
		*reply = ByteView{}
		return err
	}
	*reply = view
	return nil
}
