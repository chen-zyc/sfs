package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	flagAddr string
)

func init() {
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:9000", "the address of static file server")
}

func main() {
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Run file server(dir=%s) on %s\n", wd, flagAddr)

	http.ListenAndServe(flagAddr, http.FileServer(http.Dir(wd)))
}


