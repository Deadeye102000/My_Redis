package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
)

const ( // CommandSet is the command to set a key-value pair.
	CommandSET = "SET"
)

type Command interface {
}

type SetCommand struct {
	key, val string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err != nil {
			if err == io.EOF {
				return nil, err // disconnect, not fatal
			}
			return nil, fmt.Errorf("RESP parse error: %v", err)
		}

		if v.Type() == resp.Array {
			arr := v.Array()
			if len(arr) == 0 {
				return nil, fmt.Errorf("empty command array")
			}

			switch arr[0].String() {
			case CommandSET:
				if len(arr) != 3 {
					return nil, fmt.Errorf("SET command must have 3 elements")
				}
				return SetCommand{
					key: arr[1].String(),
					val: arr[2].String(),
				}, nil
			default:
				return nil, fmt.Errorf("unknown command: %s", arr[0].String())
			}
		}

		return nil, fmt.Errorf("expected array RESP type")
	}
}
