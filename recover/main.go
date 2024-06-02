package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime/debug"
)

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter does not support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug", debugHandlerFunc)
	mux.HandleFunc("/panic_after", panicAfterHandlerFunc)
	mux.HandleFunc("/panic", panicBeforeHandlerFunc)
	mux.HandleFunc("/", indexHandlerFunc)

	fmt.Println("Listening on port 3000")
	if err := http.ListenAndServe(":3000", recoverMw(mux, false)); err != nil {
		panic(err)
	}
}

func debugHandlerFunc(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, file)
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				if !dev {
					http.Error(w, "Oh no! The squirrels have learned how to hack into our Wi-Fi and are now plotting world domination!", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<p>Panic: %v</p><br><pre>%s<pre>", err, string(stack))
			}
		}()

		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

func panicAfterHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You are so Busted")
	mustPanic()
}

func panicBeforeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	mustPanic()
	fmt.Fprint(w, "Hi Stacy Where are you ?")
}

func indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Index Page")
}

func mustPanic() {
	panic("Uh-oh! The coffee machine just exploded and now there’s espresso everywhere, including on the cat!")
}
