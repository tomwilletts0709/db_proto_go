package main

import (
	"path/filepath"
	"testing"
)

func TestNewServer(t *testing.T) {
	tmp := t.TempDir()
	aofPath := filepath.Join(tmp, "test.aof")

	server, err := NewServer(aofPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
