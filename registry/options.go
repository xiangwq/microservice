package registry

import "time"

type Options struct {
	Addrs []string
	Timeout time.Duration
	RegistryPath string
	HeartBeat int64
}

type Option func(opt  *Options)


func WithTimeout(time time.Duration) Option{
	return func(opt *Options) {
		opt.Timeout = time
	}
}

func WithAddrs(addrs []string) Option{
	return func(opt *Options) {
		opt.Addrs = addrs
	}
}

func WithRegistryPath(path string) Option{
	return func(opt *Options) {
		opt.RegistryPath = path
	}
}

func WithHeartBeat(heartBeat int64) Option{
	return func(opt *Options) {
		opt.HeartBeat = heartBeat
	}
}