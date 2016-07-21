package teleinfo

import (
	"github.com/tarm/serial"
)

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
