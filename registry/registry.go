package registry

import "context"

type Registry interface {
	// Name 插件名字
	Name() string
	// Init 初始化
	Init(ctx context.Context, opts ...Option)(err error)
	// Register 注册服务
	Register(ctx context.Context, service *Service)(err error)
	// Unregister 服务反注册（取消注册）
	Unregister(ctx context.Context, service *Service)(err error)
	// GetService 服务发现
	GetService(ctx context.Context, name string)(service *Service, err error)
}