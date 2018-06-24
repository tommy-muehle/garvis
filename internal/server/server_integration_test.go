// +build integration

package server

import (
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestServer_AddHandler(t *testing.T) {
	s := New(":1235", zap.L())
	s.AddHandler("/foo", http.NotFoundHandler())

	go func() {
		s.ListenAndServe()
	}()

	time.Sleep(1 * time.Second)

	res, err := http.Get("http://localhost:1235/foo")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status code isn't 404: %v", res.StatusCode)
	}
}

func TestServer_Health(t *testing.T) {
	s := New(":1236", zap.L())
	s.Health("/health")

	go func() {
		s.ListenAndServe()
	}()

	defer s.Shutdown()

	time.Sleep(1 * time.Second)

	res, err := http.Get("http://localhost:1236/health")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code is not ok: %v", res.StatusCode)
	}
}
