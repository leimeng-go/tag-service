package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"tag-service/global"
	"tag-service/internal/middleware"
	"tag-service/pkg/tracer"
	pb "tag-service/proto"
)

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

func main() {
	ctx := context.Background()
	newCtx := metadata.AppendToOutgoingContext(ctx, "eddycjy", "Go语言编程之旅")
	clientConn, err := GetClientConn(newCtx, "localhost:8004", []grpc.DialOption{grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
			middleware.ClientTracing(),
		),
	)})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	var err error
	jaegerTracer, _, err := tracer.NewJaegerTracer("article-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
