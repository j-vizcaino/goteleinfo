package teleinfo

import (
	"github.com/tarm/serial"
)

// OpenPort opens and configures the serial port to read Teleinfo frames from.
func OpenPort(serialDevice string) (*serial.Port, error) {
	cfg := &serial.Config{
		Name:     serialDevice,
		Baud:     1200,
		Size:     7,
		Parity:   serial.ParityEven,
		StopBits: serial.Stop1,
	}
	return serial.OpenPort(cfg)
}
