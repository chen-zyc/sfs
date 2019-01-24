package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	flagAddr    string // 服务地址
	dumpReq     int    // 打印请求信息
	blockHeader string // 阻塞时间的请求头
	countHeader string
)

func init() {
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:9000", "the address of static file server")
	flag.IntVar(&dumpReq, "dump", 0, "1: dump request header, 2: dump request header & body")
	flag.StringVar(&blockHeader, "block-header", "", "request header representing blocking time")
	flag.StringVar(&countHeader, "count-header", "", "response header containing counter values")
}

func main() {
	flag.Parse()

	h := newHandler()
	fmt.Printf("Run file server(dir=%s) on %s\n", h.staticRoot, flagAddr)

	err := http.ListenAndServe(flagAddr, h)
	fmt.Printf("Server returned: %v\n", err)
}

type handler struct {
	staticRoot string
	fileServer http.Handler
}

func newHandler() *handler {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.Dir(wd))
	fileServer = wrapHandler(fileServer)
	return &handler{
		staticRoot: wd,
		fileServer: fileServer,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if blockHeader != "" {
		blockRequest(req)
	}
	h.fileServer.ServeHTTP(w, req)
}

func wrapHandler(h http.Handler) http.Handler {
	wrapper := func(prevHandler http.Handler, before, after http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if before != nil {
				before(w, req)
			}
			if prevHandler != nil {
				prevHandler.ServeHTTP(w, req)
			}
			if after != nil {
				after(w, req)
			}
		})
	}
	if dumpReq > 0 {
		h = wrapper(h, func(writer http.ResponseWriter, req *http.Request) {
			dumpRequest(req)
		}, nil)
	}
	if countHeader != "" {
		h = wrapper(h, func(w http.ResponseWriter, req *http.Request) {
			setCounterValueToResponse(w)
		}, nil)
	}
	return h
}

func dumpRequest(req *http.Request) {
	fmt.Println("==================== Incoming Request ====================")
	defer func() {
		fmt.Println("==================== Request Finished ====================")
	}()
	dumpData, err := httputil.DumpRequest(req, dumpReq == 2)
	if err != nil {
		fmt.Println("DUMP ERROR:", err)
		return
	}
	dumpData = bytes.TrimSpace(dumpData)
	fmt.Println(string(dumpData))
}

func blockRequest(req *http.Request) {
	blockDurationStr := req.Header.Get(blockHeader)
	if blockDurationStr == "" {
		return
	}
	blockDur, err := time.ParseDuration(blockDurationStr)
	if err != nil {
		fmt.Printf("Failed to parse header[%s]: %s, err: %v\n", blockHeader, blockDurationStr, err)
		return
	}
	start := time.Now()
	time.Sleep(blockDur)
	fmt.Printf("Block request, start: %s, duration: %v\n", start, blockDur)
}

var counter uint64

func setCounterValueToResponse(w http.ResponseWriter) {
	v := atomic.AddUint64(&counter, 1)
	w.Header().Set(countHeader, strconv.FormatUint(v, 10))
}
