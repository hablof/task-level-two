package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":23")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("new connection accepted")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("err on read: %v\n", err)
			break
		}

		log.Printf("recived: %s", s)

		rs := []rune(strings.TrimRight(s, "\r\n"))
		slices.Reverse(rs)
		s2 := string(rs)

		if _, err := conn.Write([]byte(s2 + "\n")); err != nil {
			log.Printf("err on write: %v\n", err)
			break
		}

		log.Printf("sended: %s", s2)
	}

}
