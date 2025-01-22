package modbus

import (
	"go.bug.st/serial"
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
		logger:       newLogger("modbusRTUServer => "),
	}
}
func recvData(port serial.Port) (recvData []byte, err error) {
reget:
	buff := make([]byte, 1024)
	n, err := port.Read(buff)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return
	}
	if n > 0 {
		recvData = append(recvData, buff[:n]...)
		goto reget
	}
	return recvData, nil
}

func (m *RtuMaster) Connect() error {
	if err := m.conn.Connect(); err != nil {
		return err
	}
	for {
		recvData, err := recvData(m.conn.port)
		if err != nil {
			m.logger.Errorf("recvData error:%v", err)
			return err
		}
		if len(recvData) == 0 {
			continue
		}
		sess := &MasterSession{
			m.conn.port,
			m.serverCommon,
			m.logger,
		}
		sess.frameHandler(recvData)
	}
	// for {
	// 	allLen := 0
	// reget:
	// 	newBuff := make([]byte, rtuAduMaxSize)
	// 	n, err := m.conn.port.Read(newBuff)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if n != 0 {
	// 		allLen += n
	// 		buff = append(buff, newBuff[:n]...)
	// 		if allLen < rtuAduMinSize {
	// 			goto reget
	// 		}
	// 	} else if allLen == 0 {
	// 		continue
	// 	}
	// 	// ---------------------------
	// 	data := buff[:allLen]
	// 	buff = []byte{}
	// 	sess := &MasterSession{
	// 		m.conn.port,
	// 		m.serverCommon,
	// 		m.logger,
	// 	}
	// 	sess.frameHandler(data)
	// }
}
