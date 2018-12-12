package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	flagAddr string
	dumpReq  int
)

func init() {
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:9000", "the address of static file server")
	flag.IntVar(&dumpReq, "dump", 0, "1: dump request header, 2: dump request header & body")
}

func main() {
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Run file server(dir=%s) on %s\n", wd, flagAddr)

	h := http.FileServer(http.Dir(wd))
	if dumpReq > 0 {
		prevHandler := h
		h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			dumpRequest(req)
			prevHandler.ServeHTTP(w, req)
		})
	}

	http.ListenAndServe(flagAddr, h)
}

func dumpRequest(req *http.Request) {
	fmt.Println("==================== Incoming Request ====================")
	defer func() {
		fmt.Println("==================== Dump Finished ====================")
	}()
	dumpData, err := httputil.DumpRequest(req, dumpReq == 2)
	if err != nil {
		fmt.Println("DUMP ERROR:", err)
		return
	}
	dumpData = bytes.TrimSpace(dumpData)
	fmt.Println(string(dumpData))
}
