package movement

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type Movement struct {
	Time time.Time
	Err  error
}

func StreamMovements(c chan Movement) {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		c <- Movement{Err: err}
	}

	// Lookup a pin by its number:
	pName := "GPIO17"
	p := gpioreg.ByName(pName)
	if p == nil {
		c <- Movement{Err: fmt.Errorf("failed to find %s", pName)}
	}

	// TODO figure out how to fix the linter deprecation warning here
	// fmt.Printf("%s: %s\n", p, p.Function())

	// Set it as input, with an internal pull down resistor:
	if err := p.In(gpio.PullDown, gpio.BothEdges); err != nil {
		c <- Movement{Err: err}
	}

	// Wait for edges as detected by the hardware, and print the value read:
	for {
		p.WaitForEdge(-1)
		res := p.Read()
		if res == gpio.High {
			// fmt.Printf("triggered %s\n",
			// 	time.Now().Format(time.RFC3339))
			c <- Movement{Time: time.Now(), Err: nil}
		}
		time.Sleep(1 * time.Minute)
	}
}
