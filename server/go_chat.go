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
	conn_map := make(map[string]net.Conn)

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

func handleConn(conn_map map[string]net.Conn, client net.Conn) {
	//	bufio.NewReader()创建一个默认大小的readbuf
	//	b := bufio.NewReader(client)

	for {
		//		line, err := b.ReadBytes('\n')

		conn := gochat.NewConn(client)
		line, err := conn.Recv()
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
			fmt.Println(" =>\n", c_k, c_v)
			if client == c_v {
				seek_status = true
				user = c_k + " say: "
				fmt.Println("find_user:\n", user)
				break
			}
		}

		if seek_status {
			for _, client_v := range conn_map {
				//				client_v.Write([]byte(user + data + "\n"))
				payload := []byte(user + data + "\n")
				conn_v := gochat.NewConn(client_v)
				conn_v.Send(payload)

			}
		} else {
			conn_map[data] = client
		}

	}
}
