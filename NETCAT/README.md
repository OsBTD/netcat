# TCPChat - A Simple TCP Chat Application (NetCat)

## AUTHORS
ayoub NACHTI
Zakarai DIOURI
Youssef JAOUHAR

## Overview
This project recreates the functionality of NetCat (`nc`) in a Server-Client architecture. It allows multiple clients to connect to a server and communicate via TCP, similar to how NetCat works. The server listens on a specified port, and clients can send messages that are broadcasted to everyone connected.

## Features
- **TCP connection** between server and multiple clients.
- **Name requirement** for clients.
- **Message formatting**: Timestamp and sender's name included.
- **Client join/leave notifications** to other clients.
- **Broadcasting messages** to all clients.
- **Server supports up to 10 clients**.
  
## Usage

### Running the Server:
```bash
$ go run .
Listening on the port :8989 

or

$ go run . $port
Listening on the port $port

```
