package co2

import (
	"time"

	"github.com/tarm/serial"
)

type CO2 struct {
	PPM int   `json:"ppm"`
	Err error `json:"err"`
}

func StreamLevel(c chan CO2) {
	config := &serial.Config{
		Name:        "/dev/serial0",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		c <- CO2{Err: err}
	}
	defer stream.Close()

	latestValue := 0
	for {
		// send the 'gas concentration' command to get the current reading
		if _, err := stream.Write([]byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}); err != nil {
			c <- CO2{Err: err}
		}
		response := make([]byte, 9)
		if _, err = stream.Read(response); err != nil {
			c <- CO2{Err: err}
		}

		if !checksumValidate(response) {
			c <- CO2{Err: err}
		}

		ppm := int(response[2])*256 + int(response[3])
		// logrus.Info("Reported concentration:", ppm, "ppm.")
		if latestValue != ppm {
			latestValue = ppm
			c <- CO2{PPM: ppm, Err: nil}
		}

		time.Sleep(1 * time.Second)
	}

}

func checksumValidate(b []byte) bool {
	var sum byte
	if len(b) != 9 {
		return false
	}
	for i := 1; i < len(b)-1; i++ {
		sum += b[i]
	}
	sum = 0xFF - sum
	return (sum + 1) == b[8]
}
