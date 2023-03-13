package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	username string
	password string
	conn     net.Conn
}

var clients []*User
var clientCount int32

func main() {
	// initialize server
	fmt.Println("Server listening for incoming connections...")
	listener, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		log.Fatalln(err)
	}

	// its procedure to close listener
	defer listener.Close()

	for {
		// we start accepting connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("New connection from ", conn.RemoteAddr())

		var user User
		user.conn = conn

		buffer := make([]byte, 1400)
		dataSize, err := user.conn.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		info := string(buffer[:dataSize])
		fmt.Println(info)
		userInfo := strings.SplitAfterN(info, " ", 3)
		username := userInfo[0]
		fmt.Println(username)
		password := userInfo[1]
		fmt.Println(password)
		user.username = username
		user.password = password

		//buffer = make([]byte, 1400)
		//dataSize, err = user.conn.Read(buffer)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//
		//password := buffer[:dataSize]
		//user.password = string(password)

		clients = append(clients, &user)

		// thread for listening to conns
		go listenConnection(&user)
	}
}

func listenConnection(client *User) {
	for {
		buffer := make([]byte, 1400)
		dataSize, err := client.conn.Read(buffer)
		if err != nil {
			fmt.Println("ERROR: Connection failed.")
			return
		}

		data := buffer[:dataSize]
		broadcast(client.username, data, client.conn)
		fmt.Println(bytes.NewBuffer(data).String())
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func broadcast(clientname string, msg []byte, conn net.Conn) {
	for i := 0; i < len(clients); i++ {
		// we dont want to send message to original user
		// so we skip that client
		if clients[i].conn != conn {
			_, _ = clients[i].conn.Write([]byte(clientname))
			_, _ = clients[i].conn.Write(msg)
		}
		continue
	}
}
