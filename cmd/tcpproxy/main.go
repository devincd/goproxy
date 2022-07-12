package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var localAddr *string = flag.String("local", "0.0.0.0:80", "local address")
var remoteAddr *string = flag.String("remote", "0.0.0.0:80", "remote address")

func main() {
	flag.Parse()
	fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

	listener, err := net.Listen("tcp4", *localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		log.Println("[New connection] remote addr:", conn.RemoteAddr())
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}
		go DoProxy(conn)
	}
}

func DoProxy(conn net.Conn) {
	conn2, err := net.Dial("tcp4", *remoteAddr)
	if err != nil {
		log.Println("error dialing remote addr", err)
		return
	}
	defer conn2.Close()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// send data
		io.Copy(conn2, conn)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		// receive data
		io.Copy(conn, conn2)
		wg.Done()
	}()
	wg.Wait()
	log.Println("[Connection complete] remote addr:", conn.RemoteAddr())
}
