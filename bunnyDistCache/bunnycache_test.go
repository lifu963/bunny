package bunnyDistCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestRetriever(t *testing.T) {
	var f Retriever = RetrieverFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.retrieve("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGetGroup(t *testing.T) {
	bunnyCache := NewBunnyCache()
	groupName := "scores"
	bunnyCache.NewGroup(groupName, 2<<10, RetrieverFunc(
		func(key string) (bytes []byte, err error) { return }))

	if g, err := bunnyCache.getGroup(groupName); err != nil || g.name != groupName {
		t.Fatalf("[BunnyCache]: %s not exist", groupName)
	}

	if g, err := bunnyCache.getGroup(groupName + "***"); err == nil || g != nil {
		t.Fatalf("[BunnyCache]: expect nil, but %s got", g.name)
	}
}

func TestGet(t *testing.T) {
	bunnyCache := NewBunnyCache()
	groupName := "scores"
	loadCounts := make(map[string]int, len(db))
	_ = bunnyCache.NewGroup(groupName, 2<<10, RetrieverFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for key, value := range db {
		if view, err := bunnyCache.get(groupName, key); err != nil || view.String() != value {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := bunnyCache.get(groupName, key); err != nil || loadCounts[key] > 1 {
			t.Fatalf("cache %s miss", key)
		}
	}

	if view, err := bunnyCache.get(groupName, "unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
