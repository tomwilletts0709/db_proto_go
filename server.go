package main

import (
	"fmt"
	"net"
	"strings"
)

type Server struct {
	aof *Aof
}

func NewServer(aofPath string) (*Server, error) {
	aof, err := NewAof(aofPath)
	if err != nil {
		return nil, err
	}
	return &Server{aof: aof}, nil
}

func (s *Server) Run(port string) error {
	fmt.Println("Listening on port", port)

	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	defer l.Close()
	defer s.aof.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			return
		}

		if value.typ != "array" {
			fmt.Println("invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			writer.Write(Value{typ: "error", str: "ERR unknown command '" + command + "'"})
			continue
		}

		if command == "SET" || command == "HSET" {
			s.aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
