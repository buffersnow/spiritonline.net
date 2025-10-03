package protocol

//$ https://github.com/devzspy/GameSpy-Openspy-Core/blob/master/qr/Client.h

const (
	QR2Command_Query            byte = 0x00
	QR2Command_Challenge        byte = 0x01
	QR2Command_Echo             byte = 0x02
	QR2Command_Heartbeat        byte = 0x03
	QR2Command_AddError         byte = 0x04
	QR2Command_EchoResponse     byte = 0x05
	QR2Command_ClientMessage    byte = 0x06
	QR2Command_ClientMessageAck byte = 0x07
	QR2Command_Keepalive        byte = 0x08
	QR2Command_Available        byte = 0x09
	QR2Command_ClientRegistered byte = 0x0A
)

const (
	QR1Command_Heartbeat = "heartbeat"
	QR1Command_Validate  = "validate"
	QR1Command_Echo      = "echo"
)
