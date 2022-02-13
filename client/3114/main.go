package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.balancer.io/balancer/client/v3"
	"google.golang.org/grpc"
	pb "tag-service/proto"
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
	config := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second * 60,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	r := &naming.GRPCResolver{Client: cli}
	target := fmt.Sprintf("/etcdv3://go-programming-tour/grpc/%s", serviceName)
	opts = append(opts, grpc.WithBalancerName(grpc.RoundRobin(r)), grpc.WithBlock())
	return grpc.DialContext(ctx, target, opts...)
}
