package balancer

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"strings"
	"time"
)

const (
	schema = "hahah"
)

var cli *clientv3.Client

type etcdResolver struct {
	rawAddr string
	cc      resolver.ClientConn
	cli     *clientv3.Client
}

func NewResolver(etcdAddr string) resolver.Builder {
	return &etcdResolver{
		rawAddr: etcdAddr,
	}
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	if r.cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(r.rawAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}
	r.cc = cc

	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")
	return r, nil
}
func (r *etcdResolver) watch(prefix string) {
	var addrList []resolver.Address

	getResp, err := cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), prefix)})
		}
	}
	r.cc.UpdateState(resolver.State{
		Addresses: addrList,
	})

	rch := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), prefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.UpdateState(resolver.State{
						Addresses: addrList,
					})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.UpdateState(resolver.State{
						Addresses: addrList,
					})
				}
			}
		}
	}
}
func (r etcdResolver) Scheme() string {
	return schema
}
func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	log.Println("ResolveNow")
}
func (r etcdResolver) Close() {
	log.Println("Close")
}
func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}
func remove(s []resolver.Address, add string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == add {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
