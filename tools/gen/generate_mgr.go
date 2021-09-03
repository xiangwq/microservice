package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf/test",
	"conf/product",
	"app/router",
	"app/config",
	"model",
	"generate",
	"router",
}

var genMgr *GenerateMgr = &GenerateMgr{
	genClientMap:  make(map[string]Generator),
	genServiceMap: make(map[string]Generator),
	metaData:      &ServiceMetaData{},
}

type GenerateMgr struct {
	genClientMap  map[string]Generator
	genServiceMap map[string]Generator
	metaData      *ServiceMetaData
}

func RegisterClientGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genClientMap[name]
	if ok {
		return fmt.Errorf("generator %s is exists", name)
	}
	genMgr.genClientMap[name] = gen
	return nil
}

func RegisterServiceGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genServiceMap[name]
	if ok {
		return fmt.Errorf("generator %s is exists", name)
	}
	genMgr.genServiceMap[name] = gen
	return nil
}

func (g *GenerateMgr) Run(opt *Option) error {
	err := g.initOutputDir(opt)
	if err != nil {
		return err
	}

	err = g.parseService(opt)
	if err != nil {
		return err
	}

	g.metaData.Prefix = opt.Prefix

	if opt.GenServerCode {
		err = g.CreateAllDir(opt)
		if err != nil {
			return err
		}
		for _, gen := range g.genServiceMap {
			err := gen.Run(opt, g.metaData)
			if err != nil {
				return nil
			}
		}
	}

	if opt.GenClientCode {
		for _, gen := range g.genClientMap {
			err := gen.Run(opt, g.metaData)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (g *GenerateMgr) initOutputDir(opt *Option) error {
	gopath := os.Getenv("GOPATH")
	opt.Prefix = "microservice/tools/gen/" + opt.Name
	exeFilePath, err := filepath.Abs(os.Args[0])
	lastIdx := strings.LastIndex(exeFilePath, "/")
	//exeFilePath = "/Users/xiangwenqi/go/src/microservice/tools/gen/gen"
	if err != nil {
		return err
	}
	srcPath := path.Join(gopath, "src/")
	// 判断项目是否在gopath下且不在
	if exeFilePath[:len(srcPath)] != srcPath && len(opt.Prefix) <= 0 {
		return fmt.Errorf("不在gopath下的项目需要设置包import prefix")
	}
	if exeFilePath[:len(srcPath)] != srcPath {
		opt.Output = exeFilePath[0:lastIdx] + "/" + opt.Name
		fmt.Printf("opt output:%s, prefix:%s, gopath:%s\n", opt.Output, opt.Prefix, gopath)
		return nil
	}

	// 自动构建gopath下的项目
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return err
	}
	opt.Output = exeFilePath[0:lastIdx] + "/" + opt.Name
	if srcPath[len(srcPath)-1] != '/' {
		srcPath = fmt.Sprintf("%s/", srcPath)
	}
	opt.Prefix = strings.Replace(opt.Output, srcPath, "", -1)
	fmt.Printf("opt output:%s, prefix:%s, gopath:%s\n", opt.Output, opt.Prefix, gopath)
	return nil
}

func (g *GenerateMgr) CreateAllDir(opt *Option) error {
	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err := os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err: %s", fullDir, err)
			return err
		}
		fmt.Println(fullDir, " is created")
	}
	return nil
}

func (g *GenerateMgr) parseService(opt *Option) error {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Println("open proto file failed", opt.Proto3Filename)
		return err
	}
	defer reader.Close()

	parse := proto.NewParser(reader)
	definition, err := parse.Parse()
	if err != nil {
		fmt.Println("parse proto file failed", opt.Proto3Filename)
		return err
	}

	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage))

	//fmt.Printf("parse rpc: %#v", c.rpc)
	return nil
}

func (g *GenerateMgr) handleService(s *proto.Service) {
	fmt.Println("service:", s)
	g.metaData.Service = s
}

func (g *GenerateMgr) handleMessage(m *proto.Message) {
	fmt.Println("message:", m)
	g.metaData.Message = append(g.metaData.Message, m)

}

func (g *GenerateMgr) handleRPC(r *proto.RPC) {
	fmt.Println("rpc:", r)
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GenerateMgr) handlePackage(p *proto.Package) {
	g.metaData.Package = p
}
