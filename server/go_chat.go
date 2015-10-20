package main

import (
	//	"bufio"
	"fmt"
	"net"
	"strconv"
	//	"strings"

	//"github.com/nuanri/go-chat"
	"nuanri/gochat"
)

const PORT = 9999

//var conn_map map[string]net.Conn

func main() {
	//建立 socket, 监听端口
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		fmt.Println(err)
		panic("couldn't start listening: ")
	}
	conns := clientConns(server)
	//初始化 map
	conn_map := make(map[string]*gochat.Conn)

	for {
		conn := <-conns
		go handleConn(conn_map, conn)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)

	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Println("couldn't accept: ", err)
				continue
			}

			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(conn_map map[string]*gochat.Conn, client_raw net.Conn) {

	client := gochat.NewConn(client_raw)
	pool := gochat.NewClientPool(conn_map)

	for {
		data, err := client.Recv()
		if err != nil {
			fmt.Printf("client %s was disconnected!\n", client.RemoteAddr())
			break
		}
		fmt.Println("data =", string(data))
		himsg := &gochat.HiMsg{}
		err = gochat.ParseMsg(data, himsg)
		if err != nil {
			fmt.Println("输入格式错误 :", err)
			fmt.Println("  data :", data)
			client.Send([]byte(`{"name": "我", "msg": "笨蛋"}`))
			continue
		}

		switch true {
		case himsg.Cmd == "signup":
			sigb := &gochat.SigBody{}
			err := gochat.ParseMsg([]byte(himsg.Body), sigb)
			if err != nil {
				fmt.Println("signbody格式错误")
				break
			}
			//往 map 里添加 client
			fmt.Println("sigb.Name=", sigb.Name)
			pool.Add(sigb.Name, client)

		case himsg.Cmd == "sendmessage":
			sendb := &gochat.SendBody{}
			err := gochat.ParseMsg([]byte(himsg.Body), sendb)
			if err != nil {
				fmt.Println("sendbody格式错误")
				break
			}

			if sendb.To == "all" {
				username := pool.GetByConn(client)

				r_body := &gochat.SendBody{From: username, Msg: sendb.Msg}
				r_body_b := gochat.GetJson(r_body)
				himsg.Body = r_body_b

				pool.SendToAll(himsg)

			} else {
				client_v, ok := pool.GetByUsername(sendb.To)
				if !ok {
					fmt.Println("没有找到该用户", sendb.To)
					client.Send([]byte("别烦我！"))
					continue
				}
				payload := gochat.GetJson(himsg)
				client_v.Send(payload)
			}
		}

	}
}
