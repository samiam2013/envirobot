package bme280

import (
	"time"

	"golang.org/x/exp/io/i2c"

	"github.com/quhar/bme280"
)

// AtmosData is a struct that holds the data from the BME280 sensor
type AtmosData struct {
	TempC    float64
	PressBar float64
	HumPerc  float64
	Err      error
}

func StreamData(c chan AtmosData) {
	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		c <- AtmosData{Err: err}
	}

	b := bme280.New(d)
	err = b.Init()
	if err != nil {
		c <- AtmosData{Err: err}
	}

	// fmt.Printf("Temp: %fC, Press: %fhPa, Hum: %f%%\n", t, p, h)
	latestData := AtmosData{}
	for {
		t, p, h, err := b.EnvData()
		if err != nil {
			c <- AtmosData{Err: err}
		}

		data := AtmosData{TempC: t, PressBar: p, HumPerc: h, Err: nil}
		// fmt.Println("data:", data)
		if data != latestData {
			c <- data
			latestData = data
		}
		time.Sleep(1 * time.Second)
	}
}
