package loadbalance

import (
	"context"
	"fmt"
	"microservice/registry"
	"testing"
)

func TestRandomBalance_Select(t *testing.T) {
	balance := &WeightedRoundRobin{}
	var nodes []*registry.Node
	for i := 0; i < 4; i++ {
		node := &registry.Node{
			IP:     fmt.Sprintf("127.0.0.%d", i),
			Port:   8080,
			Weight: i,
		}
		fmt.Println(node.IP, node.Weight)
		nodes = append(nodes, node)
	}

	count := make(map[string]int)

	for i := 0; i < 600; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Errorf("select failed: %v", err)
			continue
		}

		count[node.IP]++
	}

	for k, v := range count {
		fmt.Println("k", k, "v", v)
	}
}
