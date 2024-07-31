package modbus

import (
	"fmt"
	"net"
	"time"
)

type RtuMaster struct {
	conn *RTUClientProvider
	*serverCommon
	logger
}

func NewRtuMaster(r *RTUClientProvider) *RtuMaster {
	return &RtuMaster{
		conn:         r,
		serverCommon: newServerCommon(),
		logger:       newLogger("modbusTCPServer => "),
	}
}

func (m *RtuMaster) Connect() error {
	if err := m.conn.Connect(); err != nil {
		return err
	}
	if err := m.conn.port.SetReadTimeout(m.conn.TimeOut); err != nil {
		return err
	}
	fmt.Println("Connected to ", m.conn.ComName, m.conn.TimeOut)
	var tempDelay = minTempDelay // how long to sleep on accept failure
	buff := make([]byte, rtuAduMaxSize)
	for {
		n, err := m.conn.port.Read(buff)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				tempDelay <<= 1
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		if n == 0 {
			continue
		}
		newBuff := buff[:n]
		buff = make([]byte, rtuAduMaxSize)
		fmt.Printf("received [% x]\n", newBuff)
		tempDelay = minTempDelay
		sess := &MasterSession{
			m.conn.port,
			m.serverCommon,
			m.logger,
		}
		sess.frameHandler(newBuff)
	}
}
