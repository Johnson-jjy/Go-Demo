package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:7777"
	DELEMITER = '\t' // 表示一个作为数据边界的单字节字符
)

var wg sync.WaitGroup

// 为更好地记录日志而编写的辅助函数
// 隔离将来很可能发生的日志记录方式的变化，并能够避免散弹式修改
func printLog(role string, sn int, format string, args ...interface{})  {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}

// 检查数据块是否可以转化为一个int32类型的值并转换
func strToInt32(str string) (int32, error) {
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("\"%s\" is not interger", str)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, fmt.Errorf("%d is not 32-bit interger", num)
	}
	return int32(num), nil
}

// 计算立方根
func cbrt(param int32) float64 {
	return math.Cbrt(float64(param))
}

// 读取一段以数据分解符为结尾的数据
func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1) // 初始化为1:防止从连接值中读出多余的数据从而对后续的读取操作造成影响
	var buffer bytes.Buffer // 暂存当前数据块中的字节
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELEMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

// 以下代码便会造成提前读取 -> 导致一些数据不完整,甚至漏掉一些数据块
//func read(conn net.Conn) (string, error) {
//	reader := bufio.NewReader(conn)
//	readBytes, err := reader.ReadBytes(DELIMITER) // 缓存机制会预读一部分数据
//	if err != nil {
//		return "", err
//	}
//	return string(readBytes[:len(readBytes)-1]), nil
//}

func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELEMITER) // 追加一个数据分界符, 形成一个两端程序均可识别的数据块
	return conn.Write(buffer.Bytes())
}

func serverGo() {
	// 根据给定的网络协议和地址创建一个监听器
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		printServerLog("Listen Error: %s", err)
		return
	}
	defer listener.Close() // 保证在函数结束执行前关闭监听器

	printServerLog("Got listener for the server. (local address: %s)", listener.Addr())

	for {
		conn, err := listener.Accept() // 阻塞直至新连接到来
		if err != nil {
			printServerLog("Accept Error: %s", err)
			continue
		}
		printServerLog("Established a connection with a client application. (remote address: %s)", conn.RemoteAddr())
		go handleConn(conn) // 并发处理连接以避免完全串行处理众多连接
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		wg.Done()
	}()
	// 试图从连接中获取数据 -> 保证尽量及时地处理和响应请求
	for {
		// 关闭闲置连接
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF { // read未作处理,此处做等值处理
				printServerLog("The connection is closed by another side.")
			} else {
				printServerLog("Read Error: %s", err)
			}
			break
		}
		printServerLog("Receiver request: %s", strReq)
		intReq, err := strToInt32(strReq)
		if err != nil {
			n, err := write(conn, err.Error())
			printServerLog("Sent error message (written %d bytes): %s.", n, err)
			continue
		}
		floatResp := cbrt(intReq)
		respMsg := fmt.Sprintf("The cube root of %d is %f", intReq, floatResp)
		n, err := write(conn, respMsg)
		if err != nil {
			printServerLog("Write Error: %s", err)
		}
		printServerLog("Sent response (written %d bytes): %s.", n, respMsg)
	}
}

func clientGo(id int) { // id -> 在运行多个客户端程序的场景下在日志中区分它们
	defer wg.Done()
	// 与服务端程序建立连接
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		printClientLog(id, "Dial Error: %s", err)
		return
	}
	defer conn.Close()
	printClientLog(id, "Connected to server. (remote address: %s, local address: %s)", conn.RemoteAddr(), conn.LocalAddr())
	time.Sleep(200 * time.Millisecond) // 只是让两端程序看上去清晰一些

	// 发送数据 -> 另C端发送的请求数据库数量定义为5个
	requestNumber := 5
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
	for i := 0; i <requestNumber; i++ {
		req := rand.Int31() // 随机生成一个int32型
		n, err := write(conn, fmt.Sprintf("%d", req))
		if err != nil {
			printClientLog(id, "Write Error: %s", err)
			continue
		}
		printClientLog(id, "Sent request (written %d bytes): %d.", n, req)
	}

	// 接收响应数据块
	for j := 0; j < requestNumber; j++ {
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printClientLog(id, "The connection is closed by another side.")
			} else {
				printClientLog(id, "Read Error: %s", err)
			}
			break
		}
		printClientLog(id, "Received responses: %s.", strResp)
	}
}

func main()  {
	wg.Add(2)
	go serverGo()
	time.Sleep(500 * time.Millisecond)
	go clientGo(1)
	wg.Wait()
}