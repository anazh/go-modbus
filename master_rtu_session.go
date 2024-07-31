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
			sf.Errorf("painc happen,%v", err)
		}
	}()
	sf.Debugf("RX Raw[% x]", requestAdu)
	// 校验是否符合modbus-rtu协议
	slaveId, pdu, err := decodeRTUFrame(requestAdu)
	if err != nil {
		sf.Errorf("decodeRTUFrame error:%v", err)
		return err
	}
	fmt.Printf("slaveId:%v, pdu:%v\n", slaveId, pdu)
	funcCode := pdu[0]
	node, err := sf.GetNode(slaveId)
	if err != nil { // slave id not exit, ignore it
		return nil
	}
	var rspPduData []byte
	if handle, ok := sf.function[funcCode]; ok {
		rspPduData, err = handle(node, pdu)
		fmt.Printf("rspPduData:%v err:%v\n", rspPduData, err)
	} else {
		err = &ExceptionError{ExceptionCodeIllegalFunction}
	}
	fmt.Println("11111")
	if err != nil {
		fmt.Printf(" handle => err:%v\n", err)
		funcCode |= 0x80
		rspPduData = []byte{err.(*ExceptionError).ExceptionCode}
	}
	fmt.Println("2")
	sfv := protocolFrame{
		adu: []byte{},
	}
	responseAdu, err := sfv.encodeRTUFrame(slaveId, ProtocolDataUnit{
		funcCode,
		rspPduData,
	})
	fmt.Println("3")
	fmt.Printf("responseAdu:%v\n", responseAdu)
	if err != nil {
		fmt.Printf("encodeRTUFrame error:%v", err)
		sf.Errorf("encodeRTUFrame error:%v", err)
		return err
	}
	fmt.Printf("response [% x]\n", responseAdu)
	sf.Debugf("TX Raw[% x]", responseAdu)
	// write response
	return func(b []byte) error {
		_, err := sf.conn.Write(b)
		return err
	}(responseAdu)
}
