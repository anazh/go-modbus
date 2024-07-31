package modbus

import (
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
)

// SerialDefaultTimeout Serial Default timeout
const SerialDefaultTimeout = 1 * time.Second

// serialPort has configuration and I/O controller.
type serialPort struct {
	ComName string
	// Serial port configuration.
	serial.Mode
	mu      sync.Mutex
	port    serial.Port
	TimeOut time.Duration
}

// Connect try to connect the remote server
func (sf *serialPort) Connect() (err error) {
	if sf.ComName == "" {
		return fmt.Errorf("serial port name is empty")
	}
	sf.mu.Lock()
	err = sf.connect()
	sf.mu.Unlock()
	return
}

// Caller must hold the mutex before calling this method.
func (sf *serialPort) connect() error {
	if sf.port == nil {
		port, err := serial.Open(sf.ComName, &sf.Mode)
		if err != nil {
			return fmt.Errorf("failed to open serial port: %s err %v", sf.ComName, err)
		}
		port.SetReadTimeout(sf.TimeOut)
		sf.port = port
	}
	return nil
}

// IsConnected returns a bool signifying whether the client is connected or not.
func (sf *serialPort) IsConnected() (b bool) {
	sf.mu.Lock()
	b = sf.port != nil
	sf.mu.Unlock()
	return b
}

// setSerialConfig set serial config
func (sf *serialPort) setSerialConfig(commName string, config serial.Mode) {
	sf.Mode = config
	sf.ComName = commName
}

func (sf *serialPort) setTimeout(t time.Duration) {
	sf.TimeOut = t
}

func (sf *serialPort) close() (err error) {
	if sf.port != nil {
		err = sf.port.Close()
		sf.port = nil
	}
	return err
}

// Close close current connection.
func (sf *serialPort) Close() (err error) {
	sf.mu.Lock()
	err = sf.close()
	sf.mu.Unlock()
	return
}
