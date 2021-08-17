package normal

import (
	"fmt"
	"testing"
	"time"

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

	//timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancelFunc()

	m, _ := time.ParseDuration("1s")
	result := time.Now().Add(m)
	deadline, c2 := context.WithDeadline(context.Background(), result)
	defer c2()

	// 调用服务端的SayHello
	r, err := c.SayHello(deadline, &pb.HelloRequest{Name: "CN"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s !\n", r.Message)
}
