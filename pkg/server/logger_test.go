package server

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestLogResponse_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	data := []byte("{\n  \"test\": true\n}")

	LogResponseToDir(dir, data)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	if !strings.HasSuffix(entries[0].Name(), ".json") {
		t.Fatalf("expected .json file, got %s", entries[0].Name())
	}
	content, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(data) {
		t.Fatalf("content mismatch: got %q", string(content))
	}
}

func TestLogResponse_ConcurrentCallsCreateUniqueFiles(t *testing.T) {
	dir := t.TempDir()
	data := []byte("{}")
	n := 50

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			LogResponseToDir(dir, data)
		}()
	}
	wg.Wait()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != n {
		t.Fatalf("expected %d files, got %d", n, len(entries))
	}
}
