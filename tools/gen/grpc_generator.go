package main

import (
	"fmt"
	"os"
	"os/exec"
)

type GrpcGenerator struct {
}

func (g *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) error {
	// protoc --go_out=plugins=grpc:. .*.proto
	outputParams := fmt.Sprintf("--go_out=plugins=grpc:%s/generate/", opt.Output)
	fmt.Println(outputParams)
	cmd := exec.Command("protoc", outputParams, opt.Proto3Filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("grpc gen failed")
		return err
	}
	return nil
}

func init() {
	dir := &GrpcGenerator{}
	Register("grpc_generate", dir)
}
