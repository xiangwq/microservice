package loadbalance

import (
	"context"
	"math/rand"
	"microservice/registry"
)

type RandomBalance struct {
	name string
}

func NewRandomBalance() LoadBalance {
	return &RandomBalance{
		name: "random",
	}
}

func (r *RandomBalance) Name() string {
	return r.name
}

func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		return nil, ErrNotHaveNode
	}

	index := rand.Intn(len(nodes))
	node = nodes[index]
	return node, nil
}
