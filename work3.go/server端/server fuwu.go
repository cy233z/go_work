package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go_wlbc/RPC/proto"
	"io"
	"net"
)

//返回resp
type GetUserResp struct {
	UserID   int64  `json:"userId"`
	UserName string `json:"userName"`
}

func process(conn net.Conn) {
	defer conn.Close()  //关闭连接
	reader := bufio.NewReader(conn)
	for i := 0; i < 10; i++ {
		msg, err := proto.Decode(reader)//将消息解码
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed,err:", err)
			return
		}
		fmt.Println("收到client发来的数据:", msg)

		
		resp := GetUserResp{
			UserID:   13004951,
			UserName: "cyz",
		}
		a, err := json.Marshal(resp)//将返回的resp转化为json格式的字节数组
		if err != nil {
			fmt.Println("marshal failed,err:", err)
			return
		}
		fmt.Println(a)
		defer conn.Close()
		for i := 0; i < 1; i++ {
			data, err := proto.Encode(string(a))
			if err != nil {
				fmt.Println("encode failed,err:", err)
				return
			}
			conn.Write(data)
		}
	}
}
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:57073")
	if err != nil {
		fmt.Println("listen failed,err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()//与client端建立连接
		if err != nil {
			fmt.Println("accept failed,err:", err)
			continue
		}
		go process(conn)//启动一个goroutine处理连接
	}
}