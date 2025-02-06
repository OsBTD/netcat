package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clients []net.Conn
	mutex   sync.Mutex
	users   = make(map[net.Conn]string)
	history []string
)

func main() {
	port := ":8989"
	args := os.Args
	if len(args) != 2 {
		fmt.Print([]byte("		[USAGE]: ./TCPChat $port\n"))
		return
	} else {
		if strings.HasPrefix(args[1], ":") {
			port = args[1]

		} else {
			port = ":" + args[1]
		}
	}
	const Maxclients = 10
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("error listening", err)
	}
	log.Println("server started at port " + port)
	defer listener.Close()
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			log.Println("error accepting", err)
			continue
		}
		mutex.Lock()
		if len(clients) >= Maxclients {
			conn.Write([]byte("maximum number of connections reached\nplease try again..\n"))
			mutex.Unlock()
			conn.Close()
		}
		clients = append(clients, conn)

		mutex.Unlock()

		go Handle(conn)
	}
}
func Handle(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "))
a:

	bufferName := make([]byte, 1024)
	m, _ := conn.Read(bufferName)
	username := strings.TrimSpace(string(bufferName[:m]))
	var uservalid bool = true
	formattedTime := time.Now().Format("02-01-2006 15:04:05")

	if len(username) < 3 || len(username) > 30 {
		conn.Write([]byte("invalid username, try again : \nenter a new username : "))
		uservalid = false
		goto a
	}

	for _, exists := range users {
		if username == exists {
			conn.Write([]byte("username already exists, try again : \nenter a new username : "))
			uservalid = false
			goto a

		}
	}

	for _, char := range username {
		if char >= '~' || char <= ' ' {
			conn.Write([]byte("invalid username, try again : \nenter a new username : "))
			uservalid = false
			goto a

		}
	}

	users[conn] = username
	for _, entry := range history {
		conn.Write([]byte(entry))
	}
	for _, client := range clients {
		if client != conn && len(users[client]) > 0 && uservalid {
			log.Println("joined" + users[conn])
			client.Write([]byte("\n\033[32m" + users[conn] + " has joined the chat.\033[0m\n"))
			client.Write([]byte("[" + formattedTime + "][" + users[client] + "]: "))
		}
	}

	for {
		conn.Write([]byte("[" + formattedTime + "][" + users[conn] + "]: "))
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		message := strings.TrimSpace(string(buffer[:n]))

		if len(message) > 0 {
			history = append(history, "["+formattedTime+"]["+users[conn]+"]: "+message+"\n")
		}
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "wsarecv: An existing connection was forcibly closed by the remote host") {
				log.Println(users[conn] + " has left the chat.")
				mutex.Lock()
				for i, client := range clients {
					if client == conn {
						clients = append(clients[:i], clients[i+1:]...)
						break
					}
				}

				for _, client := range clients {
					if conn != client && len(users[client]) > 0 && uservalid {
						log.Println("disconnected" + users[conn])

						client.Write([]byte("\n\033[31m" + users[conn] + " has left the chat.\033[0m\n"))
						client.Write([]byte("[" + formattedTime + "][" + users[client] + "]: "))
					}
				}

				mutex.Unlock()
				conn.Close()
				break

			}

			log.Println("error reading", err)
			return
		}
		log.Println("message received", message)

		mutex.Lock()

		for _, client := range clients {
			if conn == client {
				continue
			} else if conn != client && len(users[client]) > 0 && len(message) > 0 && uservalid {
				log.Println("message" + users[conn])

				_, err3 := client.Write([]byte("\n[" + formattedTime + "] [" + users[conn] + "]: " + message + "\n[" + formattedTime + "] [" + users[client] + "]: "))
				if err3 != nil {
					log.Println("error sending message to client", err)
				}
			}
		}

		// conn.Write([]byte("[" + formattedTime + "][" + users[conn] + "]: ")) // Send prompt again to the sender after the message is processed

		mutex.Unlock()
	}
}
