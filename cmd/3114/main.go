package main

import (
	"context"
	"flag"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/credentials/insecure"
	"tag-service/pkg/balancer"

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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "tag-service/proto"
	"tag-service/server"
)

const (
	SERVICE_NAME = "tag-service"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8004", "启动端口号")
	flag.Parse()
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
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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

	//target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", SERVICE_NAME)
	//grpcproxy.Register(etcdClient, target, ":"+port, 60)
	go balancer.Register("http://localhost:2379", fmt.Sprintf("grpc/%s", SERVICE_NAME), "localhost:8004", 1)

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
