package main

import (
	"bufio"
	storage "dreyspi/godis/storage"
	"fmt"
	"net"
	"strings"
)

const (
	okStatus                    = "OK"
	errorStatus                 = "ERROR"
	emptyCommandResponse        = "Empty Command"
	notEnoughArgumentsResponse  = "Not enough arguments"
	tooManyArgumentsResponse    = "Too Many Arguments"
	emptyKeyResponse            = "Empty Key"
	storageFailureResponse      = "Storage Failure" // Todo: be more specific
	typeNotSupportedResponse    = "Type Not Supported"
	commandNotSupportedResponse = "Command Not Supported"
	ttlParseErrorResponse       = "TTL Parse Error"
	invalidDictVInputResponse   = "Invalid dict input"
)

var godisStorage = storage.New(false)

func main() {
	// Start listening on a specific port
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port 12345")

	for {
		// Accept a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection accepted from", conn.RemoteAddr())

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

const NL = '\n'
const CRLF = "\r\n"

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()

	// Create a buffer to read data from the connection
	reader := bufio.NewReader(conn)
	for {
		// Read data from the connection
		message, err := reader.ReadString(NL)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Printf("Received message: %s", makeInvisibleCharsVisible(message))
		fmt.Println()
		response, ok := routeCommandToStorage(message)
		var status string
		if ok {
			status = okStatus
		} else {
			status = errorStatus
		}

		_, err = conn.Write([]byte(status + CRLF + response + CRLF))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}

func makeInvisibleCharsVisible(input string) string {
	var builder strings.Builder
	for _, r := range input {
		switch r {
		case '\n':
			builder.WriteString("\\n")
		case '\r':
			builder.WriteString("\\r")
		case '\t':
			builder.WriteString("\\t")
		default:
			if r < 32 || r == 127 {
				// Non-printable ASCII characters
				builder.WriteString(fmt.Sprintf("\\x%02x", r))
			} else {
				builder.WriteRune(r)
			}
		}
	}
	return builder.String()
}
