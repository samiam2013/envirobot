package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/samiam2013/envirobot/bme280"
	"github.com/samiam2013/envirobot/co2"
	"github.com/samiam2013/envirobot/movement"
)

func main() {
	AtmosDataC := make(chan bme280.AtmosData)
	MovementC := make(chan movement.Movement)
	CO2C := make(chan co2.CO2)

	messageC := make(chan []byte, 3)

	go bme280.StreamData(AtmosDataC)
	go movement.StreamMovements(MovementC)
	go co2.StreamLevel(CO2C)

	// stand up the websocket server
	go func() {
		http.Handle("/envirobot",
			websocket.Handler(forwardHandlerBuilder(messageC)))
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	for {
		var atmosData bme280.AtmosData
		var movement movement.Movement
		var co2 co2.CO2

		select {
		case atmosData = <-AtmosDataC:
			if atmosData.Err != nil {
				fmt.Println("AtmosData Error:", atmosData.Err)
			} else {
				fmt.Printf("%+v\n", atmosData)
			}
		case movement = <-MovementC:
			if movement.Err != nil {
				fmt.Println("Movement Error:", movement.Err)
			} else {
				fmt.Printf("%+v\n", movement)
			}
		case co2 = <-CO2C:
			if co2.Err != nil {
				fmt.Println("CO2 Error:", co2.Err)
			} else {
				fmt.Printf("%+v\n", co2)
			}
		}

		// marshal each of the structs into json and send over the messageC channel
		b, err := json.Marshal(atmosData)
		if err != nil {
			fmt.Println("Error marshaling atmosData:", err)
		}
		messageC <- b

		b, err = json.Marshal(movement)
		if err != nil {
			fmt.Println("Error marshaling movement:", err)
		}
		messageC <- b

		b, err = json.Marshal(co2)
		if err != nil {
			fmt.Println("Error marshaling co2:", err)
		}
		messageC <- b
	}
}

func forwardHandlerBuilder(c chan []byte) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			msg := <-c
			log.Print("Forwarding message bytes:", string(msg))
			_, err := io.Copy(ws, bytes.NewReader(msg))
			if err != nil {
				log.Print("Error forwarding message bytes:", err)
			}
		}
	}
}
