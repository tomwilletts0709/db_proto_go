package main

import (
	"bytes"
	"testing"
)

func expectBad(t *testing.T, payload string) {
	t.Helper()
	_, err := NewResp(bytes.NewReader([]byte(payload))).Read()
	if err == nil {
		t.Fatalf("expected error for bad payload %q, got nil", payload)
	}
}

func expectArray(t *testing.T, payload string, wantLen int) {
	t.Helper()
	value, err := NewResp(bytes.NewReader([]byte(payload))).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.typ != "array" {
		t.Fatalf("expected array, got typ=%q", value.typ)
	}
	if len(value.array) != wantLen {
		t.Fatalf("expected array len %d, got %d", wantLen, len(value.array))
	}
}

func expectBulk(t *testing.T, payload string, wantBulk string) {
	t.Helper()
	value, err := NewResp(bytes.NewReader([]byte(payload))).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.typ != "bulk" {
		t.Fatalf("expected bulk, got typ=%q", value.typ)
	}
	if value.bulk != wantBulk {
		t.Fatalf("expected bulk %q, got %q", wantBulk, value.bulk)
	}
}

func expectInteger(t *testing.T, payload string, wantNum int) {
	t.Helper()
	value, err := NewResp(bytes.NewReader([]byte(payload))).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.typ != "integer" {
		t.Fatalf("expected integer, got typ=%q", value.typ)
	}
	if value.num != wantNum {
		t.Fatalf("expected num %d, got %d", wantNum, value.num)
	}
}

func TestRESP(t *testing.T) {
	t.Run("bad payloads", func(t *testing.T) {
		expectBad(t, "")
		expectBad(t, "*") // truncated
		expectBad(t, "*1\r\n") // array len 1 but no element
	})

	t.Run("array", func(t *testing.T) {
		expectArray(t, "*0\r\n", 0)
		expectArray(t, "*2\r\n$3\r\nSET\r\n$3\r\nkey\r\n", 2)
	})

	t.Run("bulk", func(t *testing.T) {
		expectBulk(t, "$4\r\nTEST\r\n", "TEST")
		expectBulk(t, "$0\r\n\r\n", "")
	})

	t.Run("integer", func(t *testing.T) {
		expectInteger(t, ":0\r\n", 0)
		expectInteger(t, ":42\r\n", 42)
		expectInteger(t, ":-1\r\n", -1)
	})
}
