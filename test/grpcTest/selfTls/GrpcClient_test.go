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

const (
	// OpenTLS 是否开启TLS认证
	OpenTLS = true
)

// customCredential 自定义认证
type customCredential struct{}

// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "110",
		"appKey": "token",
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

func TestGrpcClient(t *testing.T) {
	var err error
	var opts []grpc.DialOption

	if OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile("/Users/yangyue2/ca.crt", "www.yangyue.com")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	// 使用自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	//连接服务端
	conn, err := grpc.Dial(":8972", opts...)
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
