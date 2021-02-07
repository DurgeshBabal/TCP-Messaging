package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"test/models"
)

type server struct {
	ip      string
	port    string
	clients map[string]net.Conn
}

func NewServer(ip string, port string) *server {
	return &server{
		ip:      ip,
		port:    port,
		clients: make(map[string]net.Conn),
	}
}

func getClientList(clients map[string]net.Conn) (resp string) {
	for k, _ := range clients {
		resp += k + "\n"
	}

	return resp
}

func registerClient(clients map[string]net.Conn, pubKey string, c net.Conn) (resp string) {
	clients[pubKey] = c

	resp = "Client registered!"

	return resp
}

func forwardMessage(clients map[string]net.Conn, m models.Message, c net.Conn) (resp string) {
	targetConn := clients[m.Target]

	msg := models.Message{
		Operation: m.Operation,
		Value:     m.Value,
	}

	data := msg.Bytes()
	data = append(data, '~')

	_, err := targetConn.Write(data)
	if err != nil {
		log.Println(err)
	}

	return "Forwarded~"
}

func (s *server) handleOperation(c net.Conn, m models.Message) []byte {
	var resp string

	switch m.Operation {
	case "ClientList":
		resp = getClientList(s.clients)

	case "RegisterClient":
		resp = registerClient(s.clients, m.Value, c)

	case "ForwardMessage", "ClientResponse":
		resp = forwardMessage(s.clients, m, c)
		return []byte(resp)
	}

	msg := models.Message{
		Operation: m.Operation,
		Value:     resp,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	data = append(data, '~')

	return data
}

func readMessage(c net.Conn) (m models.Message, err error) {
	r := bufio.NewReader(c)

	netData, err := r.ReadString('~')
	if err != nil {
		return m, err
	}

	netData = strings.TrimSuffix(netData, "~")

	err = json.Unmarshal([]byte(netData), &m)

	return m, err
}

func (s *server) handleConnection(c net.Conn) {
	defer c.Close()
	defer func() {
		key := ""

		for k, v := range s.clients {
			if v == c {
				key = k
				break
			} 
		}

		delete(s.clients, key)
	}()
	
	for {
		message, err := readMessage(c)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}

		resp := s.handleOperation(c, message)

		_, err = c.Write(resp)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	s := NewServer("127.0.0.1", "5000")

	l, err := net.Listen("tcp", s.ip+":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TCP Server online!")

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleConnection(c)
	}
}
