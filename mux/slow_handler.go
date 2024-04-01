package mux

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"
)

// M is a type that represents a handler that introduces a delay before responding with an "OK" status.
type M struct {
	Duration int64
	Logger   *slog.Logger
}

// Mux returns a new http.ServeMux with a single handler that introduces a delay before responding with an "OK" status.
func (m M) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", m.SlowHandler)
	return mux
}

// SlowHandler is a handler function that introduces a delay before responding with an "OK" status.
// The duration of the delay is randomly generated based on the provided Duration value.
func (m M) SlowHandler(writer http.ResponseWriter, _ *http.Request) {
	n, _ := rand.Int(rand.Reader, big.NewInt(m.Duration))
	m.Logger.Info("sleeping", "seconds", n.Int64())
	time.Sleep(time.Duration(n.Int64()) * time.Second)

	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, "OK")
}
