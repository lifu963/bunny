package bunnyDistCache

import (
	"bunny/bunnyDistCache/discovery"
	"context"
	"sync"
)

type ClientController struct {
	mu      sync.Mutex
	disc    discovery.Discovery
	clients map[string]*Client
}

func NewClientController(disc discovery.Discovery) *ClientController {
	return &ClientController{disc: disc, clients: make(map[string]*Client)}
}

func (ctl *ClientController) Close() error {
	ctl.mu.Lock()
	defer ctl.mu.Unlock()
	for key, client := range ctl.clients {
		_ = client.Close()
		delete(ctl.clients, key)
	}
	return nil
}

func (ctl *ClientController) dial(rpcAddr string) (*Client, error) {
	ctl.mu.Lock()
	defer ctl.mu.Unlock()
	client, ok := ctl.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(ctl.clients, rpcAddr)
		client = nil
	}
	if client == nil {
		var err error
		client, err = Dial(rpcAddr)
		if err != nil {
			return nil, err
		}
		ctl.clients[rpcAddr] = client
	}
	return client, nil
}

func (ctl *ClientController) send(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := ctl.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Send(ctx, serviceMethod, args, reply)
}

func (ctl *ClientController) Send(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	argv := args.(Args)
	rpcAddr, err := ctl.disc.Get(argv.Key, argv.Group)
	if err != nil {
		return err
	}
	return ctl.send(rpcAddr, ctx, serviceMethod, args, reply)
}
