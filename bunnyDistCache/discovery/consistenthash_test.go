package discovery

import (
	"hash/crc32"
	"log"
	"sort"
	"testing"
)

func TestRegister(t *testing.T) {
	c := New(2, nil)
	c.Register("server1", "server2")
	if len(c.ring) != 4 {
		t.Errorf("Actual: %d\tExpect: %d\n", len(c.ring), 4)
	}

	hashValue := int(crc32.ChecksumIEEE([]byte("1server1")))
	idx := sort.SearchInts(c.ring, hashValue)
	if c.ring[idx] != hashValue {
		t.Errorf("Actual: %d\tExpect: %d\n", c.ring[idx], hashValue)
	}
}

func TestGet(t *testing.T) {
	c := New(1, nil)
	c.Register("server1", "server2")
	key := "Tom"
	keyHashValue := int(crc32.ChecksumIEEE([]byte(key)))
	log.Printf("key hash = %d\n", keyHashValue)
	for _, v := range c.ring {
		log.Printf("%d -> %s\n", v, c.hashmap[v])
	}
	peer := c.GetServer(key)
	log.Printf("Go to search -> %s\n", peer)
}
