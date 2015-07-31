package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
)

const Version = "0.0.1"

var options struct {
	localAddr  string
	remoteAddr string
	minBytes   int64
	maxBytes   int64
	seed       int64
	version    bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&options.version, "version", false, "print version and exit")
	flag.StringVar(&options.localAddr, "local", "127.0.0.1:2999", "listen for connections on specified address")
	flag.StringVar(&options.remoteAddr, "remote", "", "forward connections to specified address")
	flag.Int64Var(&options.seed, "seed", -1, "seed for random number generator")
	flag.Int64Var(&options.minBytes, "minbytes", 1, "minimum bytes before cave in")
	flag.Int64Var(&options.maxBytes, "maxbytes", 5000, "maximum bytes before cave in")
	flag.Parse()

	if options.version {
		fmt.Printf("cavein v%v\n", Version)
		os.Exit(0)
	}

	if options.remoteAddr == "" {
		log.Fatal("remote is required")
	}

	if options.seed < 0 {
		options.seed = rand.Int63()
		log.Println("seed:", options.seed)
	}

	rand.Seed(options.seed)

	ln, err := net.Listen("tcp", options.localAddr)
	if err != nil {
		log.Fatal(err)
	}

	var connCount int

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connCount += 1
		go forwardConn(conn, connCount)
	}
}

func forwardConn(localConn net.Conn, connNum int) {
	remoteConn, err := net.Dial("tcp", options.remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	byteLimit := options.minBytes + rand.Int63n(options.maxBytes-options.minBytes)
	remoteReader := &io.LimitedReader{R: remoteConn, N: byteLimit}
	log.Println("Tunnel", connNum, "established, will cave in after", remoteReader.N, "bytes")

	go io.Copy(remoteConn, localConn)
	go func() {
		n, err := io.Copy(localConn, remoteReader)
		if err != nil {
			log.Println("Tunnel", connNum, "unexpected error", err)
		}

		err = localConn.Close()
		if err != nil {
			log.Println("Tunnel", connNum, "unexpected error", err)
		}

		log.Println("Tunnel", connNum, "caved in after receiving", n, "bytes")
	}()
}
