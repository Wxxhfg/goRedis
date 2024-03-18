package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type EchoHandler struct {
	// 保存所有工作状态client的集合
	// 需要并发安全的容器
	// todo 如果不是并发安全的容器会怎样？
	activeConn sync.Map
	// 关闭状态标识位
	closing atomic.Bool
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		conn.Close()
		return
	}

	client := &Client{
		Conn: conn,
	}
	// 记住仍然存活的连接
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("connection close")
				h.activeConn.Delete(client)
			} else {
				log.Println(err)
			}

			return
		}

		client.Waiting.Add(1)
		b := []byte(msg)
		conn.Write(b)
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	log.Println("handler shutting down..")
	h.closing.Store(true)
	h.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*Client)
		client.Close()
		return true
	})

	return nil
}
