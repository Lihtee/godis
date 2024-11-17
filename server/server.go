package main

import (
	"bufio"
	storage "dreyspi/godis/storage"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
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
)

const (
	getCommand = "get"
	setCommand = "set"
)

const (
	stringType = "string"
)

const trimSet = "\n\t\r"

func routeCommandToStorage(command string) (string, bool) {
	trimmedCommand := strings.Trim(command, trimSet)
	if trimmedCommand == "" {
		return emptyCommandResponse, false
	}

	tokens := strings.Split(trimmedCommand, " ")
	if len(tokens) == 0 {
		fmt.Printf("Something extremely weird happened here: %s(command), %s(trimmed command)", command, trimmedCommand)
		return emptyCommandResponse, false
	}

	switch tokens[0] {
	case getCommand:
		if len(tokens) < 3 {
			return notEnoughArgumentsResponse, false
		}
		switch tokens[1] {
		case stringType:
			return handleGetString(tokens)
		default:
			return typeNotSupportedResponse, false
		}
	case setCommand:
		if len(tokens) < 4 {
			return notEnoughArgumentsResponse, false
		}
		switch tokens[1] {
		case stringType:
			return handleSetString(tokens)
		default:
			return typeNotSupportedResponse, false
		}
	default:
		return commandNotSupportedResponse, false
	}
}

func handleGetString(tokens []string) (string, bool) {
	if len(tokens) > 3 {
		return tooManyArgumentsResponse, false
	}
	key := tokens[2]
	if key == "" {
		return emptyKeyResponse, false
	}

	value, err := godisStorage.GetString(key)
	if err != nil {
		return storageFailureResponse, false
	}

	return value, true
}

func handleSetString(tokens []string) (string, bool) {
	if len(tokens) > 5 {
		return tooManyArgumentsResponse, false
	}

	key := tokens[2]
	if key == "" {
		return emptyKeyResponse, false
	}

	value := tokens[3]
	var ttl time.Duration
	if len(tokens) == 5 {
		var err error
		ttl, err = parseTtl(tokens[4])
		if err != nil {
			return fmt.Sprintf("%s: %v", ttlParseErrorResponse, err), false
		}
	}

	err := godisStorage.SetString(key, value, ttl)
	if err != nil {
		return storageFailureResponse, false
	}

	return "", true
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

func parseTtl(ttl string) (time.Duration, error) {
	trimmedTtl := strings.Trim(ttl, " ")
	if trimmedTtl == "" {
		return 0, nil
	}

	unitLetter := trimmedTtl[len(trimmedTtl)-1]
	var unit time.Duration
	switch unitLetter {
	case 's':
		unit = time.Second
	case 'm':
		unit = time.Minute
	case 'h':
		unit = time.Hour
	case 'd':
		unit = time.Hour * 24
	case 'w':
		unit = time.Hour * 24 * 7
	default:
		return 0, fmt.Errorf("unrecognized unit letter: %s", trimmedTtl)
	}

	amountString := trimmedTtl[:len(trimmedTtl)-1]
	var amount int
	var err error
	if len(amountString) == 0 {
		amount = 1
	} else {
		amount, err = strconv.Atoi(amountString)
		if err != nil {
			return 0, fmt.Errorf("failed to parse amount: %s", amountString)
		}
	}

	return time.Duration(amount * int(unit)), nil
}
