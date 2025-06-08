package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

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
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Read %s\n", v.Type())

		// var cmd Command
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						// panic("yikes")
						return nil, fmt.Errorf("expected 3 elements in SET command, got %d", len(v.Array()))
					}
					cmd := SetCommand{
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}

					return cmd, nil

				}

			}
		}
		return nil, fmt.Errorf("Invalid or unknown command received: %s", raw)
	}
	return nil, fmt.Errorf("Invalid or unknown command received: %s", raw)
}
