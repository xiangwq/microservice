package loadbalance

import (
	"context"
	"errors"
	"microservice/registry"
)

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}

var (
	ErrNotHaveNode = errors.New("not have node")
)
