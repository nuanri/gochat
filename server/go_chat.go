package main

import (
	//	"bufio"
	"encoding/json"
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
		user := ""

		for c_k, c_v := range conn_map {
			fmt.Println(" =>", c_k, c_v)
			if client == c_v {
				seek_status = true
				user = c_k + " say: "
				fmt.Println("find_user:", user)
				break
			}
		}

		if seek_status {
			message := ParseMsg(data)

			if message.Name == "all" {
				for _, client_v := range conn_map {
					payload := []byte(user + message.Msg + "\n")
					client_v.Send(payload)

				}
			} else {
				client_v, ok := conn_map[message.Name]
				if !ok {
					fmt.Println("没有找到该用户", message.Name)
					client.Send([]byte("别烦我！"))
					continue
				}
				payload := []byte(user + message.Msg + "\n")
				client_v.Send(payload)
			}

		} else {
			conn_map[data] = client
		}

	}
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func ParseMsg(line string) Message {

	var message Message

	data := []byte(line)
	fmt.Println("传到解析json里的字符串", data)
	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("json 解析出错", err)

	}

	return message
}
