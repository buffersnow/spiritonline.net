package net

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"github.com/luxploit/red"
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

func (n NetUtils) CreateTcpListener(port int) (TcpServer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return TcpServer{}, fmt.Errorf("net: %w", err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return TcpServer{}, fmt.Errorf("net: %w", err)
	}

	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return TcpServer{}, fmt.Errorf("net: %w", err)
	}

	logger.Action("TCP Listener", "Listening on 0.0.0.0:%d", port)
	return TcpServer{
		conn: tcpListener,
		log:  logger.Factory("TCP"),
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
	if tcp.server == nil || tcp.client == nil {
		return errors.New("net: TCP server or client not initialized")
	}

	if _, err := tcp.client.Write([]byte(data)); err != nil {
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
	if tcp.server == nil || tcp.client == nil {
		return "", errors.New("net: TCP server or client not initialized")
	}

	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	n, err := tcp.client.Read(buf)

	if err != nil {
		return "", fmt.Errorf("net: %w", err)
	}

	lastData := buf[:n]
	retstr := strings.TrimRight(string(lastData), "\r\n")

	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadText", "Data Length: %d, Data Stream: %s",
		len(lastData), retstr,
	)

	return retstr, nil
}

func (tcp TcpConnection) WriteBytes(data []byte) error {
	if tcp.server == nil || tcp.client == nil {
		return errors.New("net: TCP server or client not initialized")
	}

	if _, err := tcp.client.Write(data); err != nil {
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
	if tcp.server == nil || tcp.client == nil {
		return nil, errors.New("net: TCP server or client not initialized")
	}

	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	n, err := tcp.client.Read(buf)

	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	lastData := buf[:n]

	tcp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadBytes", "Data Length: %d, Data Stream: %s",
		len(lastData), prettyBytes(lastData),
	)

	return lastData, nil
}

func (tcp *TcpConnection) Close() error {
	if err := tcp.client.Close(); err != nil {
		return fmt.Errorf("net: %w", err)
	}

	return nil
}
