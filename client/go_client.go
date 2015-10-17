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
		conn := gochat.NewConn(conn)
		conn.Send(line)
	}
}

func Recived_back(conn net.Conn) {

	conn_c := gochat.NewConn(conn)

	for {
		r_msg, err := conn_c.Recv()
		if err != nil {
			fmt.Printf("读取返回数据错误 %s", err)
			break
		}

		message, err := gochat.ParseMsg(string(r_msg))
		if err != nil {
			fmt.Println("输入格式错误 :", err)
			fmt.Println("  r_msg :", string(r_msg))

			continue
		}
		//		color.Magenta(message.Name)
		//		fmt.Println(message.Msg)
		d := color.New(color.FgCyan, color.Bold)
		d.Printf("%s say", message.Name)
		fmt.Println(" : ", message.Msg)
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
