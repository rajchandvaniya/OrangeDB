package core

import (
	"errors"
)

func evalPING(cmd *RedisCmd) ([]byte, error) {
	if len(cmd.Args) > 1 {
		return nil, errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(cmd.Args) == 0 {
		return Encode("PONG", true), nil
	} else {
		return Encode(cmd.Args[0], false), nil
	}
}

func Eval(cmd *RedisCmd) ([]byte, error) {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd)
	default:
		return nil, errors.New("UNKNOWNCMD " + cmd.Cmd)
	}
}
