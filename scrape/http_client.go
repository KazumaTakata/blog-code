package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", "151.101.129.69:80", 5*time.Second)

	if err != nil {
		fmt.Printf("%v", err)
	}

	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: stackoverflow.com\r\n\r\n"))

	buf := make([]byte, 1000)

	_, err = conn.Read(buf)

	fmt.Printf("%s", string(buf))
}
