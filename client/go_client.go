package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"

	"nuanri/gochat"
)

var name string

func Recived(raw_conn net.Conn) {

	go Recived_back(raw_conn)

	r := bufio.NewReader(os.Stdin)

	for {
		msg, err := r.ReadString('\n')

		if err != nil {
			fmt.Printf("读取用户输入数据错误 %s", err)
			break
		}

		msg = msg[:len(msg)-1]
		sigbody := gochat.SendBody{To: "all", Msg: msg}
		s := gochat.GetJson(sigbody)

		reg := gochat.HiMsg{Cmd: "sendmessage", Body: s}
		regjson := gochat.GetJson(reg)
		//		conn.Write(line)
		conn := gochat.NewConn(raw_conn)
		conn.Send(regjson)
	}
}

func Recived_back(conn net.Conn) {

	conn_c := gochat.NewConn(conn)

	for {
		data, err := conn_c.Recv()
		if err != nil {
			fmt.Printf("读取返回数据错误 %s", err)
			break
		}

		msg := gochat.HiMsg{}
		err = gochat.ParseMsg(data, &msg)
		if err != nil {
			fmt.Println("输入格式错误 :", err)
			fmt.Println("  data :", string(data))
			continue
		}

		body := gochat.SendBody{}
		err = gochat.ParseMsg(msg.Body, &body)
		if err != nil {
			fmt.Println("输入格式错误 :", err)
			fmt.Println("  msg.Body :", msg.Body)
			continue
		}

		d := color.New(color.FgCyan, color.Bold)

		d.Print(body.From)
		fmt.Println(" : ", body.Msg)

	}
}

func Registered(raw_conn net.Conn) {

	sigbody := gochat.SigBody{Name: name}
	s := gochat.GetJson(sigbody)

	reg := gochat.HiMsg{Cmd: "signup", Body: s}
	regjson := gochat.GetJson(reg)
	//		conn.Write(line)
	conn := gochat.NewConn(raw_conn)
	conn.Send(regjson)

}

func parse_ops() {
	flag.StringVar(&name, "name", "hichat", "设置用户名")
	flag.Parse()
}

func main() {
	parse_ops()

	conn, err := net.Dial("tcp", "127.0.0.1:3540")

	if err != nil {
		fmt.Printf("Failure to connent:%s\n", err)
		return
	} else {
		fmt.Println("connected!\n")
	}

	Registered(conn)
	Recived(conn)

}
