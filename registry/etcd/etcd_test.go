package etcd

import (
	"context"
	"fmt"
	"microservice/registry"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	inst, err := registry.InitRegistry(context.TODO(), "etcd",
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("microservice/"),
		registry.WithHeartBeat(5))

	if err != nil {
		t.Errorf("init registry failed: %v", err)
		return
	}

	service := &registry.Service{
		Name: "test_service",
	}
	service.Nodes = append(service.Nodes, &registry.Node{IP: "127.0.0.1", Port: 8801}, &registry.Node{IP: "127.0.0.1", Port: 8802})
	inst.Register(context.TODO(), service)

	go func() {
		time.Sleep(10 * time.Second)
		service := &registry.Service{
			Name: "test_service",
		}
		service.Nodes = append(service.Nodes, &registry.Node{IP: "127.0.0.2", Port: 8801}, &registry.Node{IP: "127.0.0.3", Port: 8802})
		inst.Register(context.TODO(), service)
	}()

	for {
		service , err := inst.GetService(context.TODO(), "test_service" )
		if err != nil {
			t.Errorf("get service failed: %v", err)
			return
		}
		time.Sleep(time.Second)
	}
}
