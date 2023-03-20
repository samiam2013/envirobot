package main

import (
	"fmt"

	"github.com/samiam2013/envirobot/bme280"
	"github.com/samiam2013/envirobot/co2"
	"github.com/samiam2013/envirobot/movement"
)

func main() {
	AtmosDataC := make(chan bme280.AtmosData)
	MovementC := make(chan movement.Movement)
	CO2C := make(chan co2.CO2)

	go bme280.StreamData(AtmosDataC)
	go movement.StreamMovements(MovementC)
	go co2.StreamLevel(CO2C)

	for {
		select {
		case atmosData := <-AtmosDataC:
			if atmosData.Err != nil {
				fmt.Println("AtmosData Error:", atmosData.Err)
			} else {
				fmt.Printf("%+v\n", atmosData)
			}
		case movement := <-MovementC:
			if movement.Err != nil {
				fmt.Println("Movement Error:", movement.Err)
			} else {
				fmt.Printf("%+v\n", movement)
			}
		case co2 := <-CO2C:
			if co2.Err != nil {
				fmt.Println("CO2 Error:", co2.Err)
			} else {
				fmt.Printf("%+v\n", co2)
			}
		}
	}

}
