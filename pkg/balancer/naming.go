package balancer

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

func Register(etcdAddr, name string, addr string, ttl int64) error {
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return err
		}
	}
	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		for {
			getResp, err := cli.Get(context.Background(), "/"+schema+"/"+name+"/"+addr)
			if err != nil {
				log.Println(err)
			} else if getResp.Count == 0 {
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Println(err)
				}
			} else {

			}
			<-ticker.C
		}
	}()
	return nil
}

func withAlive(name string, addr string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	fmt.Printf("key: %v\n", "/"+schema+"/"+name+"/"+addr)
	_, err = cli.Put(context.Background(), "/"+schema+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	_, err = cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}
	return nil
}

func UnRegister(name string, addr string) {
	if cli != nil {
		cli.Delete(context.Background(), "/"+schema+"/"+name+"/"+addr)
	}
}
