package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	aof := &Aof{file: f, rd: bufio.NewReader(f), mu: sync.Mutex{}}

	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}

func (aof *Aof) Replay(apply func(Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(aof.file)

	for {
		resp := NewResp(reader)
		value, err := resp.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if value.typ != "array" || len(value.array) == 0 {
			continue
		}

		apply(value)
	}
	return nil
}
