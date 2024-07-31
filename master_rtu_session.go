package modbus

import (
	"fmt"

	"go.bug.st/serial"
)

// ServerSession tcp server session
type MasterSession struct {
	conn serial.Port
	*serverCommon
	logger
}

// modbus 包处理
func (sf *MasterSession) frameHandler(requestAdu []byte) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("painc happen,%v\n", err)
			sf.Errorf("painc happen,%v", err)
		}
	}()
	fmt.Printf("RX Raw[% x]\n", requestAdu)
	sf.Debugf("RX Raw[% x]", requestAdu)
	// 校验是否符合modbus-rtu协议
	slaveId, content, err := decodeRTUFrame(requestAdu)
	if err != nil {
		sf.Errorf("decodeRTUFrame error:%v", err)
		return err
	}
	fmt.Printf("slaveId:%v, pdu:%v\n", slaveId, content)
	funcCode := content[0]
	pdu := content[1:]
	node, err := sf.GetNode(slaveId)
	if err != nil { // slave id not exit, ignore it
		return nil
	}
	var rspPduData []byte
	if handle, ok := sf.function[funcCode]; ok {
		rspPduData, err = handle(node, pdu)
	} else {
		err = &ExceptionError{ExceptionCodeIllegalFunction}
	}
	if err != nil {
		funcCode |= 0x80
		rspPduData = []byte{err.(*ExceptionError).ExceptionCode}
	}
	sfv := protocolFrame{
		adu: requestAdu,
	}
	responseAdu, err := sfv.encodeRTUFrame(slaveId, ProtocolDataUnit{
		FuncCode: funcCode,
		Data:     rspPduData,
	})
	if err != nil {
		sf.Errorf("encodeRTUFrame error:%v", err)
		return err
	}
	sf.Debugf("TX Raw[% x]", responseAdu)
	// write response
	return func(b []byte) error {
		_, err := sf.conn.Write(b)
		return err
	}(responseAdu)
}
