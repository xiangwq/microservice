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
	dir := fmt.Sprintf("%s/generate/%s", opt.Output, metaData.Package.Name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("mkdir dir %s failed, err: %s", dir, err)
		return err
	}

	outputParams := fmt.Sprintf("--go_out=plugins=grpc:%s", dir)
	fmt.Println(outputParams)
	cmd := exec.Command("protoc", outputParams, opt.Proto3Filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println("grpc gen failed")
		return err
	}
	return nil
}

func init() {
	dir := &GrpcGenerator{}
	RegisterServiceGenerator("grpc_generate", dir)
	RegisterClientGenerator("grpc_generate", dir)
}
