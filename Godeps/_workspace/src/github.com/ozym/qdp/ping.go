package qdp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"time"
)

// native encoded ping packet sans crc
type Ping struct {
	Command     uint8     `json:"command"`
	Version     uint8     `json:"version"`
	DataLength  uint16    `json:"datalength"`
	Sequence    uint16    `json:"sequence"`
	Acknowledge uint16    `json:"acknowledge"`
	PingType    uint16    `json:"ping_type"`
	PingID      uint16    `json:"ping_id"`
	Data        [532]byte `json:"data"`
}

// NewSerial builds a Ping structure configured for requesting a serial packet
func NewSerial() *Ping {

	p := Ping{
		Command:     56,
		Version:     2,
		DataLength:  4,
		Sequence:    1,
		Acknowledge: 0,
		PingType:    4,
		PingID:      0,
	}

	return &p
}

// NewStatus builds a Ping structure configured for requesting a status packet
func NewStatus() *Ping {

	p := Ping{
		Command:     56,
		Version:     2,
		DataLength:  8,
		Sequence:    1,
		Acknowledge: 0,
		PingType:    2,
		PingID:      0,
	}

	binary.BigEndian.PutUint32(p.Data[0:4], 0x8f0b)

	return &p
}

func (p *Ping) String() (string, error) {

	s, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return (string)(s), nil
}

func (p *Ping) Buffer(crc uint32) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, crc)
	binary.Write(buf, binary.BigEndian, p.Command)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.DataLength)
	binary.Write(buf, binary.BigEndian, p.Sequence)
	binary.Write(buf, binary.BigEndian, p.Acknowledge)
	binary.Write(buf, binary.BigEndian, p.PingType)
	binary.Write(buf, binary.BigEndian, p.PingID)
	if !(p.DataLength < 4) {
		buf.Write(p.Data[0 : p.DataLength-4])
	}

	return buf.Bytes()
}

func (p *Ping) Crc() uint32 {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Command)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.DataLength)
	binary.Write(buf, binary.BigEndian, p.Sequence)
	binary.Write(buf, binary.BigEndian, p.Acknowledge)
	binary.Write(buf, binary.BigEndian, p.PingType)
	binary.Write(buf, binary.BigEndian, p.PingID)
	if !(p.DataLength < 4) {
		buf.Write(p.Data[0 : p.DataLength-4])
	}

	return crc(buf.Bytes())
}

// Decode unwraps the raw wire-format into a Ping structure.
// A nil will be returned if the block fails the CRC check.
func Decode(b []byte) *Ping {

	p := Ping{}
	buf := bytes.NewReader(b)

	var crc uint32

	// decode packet header
	binary.Read(buf, binary.BigEndian, &crc)
	binary.Read(buf, binary.BigEndian, &p.Command)
	binary.Read(buf, binary.BigEndian, &p.Version)
	binary.Read(buf, binary.BigEndian, &p.DataLength)
	binary.Read(buf, binary.BigEndian, &p.Sequence)
	binary.Read(buf, binary.BigEndian, &p.Acknowledge)
	binary.Read(buf, binary.BigEndian, &p.PingType)
	binary.Read(buf, binary.BigEndian, &p.PingID)

	// copy raw data
	copy(p.Data[:], b[16:12+p.DataLength])

	// check crc
	if crc != p.Crc() {
		return nil
	}

	return &p
}

func (p *Ping) Send(ipaddr string, ipport string, timeout time.Duration) (*Ping, error) {
	conn, err := net.Dial("udp", net.JoinHostPort(ipaddr, ipport))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(p.Buffer(p.Crc()))
	if err != nil {
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	b := make([]byte, 512)
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}

	return Decode(b), nil
}
