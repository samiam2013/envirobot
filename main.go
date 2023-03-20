package main

import (
	"fmt"

	"github.com/samiam2013/envirobot/bme280"
	"github.com/samiam2013/envirobot/co2"
	"github.com/samiam2013/envirobot/movement"
)

func main() {
	fmt.Println("Hello, World!")

	AtmosDataC := make(chan bme280.AtmosData)
	go bme280.StreamData(AtmosDataC)

	MovementC := make(chan movement.Movement)
	go movement.StreamMovements(MovementC)

	CO2C := make(chan co2.CO2)
	go co2.StreamLevel(CO2C)

	for {
		select {
		case atmosData := <-AtmosDataC:
			if atmosData.Err != nil {
				fmt.Println("AtmosDataC error:", atmosData.Err)
			} else {
				fmt.Println("AtmosDataC:", atmosData)
			}
		case movement := <-MovementC:
			if movement.Err != nil {
				fmt.Println("MovementC error:", movement.Err)
			} else {
				fmt.Println("MovementC:", movement)
			}
		case co2 := <-CO2C:
			if co2.Err != nil {
				fmt.Println("CO2C error:", co2.Err)
			} else {
				fmt.Println("CO2C:", co2)
			}
		}
	}

}
