package main

import "github.com/emicklei/proto"

type Generator interface {
	Run(opt *Option, metaData *ServiceMetaData) error
}

type ServiceMetaData struct {
	Service *proto.Service
	Message []*proto.Message
	Rpc     []*proto.RPC
	Package *proto.Package
	Prefix  string
}
