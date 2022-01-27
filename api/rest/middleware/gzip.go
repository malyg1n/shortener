package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("compress starting")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			fmt.Println("compress ignore")
			return
		}
		fmt.Println("compress started")
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			fmt.Println("compress error", err.Error())
			return
		}
		defer func() {
			_ = gz.Close()
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		fmt.Println("compress success")
	})
}

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("decompress starting")
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			fmt.Println("decompress ignore")
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("decompress error", err.Error())
			return
		}

		var b bytes.Buffer
		_, err = b.ReadFrom(gz)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nr := io.NopCloser(bytes.NewReader(b.Bytes()))
		rb, _ := http.NewRequest(r.Method, r.RequestURI, nr)
		_ = r.Body.Close()
		next.ServeHTTP(w, rb)
		fmt.Println("decompress success")
	})
}
