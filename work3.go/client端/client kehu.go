package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go_wlbc/RPC/proto"
	"io"
	"net"
)

//请求rep
type GetUserReq struct {
	UserID int32 `json:"userId"`
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
		fmt.Println("收到server返回的数据", msg)
	}
}

func main() {
	req := GetUserReq{
		UserID: 13004951,
	}
	conn, err := net.Dial("tcp", "127.0.0.1:57073")
	if err != nil {
		fmt.Println("dial failed, err:", err)
		return
	}

	a, err := json.Marshal(req)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	fmt.Println(a)
	defer conn.Close()
	for  {
		data, err := proto.Encode(string(a))
		if err != nil {
			fmt.Println("encode failed,err:", err)
			return
		}
		conn.Write(data)
		go process(conn)
	}
}
