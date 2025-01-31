package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dsidak/exchange/exchange/api"
	"github.com/joho/godotenv"
)

// key for the HTTP server’s address value in the context
const keyServerAddr = "serverAddr"

func main() {
	loadEnv()

	startServer()
}

func startServer() {
	// server multiplexer and http.Handler implementation
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/latest", getLatestRates)
	mux.HandleFunc("/symbols", getSymbols)

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
		log.Printf("server is closed\n")
	} else if err != nil {
		log.Printf("error listening for server: %s\n", err)
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
		log.Printf("could not read body: %s\n", err)
	}

	log.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n body=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)

	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))

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

func getLatestRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("%s: got /latest request\n", ctx.Value(keyServerAddr))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read body: %s\n", err)
	}

	log.Printf("body=\n%s\n", body)

	hasSymbols := r.URL.Query().Has("symbols")
	symbols := r.URL.Query().Get("symbols")
	hasBase := r.URL.Query().Has("base")
	base := r.URL.Query().Get("base")

	context := api.Context{
		Endpoint: api.Latest,
	}
	if hasSymbols {
		context.Symbols = symbols
	}
	if hasBase {
		context.Base = base
	}

	response := api.CallFixerIo(context)
	io.WriteString(w, fmt.Sprintf("%v\n", response))
}

func getSymbols(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("%s: got /latest request\n", ctx.Value(keyServerAddr))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read body: %s\n", err)
	}

	log.Printf("body=\n%s\n", body)

	context := api.Context{
		Endpoint: api.Symbols,
	}

	response := api.CallFixerIo(context)
	io.WriteString(w, fmt.Sprintf("%v\n", response))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("could not load .env file: %s\n", err)
	}

	log.Printf("godotenv: %s = %s \n", "FIXER_IO_API_KEY", os.Getenv("FIXER_IO_API_KEY"))
	log.Printf("godotenv: %s = %s \n", "FIXER_IO_ADDRESS", os.Getenv("FIXER_IO_ADDRESS"))
}
