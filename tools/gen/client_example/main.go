package main

import (
	"context"
	"fmt"
	"microservice/logs"
	"microservice/tools/gen/client_example/generate/client/testc"
	"microservice/tools/gen/client_example/generate/test"
	"time"
)

func main() {
	client := testc.NewTestClient("test")
	ctx := context.Background()

	resp, err := client.SayHello(ctx, &test.HelloRequest{Name: "this is test client"})
	if err != nil {
		logs.Error(ctx, "could not greet:%v", err)
		return
	}
	fmt.Println(resp.Name)
	logs.Info(ctx, "Greeting: %s", resp.Name)
	logs.Stop()
	for {
		time.Sleep(time.Second)
	}
}
