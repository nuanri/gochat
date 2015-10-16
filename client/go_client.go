package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"

	"nuanri/gochat"
)

func Recived(conn net.Conn) {

	go Recived_back(conn)

	r := bufio.NewReader(os.Stdin)

	for {
		msg, err := r.ReadString('\n')

		if err != nil {
			fmt.Printf("读取用户输入数据错误 %s", err)
			break
		}

		line := []byte(msg)
		//		conn.Write(line)
		gochat.Send(conn, line)
	}
}

func Recived_back(conn net.Conn) {

	for {
		r_msg, err := gochat.Recv(conn)
		if err != nil {
			fmt.Printf("读取返回数据错误 %s", err)
			break
		}
		color.Cyan(string(r_msg))
	}
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:3540")

	if err != nil {
		fmt.Printf("Failure to connent:%s\n", err)
		return
	} else {
		fmt.Println("connected!\n")
	}

	Recived(conn)

}
