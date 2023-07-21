package discovery

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	defaultReplicas = 5
)

type Discovery interface {
	Refresh() error
	Update(servers []string) error
	Get(group, key string) (string, error)
	GetAll() ([]string, error)
}

type MultiServersDiscovery struct {
	mu       sync.RWMutex
	consHash *Consistency
	servers  []string
}

func NewMultiServersDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{}
	d.consHash = New(defaultReplicas, nil)
	d.Update(servers)
	return d
}

func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.consHash.Register(servers...)
	return nil
}

func (d *MultiServersDiscovery) Get(group, key string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.servers) == 0 {
		return "", errors.New("[rpc discovery]: no available servers")
	}
	server := d.consHash.GetServer(fmt.Sprintf("%s.%s", group, key))
	if server != "" {
		log.Println("[rpc discovery]: get server: ", server)
		return server, nil
	}
	return "", errors.New("[rpc discovery]: get server failed")
}

func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}
