package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/dsidak/exchange/exchange/api"
)

// key for the HTTP server’s address value in the context
const keyServerAddr = "serverAddr"

func main() {
	response := api.CallFixerIo()
	fmt.Println(response)

	startServer()
}

func startServer() {
	// server multiplexer and http.Handler implementation
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/rates", getRates)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(ctx, keyServerAddr, listener.Addr().String())
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server is closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s, body=\n%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)

	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}

// As an example of how to pass message body to the handler, you can use the following curl command:
// curl -X POST -d 'This is the body' 'http://localhost:8080?first=1&second='
// Might be a good idea to use the body as JSON to request multiple exchange rates at once
// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go

func getRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /rates request\n", ctx.Value(keyServerAddr))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("body=\n%s\n", body)

	response := api.CallFixerIo()

	io.WriteString(w, fmt.Sprintf("%v\n", response))
}
