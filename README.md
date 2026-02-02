# Go Redis Rebuild

A Redis-like in-memory key-value store implementation in Go. This project is a learning exercise to understand how Redis works under the hood by rebuilding its core components from scratch.

## Overview

This project implements a simplified Redis server that:
- Listens on port 6379 (Redis default port)
- Parses and responds to RESP (REdis Serialization Protocol) commands
- Handles client connections over TCP
- Supports basic command parsing and response formatting

## Current Status

### Implemented
- **TCP Server**: Basic server that listens on port 6379 and accepts client connections
- **RESP Protocol Parser**: Full implementation of RESP protocol parsing including:
  - Arrays (`*`)
  - Bulk strings (`$`)
  - String responses (`+`)
  - Error responses (`-`)
  - Integer responses (`:`)

In Progress
- **AOF (Append Only File) Persistence**: Basic structure defined, implementation in progress
- **Command Handler**: Framework started, needs command routing logic
- **Server Architecture**: Core server structure being refined

Project Structure

```
.
├── main.go      # Entry point, TCP server setup and connection handling
├── resp.go      # RESP protocol parser implementation
├── aof.go       # AOF persistence (work in progress)
├── handler.go   # Command handler (work in progress)
└── server.go    # Server architecture (work in progress)
```

## How It Works

### RESP Protocol

The project implements the RESP (REdis Serialization Protocol) which is Redis's wire protocol. The parser can handle:

- **Arrays**: `*<number>\r\n<elements>`
- **Bulk Strings**: `$<length>\r\n<string>\r\n`
- **Simple Strings**: `+<string>\r\n`
- **Errors**: `-<string>\r\n`
- **Integers**: `:<number>\r\n`

### Server Flow

1. Server starts and listens on `0.0.0.0:6379`
2. Accepts incoming TCP connections
3. Reads RESP-formatted commands from clients
4. Parses commands using the RESP parser
5. Responds with appropriate RESP-formatted responses

## Building and Running

### Prerequisites
- Go 1.16 or later

### Build
```bash
go build -o redis-server
```

### Run
```bash
./redis-server
```

Or run directly:
```bash
go run .
```

The server will start listening on port 6379 and print:
```
Listening on port 6379
```

### Testing with Redis CLI

You can test the server using the official Redis CLI:

```bash
redis-cli -p 6379
```

Currently, the server responds with `+OK\r\n` to all commands. This is a placeholder while command handling is being implemented.

## Architecture

### RESP Parser (`resp.go`)
- `NewResp()`: Creates a new RESP parser from an `io.Reader`
- `Read()`: Main entry point that reads and parses RESP values
- `readArray()`: Parses RESP arrays recursively
- `readBulk()`: Parses RESP bulk strings
- `readInteger()`: Parses RESP integers
- `readLine()`: Low-level line reading utility

### Value Structure
The `Value` struct represents parsed RESP data:
```go
type Value struct {
    typ   string
    str   string
    num   int
    bulk  string
    array []Value
}
```

## Roadmap

- [ ] Complete AOF persistence implementation
- [ ] Implement command handler with routing
- [ ] Add support for core Redis commands (SET, GET, DEL, etc.)
- [ ] Implement in-memory data structures (hash tables, etc.)


## Learning Goals

This project serves as a hands-on way to understand:
- Network programming in Go
- Protocol design and parsing
- In-memory data structures
- Persistence mechanisms
- Concurrent connection handling

