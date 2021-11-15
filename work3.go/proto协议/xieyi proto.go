package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

//Encode将消息编码
func Encode(message string) ([]byte, error) {
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
