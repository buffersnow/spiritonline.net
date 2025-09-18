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

type UdpServer struct {
	conn *net.UDPConn
	log  log.LoggingFactory
}

type UdpConnection struct {
	server *UdpServer
	Log    log.LoggingFactory
}

type UdpPacket struct {
	Log  log.LoggingFactory
	addr *net.UDPAddr
	conn *net.UDPConn
	data []byte
}

func (n NetUtils) CreateUdpListener(port int) (*UdpServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	logger.Action("UDP Listener", "Listening on 0.0.0.0:%d", port)

	return &UdpServer{
		conn: udpConn,
		log:  logger.Factory("UDP"),
	}, nil
}

func (udp UdpPacket) GetRemoteAddress() string {
	return udp.addr.String()
}

func (udp *UdpConnection) ReadBytes() (*UdpPacket, error) {
	return udp.ReadBytesEx(time.Time{})
}

func (udp *UdpConnection) ReadBytesEx(timeout time.Time) (*UdpPacket, error) {
	if udp.server == nil || udp.server.conn == nil {
		return nil, errors.New("net: UDP server or socket not initialized")
	}

	udp.server.conn.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	n, addr, err := udp.server.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	lastData := buf[:n]

	udp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadBytes", "<IP: %s> Data Length: %d, Data Stream: %s",
		addr.String(), n, prettyBytes(lastData),
	)

	return &UdpPacket{
		addr: addr,
		data: lastData,
	}, nil
}

func (udp *UdpPacket) WriteBytes(data []byte) error {
	if udp.addr == nil {
		return errors.New("net: invalid UDP connection")
	}

	_, err := udp.conn.WriteToUDP(data, udp.addr)
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	udp.Log.Debug(log.DEBUG_TRAFFIC,
		"WriteBytes", "<IP: %s> Data Length: %d, Data Stream: %s",
		udp.addr.String(), len(data), prettyBytes(data),
	)

	return nil
}

func (udp *UdpConnection) ReadText() (*UdpPacket, error) {
	return udp.ReadTextEx(time.Time{})
}

func (udp *UdpConnection) ReadTextEx(timeout time.Time) (*UdpPacket, error) {
	if udp.server == nil || udp.server.conn == nil {
		return nil, errors.New("net: UDP server or socket not initialized")
	}

	udp.server.conn.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	n, addr, err := udp.server.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("net: %w", err)
	}

	lastData := buf[:n]
	retstr := strings.TrimRight(string(lastData), "\r\n")

	udp.Log.Debug(log.DEBUG_TRAFFIC,
		"ReadText", "<IP: %s> Data Length: %d, Data Stream: %s",
		addr.String(), n, retstr,
	)

	return &UdpPacket{
		Log:  udp.Log,
		addr: addr,
		data: []byte(retstr),
	}, nil
}

func (udp *UdpPacket) WriteText(data string) error {
	if udp.addr == nil {
		return errors.New("net: invalid UDP connection")
	}

	_, err := udp.conn.WriteToUDP([]byte(data), udp.addr)
	if err != nil {
		return fmt.Errorf("net: %w", err)
	}

	udp.Log.Debug(log.DEBUG_TRAFFIC,
		"WriteText", "<IP: %s> Data Length: %d, Data Stream: %s",
		udp.addr.String(), len(data), strings.TrimRight(data, "\r\n"),
	)

	return nil
}
