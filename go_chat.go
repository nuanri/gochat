package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const PORT = 3540

var conn_slice []net.Conn

func main() {
	//建立 socket, 监听端口
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		fmt.Println(err)
		panic("couldn't start listening: ")
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
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
			conn_slice = append(conn_slice, client)
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	//	bufio.NewReader()创建一个默认大小的readbuf
	b := bufio.NewReader(client)

	for {
		line, err := b.ReadBytes('\n')
		fmt.Println("收到:", string(line[:len(line)-1]), err)
		if err != nil { // EOF, or worse
			break
		}
		user := fmt.Sprintf("来自 %s 说: ", client.RemoteAddr())
		for _, client_v := range conn_slice {
			client_v.Write([]byte(user + string(line)))
		}
		//		client.Write(line)
	}
}
