package qdp

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

// wire format
type Info struct {
	Version    uint16
	Flags      uint16
	KMI        uint32
	SerialLow  uint32
	SerialHigh uint32
	Memory     [8]uint32
	Interface  [8]uint16
	CalErr     uint16
	SysVer     uint16
}

// decoded format
type Serial struct {
	Version   uint16    `json:"version"`
	KMI       uint32    `json:"kmi"`
	Serial    string    `json:"serial"`
	SysVer    uint16    `json:"sysver"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Serial) String() (string, error) {

	r, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}

	return (string)(r), nil
}

func (p *Ping) Serial() *Serial {

	// check correct ping type
	if p.PingType != 5 {
		return nil
	}

	i := Info{}

	// decode wire format into struct
	binary.Read(bytes.NewReader(p.Data[:]), binary.BigEndian, &i)

	// construct resultant struct
	s := Serial{
		Version:   i.Version,
		KMI:       i.KMI,
		Serial:    fmt.Sprintf("0x0%08x%08x", i.SerialLow, i.SerialHigh),
		SysVer:    i.SysVer,
		Timestamp: time.Now().UTC(),
	}

	// done
	return &s
}
