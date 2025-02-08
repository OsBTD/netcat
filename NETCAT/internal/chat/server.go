package chat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"net-cat/internal/helpers"
	"net-cat/internal/logger"
)

// Server represents the chat server which manages connections, users, and broadcasts.
type Server struct {
	Listener        net.Listener
	Users           Users
	Broadcast       chan BroadcastDetails
	HistoryMessages []string
	Shutdown        bool
	mu              sync.Mutex
	*logger.Loggers
}

// BroadcastDetails represents a message that needs to be broadcast to other users.
type BroadcastDetails struct {
	Message string
	User    string
}

// Users is a map that holds the users and their connections.
type Users map[string]net.Conn

// NewServer initializes a new Server with default values.
func NewServer() *Server {
	return &Server{
		Users:           Users{},
		Shutdown:        false,
		Broadcast:       make(chan BroadcastDetails),
		HistoryMessages: []string{},
		mu:              sync.Mutex{},
		Loggers:         logger.SetLoggers(),
	}
}

// Start begins listening on the provided port and handles incoming connections.
func (s *Server) Start(port string, kill chan os.Signal) {
	var err error
	s.Listener, err = net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		s.LogError.Println("Server failed to start, Error listening:", err.Error())
		fmt.Println("[USAGE]: ./TCPChat $port")
		kill <- os.Kill
		return
	}
	s.LogInfo.Println("Chat Server Started : server listening for connections on " + s.Listener.Addr().String())
	fmt.Println("Listening on the port :" + port)

	go s.brodcast()
	for !s.Shutdown {
		conn, err := s.Listener.Accept()
		if err != nil {
			s.LogError.Println("Error listening:", err.Error())
			continue
		}
		go s.HandleConnection(conn) // Handle each connection concurrently.
	}
	s.Listener.Close()
}

// HandleConnection handles an individual user connection.
func (s *Server) HandleConnection(conn net.Conn) {
	if s.Shutdown {
		return
	}

	s.LogInfo.Println("New connection from", conn.RemoteAddr().String())

	userName, err := s.AddUser(conn)
	if err != nil {
		s.LogError.Println("Error adding user:", err)
		return
	}

	fmt.Println(userName + " has joined the chat...")
	s.LogInfo.Println(userName + " has joined the chat.")

	s.Broadcast <- BroadcastDetails{
		Message: "\033[32m" + userName + " has joined the chat...\033[0m\n",
		User:    userName,
	}

	defer s.Removeuser(userName)

	s.HandleMessages(conn, userName)
}

// AddUser adds a new user to the server and asks for their name.
func (s *Server) AddUser(conn net.Conn) (string, error) {
	welcomeMessage := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    .       | ' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     -'       --'\n"
	conn.Write([]byte(welcomeMessage + "\n"))

	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		name = name[:len(name)-1]
		// Validate the name and ensure it's unique.
		s.mu.Lock()
		if len(s.Users) < 10 {
			err := helpers.ValidUserName(name)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				s.mu.Unlock()
				continue
			}
			if _, exists := s.Users[name]; !exists {
				// Add user to the server.
				s.Users[name] = conn
				conn.Write([]byte(strings.Join(s.HistoryMessages, "")))
				s.mu.Unlock()
				return name, nil
			}
			conn.Write([]byte("Username already exists\n"))
			s.mu.Unlock()
			continue
		} else {
			conn.Write([]byte("Server is full\n"))
			conn.Close()
			s.mu.Unlock()
			return "", errors.New("server full")
		}
	}
}

// Removeuser removes a user from the server and notifies others.
func (s *Server) Removeuser(user string) {
	s.Broadcast <- BroadcastDetails{
		Message: "\033[31m" + user + " has left the chat.\033[0m\n",
	}
	s.mu.Lock()
	s.Users[user].Close()
	delete(s.Users, user)
	s.mu.Unlock()
}

// HandleMessages listens for and processes incoming messages from a user.
func (s *Server) HandleMessages(conn net.Conn, userName string) {
	for {
		if s.Shutdown {
			conn.Close()
			break
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Log when a user leaves.
				s.LogInfo.Println(userName + " has left the chat ...")
				fmt.Println(userName + " has left the chat")
				return
			}
			fmt.Println("ERR TO READ from : ", userName, err)
			return
		}
		message, err = helpers.ValidMessage(strings.TrimSpace(message))
		if err == nil && message != "" {
			// Broadcast valid message to all users.
			s.Broadcast <- BroadcastDetails{
				Message: helpers.SetPrefix(userName) + message + "\n",
				User:    userName,
			}
		} else if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			conn.Write([]byte(helpers.SetPrefix(userName)))
		} else {
			conn.Write([]byte(helpers.SetPrefix(userName)))
		}
	}
}

// brodcast listens for messages to broadcast to all users.
func (s *Server) brodcast() {
	for brodcast := range s.Broadcast {
		s.mu.Lock()
		for user, conn := range s.Users {
			if user != brodcast.User {
				conn.Write([]byte("\033[s\033[2K\r")) // Clear line and move cursor.
				conn.Write([]byte(brodcast.Message))
			}
			conn.Write([]byte(helpers.SetPrefix(user)))
			if user != brodcast.User {
				conn.Write([]byte("\033[u\033[B")) // Restore cursor and move down.
			}
		}
		s.HistoryMessages = append(s.HistoryMessages, brodcast.Message)
		s.mu.Unlock()
	}
}

// Stop gracefully shuts down the server and notifies all users.
func (s *Server) Stop() {
	s.mu.Lock()
	s.LogInfo.Println("Stopping the server ...")
	defer s.mu.Unlock()
	s.Shutdown = true
	close(s.Broadcast)

	for username, conn := range s.Users {
		s.LogInfo.Println(username, " exited the chat server")
		conn.Write([]byte("\nServer has been stopped."))
	}
	s.LogInfo.Print("Server has been stopped\n\n")
	fmt.Println("\rServer has been stopped.")
}
