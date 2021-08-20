package main

var main_template = `package main
import (
	"google.golang.org/grpc"
	"log"
     {{ if not .Prefix}}
	"generate"
	"router"
	{{else}}
    	"{{.Prefix}}/generate"
		"{{.Prefix}}/router"
	{{end}}
		"net"
	)

	var server = &router.RouterServer{}
	const port = ":30011"

    func main() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalln("failed to listen", err.Error())
		}

		s := grpc.NewServer()
		generate.Register{{ .Service.Name }}Server(s, server)
		s.Serve(lis)
	}
`
