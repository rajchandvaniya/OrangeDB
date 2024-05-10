package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/rajchandvaniya/orangedb/config"
	"github.com/rajchandvaniya/orangedb/core"
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
			respond(cmd, c)
		}
	}

}

func readCommand(con net.Conn) (*core.RedisCmd, error) {
	buffer := make([]byte, 512)
	n, err := con.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	cmds, err := core.Decode(buffer[:n])
	tokens := cmds.([]interface{})
	if err != nil {
		return nil, err
	}
	return &core.RedisCmd{Cmd: strings.ToUpper(tokens[0].(string)), Args: tokens[1:]}, nil
}

func respond(cmd *core.RedisCmd, con net.Conn) {
	response, err := core.Eval(cmd)
	if err != nil {
		con.Write(core.EncodeError(err))
	}
	con.Write(response)
}
