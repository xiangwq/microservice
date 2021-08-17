package loadbalance

import (
	"context"
	"microservice/registry"
)

// WeightedRoundRobin 加权轮询算法
type WeightedRoundRobin struct {
	name  string
	index int
}

func NewWeightedRoundRobin() LoadBalance {
	return &RandomBalance{
		name: "WeightedRoundRobin",
	}
}

func (r *WeightedRoundRobin) Name() string {
	return r.name
}

func (r *WeightedRoundRobin) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
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

	//加权轮询
	index := r.index
	for _, v := range nodes {
		index = index - v.Weight
		if index >= 0 {
			continue
		}
		r.index++
		if r.index >= total {
			r.index = 0
		}
		return v, nil
	}

	return nil, ErrNotHaveNode
}
