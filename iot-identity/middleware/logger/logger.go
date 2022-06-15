// This code largely adapted from https://github.com/pytimer/mux-logrus

//nolint:gocritic
//MIT License
//
//Copyright (c) 2018 pytimer
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

// Package logger provides a middleware for logging with Gin using logrus
package logger

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// realClock save request times
type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// LogOptions logging middleware options
type LogOptions struct {
	EnableStarting bool
}

// LoggingMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggingMiddleware struct {
	logger         *logrus.Logger
	clock          timer
	enableStarting bool
}

// New returns a new *LoggingMiddleware, yay!
func New(l *logrus.Logger, options LogOptions) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger:         l,
		clock:          &realClock{},
		enableStarting: options.EnableStarting,
	}
}

// realIP get the real IP from http request
func realIP(req *http.Request) string {
	ra := req.RemoteAddr
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

// HandleFunc is the function for processing requests in the chain and logging
func (m *LoggingMiddleware) HandleFunc(context *gin.Context) {
	logrus.Info("HandleFunc of LoggingMiddleware")
	r := context.Request
	entry := logrus.NewEntry(logrus.StandardLogger())
	start := m.clock.Now()

	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		entry = entry.WithField("requestId", reqID)
	}

	if remoteAddr := realIP(r); remoteAddr != "" {
		entry = entry.WithField("remoteAddr", remoteAddr)
	}

	if m.enableStarting {
		entry.WithFields(logrus.Fields{
			"timestamp": time.Now().Format("2006/01/02 - 15:04:05"),
			"request":   r.RequestURI,
			"method":    r.Method,
			"requstURI": r.RequestURI,
			"token":     r.Header.Get("Authorization"),
		}).Info("started handling request")
	}

	context.Next()

	latency := m.clock.Since(start)

	entry.WithFields(logrus.Fields{
		"timestamp": time.Now().Format("2006/01/02 - 15:04:05"),
		"status":    context.Writer.Status(),
		"took":      latency,
	}).Info("completed handling request")
}
