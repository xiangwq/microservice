package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"
)

type CtrlGenerator struct {
}

type RpcMeta struct {
	Rpc *proto.RPC
	//Package *proto.Package
	//Prefix  string
	*ServiceMetaData
}

func (c *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) error {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Println("open proto file failed", opt.Proto3Filename)
		return err
	}
	defer reader.Close()

	//fmt.Printf("parse rpc: %#v", c.rpc)
	return c.generateRpc(opt, metaData)
}

func (c *CtrlGenerator) generateRpc(opt *Option, metaData *ServiceMetaData) (err error) {

	for _, rpc := range metaData.Rpc {
		var file *os.File
		tmpName := ToUnderScoreString(rpc.Name)
		filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", strings.ToLower(tmpName)))
		fmt.Printf("filename is %s\n", filename)
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("open file:%s failed, err:%v\n", filename, err)
			return
		}

		rpcMeta := &RpcMeta{}
		rpcMeta.Rpc = rpc
		rpcMeta.ServiceMetaData = metaData
		//rpcMeta.Package = metaData.Package
		//rpcMeta.Prefix = metaData.Prefix

		err = c.render(file, controller_template, rpcMeta)
		if err != nil {
			fmt.Printf("render controller failed err:%v\n", err)
			return
		}
		file.Close()
	}
	return
}

func (c *CtrlGenerator) render(file *os.File, data string, metaData *RpcMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func ToUnderScoreString(name string) string {
	var buffer bytes.Buffer
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteString("_")
			}
			buffer.WriteString(fmt.Sprintf("%c", unicode.ToLower(r)))
		} else {
			buffer.WriteString(fmt.Sprintf("%c", unicode.ToLower(r)))
		}
	}

	return buffer.String()
}

func init() {
	dir := &CtrlGenerator{}
	Register("controller_generate", dir)
}
