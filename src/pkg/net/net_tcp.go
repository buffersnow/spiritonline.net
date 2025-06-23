package net

import (
	"fmt"
	"net"
	"strings"
	"time"

	"buffersnow.com/spiritonline/pkg/util"
)

type TcpConnection struct {
	server *net.TCPListener
	client *net.TCPConn
	Log    *util.LogFactory
}

func CreateTcpListener(port int) (TcpConnection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return TcpConnection{}, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return TcpConnection{}, err
	}

	iLog := util.Log.Instance("TCP")
	iLog.Info("Listener", "Listening on 0.0.0.0:%d", port)

	return TcpConnection{
		server: tcpListener,
		Log:    iLog,
	}, nil
}

func CreateTcpConnection(server string, port int) (TcpConnection, error) {
	serverAddr := fmt.Sprintf("%s:%d", server, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return TcpConnection{}, err
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return TcpConnection{}, err
	}

	iLog := util.Log.Instance("TCP")
	iLog.ChangePostfix("<IP: %s>", tcpConn.RemoteAddr().String())

	return TcpConnection{
		client: tcpConn,
		Log:    iLog,
	}, nil
}

func (tcp *TcpConnection) AcceptIncoming() error {
	lst, err := tcp.server.AcceptTCP()
	tcp.client = lst
	tcp.Log.ChangePostfix("<IP: %s>", tcp.GetRemoteAddress())
	return err
}

func (tcp TcpConnection) GetRemoteAddress() string {
	return tcp.client.RemoteAddr().String()
}

func (tcp TcpConnection) WriteText(data string) error {
	_, err := tcp.client.Write([]byte(data))
	tcp.Log.Debug(util.LOG_DEBUG_TRAFFIC,
		"WriteText", "Data Length: %d, Data Stream: %s",
		len(data), strings.TrimRight(data, "\r\n"),
	)
	return err
}

func (tcp TcpConnection) ReadText() (data string, err error) {
	return tcp.ReadTextEx(time.Time{})
}

func (tcp TcpConnection) ReadTextEx(timeout time.Time) (data string, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		return "", err
	}

	ret := make([]byte, length)
	copy(ret, buf)

	retstr := strings.TrimRight(string(ret), "\r\n")
	tcp.Log.Debug(util.LOG_DEBUG_TRAFFIC,
		"ReadText", "Data Length: %d, Data Stream: %s",
		len(ret), retstr,
	)

	return retstr, err
}

func (tcp TcpConnection) WriteBytes(data []byte) error {
	_, err := tcp.client.Write(data)
	tcp.Log.Debug(util.LOG_DEBUG_TRAFFIC,
		"WriteBytes", "Data Length: %d, Data Stream: %s",
		len(data), util.PrettyBytes(data),
	)
	return err
}

func (tcp TcpConnection) ReadBytes() (data []byte, err error) {
	return tcp.ReadBytesEx(time.Time{})
}

func (tcp TcpConnection) ReadBytesEx(timeout time.Time) (data []byte, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		return nil, err
	}

	ret := make([]byte, length)
	copy(ret, buf)

	tcp.Log.Debug(util.LOG_DEBUG_TRAFFIC,
		"ReadBytes", "Data Length: %d, Data Stream: %s",
		len(ret), util.PrettyBytes(ret),
	)

	return ret, err
}

func (tcp *TcpConnection) Close() error {
	return tcp.client.Close()
}
