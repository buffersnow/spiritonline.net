package com

import "github.com/vmihailenco/msgpack/v5"

//% RegionID Tag Format:
//% us-east-test-2
//% + Country Code
//%   + Cardinal Region (+ Central & Global)
//%     + Service Name
//%       + Registration Index

type MessageHeader struct {
	Version     int    `msgpack:"version"`
	OpCode      string `msgpack:"opcode"`
	Sender      string `msgpack:"sender"`         // uses RegionId
	Forwardee   string `msgpack:"forwarded_from"` // uses RegionId - set to Sender unless forwarded
	Receiver    string `msgpack:"receiver"`       // uses RegionId
	ForwardedIP string `msgpack:"forwarded_ip"`   // If receiver is not router
}

type UnsignedMessage struct {
	Header  MessageHeader      `msgpack:"header"`
	Payload msgpack.RawMessage `msgpack:"payload"` // msgpack-encoded message
}

type SignedMessage struct {
	Header    MessageHeader      `msgpack:"header"`
	Payload   msgpack.RawMessage `msgpack:"payload"`         // msgpack-encoded message
	Signature []byte             `msgpack:"ecdsa_signature"` // SHA256 Hash
}
