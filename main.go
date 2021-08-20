package main // import "microservice"

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: time.Second,
	})

	if err != nil {
		log.Fatalln(err)
	}

	defer cli.Close()

	resp, err := cli.Grant(context.TODO(), 5)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = cli.Put(context.TODO(), "test", "tt", clientv3.WithLease(resp.ID))

	if err != nil {
		log.Fatalln(err)
	}

	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatalln(kaerr)
	}
	for {
		ka := <-ch
		fmt.Println(ka)
	}

}
