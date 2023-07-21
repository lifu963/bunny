package main

import (
	"bunny/bunnyDistCache"
	"bunny/bunnyDistCache/registry"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func startRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9998")
	r := registry.New()
	r.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}

func startServer(registryAddr string, wg *sync.WaitGroup) {
	cache := bunnyDistCache.NewBunnyCache()
	cache.NewGroup("default", 2<<10, bunnyDistCache.RetrieverFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	server := bunnyDistCache.NewServer()
	if err := server.Register(cache); err != nil {
		log.Fatal("register error:", err)
	}
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	wg.Done()
	server.Accept(l)
}

func startAPIServer(apiAddr string, cache *bunnyDistCache.DistCache) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			reply, err := cache.Get(context.Background(), "default", key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(reply.ByteSlice())
		}))
	log.Println("API server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	log.SetFlags(0)
	registryAddr := "http://localhost:9998/registry"
	apiAddr := "http://localhost:9999"
	var wg sync.WaitGroup
	wg.Add(1)
	go startRegistry(&wg)
	wg.Wait()

	time.Sleep(time.Second)
	wg.Add(2)
	go startServer(registryAddr, &wg)
	go startServer(registryAddr, &wg)
	wg.Wait()

	time.Sleep(time.Second)

	cache := bunnyDistCache.NewDistCache(registryAddr)
	defer func() { _ = cache.Close() }()

	startAPIServer(apiAddr, cache)
}
