package qdp

import "time"

func ReadSerial(ipaddr string, ipport string, timeout time.Duration) (*Serial, error) {
	p, err := NewSerial().Send(ipaddr, ipport, timeout)
	if err != nil {
		return nil, err
	}
	s := p.Serial()
	if s == nil {
		return nil, nil
	}
	return s, nil
}

func ReadSOH(ipaddr string, ipport string, timeout time.Duration) (*SOH, error) {
	p, err := NewStatus().Send(ipaddr, ipport, timeout)
	if err != nil {
		return nil, err
	}
	s := p.Status()
	if s == nil {
		return nil, nil
	}
	return s, nil
}
