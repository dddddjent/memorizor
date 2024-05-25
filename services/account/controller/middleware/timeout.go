package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"memorizor/services/account/util"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Only write header when it has not timed out
// Write the response only when it's finished or timed out, so cache the temporary data
type timeoutResponseWriter struct {
	gin.ResponseWriter
	// temp data
	header http.Header
	buf    bytes.Buffer // body

	mu         sync.Mutex
	isTimeOut  bool
	hasWritten bool
	code       int
}

// Basically when you call c.JSON after timeout it will not do anything.
// If we try to check the ctx, it will notify that it's time to quit, otherwise it will continue
// to run even if you can't really write anything into the response.
func Timeout(timeout time.Duration, errWhenTimeout *util.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// timeoutCtx is a child of c.Request.Context()
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		c.Request.WithContext(timeoutCtx) // sub-goroutines can quit early if they check it
		defer cancel()

		timeoutWriter := &timeoutResponseWriter{ResponseWriter: c.Writer, header: make(http.Header)}
		c.Writer = timeoutWriter

		// determine if timeout should be set
		properFinished := make(chan struct{})
		panicFinished := make(chan any, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicFinished <- p
				}
			}()
			c.Next()
			properFinished <- struct{}{}
		}()

		select {
		case <-panicFinished:
			err := util.NewInternal("Panic!")
			timeoutWriter.ResponseWriter.WriteHeader(err.HttpStatus())
			body, _ := json.Marshal(gin.H{
				"error": errWhenTimeout,
			})
			timeoutWriter.ResponseWriter.Write(body)
		case <-properFinished:
			timeoutWriter.mu.Lock()
			defer timeoutWriter.mu.Unlock()
			responseHeader := timeoutWriter.ResponseWriter.Header()
			for key, val := range timeoutWriter.Header() {
				responseHeader[key] = val
			}
			timeoutWriter.ResponseWriter.WriteHeader(timeoutWriter.code)
			timeoutWriter.ResponseWriter.Write(timeoutWriter.buf.Bytes())
		case <-timeoutCtx.Done():
			timeoutWriter.mu.Lock()

			timeoutWriter.ResponseWriter.Header().Set("Content-Type", "application/json")
			timeoutWriter.ResponseWriter.WriteHeader(errWhenTimeout.HttpStatus())
			body, _ := json.Marshal(gin.H{
				"error": errWhenTimeout,
			})
			timeoutWriter.ResponseWriter.Write(body)
			c.Abort()

			timeoutWriter.mu.Unlock()
			// Other goroutines may write into the writer here, but they won't be sent anymore
			timeoutWriter.SetTimeout()
		}
	}
}

// http.ResponseWriter's interface
func (writer *timeoutResponseWriter) Write(b []byte) (int, error) {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	if writer.isTimeOut {
		return 0, nil
	}
	return writer.buf.Write(b)
}

// http.ResponseWriter's interface
func (writer *timeoutResponseWriter) WriteHeader(code int) {
	if code < 100 || code > 999 {
		return
	}
	writer.mu.Lock()
	defer writer.mu.Unlock()
	if writer.isTimeOut || writer.hasWritten {
		return
	}
	writer.hasWritten = true
	writer.code = code
}

// http.ResponseWriter's interface
func (writer *timeoutResponseWriter) Header() http.Header {
	return writer.header
}

func (writer *timeoutResponseWriter) SetTimeout() {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.isTimeOut = true
}
