package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port 6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(conn net.Conn) {
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
					aof.Write(value)
				}

				result := handler(args)
				writer.Write(result)
			}
		}(conn)
	}
}
