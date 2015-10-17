package gochat

import (
	//"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type Conn struct {
	conn net.Conn
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		conn: conn,
	}
	return c
}

func (c *Conn) Recv() ([]byte, error) {
	var length uint32

	blen := make([]byte, 4)
	n, err := c.conn.Read(blen)
	if err != nil {
		return nil, err
	}

	if n != 4 {
		fmt.Printf("head读取错误 n=%d", n)
		return nil, errors.New("head读取错误")
	}

	buf := bytes.NewReader(blen)
	err = binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		return nil, err
	}

	if length < 0 {
		fmt.Println("长度不正确！")
		return nil, errors.New("length 错了！")
	}

	blen2 := make([]byte, length)
	n, err = c.conn.Read(blen2)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		fmt.Printf("payload 读取错误 n=%d, length=%d", n, length)
		return nil, errors.New("payload 读取错误")
	}

	return blen2[:n], nil

}

func (c *Conn) Send(payload []byte) (int, error) {
	length := uint32(len(payload))
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		fmt.Println("binary.Write failed", err)
		return 0, err
	}

	err = binary.Write(buf, binary.LittleEndian, payload)
	if err != nil {
		fmt.Println("binary.Write failed", err)
		return 0, err
	}

	return c.conn.Write(buf.Bytes())

}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
