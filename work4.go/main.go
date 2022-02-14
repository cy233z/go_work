package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"
)

func Oddparity(Code byte) []byte { //奇校验
	a := strconv.FormatInt(int64(Code), 2) //返回Code的base进制的表示
	var count = 0
	for i := 0; i < len(a); i++ { //遍历1的个数
		if a[i] == '1' {
			count++
		}
	}
	s := []byte(a)
	if count%2 == 0 { //1的个数为偶数，则奇校验位为1
		s = append(s, '1')
	} else { //1的个数为奇数，则奇校验位为0
		s = append(s, '0')
	}
	return s
}

func check(msg []byte) (bool, int64) {
	count := 0
	for i := 0; i < len(msg); i++ {
		if msg[i] == '1' {
			count++
		}
	}
	if count%2 == 0 {
		return false, 0
	}
	msg = msg[:len(msg)-1]
	result, _ := strconv.ParseInt(string(msg), 2, 64)
	return true, result
}

//Encode将消息编码
func Encode(message []byte) ([]byte, error) {
	//读取消息的长度，转换为int64类型(占8个字节)
	var length = int64(len(message))
	var pkg = new(bytes.Buffer)
	//写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	//写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

//Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	//读取消息的长度
	lengthbyte, _ := reader.Peek(8)
	lengthBuff := bytes.NewBuffer(lengthbyte)
	var length int64
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	//Buffered返回缓冲中现有的可读取的字节数
	if int64(reader.Buffered()) < length+8 {
		return "", err
	}
	//读取真正的消息数据
	pack := make([]byte, int(8+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[8:]), nil
}

func producer(out chan<- string) { //只写单向通道
	var str string
	for i := 0; i < 3; i++ {
		fmt.Scan(&str)
		out <- str //只进行发送操作
		result := []byte(str)
		for j := 0; j < len(result); j++ {
			a := Oddparity(result[j]) //奇校验
			msg, _ := Encode(a)
			file, err := os.OpenFile("./mq/mq1.mq", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) //创建 只写 追加写入
			if err != nil {
				fmt.Println("open file failed!,err:", err)
				return
			}
			defer file.Close()
			_, err = file.WriteString(string(msg)) //写入文件中
			if err != nil {
				fmt.Println("failed to writestring!,err:", err)
				return
			}
		}
	}
	close(out) //关闭通道
}

func consumer(in <-chan string) { //只读单向通道
	for i := range in {
		fmt.Println(i)
	}
	for {
		time.Sleep(1*time.Second)
		file, err := os.Open("./mq/mq1.mq")
		if err != nil {
			fmt.Println("open file failed,err:", err)
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			msg, err := Decode(reader)
			if err != nil {
				fmt.Println("decode failed,err:", err)
				return
			}
			a:=[]byte(msg)
			b,c:=check(a)
			if c==0{
				fmt.Println("奇校验的结果为:",b)
			}
			fmt.Println(c)
		}
	}
}

func main() {
	ch:=make(chan string)
	go producer(ch)
	consumer(ch)
	err:=os.Remove("./mq/mq1.mq")
	if err != nil {
		fmt.Println("remove mq1.mq failed,err:",err)
	}
}
