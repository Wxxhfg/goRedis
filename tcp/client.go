package tcp

import (
	"goRedis/lib/sync/wait"
	"net"
	"time"
)

type Client struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	// 等待数据发送完成或超时
	c.Waiting.WaitWithTimeOut(10 * time.Second)
	c.Conn.Close()
	return nil
}
