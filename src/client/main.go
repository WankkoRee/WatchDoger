package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func daemon(signalChan chan os.Signal, signalChanClosed *bool, runState *bool, wg *sync.WaitGroup) {
	defer func() {
		*runState = false
		wg.Done()
	}()

loop:
	for s := range signalChan {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Printf("程序接收到%v信号，退出中...\n", s)
			*signalChanClosed = true
			break loop
		}
	}
}

func echo(result *bufio.Reader, runState *bool, wg *sync.WaitGroup) {
	defer func() {
		*runState = false
		wg.Done()
	}()

	for *runState {
		line, err := result.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("接收回显完成，退出中...\n")
				break
			} else if errors.Is(err, net.ErrClosed) {
				log.Printf("本地连接已关闭，退出中...\n")
				break
			} else if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				log.Printf("远端连接已关闭，退出中...\n")
				break
			} else {
				log.Printf("接收终端输出失败: %v\n", err)
				break
			}
		}
		fmt.Print("| " + line)
	}
}

func key(input *bufio.Reader, conn net.Conn, runState *bool) {
	defer func() {
		*runState = false
		// 接收输入线程无需等待退出
	}()

	for *runState {
		line, err := input.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("接收用户输入完成，退出中...\n")
				break
			} else {
				log.Printf("接收用户输入失败: %v\n", err)
				break
			}
		}
		n, err := conn.Write([]byte(line))
		if err != nil {
			log.Printf("发送命令行失败: %v\n", err)
			break
		}
		if n != len(line) {
			log.Printf("发送命令行失败: %v\n", "发送长度不等于实际长度")
			break
		}
	}
}

var version = ""

func init() {
	fmt.Printf("WatchDoger Client %s\n", version)
	fmt.Printf("---- ---- ---- ----\n")
}

func main() {
	runState := true
	var err error = nil
	var wg sync.WaitGroup

	var conn net.Conn = nil

	signalChan := make(chan os.Signal)
	signalChanClosed := false
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	wg.Add(1)
	go daemon(signalChan, &signalChanClosed, &runState, &wg)

	conn, err = net.Dial("tcp", "127.0.0.1:32768")
	if err != nil {
		log.Printf("建立连接失败: %v\n", err)
		return
	}
	log.Println("建立连接成功")

	input := bufio.NewReader(os.Stdin)
	result := bufio.NewReader(conn)

	// 输出回显
	wg.Add(1)
	go echo(result, &runState, &wg)

	// 接收输入并转发
	go key(input, conn, &runState)

	for runState {

	}
	time.Sleep(1 * time.Second) // 防止信号未被处理导致`panic`
	if !signalChanClosed {
		close(signalChan)
	}
	err = conn.Close()
	if err != nil {
		log.Printf("关闭连接失败: %v\n", err)
		return
	}
	wg.Wait()
}
