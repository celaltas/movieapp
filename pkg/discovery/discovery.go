package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrNotFound = errors.New("no service addresses found")

// Registry defines a service registry.
type Registry interface {
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	ServiceAddress(ctx context.Context, serviceName string) ([]string, error)

	// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
	ReportHealtyState(instanceID string, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
