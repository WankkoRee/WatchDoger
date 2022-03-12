package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func task(pipe *bufio.Reader, needEcho *bool, result **io.PipeWriter, runState *bool, wg *sync.WaitGroup) {
	defer func() {
		*runState = false
		wg.Done()
	}()

	for {
		line, err := pipe.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("交互命令结束，退出中...\n")
				break
			} else {
				log.Printf("读取控制台结果失败: %v\n", err)
				break
			}
		}
		fmt.Print("| " + line)
		if *needEcho {
			n, err := (*result).Write([]byte(line))
			if err != nil {
				log.Printf("回显控制台结果失败: %v\n", err)
			}
			if n != len(line) {
				log.Printf("发送命令行失败: %v\n", "发送长度不等于实际长度")
			}
		}
	}
}

func listener(reader *bufio.Reader, inPipe io.WriteCloser, runState *bool, listenState *bool, wg *sync.WaitGroup) {
	defer func() {
		*listenState = false
		wg.Done()
	}()

	for *runState && *listenState {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("接收命令结束，断开连接中...\n")
				break
			} else if errors.Is(err, net.ErrClosed) {
				log.Printf("接收命令连接断开\n")
				break
			} else {
				log.Printf("接收命令失败: %v\n", err)
				break
			}
		}

		fmt.Print("> " + line)
		n, err := inPipe.Write([]byte(line))
		if err != nil {
			log.Printf("写入待执行命令失败: %v\n", err)
			break
		}
		if n != len(line) {
			log.Printf("写入待执行命令失败: %v\n", "发送长度不等于实际长度")
			break
		}
	}
}

func sender(conn net.Conn, display *bufio.Reader, runState *bool, listenState *bool, wg *sync.WaitGroup) {
	defer func() {
		*listenState = false
		wg.Done()
	}()

	for *runState && *listenState {
		line, err := display.ReadString('\n')
		if err != nil {
			if strings.Contains(err.Error(), "io: read/write on closed pipe") {
				log.Printf("同步回显结束，断开连接中...\n")
				break
			} else {
				log.Printf("同步回显失败: %v\n", err)
				break
			}
		}

		n, err := conn.Write([]byte(line))
		if err != nil {
			log.Printf("回传回显失败: %v\n", err)
			break
		}
		if n != len(line) {
			log.Printf("回传回显失败: %v\n", "发送长度不等于实际长度")
			break
		}
	}
}

func init() {
	fmt.Printf("WatchDoger Server %s\n", "v1.0.0")
	fmt.Printf("---- ---- ---- ----\n")
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Printf("请将要运行的可交互命令附于本程序命令行后再执行。\n")
		return
	}
	args = args[1:]
	log.Printf("开始执行: %s\n", strings.Join(args, " ")) // 输出用户希望执行的可交互命令
	var cmd *exec.Cmd = nil
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}
	inPipe, _ := cmd.StdinPipe()
	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	runState := true
	var err error = nil
	var wg sync.WaitGroup

	needEcho := false
	var resultReader *io.PipeReader
	var resultWriter *io.PipeWriter

	wg.Add(1)
	go task(bufio.NewReader(outPipe), &needEcho, &resultWriter, &runState, &wg)
	wg.Add(1)
	go task(bufio.NewReader(errPipe), &needEcho, &resultWriter, &runState, &wg)

	err = cmd.Start()
	if err != nil {
		log.Printf("执行失败: %v\n", err)
		return
	}

	listen, err := net.Listen("tcp", "127.0.0.1:32768")
	if err != nil {
		log.Printf("开始监听失败: %v\n", err)
		return
	}
	log.Printf("已开始监听: %v\n", "127.0.0.1:32768")

	for runState {
		log.Printf("等待连接\n")
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("建立连接失败: %v\n", err)
			continue
		}
		log.Printf("已建立连接： %v\n", conn.RemoteAddr().String())
		// 只允许一个客户端存在
		resultReader, resultWriter = io.Pipe()
		needEcho = true

		listenState := true
		var wgT sync.WaitGroup
		wgT.Add(1)
		go listener(bufio.NewReader(conn), inPipe, &runState, &listenState, &wgT)
		wgT.Add(1)
		go sender(conn, bufio.NewReader(resultReader), &runState, &listenState, &wgT)
		for runState && listenState {

		}
		needEcho = false
		err = resultWriter.Close()
		if err != nil {
			log.Printf("无法关闭回显写管道: %v\n", err)
		}
		err = resultReader.Close()
		if err != nil {
			log.Printf("无法关闭回显读管道: %v\n", err)
		}
		err = conn.Close()
		if err != nil {
			log.Printf("无法关闭连接: %v\n", err)
		}
		wgT.Wait()
	}
}
