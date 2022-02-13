package main

import (
	"context"
	"flag"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"path"
	"strings"
	"tag-service/internal/middleware"
	"tag-service/pkg/swagger"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clientv3 "go.balancer.io/balancer/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "tag-service/proto"
	"tag-service/server"

)

var port string

func init() {
	flag.StringVar(&port, "port", "8004", "启动端口号")
	flag.Parse()
}

const schema = "balancer"
const SERVICE_NAME = "tag-service"

type Resolver struct {
	endpoints []string
	service   string
	cli       *clientv3.Client
	cc        resolver.ClientConn
}

func NewResolver(endpoints []string, service string) resolver.Builder {
	return &Resolver{endpoints: endpoints, service: service}
}

func (r *Resolver) Schema() string {
	return schema + "_" + r.service
}
func (r *Resolver) ResolverNow(rn resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {

}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.endpoints,
	})
	if err!=nil{
		return nil,fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}
	r.cc=cc
	go r.
}

func (r *Resolver)watch(prefix string) {
	var addrList []resolver.Address
	getResp,err:=r.cli.Get(context.Background(),prefix,clientv3.WithPrefix())
	if err!=nil{
		log.Println(err)
	}else {
		for i:=range getResp.Kvs{
			addrList=append(addrList,resolver.Address{Addr: strings.TrimPrefix()})
		}
	}
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}
}

func RunServer(port string) error {
	httpMux := runHttpServer()
	grpcS := runGrpcServer()

	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	httpMux.Handle("/", gwmux)

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second * 60,
	})
	if err != nil {
		return err
	}
	defer etcdClient.Close()

	target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", SERVICE_NAME)
	grpcproxy.Register(etcdClient, target, ":"+port, 60)

	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func runHttpServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})

	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	serveMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)

		http.ServeFile(w, r, p)
	})

	return serveMux
}

func runGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
		)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
