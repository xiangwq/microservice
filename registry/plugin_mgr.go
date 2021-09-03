package registry

import (
	"context"
	"fmt"
	"sync"
)

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex
}

// 注册插件
func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.plugins[plugin.Name()]
	if ok {
		return fmt.Errorf("duplicate registry plugin")
	}

	p.plugins[plugin.Name()] = plugin
	return nil
}

// 实例化插件
func (p *PluginMgr) initRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	plugin, ok := p.plugins[name]
	fmt.Println(len(p.plugins))
	for k, v := range p.plugins {
		fmt.Println("k:", k, "v:", v.Name())
	}
	if !ok {
		return nil, fmt.Errorf("plugin %s is not exist", name)
	}
	registry = plugin
	err = registry.Init(ctx, opts...)
	return
}

// RegistryPlugin 注册插件
func RegistryPlugin(plugin Registry) (err error) {
	return pluginMgr.registerPlugin(plugin)
}

// InitRegistry 初始化
func InitRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	return pluginMgr.initRegistry(ctx, name, opts...)
}
