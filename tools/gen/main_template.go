package main

var main_template = `package main
import (
	"log"
	"microservice/server"
     {{ if not .Prefix}}
	"generate/{{.Package.Name}}"
	"router"
	{{else}}
    "{{.Prefix}}/generate/{{.Package.Name}}"
    "{{.Prefix}}/router"
	{{end}}
	)

	var routerServer = &router.RouterServer{}

    func main() {
		err := server.Init("{{.Package.Name }}")
		if err != nil {
			log.Fatal("init service failed, err: $v", err)
			return 
		}
	
		{{.Package.Name}}.Register{{ .Service.Name }}Server(server.GRPCServer(), routerServer)
		server.Run()
	}
`
