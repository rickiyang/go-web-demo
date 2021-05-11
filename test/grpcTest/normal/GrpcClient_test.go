package normal

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gorm-demo/models/pb"
)

func TestGrpcClient(t *testing.T) {
	// 连接服务器
	conn, err := grpc.Dial(":8972", grpc.WithInsecure())
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
