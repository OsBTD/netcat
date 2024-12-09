package main

import (
	"io"
	"log"
	"net"
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
	// args := os.Args[1:]
	// if len(args) > 0 {
	// 	var good bool

	// 	for _, char := range args[0] {
	// 		if char > '0' && char > '9' && strings.HasPrefix(args[0], ":") {
	// 			good = true
	// 		}
	// 	}
	// 	if good {
	// 		port = args[0]
	// 	} else {
	// 		log.Println("invalid port, defaulting to port :8989")
	// 		port = ":8989"
	// 	}
	// }
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
			client.Write([]byte("\n" + users[conn] + " has joined the chat.\n\nenter your message : "))
		}
	}

	for {
		conn.Write([]byte("enter your message : "))
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		message := buffer[:n]
		formattedTime := time.Now().Format("02-01-2006 15:04:05")

		history = append(history, formattedTime+" "+users[conn]+" : "+string(message)+"\n")

		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "wsarecv: An existing connection was forcibly closed by the remote host"){
				log.Println(users[conn] + " has disconneted\n\nenter your message : ")
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

						client.Write([]byte("\n" + users[conn] + " has disconneted\n\nenter your message : "))
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
			} else if conn != client && len(users[client]) > 0 && len(message) > 0 && uservalid {
				log.Println("message" + users[conn])

				_, err3 := client.Write([]byte("\n" + formattedTime + " " + users[conn] + ":  " + string(message) + "\n\nenter your message : "))
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

// format messages + time done
// history convo done
// port args + handle errors + default port
//
