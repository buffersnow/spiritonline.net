package net

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
)

type TcpServer struct {
	conn *net.TCPListener
	log  log.LoggingFactory
}

type TcpConnection struct {
	server *net.TCPListener
	client *net.TCPConn
	Log    log.LoggingFactory
}

func CreateTcpListener(port int) (TcpServer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return TcpServer{}, fmt.Errorf("net: %w", err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return TcpServer{}, fmt.Errorf("net: %w", err)
	}

	log.Global().Info("TCP Listener", "Listening on 0.0.0.0:%d", port)
	return TcpServer{
		conn: tcpListener,
		log:  log.Factory("TCP"),
	}, nil
}

// func CreateTcpConnection(server string, port int) (*TcpConnection, error) {
// 	serverAddr := fmt.Sprintf("%s:%d", server, port)
// 	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
// 	if err != nil {
// 		return nil, fmt.Errorf("net: %w", err)
// 	}
//
// 	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	if err != nil {
// 		return nil, fmt.Errorf("net: %w", err)
// 	}
//
// 	iLog := log.FactoryWithPostfix("TCP",
// 		fmt.Sprintf("<IP: %s>", tcpConn.RemoteAddr().String()),
// 	)
//
// 	iLog.
//
// 	return &TcpConnection{
// 		client: tcpConn,
// 		Log:    iLog,
// 	}, nil
// }

func (tcp *TcpServer) Accept() (*TcpConnection, error) {

	if tcp.conn == nil {
		return nil, errors.New("net: invalid call to accept incoming from tcp server")
	}

	lst, err := tcp.conn.AcceptTCP()
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	cli := &TcpConnection{
		server: tcp.conn,
		client: lst,
		Log:    tcp.log,
	}

	cli.Log.ChangePostfix("<IP: %s>", cli.GetRemoteAddress())

	return cli, nil
}

func (tcp TcpConnection) GetRemoteAddress() string {
	return tcp.client.RemoteAddr().String()
}

func (tcp TcpConnection) WriteText(data string) error {
	_, err := tcp.client.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"WriteText", "Data Length: %d, Data Stream: %s",
		len(data), strings.TrimRight(data, "\r\n"),
	)

	return nil
}

func (tcp TcpConnection) ReadText() (data string, err error) {
	return tcp.ReadTextEx(time.Time{})
}

func (tcp TcpConnection) ReadTextEx(timeout time.Time) (data string, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		return "", fmt.Errorf("net: %w", err)
	}

	ret := make([]byte, length)
	copy(ret, buf)

	retstr := strings.TrimRight(string(ret), "\r\n")
	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadText", "Data Length: %d, Data Stream: %s",
		len(ret), retstr,
	)

	return retstr, nil
}

func (tcp TcpConnection) WriteBytes(data []byte) error {
	_, err := tcp.client.Write(data)
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"WriteBytes", "Data Length: %d, Data Stream: %s",
		len(data), prettyBytes(data),
	)

	return nil
}

func (tcp TcpConnection) ReadBytes() (data []byte, err error) {
	return tcp.ReadBytesEx(time.Time{})
}

func (tcp TcpConnection) ReadBytesEx(timeout time.Time) (data []byte, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	ret := make([]byte, length)
	copy(ret, buf)

	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadBytes", "Data Length: %d, Data Stream: %s",
		len(ret), prettyBytes(ret),
	)

	return ret, nil
}

func (tcp *TcpConnection) Close() error {
	err := tcp.client.Close()
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	return nil
}
