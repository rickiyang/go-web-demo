package streamTest

import (
	"fmt"
	"io"
	"strconv"
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

	c := pb.NewBaseServiceClient(conn)
	// 客户端流式处理- Sum案例
	//sum(c)

	//客户端流式处理逻辑测试
	//clientStream(c, "我收到了服务端的请求数据拉")

	//服务端流式处理
	//serverStream(c, &pb.StreamRequest{Input: "我是一只小老虎"})

	//双向流处理
	streaming(c)

}

//sum案例- 客户端流式处理
func sum(client pb.BaseServiceClient) {
	sumCli, err := client.Sum(context.Background())
	if err != nil {
		panic("sum cli err")
	}
	sumCli.Send(&pb.SumRequest{Num: int64(1)})
	sumCli.Send(&pb.SumRequest{Num: int64(2)})
	sumCli.Send(&pb.SumRequest{Num: int64(3)})
	sumCli.Send(&pb.SumRequest{Num: int64(4)})

	recv, err := sumCli.CloseAndRecv()
	if err != nil {
		fmt.Printf("send sum request err: %v", err)
	}
	fmt.Printf("sum = : %v !\n", recv.Result)
}

// 客户端流式处理 的逻辑
func clientStream(client pb.BaseServiceClient, input string) error {
	stream, _ := client.ClientStream(context.Background())
	for _, s := range input {
		fmt.Println("Client Stream Send:", string(s))
		err := stream.Send(&pb.StreamRequest{Input: string(s)})
		if err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Client Stream Recv:", res.Output)
	return nil
}

// 服务端流式处理
func serverStream(client pb.BaseServiceClient, r *pb.StreamRequest) error {
	fmt.Println("Server Stream Send:", r.Input)
	stream, _ := client.ServerStream(context.Background(), r)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println("Server Stream Recv:", res.Output)
	}
	return nil
}

// 双向流式处理
func streaming(client pb.BaseServiceClient) error {
	stream, _ := client.Streaming(context.Background())
	for n := 0; n < 10; n++ {
		fmt.Println("Streaming Send:", n)
		err := stream.Send(&pb.StreamRequest{Input: strconv.Itoa(n)})
		if err != nil {
			return err
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println("Streaming Recv:", res.Output)
	}
	stream.CloseSend()
	return nil
}
