package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"tag-service/pkg/balancer"
	pb "tag-service/proto"
	"time"
)

func main() {
	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, "tag-service", nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, serviceName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	r := balancer.NewResolver("localhost:2378")
	resolver.Register(r)

	conn, err := grpc.Dial(r.Scheme()+"://author/project/test", grpc.WithDefaultServiceConfig("round_robin"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewTagServiceClient(conn)
	for {
		resp, err := client.GetTagList(ctx, &pb.GetTagListRequest{
			Name: "haha",
		}, grpc.WaitForReady(true))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}
		<-time.After(time.Second)
	}
}
