package main

import (
	"dreyspi/godis/storage"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	getCommand    = "get"
	setCommand    = "set"
	deleteCommand = "delete"
)

const (
	stringType = "string"
	dictType   = "dict"
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
		switch tokens[1] {
		case stringType:
			return handleGetString(tokens)
		case dictType:
			return handleGetDict(tokens)
		default:
			return typeNotSupportedResponse, false
		}
	case setCommand:
		switch tokens[1] {
		case stringType:
			return handleSetString(tokens)
		case dictType:
			return handleSetDict(tokens)
		default:
			return typeNotSupportedResponse, false
		}
	case deleteCommand:
		return handleDeleteKey(tokens)
	default:
		return commandNotSupportedResponse, false
	}
}

func handleSetDict(tokens []string) (string, bool) {
	if len(tokens) < 4 {
		return notEnoughArgumentsResponse, false
	}

	key := tokens[2]
	if key == "" {
		return emptyKeyResponse, false
	}

	lastIsTtl := !strings.Contains(tokens[len(tokens)-1], ":")
	fieldsRemaining := len(tokens) - 3
	if lastIsTtl {
		fieldsRemaining--
	}
	if fieldsRemaining <= 0 {
		return notEnoughArgumentsResponse, false
	}

	var ttl time.Duration
	if lastIsTtl {
		var err error
		ttl, err = parseTtl(tokens[len(tokens)-1])
		if err != nil {
			return fmt.Sprintf("%s: %v", ttlParseErrorResponse, err), false
		}
	}

	fields := map[string]string{}
	for _, field := range tokens[3 : 3+fieldsRemaining] {
		s := strings.Split(field, ":")
		if len(s) != 2 {
			return invalidDictVInputResponse, false
		}
		fields[s[0]] = s[1]
	}

	err := godisStorage.SetDict(key, fields, ttl)
	if err != nil {
		return storageFailureResponse, false
	}

	return "", true
}

func handleGetDict(tokens []string) (string, bool) {
	if len(tokens) < 3 {
		return notEnoughArgumentsResponse, false
	}

	key := tokens[2]
	if key == "" {
		return emptyKeyResponse, false
	}

	fields := []string{}
	if len(tokens) > 3 {
		fields = tokens[3:]
	}

	value, err := godisStorage.GetDict(key, fields...)
	if err != nil {
		return storageFailureResponse, false
	}
	if value == nil {
		return storage.NullValue, true
	}

	sb := strings.Builder{}
	for field, v := range value {
		sb.WriteString(field)
		sb.WriteString(":")
		sb.WriteString(v)
		sb.WriteString(CRLF)
	}

	return sb.String(), true
}

func handleDeleteKey(tokens []string) (string, bool) {
	if len(tokens) < 2 {
		return notEnoughArgumentsResponse, false
	}
	if len(tokens) > 2 {
		return tooManyArgumentsResponse, false
	}

	key := tokens[1]
	if key == "" {
		return emptyKeyResponse, false
	}

	err := godisStorage.DeleteKey(key)
	if err != nil {
		return storageFailureResponse, false
	}

	return "", true
}

func handleGetString(tokens []string) (string, bool) {
	if len(tokens) < 3 {
		return notEnoughArgumentsResponse, false
	}
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
	if len(tokens) < 4 {
		return notEnoughArgumentsResponse, false
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
