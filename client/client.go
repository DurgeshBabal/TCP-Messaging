package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"test/keygen"
	"test/models"
)

type client struct {
	pvtKey string
	pubKey string
}

func NewClient(pvtKey string, pubKey string) *client {
	return &client{
		pvtKey: pvtKey,
		pubKey: pubKey,
	}
}

const serverPort = ":5000"
const serverIP = "127.0.0.1"

func main() {
	// Generate new key pair
	pubKey, pvtKey, err := keygen.NewKey()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pubKey)

	// Connect to server
	client := NewClient(pvtKey, pubKey)

	c, err := net.Dial("tcp", serverIP+serverPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TCP Client online!")

	// Register client with server
	m := models.Message{
		Operation: "RegisterClient",
		Value:     client.pubKey,
	}

	writeMessage(m, c)

	go func() {
		serverReader := bufio.NewReader(c)

		for {
			m, err := readServerResponse(serverReader)
			if err != nil {
				if err == io.EOF {
					c.Close()
					os.Exit(0)
					return
				}

				log.Println(err)
				continue
			}

			client.handleOperation(m, c)
		}
	}()

	userReader := bufio.NewReader(os.Stdin)

	for {
		m, err := readUserInput(userReader)
		if err != nil {
			log.Println(err)
			continue
		}

		writeMessage(m, c)
	}

}

func readUserInput(r *bufio.Reader) (m models.Message, err error) {
	input, err := r.ReadString('~')
	if err != nil {
		return m, err
	}

	input = strings.TrimSuffix(input, "~")

	err = json.Unmarshal([]byte(input), &m)

	return m, err
}

func readServerResponse(r *bufio.Reader) (m models.Message, err error) {
	netData, err := r.ReadString('~')
	if err != nil {
		return m, err
	}

	netData = strings.TrimSuffix(netData, "~")

	fmt.Print("\nServer Response:\n" + netData)
	fmt.Print("\n>> ")

	err = json.Unmarshal([]byte(netData), &m)

	return m, err
}

func writeMessage(m models.Message, c net.Conn) {
	b := m.Bytes()
	b = append(b, '~')

	fmt.Fprintf(c, "%s", b)
}

func (client *client) handleOperation(m models.Message, c net.Conn) {
	switch m.Operation {
	case "ForwardMessage":
		msg := models.Message{
			Operation: "ClientResponse",
			Value:     client.pubKey,
			Target:    m.Value,
		}

		writeMessage(msg, c)
	}
}
