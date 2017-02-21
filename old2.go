package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Psiphon-Inc/goproxy"
)

func main() {
	fmt.Printf("HIHIHI")

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		defer func() {
			if e := recover(); e != nil {
				ctx.Logf("error connecting to remote: %v", e)
				client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
			}
			client.Close()
		}()
		fmt.Printf("HIHIHIHIHI@@")
		fmt.Printf("%+v \n %+v \n", req, ctx)
		clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
		remote, err := net.Dial("tcp", req.URL.Host)
		orPanic(err)
		remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
		for {
			req, err := http.ReadRequest(clientBuf.Reader)
			orPanic(err)
			orPanic(req.Write(remoteBuf))
			orPanic(remoteBuf.Flush())
			resp, err := http.ReadResponse(remoteBuf.Reader, req)
			orPanic(err)
			orPanic(resp.Write(clientBuf.Writer))
			orPanic(clientBuf.Flush())
		}
	})

	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8888", proxy))
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
