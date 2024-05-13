package server

import (
	"io"
	"strings"

	"github.com/rajchandvaniya/orangedb/core"
)

func readCommand(con io.ReadWriter) (*core.RedisCmd, error) {
	buffer := make([]byte, 512)
	n, err := con.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	cmds, err := core.Decode(buffer[:n])
	if err != nil {
		return nil, err
	}

	tokens := cmds.([]interface{})
	return &core.RedisCmd{Cmd: strings.ToUpper(tokens[0].(string)), Args: tokens[1:]}, nil
}

func respond(cmd *core.RedisCmd, con io.ReadWriter) {
	response, err := core.Eval(cmd)
	if err != nil {
		con.Write(core.EncodeError(err))
	}
	con.Write(response)
}
