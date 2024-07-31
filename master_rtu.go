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
	fmt.Println("Connected to ", m.conn.ComName, m.conn.TimeOut)
	var tempDelay = minTempDelay // how long to sleep on accept failure
	buff := []byte{}
	for {
		allLen := 0
	reget:
		newBuff := make([]byte, rtuAduMaxSize)
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
		if n != 0 {
			allLen += n
			buff = append(buff, newBuff[:n]...)
			goto reget
		}

		data := buff[:allLen]
		buff = []byte{}
		fmt.Printf("received [% x]\n", data)
		tempDelay = minTempDelay
		sess := &MasterSession{
			m.conn.port,
			m.serverCommon,
			m.logger,
		}
		sess.frameHandler(data)
	}
}
