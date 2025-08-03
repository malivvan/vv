package main

import (
	"encoding/hex"
	"fmt"
	"github.com/malivvan/vv/pkg/rpc"
	"net"
	"time"
)

func main() {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	go func() {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", ":4444")

		var sess rpc.Session
		check(sess.DoClient(conn))

		fmt.Println(hex.EncodeToString(sess.SharedKey()))

		sc := rpc.NewSessionConn(sess.Suite(), conn)

		for i := 0; i < 100; i++ {
			_, err = sc.Write([]byte(fmt.Sprintf("[%d] Hello from Go!", i)))
			check(err)
			check(sc.Flush())
		}

	}()

	ln, err := net.Listen("tcp", ":4444")
	check(err)
	defer ln.Close()

	conn, err := ln.Accept()
	check(err)
	defer conn.Close()

	var sess rpc.Session
	check(sess.DoServer(conn))

	fmt.Println(hex.EncodeToString(sess.SharedKey()))

	sc := rpc.NewSessionConn(sess.Suite(), conn)

	buf := make([]byte, 1024)

	for i := 0; i < 100; i++ {
		n, err := sc.Read(buf)
		check(err)

		fmt.Println("Decrypted:", string(buf[:n]))
	}
}
