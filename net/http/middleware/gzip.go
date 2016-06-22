package middleware

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// GZip provides a HTTP middleware to GZip compress the response stream.
func GZip(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		compress := strings.Contains(acceptEncoding, "gzip")
		if _, upgrade := r.Header["Upgrade"]; upgrade {
			// Avoid compressing websockets due to handling of
			// closing and flushing of connections.
			compress = false
		}
		if !compress {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Remove Accept-Encoding from the request, so that further
		// handlers don't attempt to apply their own encodings.
		r.Header.Del("Accept-Encoding")

		gz := gzip.NewWriter(w)
		gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w, gz: gz}

		defer func() {
			// Only close gz if anything was written.
			if gzw.bytesWritten != 0 {
				gz.Close()
			}
		}()

		next.ServeHTTP(gzw, r)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter

	gz            *gzip.Writer
	headerWritten bool
	bytesWritten  int
}

// WriteHeader removes the Content-Length if set.
func (w *gzipResponseWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
	w.headerWritten = true
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.Writer.Write(b)
	w.bytesWritten += n
	return n, err
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		err := fmt.Errorf("%T does not implement http.Hijacker", w.ResponseWriter)
		return nil, nil, err
	}
	return h.Hijack()
}

func (w *gzipResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		_ = w.gz.Flush()
		f.Flush()
	}
}
