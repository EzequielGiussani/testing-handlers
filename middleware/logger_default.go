package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// NewAuthTokenBasic returns a new AuthBasic
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

type DefaultLogger struct {
}

// Auth is a method that authenticates
func (a *DefaultLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Log before")

		next.ServeHTTP(w, r)

		// After

		fmt.Println("Log after")

		fmt.Println("Verb: ", r.Method)
		fmt.Println("RequestUrl: ", os.Getenv("SERVER_ADDR")+r.URL.String())
		fmt.Println("RequestBiteSize: ", r.Header.Get("Content-Length"), "bytes")
		fmt.Println("Time: ", time.Now().Format("2006-01-02 15:04:05"))
	})
}
