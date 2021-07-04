package __

import (
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// customCredential 自定义认证
type customCredential struct{}

// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if ctx.Value("authorization") == nil {
		return map[string]string{
			"authorization": "123456",
		}, nil
	}
	return map[string]string{
		"authorization": ctx.Value("authorization").(string),
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return true
}

func TestGrpcClient(t *testing.T) {
	var err error
	var opts []grpc.DialOption

	// TLS连接
	creds, err := credentials.NewClientTLSFromFile("/Users/yangyue2/ca.crt", "www.yangyue.com")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	//连接服务端
	conn, err := grpc.Dial(":8972", opts...)
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := NewTokenServiceClient(conn)
	// 调用服务端的SayHello
	r, err := c.Login(context.Background(), &LoginRequest{Username: "xiaoming", Password: "123456"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	requestToken := new(AuthToken)
	requestToken.Token = r.Token

	//连接服务端
	conn, err = grpc.Dial(":8972", grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(requestToken))
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()
	c = NewTokenServiceClient(conn)
	hello, err := c.SayHello(context.Background(), &PingMessage{Greeting: "hahah"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s, %s !\n", r.Token, hello)
}
