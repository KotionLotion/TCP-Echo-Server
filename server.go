package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	inactivityTimeout = 30 * time.Second
	maxMessageSize    = 1024
)

var (
	port = flag.Int("port", 4000, "Port to listen on")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", *port, err)
	}
	defer listener.Close()

	log.Printf("Server listening on port %d", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	logFile := createLogFile(clientAddr)
	defer logFile.Close()

	log.Printf("Client connected: %s", clientAddr)
	defer func() {
		conn.Close()
		log.Printf("Client disconnected: %s", clientAddr)
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		conn.SetDeadline(time.Now().Add(inactivityTimeout))

		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client %s closed the connection", clientAddr)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Client %s timed out", clientAddr)
				writer.WriteString("Connection timed out due to inactivity\n")
				writer.Flush()
			} else {
				log.Printf("Error reading from client %s: %v", clientAddr, err)
			}
			return
		}

		message = strings.TrimSpace(message)
		log.Printf("Received from %s: %s", clientAddr, message)
		fmt.Fprintf(logFile, "%s: %s\n", time.Now().Format(time.RFC3339), message)

		response, closeConn := processMessage(message)
		writer.WriteString(response + "\n")
		writer.Flush()

		if closeConn {
			return
		}
	}
}

func createLogFile(clientAddr string) *os.File {
	//Clean the address
	safeAddr := strings.ReplaceAll(clientAddr, ":", "_")
	safeAddr = strings.ReplaceAll(safeAddr, ".", "_")
	filename := fmt.Sprintf("client_%s.log", safeAddr)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create log file for %s: %v", clientAddr, err)
		return nil
	}

	return file
}

func processMessage(message string) (string, bool) {
	// empty message
	if message == "" {
		return "Say something...", false
	}

	// commands
	if strings.HasPrefix(message, "/") {
		return handleCommand(message)
	}

	// special messages
	switch strings.ToLower(message) {
	case "hello":
		return "Hello Wor.... i mean Hello there!", false
	case "bye":
		return "Leaving so soon? Goodbye!", true
	default:
		return message, false
	}
}

func handleCommand(cmd string) (string, bool) {
	parts := strings.SplitN(cmd, " ", 2)
	command := strings.ToLower(parts[0])

	switch command {
	case "/time":
		return time.Now().Format(time.RFC1123), false
	case "/quit":
		return "Closing connection", true
	case "/echo":
		if len(parts) > 1 {
			return parts[1], false
		}
		return "Usage: /echo <message>", false
	default:
		return fmt.Sprintf("Unknown command: %s", command), false
	}
}
