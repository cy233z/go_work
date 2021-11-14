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
	UserID   int32  `json:"userId"`
	UserName string `json:"userName"`
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for i := 0; i < 10; i++ {
		msg, err := proto.Decode(reader)
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
		a, err := json.Marshal(resp)
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
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed,err:", err)
			continue
		}
		go process(conn)
	}
}

