package core

import (
	"errors"
	"strconv"
	"time"
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

func evalGET(cmd *RedisCmd) ([]byte, error) {
	if len(cmd.Args) != 1 {
		return nil, errors.New("ERR wrong number of arguments for 'get' command")
	}

	key := cmd.Args[0]
	obj := Get(key.(string))

	// if key doesnt exist
	if obj == nil {
		return Encode(nil, false), nil
	}

	// if key is expired
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		return Encode(nil, false), nil
	}

	return Encode(obj.Value, false), nil
}

func evalSET(cmd *RedisCmd) ([]byte, error) {
	if len(cmd.Args) <= 1 {
		return nil, errors.New("ERR wrong number of arguments for 'set' command")
	}

	var key, value string
	key, value = (cmd.Args[0]).(string), (cmd.Args[1]).(string)
	var exDurationMs int64 = -1

	for i := 2; i < len(cmd.Args); i++ {
		switch cmd.Args[i] {
		case "EX", "ex":
			i++
			if i == len(cmd.Args) {
				return nil, errors.New("ERR syntax error")
			}
			exDurationSec, err := strconv.ParseInt(cmd.Args[i].(string), 10, 64)
			if err != nil {
				return nil, errors.New("ERR valye us not an integer or out of range")
			}
			exDurationMs = exDurationSec * 1000
		default:
			return nil, errors.New("ERR syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))
	return Encode("OK", true), nil
}

func evalTTL(cmd *RedisCmd) ([]byte, error) {
	if len(cmd.Args) != 1 {
		return nil, errors.New("ERR wrong number of arguments for 'ttl' command")
	}

	key := cmd.Args[0]
	obj := Get(key.(string))

	// if key doesnt exist
	if obj == nil {
		return Encode(int64(-2), false), nil
	}

	// if no expiry set
	if obj.ExpiresAt == -1 {
		return Encode(int64(-1), false), nil
	}

	// compute the remaining time
	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	// if key expired
	if durationMs < 0 {
		return Encode(int64(-2), false), nil
	}

	return Encode(int64(durationMs/1000), false), nil
}

func Eval(cmd *RedisCmd) ([]byte, error) {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd)
	case "GET":
		return evalGET(cmd)
	case "SET":
		return evalSET(cmd)
	case "TTL":
		return evalTTL(cmd)
	default:
		return nil, errors.New("UNKNOWNCMD " + cmd.Cmd)
	}
}
