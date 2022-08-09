package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	localAddr, remoteAddr string
)

func init() {
	flag.StringVar(&localAddr, "local", "0.0.0.0:80", "local address")
	flag.StringVar(&remoteAddr, "remote", "0.0.0.0:8080", "remote address")
	flag.Parse()
}

func main() {
	fmt.Println("Listening: ", localAddr)
	fmt.Println("Proxying:  ", remoteAddr)

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		panic(err)
	}
	for {
		// Accept 等待并将下一个连接返回
		// 在这个场景下返回一个 tcp 连接
		conn, err := listener.Accept()
		log.Println("New connection", conn.RemoteAddr())
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}
		go func() {
			conn1, err := net.Dial("tcp", remoteAddr)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}
			defer conn1.Close()
			closer := make(chan struct{}, 2)
			go copy(closer, conn1, conn)
			go copy(closer, conn, conn1)
			<-closer
			log.Println("Connection complete")
			conn.Close()
		}()
	}
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}
