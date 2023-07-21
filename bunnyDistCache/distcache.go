package bunnyDistCache

import (
	"bunny/bunnyDistCache/discovery"
	"bunny/bunnyDistCache/singleflight"
	"context"
	"fmt"
)

type DistCache struct {
	ctl    *ClientController
	flight *singleflight.Flight
}

func NewDistCache(registryAddr string) *DistCache {
	cache := &DistCache{
		ctl:    NewClientController(discovery.NewRegistryDiscovery(registryAddr, 0)),
		flight: &singleflight.Flight{},
	}
	return cache
}

func (cache *DistCache) Close() error {
	return cache.ctl.Close()
}

func (cache *DistCache) Get(ctx context.Context, group, key string) (ByteView, error) {
	reply, err := cache.flight.Fly(fmt.Sprintf("%s.%s", group, key), func() (interface{}, error) {
		args := Args{Group: group, Key: key}
		var view ByteView
		err := cache.ctl.Send(context.Background(), "BunnyCache.Get", args, &view)
		if err != nil {
			return ByteView{}, err
		}
		return view, nil
	})
	return reply.(ByteView), err
}
