package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"microservice/registry"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

type RegisterService struct {
	Id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

const (
	MaxServiceNum          = 8
	MaxSyncServiceInterval = time.Second * 10
)

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo)
	registry.RegistryPlugin(etcdRegistry)
	go etcdRegistry.run()
}

func (e *EtcdRegistry) Name() string {
	return "etcd"
}

func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{Endpoints: e.options.Addrs, DialTimeout: e.options.Timeout})
	if err != nil {
		return fmt.Errorf("init etcd failed, errL %v", err)
	}
	return
}

func (e *EtcdRegistry) run() {
	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceCh:
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)
				}
				registryService.registered = false
				break
			}
			registryService = &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			e.syncServiceFromEtcd()
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Second)
		}
	}
}

func (e *EtcdRegistry) registerOrKeepAlive() {
	for _, service := range e.registryServiceMap {
		if service.registered {
			e.keepAlive(service)
			continue
		}

		e.registerService(service)
	}
}

func (e *EtcdRegistry) registerService(registerService *RegisterService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		return err
	}

	registerService.Id = resp.ID

	for _, node := range registerService.service.Nodes {
		tmp := &registry.Service{
			Name: registerService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}

		data, err := json.Marshal(tmp)
		if err != nil {
			return nil
		}

		key := e.serviceNodePath(tmp)

		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			return err
		}

		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)

		if err != nil {
			return err
		}
		registerService.keepAliveCh = ch
		registerService.registered = true
	}

	return nil
}

func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
	}
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIp := fmt.Sprintf("%s:%d", service.Nodes[0].IP, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIp)
}

func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err := fmt.Errorf("register is full")
		return err
	}
	return
}

func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	// 一般情况下，都会从缓存读取
	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}
	// 如果缓存没有这个service，则从etcd中读取
	e.lock.Lock()
	defer e.lock.Unlock()
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	// 从etcd中读取指定服务
	key := e.servicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	service = &registry.Service{
		Name: name,
	}

	for _, k := range resp.Kvs {
		value := k.Value
		var tmp registry.Service
		err = json.Unmarshal(value, &tmp)
		if err != nil {
			return
		}

		for _, node := range tmp.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}

	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}
	allServiceInfoNew.serviceMap[key] = service
	e.value.Store(allServiceInfoNew)
	return
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	// 一般情况下，都会从缓存读取
	service, ok = allServiceInfo.serviceMap[name]
	return
}

func (e *EtcdRegistry) syncServiceFromEtcd() {
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	allServiceInfo := e.value.Load().(*AllServiceInfo)

	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		resp, err := e.client.Get(context.TODO(), key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, k := range resp.Kvs {
			value := k.Value
			var tmp registry.Service
			err = json.Unmarshal(value, &tmp)
			if err != nil {
				return
			}

			for _, node := range tmp.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}

		allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}
	e.value.Store(allServiceInfoNew)
	return
}

func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}
