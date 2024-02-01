package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceName string
type instanceID string

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

// NewRegistry creates a new in-memory service registry instance
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, ID string, name string, addr string) error {
	s := serviceName(name)
	i := instanceID(ID)
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[s]; !ok {
		r.serviceAddrs[s] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[s][i] = &serviceInstance{hostPort: addr, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, ID string, name string) error {
	s := serviceName(name)
	i := instanceID(ID)
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[s]; ok {
		delete(r.serviceAddrs[s], i)
	}
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry
func (r *Registry) ReportHealthyState(ID string, name string) error {
	s := serviceName(name)
	i := instanceID(ID)
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[s]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[s][i]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[s][i].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddress(ctx context.Context, name string) ([]string, error) {
	s := serviceName(name)
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.serviceAddrs[s]; !ok {
		return nil, errors.New("service is not registered yet")
	}
	if len(r.serviceAddrs[s]) == 0 {
		return nil, discovery.ErrNotFound
	}
	addrs := make([]string, 0, len(r.serviceAddrs[s]))
	for _, instance := range r.serviceAddrs[s] {
		if instance.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		addrs = append(addrs, instance.hostPort)
	}
	return addrs, nil
}


