package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	clients []net.Conn
	mutex   sync.Mutex
	users   = make(map[net.Conn]string)
)

func main() {
	const Maxclients = 10
	listener, err := net.Listen("tcp", ":8989")
	if err != nil {
		log.Println("error listening", err)
	}
	log.Println("server started at port 8989")
	defer listener.Close()
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			log.Println("error accepting", err)
			continue
		}
		mutex.Lock()
		if len(clients) >= Maxclients {
			conn.Write([]byte("maximum number of connections reached\nplease try again later..\n"))
			mutex.Unlock()
			conn.Close()
			continue
		}
		clients = append(clients, conn)

		mutex.Unlock()

		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "))
	bufferName := make([]byte, 1024)
	m, _ := conn.Read(bufferName)
	username := bufferName[:m]
	users[conn] = strings.TrimSpace(string(username))

	for _, client := range clients {
		if client != conn && len(users[client]) > 0 {
			log.Println("joined" + users[conn])
			client.Write([]byte("\n" + users[conn] + " has joined the chat\n"))
		}
	}
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		message := buffer[:n]
		if err != nil {
			if err == io.EOF {
				log.Println(users[conn] + " has disconneted\n")
				mutex.Lock()
				for i, client := range clients {
					if client == conn {
						clients = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				for _, client := range clients {
					if conn != client && len(users[client]) > 0 {
						log.Println("disconnected" + users[conn])

						client.Write([]byte("\n" + users[conn] + " has disconneted\n"))
					}
				}

				mutex.Unlock()
				conn.Close()
				break

			}

			log.Println("error reading", err)
			return
		}
		log.Println("message received", string(message))

		mutex.Lock()

		for _, client := range clients {
			if conn == client {
				continue
			} else if conn != client && len(users[client]) > 0 && len(message) > 0 {
				log.Println("message" + users[conn])

				_, err3 := client.Write([]byte("\n" + users[conn] + ":  " + string(message) + "\n"))
				if err3 != nil {
					log.Println("error sending message to client", err)
				}

			} else if conn == client && len(message) == 0 {
				client.Write([]byte("\n"))
				continue
			}
		}
		mutex.Unlock()

	}
}
