package normal

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"io/ioutil"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gorm-demo/models/pb"
)

func TestGrpcClient(t *testing.T) {

	// 加载客户端证书
	certificate, err := tls.LoadX509KeyPair(data.Path("/Users/yangyue2/client.crt"), data.Path("/Users/yangyue2/client.key"))
	if err != nil {
		fmt.Errorf("err, %v", err)
	}
	// 构建CertPool以校验服务端证书有效性
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(data.Path("/Users/yangyue2/ca.crt"))
	if err != nil {
		fmt.Errorf("err, %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Errorf("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "www.yangyue.com", // NOTE: this is required!
		RootCAs:      certPool,
	})

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
