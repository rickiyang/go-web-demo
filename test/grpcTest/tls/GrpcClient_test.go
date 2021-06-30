package normal

import (
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gorm-demo/models/pb"
)

func TestGrpcClient(t *testing.T) {
	// TLS连接
	creds, err := credentials.NewClientTLSFromFile("/Users/yangyue2/ca.crt", "www.yangyue.com")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	// 连接服务器
	conn, err := grpc.Dial(":8972", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	// 调用服务端的SayHello
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "CN"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s !\n", r.Message)
}
