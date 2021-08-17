package loadbalance

import (
	"context"
	"math/rand"
	"microservice/registry"
)

// WeightedRound 加权随机算法
type WeightedRound struct {
	name string
}

func NewWeightedRound() LoadBalance {
	return &RandomBalance{
		name: "WeightedRound",
	}
}

func (r *WeightedRound) Name() string {
	return r.name
}

func (r *WeightedRound) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		return nil, ErrNotHaveNode
	}

	total := 0
	for _, v := range nodes {
		if v.Weight <= 0 {
			total += 0
			continue
		}
		total += v.Weight
	}

	if total <= 0 {
		return nil, ErrNotHaveNode
	}

	//加权轮询
	index := rand.Intn(total)
	for _, v := range nodes {
		index = index - v.Weight
		if index < 0 {
			return v, nil
		}
	}
	return nil, ErrNotHaveNode
}
