package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/rajchandvaniya/orangedb/config"
)

func StartEchoTCPServer() {
	log.Println("starting a synchronous TCP server on ", config.Host, config.Port)
	con_clients := 0

	conn, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		c, err := conn.Accept()
		if err != nil {
			panic(err)
		}

		con_clients += 1
		log.Println("client connected with remote address", c.RemoteAddr(), "concurrent clients:", con_clients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients:", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("error:", err)
			}
			log.Println("received command:", cmd)
			err = writeCommand(c, cmd)
			if err != nil {
				log.Println("error while write:", err)
			}
		}
	}

}

func readCommand(con net.Conn) (string, error) {
	buffer := make([]byte, 512)
	n, err := con.Read(buffer[:])
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func writeCommand(con net.Conn, cmd string) error {
	_, err := con.Write([]byte(cmd))
	return err
}
