package modbus

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
	m.conn.mu.Lock()
	defer m.conn.mu.Unlock()
	err := m.conn.serialPort.connect()
	if err != nil {
		return err
	}
	

}
