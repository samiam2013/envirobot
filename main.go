package main

import (
	"database/sql"
	"log"

	"github.com/samiam2013/envirobot/bme280"
	"github.com/samiam2013/envirobot/co2"
	"github.com/samiam2013/envirobot/movement"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./sqlite.db?_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

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
				log.Println("AtmosData Error:", atmosData.Err)
			} else {
				log.Printf("%+v\n", atmosData)
				upload(db, atmosData)
			}
		case movement := <-MovementC:
			if movement.Err != nil {
				log.Println("Movement Error:", movement.Err)
			} else {
				log.Printf("%+v\n", movement)
				upload(db, movement)
			}
		case co2 := <-CO2C:
			if co2.Err != nil {
				log.Println("CO2 Error:", co2.Err)
			} else {
				log.Printf("%+v\n", co2)
				upload(db, co2)
			}
		}
	}
}

func upload(db *sql.DB, data interface{}) {
	switch data := data.(type) {
	case bme280.AtmosData:
		stmt, err := db.Prepare("INSERT INTO bme280 (temperature, pressure, humidity) VALUES (?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		if _, err = stmt.Exec(data.TempCelcius, data.PressHPa, data.Humidity); err != nil {
			log.Print(err)
		}
	case co2.CO2:
		stmt, err := db.Prepare("INSERT INTO co2 (ppm) VALUES (?)")
		if err != nil {
			log.Print(err)
		}
		if _, err = stmt.Exec(data.PPM); err != nil {
			log.Print(err)
		}
	case movement.Movement:
		stmt, err := db.Prepare("INSERT INTO movement (created_at) VALUES (?)")
		if err != nil {
			log.Print(err)
		}
		if _, err = stmt.Exec(data.Time); err != nil {
			log.Print(err)
		}
	}
}
