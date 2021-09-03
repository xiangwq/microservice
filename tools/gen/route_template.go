package main

var router_template = `package router
import(
	"context"
	"microservice/meta"
	"microservice/server"
	{{ if not .Prefix}}
	"generate/{{.Package.Name}}"
    "controller"
	{{else}}
    "{{.Prefix}}/controller"
    "{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)

type RouterServer struct{}

{{range .Rpc}}
func (s *RouterServer) {{.Name}}(ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}})(resp *{{$.Package.Name}}.{{.ReturnsType}}, err error){
    ctx = meta.InitServerMeta(ctx, "{{$.Package.Name}}", "{{.Name}}" )
	mwFunc := server.BuildServerMiddleware(mw{{.Name}})
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return 
	}
	resp = mwResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	return resp, err
}

func mw{{.Name}}(ctx context.Context, req interface{}) (resp interface{}, err error) {
	r := req.(*{{$.Package.Name}}.{{.RequestType}})
	ctrl := &controller.{{.Name }}Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
{{end}}
`
