package gochat

import (
	"fmt"
	"net"
	"testing"
)

var testData []string = []string{
	"hello1",
	"hello2",
	"hello3",
	"hello4",
	"hello6",
}

func Test_Conn(t *testing.T) {

	stop := make(chan bool, 1)

	l, err := net.Listen("tcp", "127.0.0.1:3540")
	if err != nil {
		t.Error("can not listen:", err)
		return
	}

	go func() {
		raw_conn, err := l.Accept()
		if err != nil {
			t.Error("accept client error:", err)
			return
		}
		// 处理 client 链接
		go func() {
			defer raw_conn.Close()
			conn := NewConn(raw_conn)
			for i := 0; ; i++ {
				payload, err := conn.Recv()
				if err != nil {
					t.Error("recv error:", err)
					break
				}
				if string(payload) == "quit" {
					fmt.Println("got quit signal")
					break
				}
				fmt.Println("recv:", string(payload))
				if string(payload) != testData[i] {
					t.Error("recv failed:", string(payload))
				}
			}

			stop <- true
		}()
	}()

	client_raw_conn, err := net.Dial("tcp", "127.0.0.1:3540")
	if err != nil {
		t.Error("dial error:", err)
		return
	}
	defer client_raw_conn.Close()

	client_conn := NewConn(client_raw_conn)
	for _, payload := range testData {
		client_conn.Send([]byte(payload))
	}
	client_conn.Send([]byte("quit"))

	<-stop
}
