package main

import (
	//	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	//"github.com/nuanri/go-chat"
	"nuanri/gochat"
)

const PORT = 3540

//var conn_map map[string]net.Conn

func main() {
	//建立 socket, 监听端口
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		fmt.Println(err)
		panic("couldn't start listening: ")
	}
	conns := clientConns(server)
	//初始化 mapx
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
	//	bufio.NewReader()创建一个默认大小的readbuf
	//	b := bufio.NewReader(client)
	client := gochat.NewConn(client_raw)

	for {
		//		line, err := b.ReadBytes('\n')

		line, err := client.Recv()
		fmt.Println(line)
		if err != nil {
			fmt.Printf("client %s was disconnected!\n", client.RemoteAddr())
			break
		}
		data := string(line[:len(line)-1])
		data = strings.TrimSpace(data)
		fmt.Println("recv data:", data)

		seek_status := false

		for _, c_v := range conn_map {
			if client == c_v {
				seek_status = true
				break
			}
		}

		if seek_status {
			message, err := gochat.ParseMsg(data)
			if err != nil {
				fmt.Println("输入格式错误 :", err)
				fmt.Println("  data :", data)
				client.Send([]byte(`{"name": "我", "msg": "笨蛋"}`))

				continue
			}

			if message.Name == "all" {
				for _, client_v := range conn_map {
					payload := []byte(message.JSON())
					client_v.Send(payload)

				}
			} else {
				client_v, ok := conn_map[message.Name]
				if !ok {
					fmt.Println("没有找到该用户", message.Name)
					client.Send([]byte("别烦我！"))
					continue
				}
				payload := []byte(message.JSON())
				client_v.Send(payload)
			}

		} else {
			conn_map[data] = client
		}

	}
}
