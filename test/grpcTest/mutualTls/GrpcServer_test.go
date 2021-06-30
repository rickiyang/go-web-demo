package normal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/reflection"
	pb "gorm-demo/models/pb"
	"io/ioutil"
	"net"
	"testing"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

//拦截器 - 打印日志
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	fmt.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func TestGrpcServer(t *testing.T) {
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	certificate, err := tls.LoadX509KeyPair(data.Path("/Users/yangyue2/server.crt"), data.Path("/Users/yangyue2/server.key"))
	if err != nil {
		fmt.Errorf("err, %v", err)
	}
	// 创建CertPool，后续就用池里的证书来校验客户端证书有效性
	// 所以如果有多个客户端 可以给每个客户端使用不同的 CA 证书，来实现分别校验的目的
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(data.Path("/Users/yangyue2/ca.crt"))
	if err != nil {
		fmt.Errorf("err, %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Errorf("failed to append certs")
	}

	// 构建基于 TLS 的 TransportCredentials
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{certificate},
		// 要求必须校验客户端的证书 可以根据实际情况选用其他参数
		ClientAuth: tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})

	//开启TLS认证, 注册拦截器
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LoggingInterceptor)) // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{})                                            // 在gRPC服务端注册服务

	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
