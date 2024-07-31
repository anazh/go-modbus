package modbus

import (
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
	m.logger.Debugf("connected to %s", m.conn.ComName)
	var tempDelay = minTempDelay // how long to sleep on accept failure
	for {
		buff := make([]byte, rtuAduMaxSize)
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
		buff = buff[:n]
		m.logger.Debugf("received [% x]", buff)
		tempDelay = minTempDelay
		sess := &MasterSession{
			m.conn.port,
			m.serverCommon,
			m.logger,
		}
		sess.frameHandler(buff)
	}
}
