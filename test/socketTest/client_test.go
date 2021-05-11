package socketTest

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

func TestClient(t *testing.T) {
	// 1、与服务端建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return
	}
	// 2、使用 conn 连接进行数据的发送和接收
	//input := bufio.NewReader(os.Stdin)
	//for {
	//	line, _, err := input.ReadLine()
	//	s := strings.TrimSpace(string(line))
	//	fmt.Printf("cur input is :%s", s)
	//	if strings.ToUpper(s) == "Q" {
	//		return
	//	}
	//	_, err = conn.Write([]byte(s))
	//	if err != nil {
	//		fmt.Printf("send failed, err:%v\n", err)
	//		return
	//	}
	//	// 从服务端接收回复消息
	//	var buf [1024]byte
	//	n, err := conn.Read(buf[:])
	//	if err != nil {
	//		fmt.Printf("read failed:%v\n", err)
	//		return
	//	}
	//	fmt.Printf("收到服务端回复:%v\n", string(buf[:n]))
	//}

	for i := 0; i < 300; i++ {
		s := strconv.Itoa(i) + "-【我是一条测试消息，我是一条测试消息，我是一条测试消息，我是一条测试消息，我是一条测试消息，我是一条测试消息】"
		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Printf("send failed, err:%v\n", err)
			return
		}
		//// 从服务端接收回复消息
		//var buf [1024]byte
		//n, err := conn.Read(buf[:])
		//if err != nil {
		//	fmt.Printf("read failed:%v\n", err)
		//	return
		//}
		//fmt.Printf("收到服务端回复:%v\n", string(buf[:n]))
	}
}
