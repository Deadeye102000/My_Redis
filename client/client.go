package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("SET"),
			resp.StringValue(key),
			resp.StringValue(val),
		})
	// wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue("follower"), resp.StringValue("Skyler")})
	// fmt.Printf("%s", buf.String())
	_, err = conn.Write(buf.Bytes())
	return err
}
