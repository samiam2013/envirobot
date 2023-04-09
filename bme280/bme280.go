package bme280

import (
	"time"

	"golang.org/x/exp/io/i2c"

	"github.com/quhar/bme280"
)

// AtmosData is a struct that holds the data from the BME280 sensor
type AtmosData struct {
	TempCelcius float64
	PressHPa    float64
	Humidity    float64
	Err         error
}

func StreamData(c chan AtmosData) {
	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		c <- AtmosData{Err: err}
	}

	b := bme280.New(d)
	if err := b.Init(); err != nil {
		c <- AtmosData{Err: err}
	}

	latestData := AtmosData{}
	for {
		t, p, h, err := b.EnvData()
		if err != nil {
			c <- AtmosData{Err: err}
		}

		data := AtmosData{TempCelcius: t, PressHPa: p, Humidity: h, Err: nil}
		// prssure is sensitive to sound waves probably, so only send if temp or humidity changes
		if data.TempCelcius != latestData.TempCelcius || data.Humidity != latestData.Humidity {
			c <- data
			latestData = data
		}
		time.Sleep(1 * time.Second)
	}
}
