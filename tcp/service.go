package tcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Handler interface {
	Handler(ctx context.Context, conn net.Conn)
	Close() error
}

type Config struct {
	Address string
}

// ListenAndServe 监听并提供服务，并在收到closeChan 发来的关闭通知后关闭
func ListenAndServe(listener net.Listener, handler Handler, closeChan <-chan struct{}) {
	// 监听关闭通知
	go func() {
		<-closeChan
		log.Println("shutting down..")
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		log.Println("accept link..")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()

			handler.Handler(ctx, conn)
		}()
	}

	waitDone.Wait()
}

func ListenAndServerWithSignal(cfg *Config, handler Handler) error {
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}
